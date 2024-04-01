#!/bin/bash
docker compose -f docker-compose.arm.yml build
docker compose -f docker-compose.arm.yml up -d