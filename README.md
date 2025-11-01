# DevOps FMI Project

Google Cloud Platform (GCP), Next.js, Go, PostgreSQL (Cloud SQL), Kubernetes (GKE), OpenTofu (Terraform).

---

## Frontend

### Tech Stack

- Next.js 15
- TypeScript and Tailwind CSS
- Vercel Hosting

### Description

Lightweight dashboard for monitoring backend and database health.

- `/healthz` for backend availability
- `/checkDatabase` for database connectivity

Local development:

```bash
cd frontend
npm install
npm run dev
```

Environment variable:

```bash
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
```

Runs locally on http://localhost:3000

---

## Backend

### Tech Stack

- Go 1.24
- Fiber v2
- PostgreSQL
- Docker and Kubernetes
- Terraform and GCP Artifact Registry
- Gateway API (modern replacement for Kubernetes Ingress): https://docs.cloud.google.com/kubernetes-engine/docs/concepts/gateway-api

### Description

The backend provides a simple REST API with the following endpoints:

| Endpoint         | Method | Description                       |
| ---------------- | ------ | --------------------------------- |
| `/healthz`       | GET    | Basic health check                |
| `/checkDatabase` | GET    | Verifies connection to PostgreSQL |
| `/login`         | POST   | Simulated login route             |
| `/register`      | POST   | Simulated user registration       |

---

### Environment Variables

| Variable          | Description          | Example                               |
| ----------------- | -------------------- | ------------------------------------- |
| `PORT`            | Server port          | `8080`                                |
| `ALLOWED_ORIGINS` | Allowed CORS origin  | `https://devops-chi-khaki.vercel.app` |
| `DB_HOST`         | Cloud SQL private IP | `postgres`                            |
| `DB_PORT`         | PostgreSQL port      | `5432`                                |
| `DB_NAME`         | Database name        | `postgres`                            |
| `DB_USER`         | Database username    | `postgres`                            |
| `DB_PASSWORD`     | Database password    | `postgres`                            |

---

### Local Development

Run locally using Docker Compose:

```bash
docker compose -f docker-compose.local.yaml up --build
```

This starts:

- Backend
- PostgreSQL

Test endpoints:

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/checkDatabase
```

---

## Deployment

There is currently no automatic CI/CD pipeline. Deployments are manual and require the following preconditions:

- Logged in with `gcloud`.
- Connected to the correct GKE cluster.
- Authenticated for GKE Artifact Registry.
- Infrastructure provisioned via OpenTofu (Terraform).
- Docker images built and pushed manually to Artifact Registry.
- Kubernetes manifests applied manually.

1. Build and push the backend image:

```bash
docker build -t europe-west4-docker.pkg.dev/devops-fmi-course-476112/devops-fmi-course-repo/backend-app:latest ./backend
docker push europe-west4-docker.pkg.dev/devops-fmi-course-476112/devops-fmi-course-repo/backend-app:latest
```

2. Apply Kubernetes manifests:

```bash
kubectl apply -f k8s/
```

---

## GKE and Cloud SQL Setup

- Private VPC (`10.0.0.0/16`)
- GKE cluster with Workload Identity
- Cloud SQL (PostgreSQL 15) with private IP only, accessible from GKE via `10.41.0.3`
- VPC peering through `servicenetworking.googleapis.com`
- HTTPS Gateway (Google Cloud Load Balancer) with Google-managed TLS

### Networking and HTTPS Flow

- The domain `api.users.gopherify.com` maps to a **reserved static IP** (`35.244.143.113`) in Google Cloud.
- This IP is bound to a **Gateway** resource configured with `gatewayClassName: gke-l7-global-external-managed`.
- TLS termination uses a **Google-managed certificate** attached via `networking.gke.io/pre-shared-certs`.
- The Gateway forwards HTTPS traffic through an `HTTPRoute` to the backend `Service`, which exposes Pods on port `3000`.

This replaces the legacy Kubernetes Ingress and ManagedCertificate workflow with the new **Gateway API**, offering more flexibility, scalability, and control.

---

## Security

- All sensitive environment variables are stored in **Kubernetes Secrets**.
- GKE communicates with **Cloud SQL over private IP** (no public access).
- **TLS termination** is handled by Google Cloud using a **managed certificate** for `api.users.gopherify.com`.
- **Workload Identity** is enabled for secure service-to-service communication without key files.

---

## Future Improvements

- Automate Kubernetes deployments.
- Introduce Horizontal Pod Autoscaling.
- Add monitoring and observability.

---

## Summary

| Component  | Status               | Notes                     |
| ---------- | -------------------- | ------------------------- |
| Frontend   | Deployed on Vercel   | Uses Next.js 15           |
| Backend    | Running in GKE       | HTTPS + CORS configured   |
| Database   | Cloud SQL Private IP | Connected successfully    |
| HTTPS      | Managed by GCP       | Active certificate        |
| Networking | VPC Peering          | Secure private connection |
| Gateway    | Global LB on GKE     | Using static IP + TLS     |
