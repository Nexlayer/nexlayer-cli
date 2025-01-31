package template

type Config struct {
	Name        string
	BuildCmd    string
	OutputDir   string
	DefaultPods []PodConfig
}

type PodConfig struct {
	Type       string
	Name       string
	Tag        string
	ExposeHttp bool
	Vars       []VarPair
}

type VarPair struct {
	Key   string
	Value string
}

var Registry = map[string]Config{
	"mern": {
		Name:      "MERN Stack",
		BuildCmd:  "npm install && npm run build",
		OutputDir: "build",
		DefaultPods: []PodConfig{
			{
				Type:       "mongodb",
				Name:       "database",
				Tag:        "mongo:4.4",
				ExposeHttp: false,
				Vars: []VarPair{
					{Key: "MONGO_INITDB_ROOT_USERNAME", Value: "root"},
					{Key: "MONGO_INITDB_ROOT_PASSWORD", Value: "{{password}}"},
				},
			},
			{
				Type:       "express",
				Name:       "backend",
				Tag:        "node:18",
				ExposeHttp: true,
			},
		},
	},
	"nextjs": {
		Name:      "Next.js",
		BuildCmd:  "npm install && npm run build",
		OutputDir: ".next",
		DefaultPods: []PodConfig{
			{
				Type:       "nginx",
				Name:       "web",
				Tag:        "nginx:1.23",
				ExposeHttp: true,
			},
		},
	},
}

func GetTemplate(name string) (Config, bool) {
	cfg, exists := Registry[name]
	return cfg, exists
}
