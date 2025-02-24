
# Nexlayer YAML Schema Deployment Template Documentation (v1.0)

## 🚀 What is Nexlayer?

Nexlayer is an AI-powered Kubernetes cloud platform that enables seamless scalable deployments through a specialized YAML template schema. It's designed for startups, developers, creators and AI agents, making complex production deployments simple and accessible without deep cloud infrastructure expertise.

## ✨ Key Features

- Pod Management: Define individual containers for your application
- Storage Configuration: Persistent volumes for data retention
- Security Controls: Secure storage of secrets, API keys, and sensitive data
- Environment Management: Application-specific configuration settings
- Network Configuration: Service ports for inter-container and external communication
- Registry Authentication: Private registry login support
- Auto-Discovery: Automatic networking between components
- Container Controls: Customize startup behavior with entrypoints and commands
- AI/ML Support: Specialized configurations for AI-powered applications

## 📖 Table of Contents

1. [Overview](#-what-is-nexlayer)
2. [Key Features](#-key-features)
3. [Basic Structure](#basic-structure)
4. [Application Configuration](#application-configuration)
5. [Pod Configuration](#pod-configuration)
6. [Environment Variables](#environment-variables)
7. [Service Ports](#service-ports)
8. [Secrets Management](#secrets-management)
9. [Registry Authentication](#registry-authentication)
10. [Examples](#examples)
11. [Best Practices](#best-practices)
12. [Need Help?](#need-help)

## Basic Structure

Nexlayer base YAML template structure follows this pattern:

```yaml
application:
  name: The name of the deployment
  url: Permanent domain URL (optional). No need to add this key if this is not going to be a permanent deployment.
  registryLogin:
    registry: The registry where private images are stored.
    username: Registry username.
    personalAccessToken: Read-only registry Personal Access Token.
  pods:
    - name: Pod name (must start with a lowercase letter and can include only alphanumeric characters, '-', '.')
      path: Path to render pod at (such as '/' for frontend). Only required for forward-facing pods.
      image: Docker image for the pod. 
        # For private images, use the following schema exactly as shown: '<% REGISTRY %>/some/path/image:tag'.
        # Images will be tagged as private if they include '<% REGISTRY %>', which will be replaced with the registry specified above.
      entrypoint: command to replace ENTRYPOINT of image
      command: command to replace CMD of image
      volumes:
        # Array of volumes to be mounted for this pod. Example:
        - name: Name of the volume (lowercase, alphanumeric, '-')
          size: 1Gi  # Required: Volume size (e.g., "1Gi", "500Mi").
          mountPath: /var/some/directory  # Required: Must start with '/'.
      secrets:
        # Array of secret files for this pod. Example:
        - name: Secret name (lowercase, alphanumeric, '-')
          data: Raw text or Base64-encoded string for the secret (e.g., JSON files should be encoded).
          mountPath: Mount path where the secret file will be stored (must start with '/').
          fileName: Name of the secret file (e.g., "secret-file.txt"). 
            # This will be available at "/var/secrets/my-secret-volume/secret-file.txt".
      vars:
        # Array of environment variables for this pod. Example:
        - key: ENV_VAR_NAME
          value: Value of the environment variable.
        # Can use <pod-name>.pod to reference other pods dynamically. Example:
        - key: API_URL
          value: http://express.pod:3000  # Where 'express' is the name of another pod.
        # Can use <% URL %> to reference the deployment's base URL dynamically. Example:
        - key: API_URL
          value: <% URL %>/api
    servicePorts:
      # Array of ports to expose for this pod. Example:
      - 3000  # Exposing port 3000.
  entrypoint: Custom container entrypoint (optional).
  command: Custom container command (optional).
```

# Nexlayer YAML Template Overview

Below is an overview of the Nexlayer base YAML template structure from above but in table format for better readability:

---

### Application Configuration

| Key   | Description                                                                                             |
|-------|---------------------------------------------------------------------------------------------------------|
| `name` | The name of the deployment.                                                                            |
| `url`  | Permanent domain URL (optional). Do not include if not using a permanent deployment.                   |

---

### Registry Login Configuration

| Key                   | Description                                                     |
|-----------------------|-----------------------------------------------------------------|
| `registry`            | The registry where private images are stored.                   |
| `username`            | Registry username.                                              |
| `personalAccessToken` | Read-only registry Personal Access Token.                       |

---

### Pod Configuration

| Key          | Description                                                                                                                                                       |
|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `name`       | Pod name (must start with a lowercase letter and can include only alphanumeric characters, `-`, or `.`).                                                             |
| `path`       | Path to render pod at (e.g., `/` for frontend). Only required for forward-facing pods.                                                                            |
| `image`      | Docker image for the pod. For private images, use the schema: `<% REGISTRY %>/some/path/image:tag`. The `<% REGISTRY %>` token will be replaced dynamically.    |
| `entrypoint` | Command to replace the image’s ENTRYPOINT.                                                                                                                        |
| `command`    | Command to replace the image’s CMD.                                                                                                                               |

---

### Pod Volumes Configuration

| Key         | Description                                                                                      |
|-------------|--------------------------------------------------------------------------------------------------|
| `name`      | Volume name (lowercase, alphanumeric, and `-`).                                                  |
| `size`      | Required volume size (e.g., `1Gi`, `500Mi`).                                                     |
| `mountPath` | Mount path for the volume (must start with `/`).                                                |

---

### Pod Secrets Configuration

| Key         | Description                                                                                                                                       |
|-------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| `name`      | Secret name (lowercase, alphanumeric, and `-`).                                                                                                   |
| `data`      | Raw text or Base64-encoded string for the secret (e.g., encoded JSON).                                                                            |
| `mountPath` | Mount path where the secret file will be stored (must start with `/`).                                                                            |
| `fileName`  | Name of the secret file (e.g., `secret-file.txt`). The file will be available at `<mountPath>/<fileName>`.                                         |

---

### Pod Environment Variables Configuration

| Key    | Description                                                                                                                                                                         |
|--------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `key`  | Environment variable name.                                                                                                                                                          |
| `value`| Value of the environment variable. Can reference other pods using `<pod-name>.pod` (e.g., `http://express.pod:3000`) or the deployment’s base URL using `<% URL %>` (e.g., `<% URL %>/api`). |

---



## Application Configuration

### Required Fields

| Field | Description | Required |
|-------|-------------|----------|
| application.name | Your application name | ✅ |
| pods | List of pod configurations | ✅ |

### Optional Fields

| Field | Description | Default |
|-------|-------------|---------|
| application.url | Application domain | Auto-generated |
| registryLogin | Private registry config | None |

## Pod Configuration

Each pod represents a container in your application. Here's a detailed pod configuration:

```yaml
pods:
  - name: web-service
    image: nginx:latest
    path: /
    servicePorts:
      - 80
    vars:
      API_URL: http://api.pod:3000
    volumes:
      - name: data
        size: 1Gi
        mountPath: /data
```

### Pod Fields Explained

| Field | Description | Required |
|-------|-------------|----------|
| name | Pod identifier | ✅ |
| image | Container image | ✅ |
| path | URL path for routing | ❌ |
| servicePorts | Exposed ports | ✅ |
| vars | Environment variables | ❌ |
| volumes | Persistent storage | ❌ |

## Environment Variables

Define environment variables using the `vars` field:

```yaml
vars:
  DATABASE_URL: postgresql://user:pass@db.pod:5432/dbname
  NODE_ENV: production
```

## Service Ports

Configure network ports for your containers:

```yaml
servicePorts:
  - 8080
  - 8081
```

## Secrets Management

Store sensitive file information securely:

```yaml
secrets:
  - name: firebase-config
    data: base64_encoded_json_string
    mountPath: /var/secrets/firebase
    fileName: firebase_config.json
```

## Registry Authentication

Configure private registry access:

```yaml
registryLogin:
  registry: registry.example.com
  username: your_username
  personalAccessToken: pat_token_here
```

## Examples

### Simple Hello World App using Next.js

```yaml
application:
  name: "Hello World NextJS App"
  pods:
  - name: nextjs-nginx 
    path: /            
    image: ghcr.io/nexlayer/hello-world-nextjs:v0.0.1
    servicePorts:
    - 80
```

### Fullstack App using Next.js, Prisma, and Postgres with Private Images

```yaml
application:
  name: "Todo Nextjs"         # Required: Application name
  url: www.todo.com           # Optional: Custom domain URL (omit if not needed)
  registryLogin:              # Optional: Required only for private images
    registry: myregistry.com  # Required if registryLogin present: Registry hostname
    username: MyUsername      # Required if registryLogin present: Registry username
    personalAccessToken: myaccesstoken  # Required if registryLogin present: Read-only registry PAT

  pods:  # Required: List of pod configurations
  - name: web-app           # Required: Pod name (lowercase, alphanumeric, '-', '.')
    path: /
    image: <% REGISTRY %>/nextjs/app:latest  # Required: Full image URL following <% REGISTRY %>/repo:<tag> format for private images.  Must include '<% REGISTRY %>' exactly as shown.
    volumes:  # Optional: List of persistent storage volumes
    - name: nextjs-cache    # Required: Volume name (lowercase, alphanumeric, '-')
      size: 1Gi             # Required: Volume size (e.g., "1Gi", "500Mi")
      mountPath: /app/.next/cache  # Required: Volume mount path (must start with '/')
    
    secrets:  # Optional: List of secret configurations
    - name: nextauth-secret  # Required: Secret name (lowercase, alphanumeric, '-')
      data: myrandomsecret   # Required: Raw text or Base64-encoded secret content
      mountPath: /var/secrets/nextauth  # Required: Secret mount path (must start with '/')
      fileName: secret.txt   # Required: Secret file name
    
    vars:  # Optional: Environment variables
    - key: API_URL
      value: http://backend.pod:3001  # Reference another pod dynamically
    - key: NEXTAUTH_URL
      value: <% URL %>
    - key: DATABASE_URL
      value: postgresql://postgres:postgres@database.pod:5432/mydb
    - key: NEXTAUTH_SECRET
      value: <% URL %>/secrets/nextauth
    - key: GITHUB_CLIENT_ID
      value: Github_Client_ID         # GitHub OAuth Client ID
    - key: GITHUB_CLIENT_SECRET
      value: Github_Client_Secret     # GitHub OAuth Client Secret
    
    servicePorts:  # Required: List of port configurations (shorthand supported)
    - 3000  # Exposing Next.js frontend on port 3000

  - name: backend  # Required: Backend service (Next.js API routes)
    image: <% REGISTRY %>/node/api:latest # Private image
    volumes:
    - name: prisma-migrations
      size: 1Gi
      mountPath: /app/prisma/migrations
    vars:
    - key: DATABASE_URL
      value: postgresql://postgres:postgres@database.pod:5432/mydb
    - key: NEXTAUTH_SECRET
      value: <% URL %>/secrets/nextauth
    - key: GITHUB_CLIENT_ID
      value: Github_Client_ID
    - key: GITHUB_CLIENT_SECRET
      value: Github_Client_Secret
    servicePorts:
    - 3001  # Exposing backend API on port 3001

  - name: database  # Required: PostgreSQL database pod
    image: postgres:latest  # Public image
    volumes:
    - name: postgres-data
      size: 5Gi
      mountPath: /var/lib/postgresql
    vars:
    - key: POSTGRES_USER
      value: postgres
    - key: POSTGRES_PASSWORD
      value: postgres
    - key: POSTGRES_DB
      value: mydb
    servicePorts:
    - 5432  # Exposing PostgreSQL database on port 5432
```

### Example Fullstack AI-Powered Web Application Architecture

1. Web Application Pod (Frontend & API)

```yaml
pods:
  - name: web-app
    image: mycompany/web-app:latest
    path: /
    servicePorts:
      - 3000
    vars:
      DATABASE_URL: postgresql://postgres:postgres@postgres.pod:5432/postgres
      AI_API_URL: http://ai-model.pod:5000
      REDIS_HOST: redis.pod
      REDIS_PORT: "6379"
      REDIS_AUTH: myredissecret
```

2. AI Model Inference Pod

```yaml
pods:
  - name: ai-model
    image: huggingface/transformers:latest
    servicePorts:
      - 5000
    vars:
      MODEL_NAME: bert-base-uncased
      GPU_ENABLED: "true"
```

3. Data Processing Worker Pod

```yaml
pods:
  - name: data-worker
    image: mycompany/data-processor:latest
    servicePorts:
      - 4000
    vars:
      VECTOR_DB_URL: http://vector-db.pod:8080
      REDIS_HOST: redis.pod
```

4. Vector Database Pod

```yaml
pods:
  - name: vector-db
    image: weaviate/weaviate:latest
    servicePorts:
      - 8080
    vars:
      PERSISTENCE_MODE: disk
      ENABLE_MODULES: text2vec-transformers
```

5. Redis Pod (Caching & Message Queue)

```yaml
pods:
  - name: redis
    image: redis:7
    command: redis-server --requirepass myredissecret
    servicePorts:
      - 6379
    vars:
      REDIS_AUTH: myredissecret
```

6. PostgreSQL Pod (Relational Database)

```yaml
pods:
  - name: postgres
    image: postgres:latest
    servicePorts:
      - 5432
    vars:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: ai_application
    volumes:
      - name: postgres-data
        size: 10Gi
        mountPath: /var/lib/postgresql
```

7. MinIO Pod (Object Storage)

```yaml
pods:
  - name: minio
    image: minio/minio:latest
    entrypoint: /bin/sh
    command: -c "mkdir -p /data/appdata && minio server --address ':9000' --console-address ':9001' /data"
    servicePorts:
      - 9000
      - 9001
    vars:
      MINIO_ROOT_USER: minio
      MINIO_ROOT_PASSWORD: miniosecret
    volumes:
      - name: minio-data
        size: 10Gi
        mountPath: /data
```

#### Pod Communication Flow

- Web App → AI Model: Sends inference requests (web-app → ai-model.pod:5000)
- AI Model → Vector DB: Retrieves embeddings (ai-model → vector-db.pod:8080)
- Web App → PostgreSQL: Stores application data (web-app → postgres.pod:5432)
- Web App → Redis: Caches results (web-app → redis.pod:6379)
- Data Worker → MinIO: Processes large datasets (data-worker → minio.pod:9000)

## Best Practices

### 🔐 Security

- Always use secrets for sensitive data
- Use specific image tags instead of latest
- Implement least privilege access

### 🚀 Performance

- Configure appropriate resource limits
- Use persistent volumes for stateful applications
- Implement health checks

### 🔄 Maintenance

- Document all environment variables
- Use meaningful pod names
- Keep configurations version controlled

## Need Help?

- 📚 [Official Documentation](https://github.com/Nexlayer/templates/)
- 💬 [Community](https://github.com/orgs/Nexlayer/discussions)
- 🐛 [Issues/Feedback](https://github.com/Nexlayer/templates/issues)

---

Made with ❤️ by the Nexlayer Team
```