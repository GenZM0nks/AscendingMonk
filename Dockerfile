# Stage 1: Build environment
FROM golang:1.26-bookworm AS builder

WORKDIR /src

# go-sqlite3 requires a C compiler.
RUN apt-get update \
    && apt-get install -y --no-install-recommends gcc libc6-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build -o /out/ascendingmonk ./cmd/main


# Stage 2: Runtime environment
FROM debian:bookworm-slim

WORKDIR /app

# Runtime libraries, TLS certificates, and an unprivileged user.
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 10001 ascendingmonk \
    && useradd --system --uid 10001 \
       --gid ascendingmonk \
       --no-create-home ascendingmonk \
    && mkdir -p /data \
    && chown ascendingmonk:ascendingmonk /data

COPY --from=builder /out/ascendingmonk /app/ascendingmonk
COPY --from=builder /src/web/templates /app/web/templates
COPY --from=builder /src/web/static /app/web/static

ENV DB_PATH=/data/ascendingmonk.db

USER ascendingmonk

EXPOSE 8080

CMD ["/app/ascendingmonk"]