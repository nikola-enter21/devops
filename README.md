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

### Deployment Flow

There is **no automatic deployment pipeline configured yet**.  
Before running any deployment commands, the following conditions must be met:

- You are **logged in with gcloud**.
- You are **connected to the correct GKE cluster**.
- You are **authenticated for the GKE Artifact Registry**.
- GCP infrastructure has been **provisioned with OpenTofu**.
- Docker images have been **built and pushed manually** to Artifact Registry.
- Kubernetes manifests will be **applied manually** to GKE.

1. Build and push backend image:

```bash
docker build -t europe-west4-docker.pkg.dev/devops-fmi-course-476112/devops-fmi-course-repo/backend-app:latest ./backend
docker push europe-west4-docker.pkg.dev/devops-fmi-course-476112/devops-fmi-course-repo/backend-app:latest
```

2. Deploy to GKE manually:

```bash
kubectl apply -f k8s/users-deployment.yml
kubectl apply -f k8s/users-service.yml
kubectl apply -f k8s/users-backendconfig.yml
kubectl apply -f k8s/managed-cert.yml
kubectl apply -f k8s/users-ingress.yml
```

---

### GKE and Cloud SQL Setup

- Private VPC (10.0.0.0/16)
- GKE cluster with Workload Identity (Google Cloud’s modern, secure way for Kubernetes Pods to access Google Cloud resources (like Cloud SQL, Pub/Sub, Storage, etc.) without using service account keys)
- Cloud SQL (PostgreSQL 15) with private IP only  
  Accessible from GKE over 10.41.0.3
- VPC peering via servicenetworking.googleapis.com
- HTTPS Load Balancer with managed TLS certificate

### Networking and HTTPS Flow

- **A-Record** maps `api.users.gopherify.com` to the **static external IP** (`35.244.143.113`) reserved in Google Cloud.  
  This IP is attached to the **GKE Ingress Load Balancer** that Kubernetes creates automatically.

- The **ManagedCertificate** automatically provisions a valid TLS certificate for `api.users.gopherify.com`, enabling secure HTTPS connections.
  After validation, Google’s global load balancer handles HTTPS traffic and forwards decrypted requests to the backend service inside GKE.

- The **Ingress (GCE)** accepts traffic on ports **80 (HTTP)** and **443 (HTTPS)** specifically for `api.users.gopherify.com`.  
  Inside the cluster, the Ingress forwards those requests to `users-service:80`, which then routes internally to backend Pods running on port `3000`.

---

### Security

- Database credentials are currently **hardcoded in the Deployment manifest** for simplicity.  
  In a production environment, these must be stored in a **Kubernetes Secret** or **Google Secret Manager** instead.

- GKE communicates with **Cloud SQL over a private VPC IP**, ensuring no external exposure.

- The **Cloud SQL instance** is private-only --- it cannot be accessed from the public Internet.

- **HTTPS** traffic is managed entirely by a **Google-managed TLS certificate** attached to the global load balancer.

- **Workload Identity** is enabled, replacing legacy service account keys for secure GCP service authentication.

---

### Future Improvements

- Automate Kubernetes deployment

---

### Summary

| Component  | Status               | Notes                     |
| ---------- | -------------------- | ------------------------- |
| Frontend   | Deployed on Vercel   | Uses Next.js 15           |
| Backend    | Running in GKE       | HTTPS + CORS configured   |
| Database   | Cloud SQL Private IP | Connected successfully    |
| HTTPS      | Managed by GCP       | Active certificate        |
| Networking | VPC Peering          | Secure private connection |
