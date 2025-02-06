# Nexlayer AI Training

> Copyright (c) 2025 Nexlayer. All rights reserved.  
> Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

This directory contains essential files for training and guiding Nexlayer's AI capabilities. It serves as the authoritative source for AI behavior, template generation, and deployment patterns.

## Directory Structure

```
/ai_training/
├── assets/                          # Contains supplementary files (images, diagrams, etc.)
│   ├── logo.png
│   └── architecture-diagram.png
├── schema/                          # Contains authoritative schema/reference files
│   ├── nexlayer_template_reference_v1.0.yaml
│   └── nexlayer_template_reference_v1.0.json
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

Available formats:
- YAML: Human-readable reference format
- JSON: Machine-readable schema for validation

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
