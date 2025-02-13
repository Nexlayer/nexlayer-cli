# Nexlayer AI Training

## Overview
This directory contains examples and resources used to train Nexlayer's AI assistant in generating and validating deployment templates. The official schema and template files are located in `/docs/reference/schemas/yaml/`, which serve as the central source of truth for our YAML configuration format.

## Schema Files

### Schema and Template Files

The following files in `/docs/reference/schemas/yaml/` define our YAML format:
This YAML file serves as a comprehensive reference for our template structure. It includes:
- Detailed descriptions of each field and its purpose
- Examples of proper usage
- Best practices and conventions
- Component type definitions and their standard configurations

This file is particularly useful for:
1. Training the AI to understand our template structure
2. Providing examples of common deployment patterns
3. Documenting best practices for different component types

### `schema/nexlayer_template_reference_v1.0.json`
JSON version of the reference document, used for:
- Programmatic access during AI training
- Schema validation in development tools
- Documentation generation

## How It Works

1. **AI Training**
   - The AI uses these files to learn about:
     - Valid component types (frontend, backend, database, etc.)
     - Standard port configurations
     - Environment variable patterns
     - Common deployment patterns

2. **Template Generation**
   - When users request new templates, the AI refers to these files to:
     - Select appropriate component types
     - Configure correct ports and environment variables
     - Apply best practices

3. **Validation**
   - While these files guide template generation, the actual validation is done using:
     - The official schema in `/pkg/validation/schema/template.v2.schema.json`
     - Runtime validators in the `template` package

## Example Usage

```yaml
# Example of AI-generated template following our standards
application:
  template:
    name: MyApp
    deploymentName: my-app
    pods:
      - type: llm
        name: ollama
        image: us-east1-docker.pkg.dev/my-registry/models/ollama:latest
        ports:
          - containerPort: 11434
            servicePort: 11434
            name: ollama
      - type: frontend
        name: web-ui
        image: us-east1-docker.pkg.dev/my-registry/apps/web-ui:v1
        vars:
          - key: REACT_APP_API_URL
            value: http://CANDIDATE_DEPENDENCY_URL_0:11434
        ports:
          - containerPort: 3000
            servicePort: 80
            name: web
```

## Contributing
When updating the training schema:
1. Update both YAML and JSON versions
2. Add new examples if introducing new patterns
3. Keep the documentation comprehensive and up-to-date
4. Ensure examples follow our best practices

> Copyright (c) 2025 Nexlayer. All rights reserved.  
> Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

This directory contains essential files for training and guiding Nexlayer's AI capabilities. It serves as the authoritative source for AI behavior, template generation, and deployment patterns.

## Directory Structure

```
/ai_training/
├── schema/                          # Contains authoritative schema/reference files
│   ├── nexlayer_template_reference_v1.0.yaml  # YAML reference for AI training
│   └── nexlayer_template_reference_v1.0.json  # JSON version for programmatic use
└── README.md                        # Documentation for the AI training folder
```

## Schema Files

### Template Reference
The template reference files (`nexlayer_template_reference_v1.0.*`) define the structure and validation rules for Nexlayer deployment templates. They include:

- Component type definitions
- Required and optional fields
- Environment variable patterns
- Port configurations
- Best practices
- Example templates

## Schema Formats

### YAML Reference (`nexlayer_template_reference_v1.0.yaml`)
- Human-readable reference format
- Contains comprehensive documentation and examples
- Includes best practices and component configurations
- Used primarily for AI training and documentation

### JSON Schema (`nexlayer_template_reference_v1.0.json`)
- Machine-readable schema for validation
- Follows JSON Schema specification
- Includes detailed type definitions and validation rules
- Used for programmatic validation and tooling

**Note:** These files need to be synchronized. Currently there are some discrepancies:
1. Version numbers don't match (YAML: 2.0, JSON: 1.0)
2. Port configuration structure differs
3. Schema validation rules are more detailed in the JSON version

## Version Control

Schema files are versioned to track changes and maintain compatibility:
- Major version changes (v1.0 → v2.0): Breaking changes
- Minor version changes (v1.0 → v1.1): Backward-compatible additions
- Patch changes tracked in git history

## Usage Guidelines

1. AI agents must strictly adhere to these specifications
2. Do not generate fields or values not defined in the schema
3. Follow best practices outlined in the reference files
4. Use example templates as guidance for common deployment patterns

## Contributing

When updating AI training materials:
1. Create new versioned files for major changes
2. Document changes in commit messages
3. Update this README if new file types are added
4. Ensure backward compatibility when possible
