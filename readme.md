# Go Sparkplug

A Go implementation of Eclipse Sparkplug B specification for MQTT-based industrial IoT communication.

## Overview

This library provides a complete implementation of the [Eclipse Sparkplug B specification](https://www.eclipse.org/tahu/spec/Sparkplug%20Topic%20Namespace%20and%20State%20ManagementV2.2-with%20appendix%20B%20format%20-%20Eclipse.pdf), enabling seamless integration with industrial MQTT systems. Sparkplug B is a specification that defines MQTT topic namespace, payload format, and session state management for industrial IoT applications.

## Features

- ✅ **Complete Sparkplug B Implementation**: Full support for all message types (NBIRTH, NDEATH, DBIRTH, DDEATH, NDATA, DDATA)
- ✅ **Node Management**: Automatic birth/death certificate handling with sequence numbering
- ✅ **Device Management**: Support for multiple devices per node
- ✅ **Metric Support**: Built-in support for all Sparkplug data types (Int, Float, String, Boolean, Bytes, etc.)
- ✅ **Command Handling**: Built-in support for Node Control commands (Rebirth, Reboot)
- ✅ **Thread-Safe**: Concurrent-safe sequence number management
- ✅ **Auto-Reconnection**: Automatic MQTT reconnection with proper state management
- ✅ **Protocol Buffers**: Uses official Sparkplug B protobuf definitions

## Installation

```bash
go get github.com/tjeumaster/go-sparkplug
```

## Quick Start

### Basic Node Setup

```go
package main

import (
    "log"
    "time"
    "github.com/tjeumaster/go-sparkplug/sparkplug"
)

func main() {
    // Configure the Sparkplug client
    config := &sparkplug.Config{
        Host:     "localhost",
        Port:     1883,
        Username: "your-username",
        Password: "your-password",
        ClientID: "go-sparkplug-client",
        GroupID:  "MyGroup",
        NodeID:   "MyNode",
    }

    // Create and start the client
    client := sparkplug.NewClient(config)
    client.Start()
    defer client.Stop()

    // Publish node data periodically
    for {
        metrics := map[string]any{
            "temperature": 25.5,
            "humidity":    60,
            "status":      "running",
        }
        
        if err := client.PublishNDATA(metrics); err != nil {
            log.Printf("Error publishing NDATA: %v", err)
        }
        
        time.Sleep(10 * time.Second)
    }
}
```

### Device Management

```go
// Implement the Device interface
type MyDevice struct {
    ID string
}

func (d *MyDevice) GetId() string {
    return d.ID
}

func (d *MyDevice) GetMetricValues() map[string]any {
    return map[string]any{
        "voltage":     3.3,
        "current":     1.2,
        "power_state": true,
    }
}

// Usage
device := &MyDevice{ID: "sensor-001"}

// Publish device birth certificate
client.PublishDBIRTH(device)

// Publish device data
deviceMetrics := map[string]any{
    "voltage": 3.4,
    "current": 1.1,
}
client.PublishDDATA(device, deviceMetrics)

// Publish device death certificate
client.PublishDDEATH(device)
```

## API Reference

### Client

#### Configuration

```go
type Config struct {
    Host     string  // MQTT broker hostname
    Port     int     // MQTT broker port
    Username string  // MQTT username
    Password string  // MQTT password
    ClientID string  // MQTT client ID
    GroupID  string  // Sparkplug Group ID
    NodeID   string  // Sparkplug Node ID
}
```

#### Methods

- `NewClient(config *Config) *Client` - Create a new Sparkplug client
- `Start()` - Connect to MQTT broker and publish NBIRTH
- `Stop()` - Publish NDEATH and disconnect from broker
- `PublishNBIRTH() error` - Publish Node Birth certificate
- `PublishNDEATH() error` - Publish Node Death certificate
- `PublishNDATA(metrics map[string]any) error` - Publish Node Data
- `PublishDBIRTH(device Device) error` - Publish Device Birth certificate
- `PublishDDEATH(device Device) error` - Publish Device Death certificate
- `PublishDDATA(device Device, metrics map[string]any) error` - Publish Device Data

### Device Interface

Implement this interface for your devices:

```go
type Device interface {
    GetId() string
    GetMetricValues() map[string]any
}
```

### Supported Data Types

The library supports all Sparkplug B data types:

- `int`, `int32`, `int64` → Int32/Int64
- `uint32`, `uint64` → UInt32/UInt64
- `float32` → Float
- `float64` → Double
- `string` → String
- `bool` → Boolean
- `[]byte` → Bytes

## Message Types

### Node Messages

- **NBIRTH**: Published when node comes online, includes birth sequence number and node control metrics
- **NDEATH**: Published when node goes offline (via last will testament)
- **NDATA**: Published to send node-level metric data

### Device Messages

- **DBIRTH**: Published when a device comes online
- **DDEATH**: Published when a device goes offline
- **DDATA**: Published to send device-level metric data

### Command Messages

The client automatically handles these command types:

- **Node Control/Rebirth**: Triggers republishing of NBIRTH message
- **Node Control/Reboot**: Ready for custom reboot logic implementation

## Topic Structure

The library follows the standard Sparkplug B topic structure:

```
spBv1.0/{GroupID}/{MessageType}/{NodeID}[/{DeviceID}]
```

Examples:
- `spBv1.0/MyGroup/NBIRTH/MyNode`
- `spBv1.0/MyGroup/DDATA/MyNode/Sensor001`

## Requirements

- Go 1.21 or later
- MQTT broker (e.g., Eclipse Mosquitto, HiveMQ, AWS IoT Core)

## Dependencies

- [Eclipse Paho MQTT Go Client](https://github.com/eclipse/paho.mqtt.golang)
- [Protocol Buffers](https://google.golang.org/protobuf)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## References

- [Eclipse Sparkplug B Specification](https://www.eclipse.org/tahu/spec/Sparkplug%20Topic%20Namespace%20and%20State%20ManagementV2.2-with%20appendix%20B%20format%20-%20Eclipse.pdf)
- [Eclipse Tahu Project](https://github.com/eclipse/tahu)
- [MQTT Specification](https://mqtt.org/mqtt-specification/)

## Support

For questions, issues, or contributions, please visit the [GitHub repository](https://github.com/tjeumaster/go-sparkplug).