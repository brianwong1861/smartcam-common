# SmartCam Shared Proto

This directory contains shared Protocol Buffer definitions for the SmartCam project.

## Structure

```
shared-proto/
├── iam/                    # Identity and Access Management proto files
│   ├── iam.proto          # IAM service definitions
│   ├── iam.pb.go          # Generated Go code
│   └── iam_grpc.pb.go     # Generated gRPC Go code
├── camera/                 # Camera service proto files
│   ├── camera.proto       # Camera service definitions
│   ├── camera.pb.go       # Generated Go code
│   └── camera_grpc.pb.go  # Generated gRPC Go code
├── device/                 # Device Management service proto files
│   ├── device.proto       # Device service definitions
│   ├── device.pb.go       # Generated Go code
│   └── device_grpc.pb.go  # Generated gRPC Go code
├── dealer/                 # Dealer Management service proto files
│   ├── dealer.proto       # Dealer service definitions
│   ├── dealer.pb.go       # Generated Go code
│   └── dealer_grpc.pb.go  # Generated gRPC Go code
├── common/                 # Common proto definitions
│   ├── common.proto       # Common message types
│   └── common.pb.go       # Generated Go code
├── go.mod                  # Go module definition
├── generate.sh            # Proto generation script
└── README.md              # This file
```

## Usage

### Generating Proto Files

Run the generation script to create Go code from proto files:

```bash
./generate.sh
```

### Dependencies

This module requires:
- Protocol Buffer compiler (`protoc`)
- Go protobuf plugin (`protoc-gen-go`)
- Go gRPC plugin (`protoc-gen-go-grpc`)

### Integration

All SmartCam microservices use this shared proto module:

```go
import "smartcam-proto/iam"
import "smartcam-proto/camera"
import "smartcam-proto/device"
import "smartcam-proto/dealer"
import "smartcam-proto/common"
```

## Benefits

1. **Unified Standards**: Ensures both services use the same data structures
2. **Code Reuse**: Avoids duplicate proto definitions
3. **Version Control**: Centralized management of API contracts
4. **Consistency**: Guarantees compatibility between services

## Development

When modifying proto files:

1. Update the `.proto` file
2. Run `./generate.sh` to regenerate Go code
3. Update both services to use the new definitions
4. Test the integration

## Services Using This Module

- **iam-service**: Uses `smartcam-proto/iam` for gRPC service definitions
- **camera-api-server**: Uses `smartcam-proto/iam` for gRPC client communication
- **device-management-service**: Uses `smartcam-proto/device` for gRPC service definitions
- **dealer-management-service**: Uses `smartcam-proto/dealer` for gRPC service definitions
- **All services**: Use `smartcam-proto/common` for shared message types
