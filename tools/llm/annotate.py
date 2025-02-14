#!/usr/bin/env python3
"""
Annotates YAML templates with LLM-friendly comments and explanations.
This helps LLMs understand the purpose and structure of Nexlayer templates.
"""

import sys
import yaml
from typing import Any, Dict, Optional

def add_llm_annotations(data: Dict[Any, Any]) -> Dict[Any, Any]:
    """
    Recursively add LLM-friendly annotations to YAML data based on the Nexlayer schema.
    """
    # Define annotations for each key based on the new Nexlayer YAML template.
    annotations = {
        "application": {
            "_llm_description": "Root configuration for a Nexlayer application deployment",
            "_llm_examples": [
                "Simple web service",
                "Multi-pod application",
                "Stateful application with volumes"
            ],
            "_llm_common_patterns": [
                "Frontend with backend services",
                "API with database",
                "Static website"
            ],
            "name": {
                "_llm_description": "The name of the deployment (must be lowercase, alphanumeric, '-', '.')"
            },
            "url": {
                "_llm_description": "Permanent domain URL (optional). Omit if not a permanent deployment."
            },
            "registryLogin": {
                "_llm_description": "Registry login information for private image storage",
                "registry": {
                    "_llm_description": "The registry where private images are stored"
                },
                "username": {
                    "_llm_description": "Registry username"
                },
                "personalAccessToken": {
                    "_llm_description": "Read-only registry Personal Access Token"
                }
            },
            "pods": {
                "_llm_description": "List of containers to deploy",
                "_llm_pod_types": [
                    "Web servers",
                    "Application servers",
                    "Databases",
                    "Background workers"
                ],
                "_llm_best_practices": [
                    "Use descriptive pod names",
                    "Configure appropriate resource limits",
                    "Set health checks for reliability"
                ],
                # Annotations for each pod in the pods list
                "_llm_item_annotations": {
                    "name": {
                        "_llm_description": "Pod name (must start with a lowercase letter and can include only alphanumeric characters, '-', '.')"
                    },
                    "path": {
                        "_llm_description": "Path to render pod at (e.g., '/' for frontend). Only required for forward-facing pods."
                    },
                    "image": {
                        "_llm_description": (
                            "Docker image for the pod. For private images, use the schema "
                            "'<% REGISTRY %>/some/path/image:tag'. Images including '<% REGISTRY %>' "
                            "will have the registry value replaced."
                        )
                    },
                    "volumes": {
                        "_llm_description": "List of volumes to be mounted for this pod",
                        "_llm_item_annotations": {
                            "name": {
                                "_llm_description": "Volume name (lowercase, alphanumeric, '-')"
                            },
                            "size": {
                                "_llm_description": 'Volume size (e.g., "1Gi", "500Mi")'
                            },
                            "mountPath": {
                                "_llm_description": "Mount path for the volume (must start with '/')"
                            }
                        }
                    },
                    "secrets": {
                        "_llm_description": "List of secret files for this pod",
                        "_llm_item_annotations": {
                            "name": {
                                "_llm_description": "Secret name (lowercase, alphanumeric, '-')"
                            },
                            "data": {
                                "_llm_description": "Secret data (raw or Base64-encoded)"
                            },
                            "mountPath": {
                                "_llm_description": "Mount path where the secret file will be stored (must start with '/')"
                            },
                            "fileName": {
                                "_llm_description": "Name of the secret file (e.g., 'secret-file.txt')"
                            }
                        }
                    },
                    "vars": {
                        "_llm_description": "Environment variables for this pod",
                        "_llm_item_annotations": {
                            "key": {
                                "_llm_description": "Environment variable name"
                            },
                            "value": {
                                "_llm_description": "Value of the environment variable"
                            }
                        }
                    },
                    "servicePorts": {
                        "_llm_description": "Ports to expose for this pod",
                        "_llm_examples": ["3000"]
                    }
                }
            },
            "entrypoint": {
                "_llm_description": "Custom container entrypoint (optional)"
            },
            "command": {
                "_llm_description": "Custom container command (optional)"
            }
        }
    }

    def annotate_data(d: Any, ann: Optional[Dict[Any, Any]]) -> Any:
        """
        Recursively annotate the data `d` using the provided annotations `ann`.
        """
        if isinstance(d, dict):
            result = {}
            for key, value in d.items():
                # If there is an annotation for this key, attach it.
                if ann is not None and key in ann:
                    result[f"_{key}_annotations"] = ann[key]
                    # Prepare nested annotations: only those keys that don't start with '_'
                    nested_ann = {k: v for k, v in ann[key].items() if not k.startswith("_")}
                    result[key] = annotate_data(value, nested_ann)
                else:
                    # Process without annotations if no matching annotation exists.
                    if isinstance(value, dict):
                        result[key] = annotate_data(value, None)
                    elif isinstance(value, list):
                        result[key] = [
                            annotate_data(item, None) if isinstance(item, dict) else item
                            for item in value
                        ]
                    else:
                        result[key] = value
            return result
        elif isinstance(d, list):
            # If annotations for list items exist, they should be under '_llm_item_annotations'
            item_ann = ann.get("_llm_item_annotations") if (ann and isinstance(ann, dict)) else None
            return [annotate_data(item, item_ann) if isinstance(item, dict) else item for item in d]
        else:
            return d

    return annotate_data(data, annotations)

def main():
    if len(sys.argv) != 3:
        print("Usage: annotate.py <input_yaml> <output_yaml>")
        sys.exit(1)
        
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    # Read and parse input YAML.
    with open(input_file, "r") as f:
        data = yaml.safe_load(f)
    
    # Add LLM annotations.
    annotated_data = add_llm_annotations(data)
    
    # Write annotated YAML.
    with open(output_file, "w") as f:
        yaml.dump(annotated_data, f, sort_keys=False, allow_unicode=True)

if __name__ == "__main__":
    main()
