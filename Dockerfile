# Stage 1: Build stage
FROM golang:1.21.5 AS builder

WORKDIR /backend/

COPY ./cmd /backend/cmd
COPY ./internal /backend/internal
COPY ./dev /backend/dev

COPY ./go.mod /backend/go.mod
COPY ./go.sum /backend/go.sum

RUN go build -o /backend/build ./cmd/L0/

# Stage 2: Final stage
FROM ubuntu:22.04

WORKDIR /backend/

COPY --from=builder /backend/build /backend/build
COPY --from=builder /backend/dev/.env /backend/dev/.env
COPY --from=builder /backend/internal/app/migrations /backend/internal/app/migrations
COPY --from=builder /backend/internal/static /backend/internal/static
COPY --from=builder /backend/internal/templates /backend/internal/templates

CMD [ "/backend/build" ]
