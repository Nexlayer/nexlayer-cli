# Nexlayer YAML Schema Documentation (v1.2)

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
- **Pod-to-pod discovery** (connecting your app's parts together)

It's designed for:
✅ Developers (easy to understand and edit)  
✅ AI Systems (clear structure for automation)  
✅ Machines (for CI/CD and deployment tools)

---

## 🚀 YAML Template Breakdown

### 📌 Basic App Information

```yaml
application:
  name: Example App  # Name of the app
  url: www.example.ai  # (Optional) Permanent domain
```

| Key      | Description |
|----------|------------|
| `name`   | The name of the application |
| `url`    | A permanent domain (optional) |

---

### 🔐 Private Registry Login

```yaml
registryLogin:
  registry: ghcr.io
  username: SomeUser1234
  personalAccessToken: ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxx
```

| Key                  | Description |
|----------------------|-------------|
| `registry`           | The registry storing the private images (e.g., `docker.io`, `ghcr.io`) |
| `username`           | Login username |
| `personalAccessToken`| Read-only token for access |

🛑 **Is This Required?**

| Image Type       | Is `registryLogin` Needed? |
|------------------|--------------------------|
| **Private Images** | ✅ Yes |
| **Public Images**  | ❌ No, remove this section |

---

### 📦 Defining Pods (Services)

```yaml
pods:
  - name: react  # Pod name
```

| Key   | Description |
|-------|-------------|
| `name` | The pod's unique name (lowercase, can use `-` or `.`) |

Each **pod** represents a different part of your app (frontend, backend, database, etc.).

---

### 🌎 Routing for Frontend Pods

```yaml
path: /
```

| Key   | Description |
|-------|-------------|
| `path` | Defines the public-facing route (use `/` for frontend) |

✅ **Examples:**
- `/` → Main website
- `/api` → Backend API

---

### 📌 Docker Image Definition

```yaml
image: <% REGISTRY %>/someUser1234/image:tag
```

| Key   | Description |
|-------|-------------|
| `image` | Docker image for the pod (use `<% REGISTRY %>` for private images) |

✅ **Best Practice:**
- Use `repo/image:tag` for **public images**
- Use `ghcr.io/repo/image:tag` for **private images**

---

### 💾 Persistent Storage (Volumes)

```yaml
volumes:
  - name: volume
    size: 1Gi
    mountPath: /var/some/directory
```

| Key        | Description |
|------------|-------------|
| `name`     | Volume name |
| `size`     | Size (e.g., `1Gi`) |
| `mountPath` | Where it's mounted inside the pod |

✅ **Use This For:**
- Storing database files
- Logs
- Anything that should persist between restarts

---

### 🔑 Secrets (For Secure Data)

```yaml
secrets:
  - name: my-secret
    data: "My secret text"
    mountPath: /var/secrets/my-secret-volume
    fileName: secret-file.txt
```

| Key        | Description |
|------------|-------------|
| `name`     | Secret name |
| `data`     | Plain text or Base64-encoded value |
| `mountPath` | Where it's stored in the pod |
| `fileName`  | File name in the pod |

✅ **Use This For:**
- API keys
- Passwords
- Configurations that must be secure

---

### 🌍 Environment Variables

```yaml
vars:
  - key: API_URL
    value: http://express.pod:3000
  - key: SITE_URL
    value: <% URL %>/api
```

| Key   | Description |
|-------|-------------|
| `key` | Environment variable name |
| `value` | Value (must be a string, supports `<pod-name>.pod`) |

✅ **Best Practices:**
- Use `postgres.pod:5432/dbname` instead of just `postgres:5432/dbname`
- Use `<% URL %>` for auto-generating site links

---

### 📡 Ports for Services

```yaml
servicePorts:
  - 3000
```

| Key   | Description |
|-------|-------------|
| `servicePorts` | List of ports exposed by the pod |

✅ **Examples:**
- `3000` → Web server
- `5432` → PostgreSQL database
- `6379` → Redis

---

### **📌 Common Fixes & Best Practices**

| Problem | Solution |
|---------|---------|
| Service names don't resolve | Use `<pod-name>.pod` for internal communication |
| Wrong image format | Use `repo/image:tag` for public, `ghcr.io/repo/image:tag` for private |
| Manual Kubernetes setup | Just run `nexlayer deploy`, no need for Kubernetes YAML |

✅ **Why Use Nexlayer?**
- No Kubernetes experience needed 🚀
- Built-in **service discovery** (no `depends_on` needed!)
- **Automated storage** (just define `size` and `mountPath`)
- **Environment variable templating** (`<% URL %>` for easy links)

💡 **Deploy with one command:**
```sh
nexlayer deploy
```

Your app is live in seconds! 🎉

---

## 🛠 Troubleshooting

| Issue | Solution |
|-------|---------|
| My pod isn't reachable by its service name | Ensure you reference the pod with `<pod-name>.pod` in your endpoints. |
| Why isn't my private image being pulled? | Confirm your image tag starts with `<%registry%>` to trigger the private image workflow. |
| Hardcoded URLs are causing configuration issues | Replace any instances of `localhost:3000` with `<%url%>` for proper routing. |
