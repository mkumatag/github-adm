package cmd

import (
	"os"

	synclabels "github.com/mkumatag/github-adm/cmd/sync-labels"
	"github.com/mkumatag/github-adm/pkg"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var rootCmd = &cobra.Command{
	Use:   "github-adm",
	Short: "github-adm command line tool",
	Long: "",
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if token := os.Getenv("GH_TOKEN"); token != "" {
			pkg.GlobalOptions.ApiKey = token
		}
	},
}

func init() {
	rootCmd.AddCommand(synclabels.Cmd)

	rootCmd.PersistentFlags().StringVar(&pkg.GlobalOptions.BaseURL, "base-url", "", "GH Base URL if enterprise account")
	rootCmd.PersistentFlags().StringVar(&pkg.GlobalOptions.UploadURL, "upload-url", "", "GH Upload URL if enterprise account")
	rootCmd.PersistentFlags().StringVarP(&pkg.GlobalOptions.ApiKey, "gh-token", "t", "", "GH Token(env: GH_TOKEN)")
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		klog.Errorln(err)
		os.Exit(1)
	}
}
