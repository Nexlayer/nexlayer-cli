# Nexlayer Template Builder Plugin

A plugin for the Nexlayer CLI that automatically generates deployment templates by analyzing your local codebase and applying best practices.

## Features

- Auto-detects your project's tech stack (Node.js, Python, etc.)
- Fetches appropriate base templates from official Nexlayer repositories
- Customizes templates with your project's details
- Optional AI-powered template refinement (if OPENAI_API_KEY or CLAUDE_API_KEY is set)
- Generates instantly deployable Nexlayer templates

## Installation

1. Build the plugin:
```bash
go build -o template-builder
```

2. Move the binary to your Nexlayer plugins directory:
```bash
mkdir -p ~/.nexlayer/plugins
mv template-builder ~/.nexlayer/plugins/
chmod +x ~/.nexlayer/plugins/template-builder
```

## Usage

From your project's root directory, run:
```bash
nexlayer template:generate
```

For a dry run (preview without writing files):
```bash
nexlayer template:generate --dry-run
```

## How It Works

1. **Stack Detection**: The plugin scans your project directory for common files:
   - `package.json` → Node.js projects
   - `requirements.txt` → Python projects
   - `.env` → Database and environment configurations

2. **Template Selection**: Based on the detected stack, it selects the most appropriate template from official Nexlayer repositories.

3. **Customization**: The template is customized with:
   - Project name (from directory name)
   - Detected technologies
   - Environment configuration

4. **AI Enhancement** (Optional): If an AI API key is configured, the template can be refined with:
   - Best practices validation
   - Resource optimization suggestions
   - Configuration improvements

5. **Output**: Generates a `<projectName>-nexlayer-template.yaml` file ready for deployment.

## Supported Stacks

- MERN (MongoDB, Express, React, Node.js)
- MEVN (MongoDB, Express, Vue.js, Node.js)
- Django + PostgreSQL
- More coming soon!

## Environment Variables

The plugin uses the same environment variables as the main Nexlayer CLI:
- `OPENAI_API_KEY`: For OpenAI-powered template refinement
- `CLAUDE_API_KEY`: For Claude AI-powered template refinement

## Contributing

Feel free to contribute by:
1. Adding new template detectors
2. Improving stack detection logic
3. Adding new official templates
4. Enhancing AI refinement capabilities
