# Nexlayer YAML Reference (v1.2)

Nexlayer Cloud makes deploying applications easier by handling the complicated parts of Kubernetes for you. Instead of setting up everything manually, you can use a simple `nexlayer.yaml` file to define your app.

## Quick Start

```yaml
application:
  name: my-app
  pods:
    - name: web
      type: react
      image: ghcr.io/myorg/web:latest
      ports:
        - containerPort: 3000
          servicePort: 80
          name: web
```

## Schema Reference

### Basic App Information
```yaml
application:
  name: Example App  # Name of the app
  url: www.example.ai  # (Optional) Permanent domain
```

### Private Registry Login (Optional)
```yaml
registryLogin:
  registry: ghcr.io
  username: SomeUser1234
  personalAccessToken: ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Pod Configuration
```yaml
pods:
  - name: web  # Pod name (lowercase, can use - or .)
    type: react  # Pod type (react, django, postgres, etc.)
    path: /  # Public route (optional)
    image: ghcr.io/myorg/web:latest  # Docker image
    ports:
      - containerPort: 3000  # Container's port
        servicePort: 80      # External port
        name: web           # Port name
    volumes:  # Optional persistent storage
      - name: data
        size: 1Gi
        mountPath: /data
    secrets:  # Optional secure data
      - name: api-key
        data: "secret-value"
        mountPath: /secrets
        fileName: api-key.txt
    vars:  # Optional environment variables
      - key: DATABASE_URL
        value: postgres://db.pod:5432
```

## Best Practices

1. **Image Names**
   - Public images: `repo/image:tag`
   - Private images: `ghcr.io/repo/image:tag`

2. **Service Discovery**
   - Use `<pod-name>.pod` for internal connections
   - Example: `postgres://db.pod:5432`

3. **Port Numbers**
   - Web servers: 3000, 8000, 8080
   - Databases: 5432 (Postgres), 6379 (Redis)

## Common Issues

| Problem | Solution |
|---------|----------|
| Service not found | Use `<pod-name>.pod` for internal URLs |
| Image pull error | Check registry login for private images |
| Storage issues | Ensure volume paths start with `/` |

## Command Reference

```sh
# Deploy your app
nexlayer deploy

# Check deployment status
nexlayer status

# View logs
nexlayer logs -f
