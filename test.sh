#!/bin/bash

# Timeout in seconds
timeout=300  # Adjust the timeout value as needed

# Start the services
docker compose build
docker compose up -d

# Function to check if a service is healthy
is_service_healthy() {
    service_name="$1"
    docker-compose ps -q "$service_name" | xargs docker inspect --format='{{.State.Health.Status}}' | grep -q 'healthy'
}

# Function to check if the timeout has expired
is_timeout_expired() {
    (( $(date +%s) - start_time >= timeout ))
}

start_time=$(date +%s)

# Wait for all services to be healthy within the specified timeout
for service in $(docker-compose config --services); do
    until is_service_healthy "$service" || is_timeout_expired; do
        echo "Waiting for $service to be healthy..."
        sleep 1
    done

    # Break out of the loop if timeout has expired
    if is_timeout_expired; then
        echo "Timed out waiting for services to be healthy."
        exit 1
    fi
done

# Ping the healthz endpoint
if curl --fail http://localhost:8080/mock-api/v1/mock/healthz; then
    echo "" 
    echo "---------------------------------"
    echo "|  API health check successful  |"
    echo "---------------------------------"
    docker-compose -f docker-compose.amd64.yml down
else
    echo ""
    echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
    echo "!  API health check failed  !"
    echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
    docker-compose -f docker-compose.amd64.yml down
    exit 1
fi

echo "All services are healthy and API health check successful."
