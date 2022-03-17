FROM 1.18.0-alpine3.15 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies.
RUN go mod download

# Copy the source to build container
COPY . .

ARG VERSION=Develop

# Build the Go app
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/alecbcs/caffeinate-workflows/main.Version=${VERSION}" -o caffeinate-workflows .

FROM alpine:3.15.0

COPY --from=builder /app/caffeinate-workflows /caffeinate-workflows

ENTRYPOINT ["/caffeinate-workflows"]