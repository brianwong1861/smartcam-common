#!/bin/bash

# Proto generation script for shared-proto
# This script generates Go code from proto files

set -e

echo "Generating proto files..."

# Create output directories
mkdir -p iam
mkdir -p camera
mkdir -p common
mkdir -p device
mkdir -p dealer
mkdir -p realtime
mkdir -p videostream

# Generate Go code for each proto file
echo "Generating IAM proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    iam/iam.proto

echo "Generating Camera proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    camera/camera.proto

echo "Generating Common proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    common/common.proto

echo "Generating Device proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    device/device.proto

echo "Generating Dealer proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    dealer/dealer.proto

echo "Generating Financial proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    financial/financial.proto

echo "Generating Quota proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    quota/quota.proto

echo "Generating Realtime proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    realtime/realtime.proto

echo "Generating Videostream proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    videostream/videostream.proto

echo "Proto generation completed!"

echo "Generating Notification proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    notification/notification.proto

echo "Proto generation completed!"