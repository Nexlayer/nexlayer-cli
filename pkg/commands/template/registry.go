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

type TemplateCategory struct {
	Name        string
	Description string
	Templates   []Template
}

type Template struct {
	ID          string
	Name        string
	Description string
	Type        string
}

func GetTemplate(name string) (Config, bool) {
	cfg, exists := Registry[name]
	return cfg, exists
}

func GetCategories() []TemplateCategory {
	return []TemplateCategory{
		{
			Name:        "Web Applications",
			Description: "Templates for web applications",
			Templates: []Template{
				{
					ID:          "react",
					Name:        "React",
					Description: "React web application",
					Type:        "frontend",
				},
				{
					ID:          "vue",
					Name:        "Vue.js",
					Description: "Vue.js web application",
					Type:        "frontend",
				},
				{
					ID:          "angular",
					Name:        "Angular",
					Description: "Angular web application",
					Type:        "frontend",
				},
			},
		},
		{
			Name:        "Backend Services",
			Description: "Templates for backend services",
			Templates: []Template{
				{
					ID:          "express",
					Name:        "Express.js",
					Description: "Express.js backend service",
					Type:        "backend",
				},
				{
					ID:          "fastapi",
					Name:        "FastAPI",
					Description: "FastAPI backend service",
					Type:        "backend",
				},
				{
					ID:          "django",
					Name:        "Django",
					Description: "Django backend service",
					Type:        "backend",
				},
			},
		},
		{
			Name:        "Machine Learning",
			Description: "Templates for machine learning applications",
			Templates: []Template{
				{
					ID:          "kubeflow",
					Name:        "Kubeflow",
					Description: "Kubeflow pipeline template",
					Type:        "llm",
				},
			},
		},
	}
}

func GetTemplateByID(id string) *Template {
	categories := GetCategories()
	for _, category := range categories {
		for _, template := range category.Templates {
			if template.ID == id {
				return &template
			}
		}
	}
	return nil
}
