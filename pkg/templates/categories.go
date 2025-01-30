package templates

// TemplateCategory represents a category of templates
type TemplateCategory struct {
	Name        string
	Description string
	Templates   []Template
}

// Template represents a template configuration
type Template struct {
	ID          string
	Name        string
	Description string
	Type        string
}

// GetCategories returns all available template categories
func GetCategories() []TemplateCategory {
	return []TemplateCategory{
		{
			Name:        "Web Applications",
			Description: "Traditional web application stacks",
			Templates: []Template{
				{ID: "mern", Name: "MERN Stack", Description: "MongoDB, Express, React, Node.js", Type: "web"},
				{ID: "mean", Name: "MEAN Stack", Description: "MongoDB, Express, Angular, Node.js", Type: "web"},
				{ID: "mevn", Name: "MEVN Stack", Description: "MongoDB, Express, Vue.js, Node.js", Type: "web"},
				{ID: "pern", Name: "PERN Stack", Description: "PostgreSQL, Express, React, Node.js", Type: "web"},
				{ID: "mnfa", Name: "MNFA Stack", Description: "MongoDB, Neo4j, FastAPI, Angular", Type: "web"},
				{ID: "pdn", Name: "PDN Stack", Description: "PostgreSQL, Django, Node.js", Type: "web"},
			},
		},
		{
			Name:        "Machine Learning",
			Description: "ML pipeline and model serving templates",
			Templates: []Template{
				{ID: "kubeflow", Name: "Kubeflow Pipeline", Description: "ML pipeline with Kubeflow", Type: "ml"},
				{ID: "mlflow", Name: "MLflow Stack", Description: "MLflow with tracking server", Type: "ml"},
				{ID: "tensorflow-serving", Name: "TensorFlow Serving", Description: "Model serving with TF Serving", Type: "ml"},
				{ID: "triton", Name: "Triton Server", Description: "NVIDIA Triton Inference Server", Type: "ml"},
			},
		},
		{
			Name:        "AI/LLM",
			Description: "AI and Large Language Model templates",
			Templates: []Template{
				{ID: "langchain-nextjs", Name: "LangChain.js", Description: "LangChain.js with Next.js", Type: "ai"},
				{ID: "langchain-fastapi", Name: "LangChain Python", Description: "LangChain Python with FastAPI", Type: "ai"},
				{ID: "openai-node", Name: "OpenAI Node.js", Description: "OpenAI with Express and React", Type: "ai"},
				{ID: "openai-py", Name: "OpenAI Python", Description: "OpenAI with FastAPI and Vue", Type: "ai"},
				{ID: "llama-node", Name: "Llama Node.js", Description: "Llama.cpp with Next.js", Type: "ai"},
				{ID: "llama-py", Name: "Llama Python", Description: "Llama.cpp with FastAPI", Type: "ai"},
				{ID: "vertex-ai", Name: "Vertex AI", Description: "Google Vertex AI with Flask", Type: "ai"},
				{ID: "huggingface", Name: "Hugging Face", Description: "Hugging Face with FastAPI", Type: "ai"},
				{ID: "anthropic-py", Name: "Anthropic Python", Description: "Anthropic Claude with FastAPI", Type: "ai"},
				{ID: "anthropic-js", Name: "Anthropic Node.js", Description: "Anthropic Claude with Next.js", Type: "ai"},
			},
		},
	}
}

// GetTemplateByID returns a template by its ID
func GetTemplateByID(id string) *Template {
	for _, category := range GetCategories() {
		for _, template := range category.Templates {
			if template.ID == id {
				return &template
			}
		}
	}
	return nil
}
