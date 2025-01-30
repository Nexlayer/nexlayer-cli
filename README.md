# Nexlayer CLI

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/Nexlayer/nexlayer-cli)](https://goreportcard.com/report/github.com/Nexlayer/nexlayer-cli)
[![GoDoc](https://godoc.org/github.com/Nexlayer/nexlayer-cli?status.svg)](https://godoc.org/github.com/Nexlayer/nexlayer-cli)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Deploy AI applications in seconds**

[Quick Start](#quick-start) • [Templates](#templates) • [Examples](#stack-examples) • [Docs](https://docs.nexlayer.com)

</div>

---

## Prerequisites (Recommended)

- **Go**: You'll need Go 1.18+ installed to run `go install`.
- **Docker**: Nexlayer uses Docker for containerizing your applications.
- **A GitHub Account** (optional): Needed if you plan to push to GHCR (GitHub Container Registry).

*(If you need detailed setup steps for Go, Docker, or private registries, see [Nexlayer Docs](https://docs.nexlayer.com).)*

---

## Quick Start

1. **Install the CLI**  
   ```bash
   go install github.com/Nexlayer/nexlayer-cli@latest
   ```
   Make sure `$GOPATH/bin` is in your PATH so that the `nexlayer` command is recognized.

2. **Initialize a New Project**
   ```bash
   nexlayer init myapp
   ```
   This automatically detects your stack, creates a `nexlayer.yaml` configuration file, and prepares the project for deployment.

3. **Deploy Your App**
   ```bash
   nexlayer deploy
   ```
   That's it! Your app goes live in seconds.

### Next Steps
- Check Status: `nexlayer status` to view current deployment state.
- View Logs: `nexlayer logs -f [podName]` to stream logs.
- Add a Custom Domain: `nexlayer domain add yourdomain.com`
- (Everything else can be fine-tuned in `nexlayer.yaml` or by choosing a template.)

## Templates

Nexlayer provides a variety of templates to help you get started quickly. Templates are organized into three categories:

### Web Applications
Traditional web application stacks:
- `mern`: MongoDB, Express, React, Node.js
- `mean`: MongoDB, Express, Angular, Node.js
- `mevn`: MongoDB, Express, Vue.js, Node.js
- `pern`: PostgreSQL, Express, React, Node.js
- `mnfa`: MongoDB, Neo4j, FastAPI, Angular
- `pdn`: PostgreSQL, Django, Node.js

### Machine Learning
ML pipeline and model serving templates:
- `kubeflow`: ML pipeline with Kubeflow
- `mlflow`: MLflow with tracking server
- `tensorflow-serving`: Model serving with TF Serving
- `triton`: NVIDIA Triton Inference Server

### AI/LLM
AI and Large Language Model templates:
- `langchain-nextjs`: LangChain.js with Next.js
- `langchain-fastapi`: LangChain Python with FastAPI
- `openai-node`: OpenAI with Express and React
- `openai-py`: OpenAI with FastAPI and Vue
- `llama-node`: Llama.cpp with Next.js
- `llama-py`: Llama.cpp with FastAPI
- `vertex-ai`: Google Vertex AI with Flask
- `huggingface`: Hugging Face with FastAPI
- `anthropic-py`: Anthropic Claude with FastAPI
- `anthropic-js`: Anthropic Claude with Next.js

### Using Templates

There are two ways to use templates:

1. Interactive Selection:
```bash
nexlayer init my-project
```
This will prompt you to:
1. Select a template category
2. Choose a specific template
3. Configure your project

2. Direct Selection:
```bash
nexlayer init my-project -t mern
```
Replace `mern` with any template ID from the list above.

### Listing Templates
To see all available templates:
```bash
nexlayer templates list
```

### Template Structure
All templates follow this structure:
```yaml
application:
  template:
    name: "project-name"
    deploymentName: "project-name"
    registryLogin:
      registry: ghcr.io
      username: <username>
      personalAccessToken: <token>
    pods:
      - type: backend|frontend|database|nginx|llm
        name: <pod-name>
        tag: <image-tag>
        vars:
          - key: VAR_NAME
            value: VALUE
```

### Supported Pod Types
- Frontend: `react`, `angular`, `vue`
- Backend: `express`, `django`, `fastapi`
- Database: `mongodb`, `postgres`, `redis`, `neo4j`
- Others: `nginx` (load balancing/static assets), `llm` (custom workloads)

### Standard Environment Variables
- Database: `DATABASE_CONNECTION_STRING`
- Frontend/Backend: `FRONTEND_CONNECTION_URL`, `BACKEND_CONNECTION_URL`
- LLM: `LLM_CONNECTION_URL`

## Stack Examples

### Full-Stack AI (Next.js & TypeScript)
```yaml
# nexlayer.yaml
application:
  template:
    name: fullstack-ai
    deploymentName: my-ai-app
    registryLogin:
      registry: ghcr.io
      username: your-username
      personalAccessToken: your-pat
  pods:
    - type: frontend
      name: next-app
      tag: node:18
      exposeHttp: true
      vars:
        - key: TOGETHER_API_KEY
          value: your-key
        - key: CLERK_SECRET_KEY
          value: your-key
        - key: DATABASE_URL
          value: your-neon-db-url
        - key: MIXPANEL_TOKEN
          value: your-token
    - type: database
      name: postgres
      tag: postgres:15
      vars:
        - key: POSTGRES_DB
          value: aiapp
        - key: POSTGRES_USER
          value: postgres
```

### Python ML Stack
```yaml
# nexlayer.yaml
application:
  template:
    name: ml-python
    deploymentName: my-ml-app
    registryLogin:
      registry: ghcr.io
      username: your-username
      personalAccessToken: your-pat
  pods:
    - type: backend
      name: fastapi
      tag: python:3.9
      exposeHttp: true
      vars:
        - key: AWS_ACCESS_KEY_ID
          value: your-key
        - key: AWS_SECRET_ACCESS_KEY
          value: your-key
        - key: DATABASE_CONNECTION_STRING
          value: postgresql://postgres:password@postgres:5432/mlapp
    - type: database
      name: postgres
      tag: postgres:15
    - type: database
      name: redis
      tag: redis:7
    - type: frontend
      name: react-app
      tag: node:18
      exposeHttp: true
```

### Browser-Based AI
```yaml
# nexlayer.yaml
application:
  template:
    name: browser-ai
    deploymentName: my-tfjs-app
    registryLogin:
      registry: ghcr.io
      username: your-username
      personalAccessToken: your-pat
  pods:
    - type: frontend
      name: react-app
      tag: node:18
      exposeHttp: true
    - type: backend
      name: express
      tag: node:18
      exposeHttp: true
      vars:
        - key: DATABASE_CONNECTION_STRING
          value: mongodb://mongodb:27017/tfjs
    - type: database
      name: mongodb
      tag: mongodb:6
```

### LangChain Chat App
```yaml
# nexlayer.yaml
application:
  template:
    name: langchain-nextjs
    deploymentName: my-chat-app
    registryLogin:
      registry: ghcr.io
      username: your-username
      personalAccessToken: your-pat
  pods:
    - type: frontend
      name: next-app
      tag: node:18
      exposeHttp: true
      vars:
        - key: OPENAI_API_KEY
          value: your-key
        - key: LANGCHAIN_TRACING_V2
          value: "true"
```

### LangChain RAG App
```yaml
# nexlayer.yaml
application:
  template:
    name: langchain-fastapi
    deploymentName: my-rag-app
    registryLogin:
      registry: ghcr.io
      username: your-username
      personalAccessToken: your-pat
  pods:
    - type: backend
      name: fastapi
      tag: python:3.9
      exposeHttp: true
      vars:
        - key: OPENAI_API_KEY
          value: your-key
        - key: PINECONE_API_KEY
          value: your-key
        - key: PINECONE_ENVIRONMENT
          value: gcp-starter
```

### Enterprise AI SaaS
```yaml
# nexlayer.yaml
application:
  template:
    name: enterprise-ai
    deploymentName: my-enterprise-app
    registryLogin:
      registry: ghcr.io
      username: your-username
      personalAccessToken: your-pat
  pods:
    - type: frontend
      name: react-app
      tag: node:18
      exposeHttp: true
      vars:
        - key: BACKEND_CONNECTION_URL
          value: http://django:8000
        - key: OKTA_CLIENT_ID
          value: your-client-id
    - type: backend
      name: django
      tag: python:3.9
      exposeHttp: true
      vars:
        - key: DATABASE_CONNECTION_STRING
          value: postgresql://postgres:password@postgres:5432/enterprise
        - key: AWS_BEDROCK_ACCESS_KEY
          value: your-key
    - type: database
      name: postgres
      tag: postgres:15
    - type: nginx
      name: nginx
      tag: nginx:1.25-alpine
      exposeHttp: true
      vars:
        - key: FRONTEND_CONNECTION_URL
          value: http://react-app:3000
```

### Kubeflow AI Pipelines
```yaml
# nexlayer.yaml
application:
  template:
    name: kubeflow
    deploymentName: my-ml-pipeline
  pods:
    - type: ml-workflow
      name: kubeflow-pipeline
      tag: kubeflow/pipelines
      vars:
        - key: DATASET_PATH
          value: gs://my-data
        - key: MODEL_STORAGE
          value: gs://my-models
    - type: backend
      name: fastapi
      tag: python:3.9
      exposeHttp: true
      vars:
        - key: API_KEY
          value: your-key
```

Deploy Kubeflow with a single command:
```bash
nexlayer deploy kubeflow
```
> No need to specify compute—Nexlayer auto-handles resources.

### AI Model Images for Training & Serving

#### 1. AI Model Training Images

| Framework | Base Image | Usage |
|-----------|------------|--------|
| TensorFlow | tensorflow/tensorflow:latest | General ML/DL training |
| PyTorch | pytorch/pytorch:latest | Training for PyTorch models |
| XGBoost | dmlc/xgboost:latest | Gradient boosting training |
| Scikit-Learn | python:3.9 + scikit-learn | Traditional ML models |
| FastAI | fastai/fastai:latest | Deep learning training |

Example YAML for Training:
```yaml
pods:
  - type: ml-training
    name: tensorflow-trainer
    tag: tensorflow/tensorflow:latest
    vars:
      - key: DATA_PATH
        value: gs://my-dataset
      - key: MODEL_OUTPUT
        value: gs://my-models
```
> 💡 Future-Ready: Can later swap with tensorflow/tensorflow:latest-gpu when GPU support is added. Request GPU support via [GitHub Issues](https://github.com/Nexlayer/nexlayer-cli/issues)

#### 2. AI Model Serving (Inference) Images

| Model Format | Serving Image | Usage |
|--------------|---------------|--------|
| TensorFlow SavedModel | tensorflow/serving:latest | Serving TensorFlow models |
| ONNX Models | microsoft/onnxruntime:latest | Optimized ONNX inference |
| PyTorch TorchServe | pytorch/torchserve:latest | Serving PyTorch models |
| Hugging Face Transformers | huggingface/transformers-pipeline:latest | NLP model inference |

Example YAML for Model Deployment:
```yaml
pods:
  - type: ml-inference
    name: model-serving
    tag: tensorflow/serving:latest
    vars:
      - key: MODEL_PATH
        value: gs://my-models/tf-model
```
> 🚀 Scales dynamically based on request load.

#### 3. AI Pipeline & Workflow Images

| Pipeline Task | Base Image | Usage |
|---------------|------------|--------|
| Data Processing | python:3.9 + Pandas, NumPy | Prepares data before training |
| Hyperparameter Tuning | kubeflowkatib/katib:latest | Runs AutoML optimization |
| Model Evaluation | python:3.9 + SciPy, Matplotlib | Model performance analysis |

Example YAML for Kubeflow Pipeline:
```yaml
pods:
  - type: ml-workflow
    name: preprocess-data
    tag: python:3.9
    command: ["python", "preprocess.py"]
```

### AI Model Monitoring & Logging

#### 1. Nexlayer Real-Time Logs

Get Deployment Status:
```bash
nexlayer status
```

View Real-Time Logs:
```bash
nexlayer logs -f model-serving
```

What You Can Track with Nexlayer Logs:
- Build & Deployment Logs: Track progress from build to live deployment
- Pod Activity: Monitor AI pipeline execution in real time
- Model Serving Requests: See live inference requests and responses
- Errors & Failures: Identify and debug model issues quickly

#### 2. External AI Monitoring & Observability

Nexlayer does not provide model performance monitoring, but you can integrate with third-party observability tools:

| Category | Tool | Use Case |
|----------|------|----------|
| LLM Observability | Helicone | Logs & monitors OpenAI/Anthropic model usage |
| Application Monitoring | Datadog | Real-time app performance tracking |
| ML Model Performance | Arize AI | Detects model drift, bias, and degradation |
| Log Management | LogDNA | Aggregates & analyzes system logs |
| Experiment Tracking | MLflow | Logs & versions model experiments |

> 📝 Nexlayer handles logs; you can integrate with any third-party tool for model performance tracking.

## Template Configuration

Each Nexlayer deployment requires a YAML configuration file that defines your application structure. Here's how to configure it:

### Basic Structure
```yaml
application:
  template:
    name: my-app-stack          # Identifier for your app stack
    deploymentName: my-app      # Your deployment name
    registryLogin:              # Optional: for private registries
      username: user
      password: pass

  pods:                         # Define your app components
    - type: react              # Pod type (database/frontend/backend/etc)
      name: frontend           # Specific name for the pod
      tag: node:14-alpine      # Docker image
      privateTag: false        # Is it from a private registry?
      vars:                    # Environment variables
        - name: PORT
          value: "3000"
      exposeHttp: true        # Make pod accessible via HTTP
```

### Supported Pod Types
- **Database**: `postgres`, `mysql`, `neo4j`, `redis`, `mongodb`
- **Frontend**: `react`, `angular`, `vue`
- **Backend**: `django`, `fastapi`, `express`
- **Others**: `nginx`, `llm` (custom naming allowed)

### Environment Variables
Nexlayer automatically provides these environment variables to your pods:

| Variable | Description | Example |
|----------|-------------|---------|
| `PROXY_URL` | Your Nexlayer site URL | `https://your-site.alpha.nexlayer.ai` |
| `PROXY_DOMAIN` | Your Nexlayer site domain | `your-site.alpha.nexlayer.ai` |
| `DATABASE_HOST` | Database hostname | - |
| `DATABASE_CONNECTION_STRING` | Database connection string | `postgresql://user:pass@host:port/db` |
| `FRONTEND_CONNECTION_URL` | Frontend URL (with http://) | - |
| `BACKEND_CONNECTION_URL` | Backend URL (with http://) | - |
| `LLM_CONNECTION_URL` | LLM URL (with http://) | - |
| `FRONTEND_CONNECTION_DOMAIN` | Frontend domain (no prefix) | - |
| `BACKEND_CONNECTION_DOMAIN` | Backend domain (no prefix) | - |
| `LLM_CONNECTION_DOMAIN` | LLM domain (no prefix) | - |

### GitHub Actions Integration
Create `.github/workflows/docker-publish.yml`:

```yaml
name: Build and Push Docker Image

on:
  push:
    branches:
      - main

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
    - uses: actions/checkout@v2
    - uses: docker/login-action@v2
      with:
        registry: ghcr.io
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    - run: echo "owner_lowercase=$(echo '${{ github.repository_owner }}' | tr '[:upper:]' '[:lower:]')" >> $GITHUB_ENV
    - uses: docker/build-push-action@v5
      with:
        context: .
        push: true
        tags: ghcr.io/${{ env.owner_lowercase }}/my-image-name:v0.0.1
```

## Features

- **Smart Detection**: Automatically detects your stack and configures everything
- **Simple Controls**: One command to initialize, one to deploy
- **Fast Cold Starts**: Sub-second startup times
- **Zero Config**: Sensible defaults for every stack
- **GPU Ready**: Built-in support for GPU acceleration
- **Cost Efficient**: Scale to zero when idle
- **Progress Feedback**: Visual progress indicators during operations
- **Error Handling**: Clear error messages and validation

## Plugins

Nexlayer supports plugins to extend its functionality. Plugins are Go shared libraries (.so files) that implement the Plugin interface.

### Using Plugins

```bash
# List installed plugins
nexlayer plugin list

# Run a plugin
nexlayer plugin run hello --name "John"

# Install a plugin
nexlayer plugin install ./my-plugin.so
```

### Creating Plugins

1. Create a new Go file for your plugin:

```go
package main

type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-plugin"
}

func (p *MyPlugin) Description() string {
    return "Description of what my plugin does"
}

func (p *MyPlugin) Run(opts map[string]interface{}) error {
    // Plugin logic here
    return nil
}

// Export the plugin
var Plugin MyPlugin
```

2. Build your plugin as a shared library:

```bash
go build -buildmode=plugin -o my-plugin.so my-plugin.go
```

3. Install your plugin:

```bash
nexlayer plugin install my-plugin.so
```

### Plugin Directory

Plugins are stored in `~/.nexlayer/plugins/`. Each plugin is a `.so` file that implements the Plugin interface.

### Plugin Interface

```go
type Plugin interface {
    // Name returns the name of the plugin
    Name() string
    
    // Description returns a description of what the plugin does
    Description() string
    
    // Run executes the plugin with the given options
    Run(opts map[string]interface{}) error
}

## Usage

```bash
# Deployment
nexlayer deploy          # Deploy your application
nexlayer status         # Check deployment status

# Configuration
nexlayer domain add     # Add custom domain

# AI-Powered Features
nexlayer init myapp     # Initialize a new app with AI-generated config
nexlayer ai detect      # Detect available AI assistants
nexlayer ai debug       # Get AI-powered deployment debugging
nexlayer ai scale       # AI-driven scaling recommendations
```

## AI Integration

Nexlayer CLI integrates with your IDE's AI capabilities to provide enhanced features:

### Automatic AI Detection
- Detects supported AI tools (GitHub Copilot, JetBrains AI, Cursor, Windsurf, Cline)
- Caches detection results in `~/.nexlayer/config.yaml`
- Runs automatically during installation or first `init`

```bash
$ nexlayer ai detect
✅ Detected AI Models:
   - GitHub Copilot (VS Code)
   - Cursor AI
```

### Smart YAML Generation
When using `nexlayer init`, the CLI:
- Analyzes your project structure
- Detects frameworks and dependencies
- Generates optimized deployment configuration

Example generated YAML:
```yaml
application:
  template:
    name: myapp
    deploymentName: myapp
    pods:
      - type: backend
        name: Node.js API
        tag: node:14
      - type: frontend
        name: React
        tag: nginx:latest
```

### AI-Powered Debugging
Debug deployment issues with AI assistance:
```bash
$ nexlayer ai debug --app myapp
❌ Deployment Error:
   - Issue: Missing environment variable `DATABASE_URL`
   - Suggested Fix: Add `DATABASE_URL` to your YAML under the `backend` pod

Suggested YAML Fix:
application:
  template:
    pods:
      - type: backend
        name: Node.js API
        vars:
          - key: DATABASE_URL
            value: mongodb://mongo-service
```

### Intelligent Scaling
Get AI-driven scaling recommendations:
```bash
$ nexlayer ai scale --app myapp
✅ Scaling Recommendation:
   - Current replicas: 2
   - Recommended replicas: 5 (based on traffic patterns)
```

## Testing

Run the test suite:

```bash
# Run all tests
./test/cli_test.sh

# Test specific functionality
nexlayer init myapp -t langchain-nextjs    # Test template initialization
nexlayer init myapp                        # Test auto-detection
```

The test suite covers:
- Command validation
- Template handling
- Project initialization
- Auto-detection
- Error scenarios
- Performance
- Concurrent operations

## Support
- [Documentation](https://docs.nexlayer.com)
- [GitHub Issues](https://github.com/Nexlayer/nexlayer-cli/issues)

## License

MIT

---

### Potential Missing Pieces

- **Local vs Cloud Deploy**: Depending on your environment, you might need additional login/credentials. Check out the [Nexlayer Docs](https://docs.nexlayer.com) for cloud deployments, secrets management, and advanced config.
- **Logging & Monitoring**: For deep observability, you may want to integrate with existing logging solutions (Datadog, Sentry, etc.).
- **Custom Domains & SSL**: See `nexlayer domain add` and the docs for info on SSL certificates and custom domain mappings.

Happy Deploying! 🚀
