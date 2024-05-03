ARG NGINX_VERSION=stable-alpine3.17-slim

FROM nginx:${NGINX_VERSION}

# Replace with our own nginx.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Import static files
COPY api_errors/ /api_errors