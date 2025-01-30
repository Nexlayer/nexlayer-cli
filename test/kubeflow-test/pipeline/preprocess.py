import kfp
from kfp import dsl
from kfp.components import create_component_from_func
from typing import NamedTuple
from kfp.components import OutputPath, InputPath
import json

@create_component_from_func
def preprocess_data(data_path: str) -> str:
    import pandas as pd
    import numpy as np
    import os
    
    # Simulate data preprocessing
    print(f"Loading data from {data_path}")
    # Create dummy data
    data = pd.DataFrame({
        'feature1': np.random.rand(1000),
        'feature2': np.random.rand(1000),
        'target': np.random.randint(0, 2, 1000)
    })
    
    # Save processed data
    os.makedirs('/tmp/data', exist_ok=True)
    output_path = '/tmp/data/processed.csv'
    data.to_csv(output_path, index=False)
    print(f"Saved processed data to {output_path}")
    return output_path

@create_component_from_func
def train_model(data_path: str) -> str:
    import pandas as pd
    from sklearn.model_selection import train_test_split
    from sklearn.ensemble import RandomForestClassifier
    import joblib
    import os
    import json
    
    # Load data
    data = pd.read_csv(data_path)
    X = data[['feature1', 'feature2']]
    y = data['target']
    
    # Split data
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2)
    
    # Train model
    model = RandomForestClassifier()
    model.fit(X_train, y_train)
    
    # Create output directories
    os.makedirs('/tmp/model', exist_ok=True)
    os.makedirs('/tmp/test_data', exist_ok=True)
    
    # Save outputs
    model_path = '/tmp/model/model.joblib'
    x_test_path = '/tmp/test_data/x_test.csv'
    y_test_path = '/tmp/test_data/y_test.csv'
    
    joblib.dump(model, model_path)
    X_test.to_csv(x_test_path, index=False)
    y_test.to_csv(y_test_path, index=False)
    
    # Return paths as JSON
    output_paths = {
        'model_path': model_path,
        'x_test_path': x_test_path,
        'y_test_path': y_test_path
    }
    return json.dumps(output_paths)

@create_component_from_func
def evaluate_model(paths_json: str) -> float:
    import json
    from sklearn.metrics import accuracy_score
    import joblib
    import pandas as pd
    
    # Parse paths
    paths = json.loads(paths_json)
    model_path = paths['model_path']
    x_test_path = paths['x_test_path']
    y_test_path = paths['y_test_path']
    
    # Load model and test data
    model = joblib.load(model_path)
    X_test = pd.read_csv(x_test_path)
    y_test = pd.read_csv(y_test_path)
    
    # Make predictions
    y_pred = model.predict(X_test)
    
    # Calculate accuracy
    accuracy = accuracy_score(y_test.iloc[:, 0], y_pred)
    print(f"Model accuracy: {accuracy}")
    
    return accuracy

@dsl.pipeline(
    name='Simple ML Pipeline',
    description='A simple ML pipeline for testing Kubeflow'
)
def ml_pipeline(data_path: str = 'data/input.csv'):
    # Preprocess data
    preprocess_op = preprocess_data(data_path)
    
    # Train model
    train_op = train_model(preprocess_op.output)
    
    # Evaluate model
    evaluate_op = evaluate_model(train_op.output)

if __name__ == '__main__':
    # Compile the pipeline
    kfp.compiler.Compiler().compile(ml_pipeline, 'ml_pipeline.yaml')
