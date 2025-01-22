package vars

// Common variables
var (
	AppName     string
	ServiceName string
	EnvVars     []string
	APIURL      string
)

// CI variables
var (
	BuildContext string
	ImageTag     string
)

// Graph variables
var (
	Depth        int
	OutputFormat string
	OutputFile   string
	EnvPairs     []string
)
