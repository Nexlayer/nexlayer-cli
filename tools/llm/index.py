#!/usr/bin/env python3
"""
Creates a semantic search index from LLM metadata for faster context retrieval.
This helps LLMs quickly find relevant deployment patterns, examples, and annotated templates.
"""

import os
import sys
import json
import yaml
from typing import Dict, List, Any

def load_metadata(filepath: str) -> Dict[Any, Any]:
    """
    Load metadata from a JSON or YAML file based on the file extension.
    """
    with open(filepath, 'r') as f:
        if filepath.lower().endswith(('.yaml', '.yml')):
            return yaml.safe_load(f)
        else:
            return json.load(f)

def write_index(index: Dict[str, Any], filepath: str) -> None:
    """
    Write the semantic index to a file in JSON or YAML format based on the file extension.
    """
    with open(filepath, 'w') as f:
        if filepath.lower().endswith(('.yaml', '.yml')):
            yaml.dump(index, f, sort_keys=False, allow_unicode=True)
        else:
            json.dump(index, f, indent=2)

def create_semantic_index(metadata: Dict[Any, Any]) -> Dict[str, List[Dict[str, Any]]]:
    """
    Create a semantic search index from metadata.
    
    The index includes:
      - 'intents': User intents with context.
      - 'patterns': Deployment patterns with associated details.
      - 'examples': Deployment examples (if provided).
      - 'api_usage': API endpoint usage examples.
      - 'templates': Annotated template information (if provided).
    """
    index = {
        'intents': [],
        'patterns': [],
        'examples': [],
        'api_usage': [],
        'templates': []
    }
    
    # Index user intents with context.
    for intent in metadata.get('user_intents', []):
        index['intents'].append({
            'text': intent.get('intent', ''),
            'keywords': intent.get('keywords', []),
            'context': {
                'actions': intent.get('actions', []),
                'examples': intent.get('examples', []),
                'suggestions': intent.get('suggestions', [])
            }
        })
    
    # Index deployment patterns.
    for pattern in metadata.get('deployment_patterns', []):
        index['patterns'].append({
            'text': pattern.get('description', ''),
            'keywords': pattern.get('keywords', []),
            'context': {
                'name': pattern.get('name', ''),
                'template': pattern.get('template', ''),
                'explanation': pattern.get('explanation', ''),
                'use_case': pattern.get('use_case', '')
            }
        })
    
    # Index deployment examples if available.
    for example in metadata.get('deployment_examples', []):
        index['examples'].append({
            'text': example.get('description', ''),
            'keywords': example.get('keywords', []),
            'context': example
        })
    
    # Index API usage examples.
    for endpoint in metadata.get('api_endpoints', []):
        for example in endpoint.get('usage_examples', []):
            index['api_usage'].append({
                'text': example,
                'context': {
                    'path': endpoint.get('path', ''),
                    'method': endpoint.get('method', ''),
                    'description': endpoint.get('description', ''),
                    'patterns': endpoint.get('common_patterns', [])
                }
            })
    
    # Optionally index annotated template information if available.
    if 'annotated_template' in metadata:
        app = metadata['annotated_template'].get('application', {})
        template_entry = {
            'text': f"Deployment template for {app.get('name', 'unknown')}",
            'keywords': [app.get('name', ''), app.get('url', '')] if app else [],
            'context': metadata['annotated_template']
        }
        index['templates'].append(template_entry)
    
    return index

def main():
    if len(sys.argv) != 3:
        print("Usage: index.py <metadata_file.(json|yaml|yml)> <output_index.(json|yaml|yml)>")
        sys.exit(1)
        
    metadata_file = sys.argv[1]
    output_file = sys.argv[2]
    
    # Load metadata from JSON or YAML.
    metadata = load_metadata(metadata_file)
    
    # Create semantic index.
    index = create_semantic_index(metadata)
    
    # Write the semantic index to output file.
    write_index(index, output_file)

if __name__ == '__main__':
    main()
