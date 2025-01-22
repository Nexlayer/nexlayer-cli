package commands

// CI command variables
var (
	Stack        string
	Registry     string
	ImageName    string
	ImageTag     string
	BuildContext string
	Token        string
)

// Service command variables
var (
	AppName      string
	Service      string
	OutputFormat string
	OutputFile   string
	EnvPairs     []string
)
