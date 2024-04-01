#!/bin/bash
docker compose -f docker-compose.amd64.yml build
docker compose -f docker-compose.amd64.yml up -d