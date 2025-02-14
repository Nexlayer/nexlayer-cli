## Nexlayer YAML Schema Documentation (v1.2)

Nexlayer Cloud makes deploying applications easier by handling the complicated parts of Kubernetes for you. Instead of setting up everything manually, you can use a simple `nexlayer.yaml` file to define your app.

---

### 📜 What This Does

This YAML template helps you define:

- **Pods** (individual containers that run your app)
- **Storage** (saving data between restarts)
- **Secrets** (keeping passwords and keys safe)
- **Environment variables** (app settings)
- **Service ports** (how pods communicate)
- **Private registry login** (for private Docker images)
- **Pod-to-pod discovery** (connecting your app’s parts together)
- **Container command overrides** (customizing startup behavior)

It’s designed for: ✅ Developers (easy to understand and edit)\
✅ AI Systems (clear structure for automation)\
✅ Machines (for CI/CD and deployment tools)

---

## 🚀 YAML Template Breakdown

### 📌 Full YAML Structure

```yaml
application:
  name: string       # Required: Application name (lowercase, alphanumeric, '-', '.')
  url: string       # Optional: Custom domain URL
  registryLogin:    # Optional: Private registry authentication
    registry: string        # Required if registryLogin present: Registry hostname
    username: string        # Required if registryLogin present: Registry username
    personalAccessToken: string  # Required if registryLogin present: Registry PAT
  pods:            # Required: List of pod configurations
    - name: string        # Required: Pod name (lowercase, alphanumeric, '-', '.')
      path: string      # Optional: Mount path (must start with '/')
      image: string     # Required: Full image URL including registry and tag
      volumes:         # Optional: List of persistent storage volumes
        - name: string     # Required: Volume name (lowercase, alphanumeric, '-')
          size: string     # Required: Volume size (e.g., "1Gi", "500Mi")
          mountPath: string # Required: Volume mount path (must start with '/')
      secrets:         # Optional: List of secret configurations
        - name: string     # Required: Secret name (lowercase, alphanumeric, '-')
          data: string     # Required: Raw or Base64-encoded secret content
          mountPath: string # Required: Secret mount path (must start with '/')
          fileName: string  # Required: Secret file name
      ports:           # Required: List of port configurations
        - containerPort: int  # Required: Port inside the container
          servicePort: int    # Required: Port exposed to other services
          name: string        # Required: Unique name for the port
      vars:            # Optional: List of environment variables
        DATABASE_URL: postgresql://...
        SALT: mysalt
      servicePorts:           # Required: List of port configurations (shorthand supported)
        - 3030
      entrypoint: string      # Optional: Custom container entrypoint
      command: string         # Optional: Custom container command
```

---

### 📌 Supported Pod Types

Each pod has a `type` field that defines its role. Nexlayer supports the following pod types:

| Type       | Description                                                  |
| ---------- | ------------------------------------------------------------ |
| `frontend` | UI components (e.g., React, Vue, Angular)                    |
| `backend`  | APIs and server logic (e.g., Express, FastAPI, Django)       |
| `database` | Databases (e.g., PostgreSQL, MongoDB, Redis)                 |
| `proxy`    | Reverse proxy/load balancer (e.g., Nginx, Traefik)           |
| `worker`   | Background jobs or task queues (e.g., Celery, Sidekiq)       |
| `llm`      | AI/LLM models for inference (e.g., OpenAI, Llama, LangChain) |

✅ **Example: Defining a Backend Pod**

```yaml
pods:
  - name: api
    image: my-org/my-backend:latest
    servicePorts:
      - 8080
```

Note: The `type` field has become optional in practice.

---

### 📌 Basic App Information

```yaml
application:
  name: Example App  # Name of the app
  url: www.example.ai  # (Optional) Permanent domain
```

| Key    | Description                   |
| ------ | ----------------------------- |
| `name` | The name of the application   |
| `url`  | A permanent domain (optional) |

---

### 🔐 Private Registry Login

```yaml
registryLogin:
  registry: ghcr.io
  username: SomeUser1234
  personalAccessToken: ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxx
```

| Key                   | Description                                                            |
| --------------------- | ---------------------------------------------------------------------- |
| `registry`            | The registry storing the private images (e.g., `docker.io`, `ghcr.io`) |
| `username`            | Login username                                                         |
| `personalAccessToken` | Read-only token for access                                             |

🛑 **Is This Required?**

| Image Type         | Is `registryLogin` Needed? |
| ------------------ | -------------------------- |
| **Private Images** | ✅ Yes                      |
| **Public Images**  | ❌ No, remove this section  |

---

### 📦 Defining Pods (Services)

```yaml
pods:
  - name: react  # Pod name
```

| Key    | Description                                           |
| ------ | ----------------------------------------------------- |
| `name` | The pod’s unique name (lowercase, can use `-` or `.`) |

Each **pod** represents a different part of your app (frontend, backend, database, etc.).

---

💡 **Deploy with one command:**

```sh
nexlayer deploy
```

Your app is live in seconds! 🎉

