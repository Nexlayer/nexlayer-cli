#!/usr/bin/env python3
"""
Creates a semantic search index from LLM metadata for faster context retrieval.
This helps LLMs quickly find relevant deployment patterns and examples.
"""

import json
import sys
from typing import Dict, List, Any

def create_semantic_index(metadata: Dict[Any, Any]) -> Dict[str, List[Dict[str, Any]]]:
    """Create a semantic search index from metadata."""
    index = {
        'intents': [],
        'patterns': [],
        'examples': [],
        'api_usage': []
    }
    
    # Index user intents with context
    for intent in metadata.get('user_intents', []):
        index['intents'].append({
            'text': intent['intent'],
            'keywords': intent['keywords'],
            'context': {
                'actions': intent['actions'],
                'examples': intent['examples'],
                'suggestions': intent['suggestions']
            }
        })
    
    # Index deployment patterns
    for pattern in metadata.get('deployment_patterns', []):
        index['patterns'].append({
            'text': pattern['description'],
            'keywords': pattern['keywords'],
            'context': {
                'name': pattern['name'],
                'template': pattern['template'],
                'explanation': pattern['explanation'],
                'use_case': pattern['use_case']
            }
        })
    
    # Index API usage examples
    for endpoint in metadata.get('api_endpoints', []):
        for example in endpoint['usage_examples']:
            index['api_usage'].append({
                'text': example,
                'context': {
                    'path': endpoint['path'],
                    'method': endpoint['method'],
                    'description': endpoint['description'],
                    'patterns': endpoint['common_patterns']
                }
            })
    
    return index

def main():
    if len(sys.argv) != 3:
        print("Usage: index.py <metadata_json> <output_index>")
        sys.exit(1)
        
    metadata_file = sys.argv[1]
    output_file = sys.argv[2]
    
    # Read metadata
    with open(metadata_file, 'r') as f:
        metadata = json.load(f)
    
    # Create semantic index
    index = create_semantic_index(metadata)
    
    # Write index
    with open(output_file, 'w') as f:
        json.dump(index, f, indent=2)

if __name__ == '__main__':
    main()
