#!/usr/bin/env python3
"""
Annotates YAML templates with LLM-friendly comments and explanations.
This helps LLMs understand the purpose and structure of Nexlayer templates.
"""

import sys
import yaml
from typing import Dict, Any

def add_llm_annotations(data: Dict[Any, Any]) -> Dict[Any, Any]:
    """Add LLM-friendly annotations to YAML data."""
    annotations = {
        'application': {
            '_llm_description': 'Root configuration for a Nexlayer application deployment',
            '_llm_examples': [
                'Simple web service',
                'Multi-pod application',
                'Stateful application with volumes'
            ],
            '_llm_common_patterns': [
                'Frontend with backend services',
                'API with database',
                'Static website'
            ]
        },
        'pods': {
            '_llm_description': 'List of containers to deploy',
            '_llm_pod_types': [
                'Web servers',
                'Application servers',
                'Databases',
                'Background workers'
            ],
            '_llm_best_practices': [
                'Use descriptive pod names',
                'Configure appropriate resource limits',
                'Set health checks for reliability'
            ]
        }
    }
    
    # Add annotations while preserving original data
    result = {}
    for key, value in data.items():
        if key in annotations:
            result[f'_{key}_annotations'] = annotations[key]
        result[key] = value
        
    return result

def main():
    if len(sys.argv) != 3:
        print("Usage: annotate.py <input_yaml> <output_yaml>")
        sys.exit(1)
        
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    # Read and parse input YAML
    with open(input_file, 'r') as f:
        data = yaml.safe_load(f)
    
    # Add LLM annotations
    annotated_data = add_llm_annotations(data)
    
    # Write annotated YAML
    with open(output_file, 'w') as f:
        yaml.dump(annotated_data, f, sort_keys=False, allow_unicode=True)

if __name__ == '__main__':
    main()
