# Sparkplug B Go Client

A Go implementation of the Eclipse Sparkplug B MQTT protocol specification for IIoT (Industrial Internet of Things) applications.

## Overview

This library provides a client implementation for the Sparkplug B protocol, enabling edge nodes and devices to communicate with MQTT-based SCADA/ICS systems using the standardized Sparkplug B specification.

## Features

- ✅ Full Sparkplug B protocol support
- ✅ Node birth/death certificates (NBIRTH/NDEATH)
- ✅ Device birth/death certificates (DBIRTH/DDEATH)
- ✅ Node and device data messages (NDATA/DDATA)
- ✅ Command handling (NCMD/DCMD)
- ✅ Auto-reconnection with proper state management
- ✅ Sequence number tracking
- ✅ Last Will and Testament (LWT) support
- ✅ Built-in support for multiple data types (int, uint, float, string, bool, bytes)
- ✅ Thread-safe operations

## Installation

```bash
go get github.com/tjeumaster/go-sparkplug
```

## Dependencies

- [Eclipse Paho MQTT Go Client](https://github.com/eclipse/paho.mqtt.golang)
- [Protocol Buffers](https://pkg.go.dev/google.golang.org/protobuf)

## Usage

### Basic Setup

```go
package main

import (
    "log"
    "github.com/tjeumaster/sparkplug-b/spb"
)

func main() {
    // Create client configuration
    config := spb.Config{
        Host:     "localhost",
        Port:     1883,
        Username: "your-username",
        Password: "your-password",
        ClientID: "edge-node-1",
        GroupID:  "group1",
        NodeID:   "node1",
    }

    // Create new Sparkplug B client
    client := spb.NewClient(config)

    // Connect to MQTT broker
    if err := client.Connect(); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer client.Disconnect()

    log.Println("Connected successfully!")
}
```

### Publishing Node Data

```go
// Publish node data
metrics := map[string]any{
    "temperature": 25.5,
    "pressure":    101.3,
    "status":      "online",
}

if err := client.PublishNDATA(metrics); err != nil {
    log.Printf("Failed to publish NDATA: %v", err)
}
```

### Working with Devices

```go
// Implement the Device interface
type MyDevice struct {
    id string
}

func (d *MyDevice) GetId() string {
    return d.id
}

func (d *MyDevice) GetMetricValues() map[string]any {
    return map[string]any{
        "sensor1": 42.0,
        "sensor2": true,
        "name":    "Device A",
    }
}

// Create and register a device
device := &MyDevice{id: "device-001"}

// Publish device birth certificate
if err := client.PublishDBIRTH(device); err != nil {
    log.Printf("Failed to publish DBIRTH: %v", err)
}

// Publish device data
deviceMetrics := map[string]any{
    "sensor1": 43.5,
}
if err := client.PublishDDATA(device, deviceMetrics); err != nil {
    log.Printf("Failed to publish DDATA: %v", err)
}

// Publish device death certificate
if err := client.PublishDDEATH(device); err != nil {
    log.Printf("Failed to publish DDEATH: %v", err)
}
```

### Supported Data Types

The library automatically maps Go types to Sparkplug B data types:

| Go Type    | Sparkplug B Type |
|-----------|------------------|
| `int`     | Int64            |
| `int32`   | Int32            |
| `int64`   | Int64            |
| `uint32`  | UInt32           |
| `uint64`  | UInt64           |
| `float32` | Float            |
| `float64` | Double           |
| `string`  | String           |
| `bool`    | Boolean          |
| `[]byte`  | Bytes            |

## Architecture

### Project Structure

```
sparkplug-b/
├── spb/
│   ├── client.go      # Main client implementation
│   ├── payload.go     # Payload builders (NBIRTH, NDEATH, DBIRTH, etc.)
│   └── metric.go      # Metric conversion utilities
├── sproto/
│   ├── sparkplug_b.proto    # Protocol Buffer definition
│   └── sparkplug_b.pb.go    # Generated protobuf code
├── go.mod
└── README.md
```

### Key Components

- **Client**: Main entry point for Sparkplug B operations
- **Config**: Configuration for MQTT connection and Sparkplug B identity
- **Device Interface**: Contract for device implementations
- **Payload Builders**: Internal methods for constructing Sparkplug B messages
- **Metric Converters**: Utilities for converting Go types to Sparkplug B metrics

## Command Handling

The client automatically handles standard Sparkplug B commands:

- **Node Control/Rebirth**: Triggers republishing of NBIRTH
- **Node Control/Reboot**: Placeholder for reboot logic (implement as needed)

Custom command handling can be extended by modifying the `handleCommandMetric` method.

## Best Practices

1. **Always disconnect gracefully**: Use `defer client.Disconnect()` to ensure proper NDEATH publication
2. **Handle errors**: Check return values from all publish methods
3. **Implement Device interface properly**: Ensure `GetMetricValues()` returns current state
4. **Use appropriate data types**: Match your metrics to the correct Go types for proper Sparkplug B encoding
5. **Monitor sequence numbers**: The client handles this automatically, but be aware of the 0-255 range

## Thread Safety

The client uses mutex locks to ensure thread-safe sequence number management. All public methods can be safely called from multiple goroutines.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.

## License

MIT License

## References

- [Eclipse Sparkplug Specification](https://sparkplug.eclipse.org/)
- [Sparkplug B Protocol Documentation](https://www.eclipse.org/tahu/spec/Sparkplug%20Topic%20Namespace%20and%20State%20ManagementV2.2-with%20appendix%20B%20format%20-%20Eclipse.pdf)
- [MQTT Protocol](https://mqtt.org/)

## Support

For questions and support, please open an issue on GitHub.
