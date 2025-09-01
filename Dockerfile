# Stage 1: The build environment, named "builder"
# We use a specific Go version on Alpine Linux for a smaller build image.
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/server -ldflags="-w -s" .


# Stage 2: The final, minimal production image
# 'scratch' is an empty image, providing the smallest possible base.
FROM alpine:latest

WORKDIR /

COPY --from=builder /app/server /server

COPY quotes.txt .

EXPOSE 8080

# The command to run when the container starts.
ENTRYPOINT ["/server"]