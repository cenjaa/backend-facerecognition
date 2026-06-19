# VPS Deployment Guide — Backend Stack

This guide covers deploying the **Go backend** and all infrastructure
(PostgreSQL, MinIO, Redis) from this `deploy/` folder.

The **ML service** is deployed separately from `ml-service/deploy/`.
Both stacks share a Docker network so they can communicate by container name.

---

## Architecture

```
Internet
   │
 Nginx (80 / 443)
   │
   └── → backend (127.0.0.1:5000)
              │
              │  http://ml-service:8001  ← shared Docker network
              ▼
        ml-service (127.0.0.1:8001 — internal only)

Backend Compose Stack (backend-facerecognition/deploy/)
  ├── backend          (Go/Fiber API)
  ├── postgres         (PostgreSQL 16)
  ├── minio            (Object Storage)
  ├── redis            (Cache)
  └── minio-init       (one-shot: creates bucket)

ML Service Compose Stack (ml-service/deploy/)
  └── ml-service       (Python FastAPI)

Shared Docker network: app_network
  (both stacks join this network)
```

---

## PART 1 — One-Time VPS Setup

### Step 1 — Install Docker & Nginx

```bash
sudo apt-get update && sudo apt-get upgrade -y

# Install Docker
sudo apt-get install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | \
    sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
  https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Run Docker without sudo
sudo usermod -aG docker $USER && newgrp docker

# Install Nginx
sudo apt-get install -y nginx
sudo systemctl enable nginx
```

### Step 2 — Configure Firewall

```bash
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
sudo ufw status
# Only ports 22, 80, 443 should be open publicly
```

### Step 3 — Create the Shared Docker Network

```bash
docker network create app_network
```

> The `deploy.sh` script creates this automatically if it doesn't exist.

---

## PART 2 — Transfer Files to VPS (Termius + SFTP)

Use **Termius** to connect to your VPS via SSH for terminal commands,
and its built-in **SFTP** panel to drag-and-drop files.

### Step 1 — Connect in Termius

1. Open Termius → **New Host**
2. Fill in your VPS IP, SSH port (22), username, and password/key
3. Click **Connect** — you now have a terminal on your VPS

### Step 2 — Create target directories on VPS

In the Termius terminal, run:

```bash
sudo mkdir -p /opt/backend
sudo mkdir -p /opt/ml-service
sudo chown -R $USER:$USER /opt/backend /opt/ml-service
```

### Step 3 — Upload files via SFTP

1. In Termius, open the **SFTP** tab (or right-click the host → SFTP)
2. Left panel = your local machine, Right panel = VPS

**Upload Backend:**
- Left panel: navigate to
  `C:\Users\Farchan Indrianto\Documents\IPB\Semester 8\TA\backend-facerecognition`
- Right panel: navigate to `/opt/backend`
- Select all files/folders **except**: `venv/`, `.git/`, `dataset/`, `videos/`, `tmp.exe`
- Drag-and-drop to upload

**Upload ML Service:**
- Left panel: navigate to
  `C:\Users\Farchan Indrianto\Documents\IPB\Semester 8\TA\ml-service`
- Right panel: navigate to `/opt/ml-service`
- Select all files/folders **except**: `venv/`, `.git/`, `dataset/`, `debug_output/`, `__pycache__/`
- Drag-and-drop to upload

### Step 4 — Upload the SQL dump

- Left panel: navigate to
  `C:\Users\Farchan Indrianto\Documents\IPB\Semester 8\TA\backend-facerecognition\deploy\init-db`
- Right panel: navigate to `/opt/backend/deploy/init-db`
- Upload `dump-ATTENDANCE_DB.sql`

> **Tip:** On subsequent deploys you only need to re-upload files that changed.
> The `venv/` and `node_modules/` equivalents are rebuilt inside Docker — no need to upload them.

---

## PART 3 — Configure Secrets on VPS

SSH into the VPS, then:

```bash
cd /opt/backend/deploy

# Create .env from template
cp .env.example .env
nano .env
```

Fill in real passwords:
```env
DATABASE_NAME=ATTENDANCE_DB
DATABASE_USER=admin
DATABASE_PASSWORD=StrongPass123!

MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=StrongMinioSecret456!
MINIO_BUCKET=attendance

REDIS_PASSWORD=StrongRedisPass789!

JIRA_TOKEN=your_real_jira_token
```

```bash
# Create config.yaml from template
cp config.yaml.example config.yaml
nano config.yaml
# Only replace CHANGE_ME values — service hostnames (postgres, minio, redis) are correct as-is
```

---

## PART 4 — Deploy Backend Stack

