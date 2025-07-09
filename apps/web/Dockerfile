# ---- Build Stage ----
FROM node:20-slim AS builder

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app

ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .
RUN pnpm build

# ---- Serve Stage ----
FROM nginx:alpine AS serve

COPY --from=builder /app/dist /usr/share/nginx/html

COPY ./config/nginx/nginx.conf /etc/nginx/conf.d

EXPOSE 80