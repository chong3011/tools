# Documentation

This directory contains additional documentation and examples for various topics.

## GDBus Communication

Documentation explaining the differences between Session Bus and System Bus in GDBus communication:

- **[gdbus-bus-types.md](gdbus-bus-types.md)** - Comprehensive English documentation on Session Bus vs System Bus
- **[gdbus-bus-types-zh.md](gdbus-bus-types-zh.md)** - Chinese version (会话总线 vs 系统总线)
- **[gdbus-example.go](gdbus-example.go)** - Practical Go examples demonstrating both bus types

### Quick Summary

**Session Bus (会话总线)**:
- Per-user session scope
- For desktop applications and user services
- Permissive within user session
- Examples: notifications, media players

**System Bus (系统总线)**:
- System-wide scope
- For system services and hardware management
- Strict security policies
- Examples: NetworkManager, systemd

### Running the Example

To run the Go example (requires `godbus` package):

```bash
# Install dependency
go get github.com/godbus/dbus/v5

# Run the example
go run docs/gdbus-example.go
```

Note: The example requires a running D-Bus session and system bus to fully demonstrate all features.
