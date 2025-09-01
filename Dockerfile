# Stage 1: The build environment, named "builder"
# We use a specific Go version on Alpine Linux for a smaller build image.
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod file and download dependencies first.
# This leverages Docker's layer caching. Dependencies are only re-downloaded
# if go.mod changes.
COPY go.mod ./
RUN go mod download

COPY . ./

# Build the Go application as a static binary.
# CGO_ENABLED=0 is critical for building a binary that can run in a minimal 'scratch' image.
# -ldflags="-w -s" strips debug symbols, reducing the binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/server -ldflags="-w -s" .


# Stage 2: The final, minimal production image
# 'scratch' is an empty image, providing the smallest possible base.
FROM scratch

# Set the working directory
WORKDIR /

# Copy only the compiled application binary from the 'builder' stage.
COPY --from=builder /app/server /server

# Copy the quotes.txt file that the application needs to run.
COPY quotes.txt .

# Expose port 8080 to the outside world.
EXPOSE 8080

# The command to run when the container starts.
ENTRYPOINT ["/server"]