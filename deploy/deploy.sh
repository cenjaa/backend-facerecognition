#!/usr/bin/env bash
# =============================================================
#  deploy/deploy.sh — Backend Stack
#  Run this on your VPS to deploy or update the backend stack.
#
#  Usage:
#    bash deploy.sh            # Start / update (no image rebuild)
#    bash deploy.sh --build    # Force rebuild of Docker images
#    bash deploy.sh --down     # Stop containers (volumes kept)
# =============================================================

set -euo pipefail

DEPLOY_DIR="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="$DEPLOY_DIR/docker-compose.yml"
BUILD_FLAG=""
DOWN_MODE=false

for arg in "$@"; do
  case $arg in
    --build) BUILD_FLAG="--build" ;;
    --down)  DOWN_MODE=true ;;
  esac
done

echo "=============================================="
echo "  Backend Stack — VPS Deploy"
echo "=============================================="

# ── Sanity checks ──────────────────────────────────────────────
if [ ! -f "$DEPLOY_DIR/.env" ]; then
  echo "ERROR: .env not found in deploy/!"
  echo "  cp $DEPLOY_DIR/.env.example $DEPLOY_DIR/.env"
  echo "  Then fill in your secrets."
  exit 1
fi

if [ ! -f "$DEPLOY_DIR/config.yaml" ]; then
  echo "ERROR: config.yaml not found in deploy/!"
  echo "  cp $DEPLOY_DIR/config.yaml.example $DEPLOY_DIR/config.yaml"
  echo "  Then fill in your secrets."
  exit 1
fi

# ── Ensure shared network exists ───────────────────────────────
if ! docker network inspect app_network &>/dev/null; then
  echo "[0/4] Creating shared Docker network: app_network"
  docker network create app_network
fi

# ── Tear down ─────────────────────────────────────────────────
if [ "$DOWN_MODE" = true ]; then
  echo "[DOWN] Stopping backend stack..."
  docker compose -f "$COMPOSE_FILE" down
  echo "✅  Backend stack stopped. Volumes preserved."
  exit 0
fi

echo "[1/4] Pulling latest code..."
git -C "$DEPLOY_DIR/.." pull --ff-only 2>/dev/null || echo "  Skipping git pull."

echo "[2/4] Building / updating Docker images..."
docker compose -f "$COMPOSE_FILE" pull --ignore-pull-failures 2>/dev/null || true
docker compose -f "$COMPOSE_FILE" up -d $BUILD_FLAG

echo "[3/4] Waiting for containers to stabilise (10s)..."
sleep 10

echo "[4/4] Container status:"
docker compose -f "$COMPOSE_FILE" ps

echo ""
echo "✅  Backend deploy complete!"
echo "  Logs:  docker compose -f $COMPOSE_FILE logs -f backend"
echo "  DB:    docker compose -f $COMPOSE_FILE logs -f postgres"
