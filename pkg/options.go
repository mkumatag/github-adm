package pkg

var (
	GlobalOptions     = &globalOptions{}
	SyncLabelsOptions = &syncLabelsOptions{}
)

type globalOptions struct {
	ApiKey    string
	BaseURL   string
	UploadURL string
}

type syncLabelsOptions struct {
	Manifest        string
	Org             string
	Repo            string
	DeleteOutOfSync bool
}
