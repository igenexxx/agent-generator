---
name: go-adk-deploy
description: >
  Go ADK v2.0 Containerization and Cloud Run deployment guidelines.
  Provides optimized multi-stage Dockerfiles, port binding setups, and GCP deployment commands.
metadata:
  author: Antigravity
  license: Apache-2.0
  version: 2.0.0
---

# Containerizing & Deploying Go ADK Agents

This skill provides the deployment guide and infrastructure templates for packaging and deploying Go ADK v2.0 agents on Google Cloud Platform.

---

## 1. Multi-Stage Dockerfile (Optimized for Go)

Use a multi-stage Docker build to keep container size to a minimum (~20MB) and ensure a secure, non-root runtime environment. Save this as `Dockerfile`:

```dockerfile
# Stage 1: Build the static Go binary
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# Copy dependencies lists
COPY go.mod go.sum ./
RUN go mod download

# Copy source code files
COPY . .

# Compile static binary with flags to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o agent-bin main.go

# Stage 2: Packaging using a secure scratch or distroless base
FROM gcr.io/distroless/static-debian12:latest-amd64

WORKDIR /

# Copy the compiled binary from Stage 1
COPY --from=builder /app/agent-bin /agent-bin

# Expose port (default Cloud Run port is 8080)
EXPOSE 8080

# Run as non-root user
USER 65532:65532

# Set entrypoint
ENTRYPOINT ["/agent-bin", "web", "api"]
```

---

## 2. Port Binding Convention

The launcher needs to bind to the port specified in the `$PORT` environment variable (which is dynamically injected by Cloud Run). Ensure your `main.go` or startup launcher configuration supports this:

```go
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}

// Bind using launcher config flag
// Example launcher execution:
// ./agent-bin web -port <port> api
```

---

## 3. Deploying to Cloud Run

Deploy your container directly to Cloud Run using the `gcloud` CLI:

```bash
# 1. Build and push container to Artifact Registry
gcloud builds submit --tag gcr.io/YOUR_PROJECT_ID/go-adk-agent:latest

# 2. Deploy to Cloud Run
gcloud run deploy go-adk-agent \
    --image gcr.io/YOUR_PROJECT_ID/go-adk-agent:latest \
    --platform managed \
    --region us-central1 \
    --allow-unauthenticated
```

---

## 4. Secret Management (Google Cloud Secret Manager)

Do not ship `.env` files inside your Docker containers. Mount your API keys securely using GCP Secret Manager:

```bash
# Mount the secret to GOOGLE_API_KEY environment variable on Cloud Run
gcloud run deploy go-adk-agent \
    --image gcr.io/YOUR_PROJECT_ID/go-adk-agent:latest \
    --update-secrets=GOOGLE_API_KEY=gemini-api-key-secret:latest
```
This mounts the secret value dynamically, injecting it directly into the process environment variables at startup.
