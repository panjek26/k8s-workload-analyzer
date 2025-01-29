#!/bin/sh
# entrypoint for ./Dockerfile and ./Dockerfile.rabbitmq
set -e
SERVICE_LOG_DIR="${SERVICE_LOG_DIR:-/tmp/log}"
SERVICE_NAME="${SERVICE_NAME:-service}"
SERVICE_EXECUTABLE="${SERVICE_EXECUTABLE:-/app/server}"
# ----------------------------------------
$SERVICE_EXECUTABLE
