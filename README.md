# DevOps

Google Cloud Platform (GCP), Next.js, Go, PostgreSQL (Cloud SQL), Kubernetes (GKE), OpenTofu (Terraform), Open Policy Agent (OPA), Buf

---

## Frontend

### Tech Stack

- Next.js 15
- TypeScript and Tailwind CSS
- Vercel Hosting

### Description

Lightweight dashboard for monitoring backend and database health.

- `/api/v1/healthz` → backend availability
- `/api/v1/checkDatabase` → database connectivity

#### Local development

```bash
cd frontend
npm install
npm run dev
```

Environment variable:

```bash
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
```

Runs locally at [http://localhost:3000](http://localhost:3000)

---

## Backend

### Tech Stack

- Go 1.25.3
- gRPC + gRPC-Gateway v2
- Buf - modern Protocol Buffers toolchain
- Open Policy Agent (OPA)
- PostgreSQL (Cloud SQL Private IP)
- Docker and Kubernetes (GKE)
- Terraform and GCP Artifact Registry
- Gateway API (modern replacement for Kubernetes Ingress): https://docs.cloud.google.com/kubernetes-engine/docs/concepts/gateway-api

---

### Description

The backend now exposes both **gRPC** and **REST** endpoints using **gRPC-Gateway**, secured with **OPA Rego policies** for fine-grained access control.

| Endpoint                | Method | Description                       |
| ----------------------- | ------ | --------------------------------- |
| `/api/v1/healthz`       | GET    | Basic health check                |
| `/api/v1/checkDatabase` | GET    | Verifies connection to PostgreSQL |
| `/api/v1/login`         | POST   | Simulated login route             |
| `/api/v1/register`      | POST   | Simulated user registration       |

---

### Environment Variables

| Variable          | Description          | Example                               |
| ----------------- | -------------------- | ------------------------------------- |
| `HTTP_PORT`       | HTTP Gateway port    | `8080`                                |
| `GRPC_PORT`       | gRPC Server port     | `8079`                                |
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

- gRPC + Gateway backend
- PostgreSQL

Test endpoints:

```bash
curl http://localhost:8080/api/v1/healthz
curl http://localhost:8080/api/v1/checkDatabase
```

---

## Deployment

Manual deployment (CI/CD not yet automated).

1. Build and push the backend image:

- Logged in with `gcloud`
- Connected to correct GKE cluster
- Authenticated for Artifact Registry
- Infrastructure provisioned with OpenTofu
- Docker image built + pushed manually
- Kubernetes manifests applied manually

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
- The Gateway forwards HTTPS traffic through an `HTTPRoute` to the backend `Service`, which exposes Pods on port `8080`.

This replaces the legacy Kubernetes Ingress and ManagedCertificate workflow with the new **Gateway API**, offering more flexibility, scalability, and control.

---

## Security

- **OPA**: fine-grained RBAC via Rego policies
- All sensitive environment variables are stored in **Kubernetes Secrets**.
- GKE communicates with **Cloud SQL over private IP** (no public access).
- **TLS termination** is handled by Google Cloud using a **managed certificate** for `api.users.gopherify.com`.
- **Workload Identity** is enabled for secure service-to-service communication without key files.

---

## Future Improvements

- Automate Kubernetes deployments
- Add **Horizontal Pod Autoscaling (HPA)**
- Integrate **Ambient Service Mesh**
- Deploy **ArgoCD for GitOps** automation
- Extend **monitoring & observability**

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