```bash
cd /opt/backend/deploy

# First deploy — build images and start everything
bash deploy.sh --build

# Watch startup sequence
docker compose logs -f
```

Expected startup output:
```
postgres    | database system is ready to accept connections  ✅
minio       | API: http://0.0.0.0:9000                        ✅
minio-init  | MinIO bucket ready.                             ✅
redis       | Ready to accept connections                     ✅
backend     | Database Connection Success.                    ✅
backend     | Minio Connection (minio:9000) Success.          ✅
backend     | Redis Connection Success.                       ✅
backend     | HTTP VPS server is running  port=5000           ✅
```

---

## PART 5 — Configure Domain & Nginx with API Key

Your backend API is now secured by an API Key at the Nginx level.
All requests to `https://api.chanfolio.my.id/` must include the header `X-API-Key`.

### Step 1 — DNS Configuration
Before configuring Nginx, log in to your domain registrar and create an **A Record**:
- **Name/Host:** `api`
- **Target/IP:** `20.2.86.39`

### Step 2 — Verify Your API Key
The API key has already been configured locally in `backend-facerecognition/deploy/nginx.conf`:

```nginx
if ($http_x_api_key != "f663f96eb5bacad610d978e1e62de8a96fd2f186b638e87376dbb4826154f65a") {
```

Make sure you upload this newly updated `nginx.conf` to the VPS (`/opt/backend/deploy/nginx.conf`) via Termius SFTP before proceeding.

### Step 3 — Apply Nginx Config
```bash
# Copy nginx config
sudo cp /opt/backend/deploy/nginx.conf /etc/nginx/sites-available/backend

# Enable and reload
sudo ln -s /etc/nginx/sites-available/backend /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl reload nginx
```

### Step 3 — Enable HTTPS (SSL Certificate)
Run Certbot to secure `api.chanfolio.my.id` with HTTPS.

```bash
sudo apt-get install -y certbot python3-certbot-nginx
sudo certbot --nginx -d api.chanfolio.my.id
```
Certbot will automatically modify your Nginx config to add the SSL certificates.

---

## PART 6 — Deploy ML Service

```bash
cd /opt/ml-service/deploy
cp .env.example .env
nano .env
# Use the same MinIO credentials as backend

bash deploy.sh --build
```

---

## PART 7 — Verify Everything

```bash
# All containers running
docker compose -f /opt/backend/deploy/docker-compose.yml ps
docker compose -f /opt/ml-service/deploy/docker-compose.yml ps

# Backend can reach ML service
docker exec $(docker ps -qf "name=backend") \
    sh -c "wget -qO- http://ml-service:8001/ && echo '✅ ML reachable'"

# Public endpoint responds
curl -s -o /dev/null -w "HTTP %{http_code}\n" http://YOUR_VPS_IP/

# MinIO console (via SSH tunnel)
# ssh -L 9001:127.0.0.1:9001 YOUR_USER@YOUR_VPS_IP
# Then open: http://localhost:9001
```

---

## Day-to-Day Operations

### Update backend code
```bash
cd /opt/backend/deploy
bash deploy.sh --build
```

### Update ML service code
```bash
cd /opt/ml-service/deploy
bash deploy.sh --build
```

### View live logs
```bash
docker compose -f /opt/backend/deploy/docker-compose.yml logs -f
docker compose -f /opt/ml-service/deploy/docker-compose.yml logs -f
```

### Stop everything
```bash
docker compose -f /opt/backend/deploy/docker-compose.yml down
docker compose -f /opt/ml-service/deploy/docker-compose.yml down
```

> ⚠️ Use `down` without `--volumes` to **keep** your data.
> Add `--volumes` only to **wipe** all data and start fresh.

---

## (Optional) Enable HTTPS

Requires a domain name pointing to your VPS IP:

```bash
sudo apt-get install -y certbot python3-certbot-nginx
sudo certbot --nginx -d your.domain.com
sudo certbot renew --dry-run
```

---

## Security Checklist

- [ ] `.env` files are on VPS only — never committed to git
- [ ] `config.yaml` is on VPS only — never committed to git
- [ ] Firewall: only ports 22, 80, 443 open (`sudo ufw status`)
- [ ] Backend port 5000 and ML port 8001 are `127.0.0.1` only
- [ ] MinIO port 9000 has no public port binding (internal only)
- [ ] MinIO console port 9001 bound to `127.0.0.1` (use SSH tunnel to access)

**Docker log rotation:**
```bash
sudo tee /etc/docker/daemon.json > /dev/null <<'EOF'
{
  "log-driver": "json-file",
  "log-opts": { "max-size": "10m", "max-file": "3" }
}
EOF
sudo systemctl restart docker
```
