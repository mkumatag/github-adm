package sync_labels

import (
	"encoding/json"
	"github.com/google/go-github/v33/github"
	"github.com/mkumatag/github-adm/pkg"
	"github.com/mkumatag/github-adm/pkg/client"
	"github.com/spf13/cobra"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/url"
	"sync"
)

type Label struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description,omitempty"`
}

const (
	delete = "DELETE"
)

var DeleteOutOfSyncActions = map[bool]string{
	true:  "delete",
	false: "none",
}

func searchManifestLabels(name string, labels []Label) *Label {
	for i := range labels {
		if name == labels[i].Name {
			return &labels[i]
		}
	}
	return nil
}

func searchGitHubLabels(name string, labels []*github.Label) *github.Label {
	for i := range labels {
		if name == *labels[i].Name {
			return labels[i]
		}
	}
	return nil
}

var Cmd = &cobra.Command{
	Use:   "sync-labels",
	Short: "Sync the labels",
	Long: `
examples:

  GH_TOKEN=<GH_TOKEN> github-adm sync-labels --base-url https://github.ibm.com/api/v3 --upload-url https://uploads.github.ibm.com/ --org org --repo repo --manifest labels.json --delete-out-of-sync`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		globalOpt := pkg.GlobalOptions
		opt := pkg.SyncLabelsOptions
		klog.Infof("Manifest file is : %s", opt.Manifest)
		manifest, err := ioutil.ReadFile(opt.Manifest)
		if err != nil {
			return err
		}
		var labels []Label
		err = json.Unmarshal(manifest, &labels)
		if err != nil {
			return err
		}
		klog.Infof("manifest: %v", labels)
		gh, err := client.NewGithub(globalOpt.BaseURL, globalOpt.UploadURL, globalOpt.ApiKey)
		if err != nil {
			return err
		}
		repoLabels, err := gh.ListLabels(opt.Org, opt.Repo)
		if err != nil {
			return err
		}

		klog.Infof("Deleting the labels from the repo which is not sync with the manifest file")
		var wg sync.WaitGroup
		for i := range repoLabels {
			if l := searchManifestLabels(*repoLabels[i].Name, labels); l == nil {
				klog.Infof("out of sync: label not found in the manifest, %s action taken on the label: %s", DeleteOutOfSyncActions[opt.DeleteOutOfSync], *repoLabels[i].Name)
				if opt.DeleteOutOfSync {
					wg.Add(1)
					go func(wg *sync.WaitGroup, label string) {
						defer wg.Done()
						klog.Infof("deleting label: %s", label)
						resp, err := gh.DeleteLabel(opt.Org, opt.Repo, url.QueryEscape(label))
						if err != nil {
							klog.Errorf("failed to delete label: %s, err: %v", label, err)
                                                        return
						}
						klog.Infof("deleted label: %s, %v %v", label, resp.StatusCode, resp.Body)
					}(&wg, *repoLabels[i].Name)
				}
			}
		}

		klog.Infof("Deleting the labels from the repo which is not sync with the manifest file.. DONE")

		klog.Infof("Adding any missing label from manifest file to repository")
		for i := range labels {
			if l := searchGitHubLabels(labels[i].Name, repoLabels); l == nil {
				klog.Infof("out of sync: label not found in the github repo, create action taken on the label: %s", labels[i].Name)
				wg.Add(1)
				label := github.Label{
					Name:  &labels[i].Name,
					Color: &labels[i].Color,
				}
				go func(wg *sync.WaitGroup, label *github.Label) {
					defer wg.Done()
					klog.Infof("adding label: %s", label)
					_, _, err := gh.CreateLabel(opt.Org, opt.Repo, label)
					if err != nil {
						klog.Errorf("failed to delete label: *repoLabels[i].Name, err: %v", err)
					}
					klog.Infof("added label: %s", label)
				}(&wg, &label)
			}
		}
		klog.Infof("Adding any missing label from manifest file to repository... DONE")

		wg.Wait()
		return nil
	},
}

func init() {
	Cmd.Flags().StringVar(&pkg.SyncLabelsOptions.Org, "org", "", "GH Organization/User")
	Cmd.Flags().StringVar(&pkg.SyncLabelsOptions.Repo, "repo", "", "GH repository")
	Cmd.Flags().StringVar(&pkg.SyncLabelsOptions.Manifest, "manifest", "", "JSON manifest file")
	Cmd.Flags().BoolVar(&pkg.SyncLabelsOptions.DeleteOutOfSync, "delete-out-of-sync", false, "Delete the labels which aren't present in the manifest JSON")
	Cmd.Flags().SortFlags = false
	Cmd.PersistentFlags().SortFlags = false
}
