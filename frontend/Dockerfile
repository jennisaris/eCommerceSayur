# ==============================================================================
# Multi-Stage Dockerfile for Nuxt 3 Production Application
# ==============================================================================

FROM node:20-alpine AS builder

WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY package*.json ./

RUN npm ci

COPY . .

ENV NODE_ENV=production
RUN npm run build

# Stage 2: Runtime Image
FROM node:20-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0

COPY --from=builder --chown=node:node /app/.output ./

USER node

EXPOSE 3000

CMD ["node", "server/index.mjs"]
