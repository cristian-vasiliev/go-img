# Image for building the Lambda function
FROM --platform=linux/amd64 golang:1.19 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the host to the container
COPY . .

# Download dependencies
RUN go mod download

# Install required development libraries for AVIF and WEBP support
RUN apt-get update && \
    apt-get install -y libaom-dev libwebp-dev

# Build the application for x86_64
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main .

# Image for running the Lambda function
FROM public.ecr.aws/lambda/provided:al2

# Install the required runtime libraries for AVIF and WEBP support
RUN yum -y install libaom libwebp

# Copy custom runtime bootstrap
COPY bootstrap ${LAMBDA_RUNTIME_DIR}

# Copy function code
COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}

# Copy the static files from the builder image
COPY --from=builder /app/static ${LAMBDA_TASK_ROOT}/static

# Set an empty CMD
CMD ["main"]
