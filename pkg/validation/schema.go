package validation

// SchemaV2 contains the JSON Schema for validating Nexlayer templates
const SchemaV2 = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["application"],
  "properties": {
    "application": {
      "type": "object",
      "required": ["name", "pods"],
      "properties": {
        "name": {
          "type": "string",
          "description": "REQUIRED: The name of the deployment (must be unique)"
        },
        "url": {
          "type": "string",
          "description": "OPTIONAL: Permanent domain (only include if needed)"
        },
        "registryLogin": {
          "type": "object",
          "required": ["registry", "username", "personalAccessToken"],
          "properties": {
            "registry": {
              "type": "string",
              "description": "REQUIRED for private images: The registry for images"
            },
            "username": {
              "type": "string",
              "description": "REQUIRED if using a private registry"
            },
            "personalAccessToken": {
              "type": "string",
              "description": "REQUIRED for read-only registry authentication"
            }
          }
        },
        "pods": {
          "type": "array",
          "items": {
            "type": "object",
            "required": ["name", "image", "servicePorts"],
            "properties": {
              "name": {
                "type": "string",
                "pattern": "^[a-z][a-z0-9\\.\\-]*$",
                "description": "REQUIRED: Pod name (must start with a lowercase letter, only alphanumeric, '-', or '.')"
              },
              "path": {
                "type": "string",
                "description": "OPTIONAL: Route path for frontend (e.g., '/' for web apps)"
              },
              "image": {
                "type": "string",
                "description": "REQUIRED: Docker image path (supports '<% REGISTRY %>' for private images)"
              },
              "volumes": {
                "type": "array",
                "items": {
                  "type": "object",
                  "required": ["name", "size", "mountPath"],
                  "properties": {
                    "name": {
                      "type": "string",
                      "description": "REQUIRED: Name of the volume"
                    },
                    "size": {
                      "type": "string",
                      "pattern": "^\\d+[KMGT]i$",
                      "description": "REQUIRED: Volume size (e.g., '1Gi')"
                    },
                    "mountPath": {
                      "type": "string",
                      "description": "REQUIRED: Path inside the container"
                    }
                  }
                }
              },
              "secrets": {
                "type": "array",
                "items": {
                  "type": "object",
                  "required": ["name", "data", "mountPath", "fileName"],
                  "properties": {
                    "name": {
                      "type": "string",
                      "description": "REQUIRED: Secret name"
                    },
                    "data": {
                      "type": "string",
                      "description": "REQUIRED: Base64-encoded or raw secret value"
                    },
                    "mountPath": {
                      "type": "string",
                      "description": "REQUIRED: Directory where the secret file will be stored"
                    },
                    "fileName": {
                      "type": "string",
                      "description": "REQUIRED: File name for the secret (e.g., 'config.json')"
                    }
                  }
                }
              },
              "vars": {
                "type": "array",
                "items": {
                  "type": "object",
                  "required": ["key", "value"],
                  "properties": {
                    "key": {
                      "type": "string",
                      "description": "REQUIRED: Environment variable key"
                    },
                    "value": {
                      "type": "string",
                      "description": "REQUIRED: Value (Supports: pod references, '<% URL %>', etc.)"
                    }
                  }
                }
              },
              "servicePorts": {
                "type": "array",
                "items": {
                  "type": "integer",
                  "minimum": 1,
                  "maximum": 65535,
                  "description": "REQUIRED: Port to expose (e.g., 3000)"
                },
                "minItems": 1
              }
            }
          },
          "minItems": 1
        }
      }
    }
  }
}`
