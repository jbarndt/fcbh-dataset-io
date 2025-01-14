FROM public.ecr.aws/docker/library/golang:1.23.4-alpine as builder

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application
RUN go build -o taskService controller/ecs_entry.go

# Deployment stage
FROM public.ecr.aws/amazonlinux/amazonlinux:latest

WORKDIR /root/


# Install necessary tools
RUN yum update -y && \
    # yum install -y wget tar xz && \
    yum clean all


# TODO: include build instructions from docs/user_data.sh

# Copy the built binary
COPY --from=builder /app/taskService .

# Ensure the binary is executable
RUN chmod +x ./taskService

# Use environment variable to determine the service to run
ENTRYPOINT ["./taskService"]
