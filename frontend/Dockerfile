FROM oven/bun:alpine AS builder
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install
COPY . .
RUN bun run build

FROM steebchen/nginx-spa:stable AS runtime
COPY --from=builder /app/dist /app
EXPOSE 80
CMD ["nginx"]
