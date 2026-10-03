# Stage 1: Install dependencies and build
FROM node:20-alpine AS build
WORKDIR /app

RUN npm install -g pnpm

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY frontend/package.json ./frontend/package.json
RUN pnpm install --frozen-lockfile --config.ignore-scripts=false

COPY frontend ./frontend

# Build-time variables — Vite bakes these into the static bundle.
# Pass them via docker-compose build.args or --build-arg.
ARG VITE_BASE_URL
ARG VITE_OIDC_CLIENT_ID
ENV VITE_BASE_URL=$VITE_BASE_URL
ENV VITE_OIDC_CLIENT_ID=$VITE_OIDC_CLIENT_ID

RUN pnpm --filter @authforest/frontend build

# Stage 2: Serve with nginx
FROM nginx:alpine

COPY --from=build /app/frontend/build /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
