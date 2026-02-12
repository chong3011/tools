// Package dbus demonstrates the differences between session bus and system bus
// in D-Bus communication for Go applications.
//
// This example shows how to connect to both bus types and interact with services.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/godbus/dbus/v5"
)

// SessionBusExample demonstrates connecting to and using the session bus
func SessionBusExample() {
	fmt.Println("=== Session Bus Example ===")
	
	// Connect to session bus
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Printf("Failed to connect to session bus: %v", err)
		return
	}
	defer conn.Close()

	fmt.Println("✓ Connected to session bus")

	// Example 1: List all names on session bus
	var names []string
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
	if err != nil {
		log.Printf("Failed to list names: %v", err)
		return
	}

	fmt.Printf("Found %d services on session bus\n", len(names))
	fmt.Println("Sample services:")
	count := 0
	for _, name := range names {
		if !isUniqueName(name) {
			fmt.Printf("  - %s\n", name)
			count++
			if count >= 5 {
				break
			}
		}
	}

	// Example 2: Send a notification (common session bus use case)
	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0,
		"ExampleApp",           // app_name
		uint32(0),              // replaces_id
		"",                     // app_icon
		"Session Bus Example",  // summary
		"This notification was sent via the session bus", // body
		[]string{},             // actions
		map[string]dbus.Variant{}, // hints
		int32(5000),            // expire_timeout (5 seconds)
	)
	if call.Err != nil {
		log.Printf("Note: Notification failed (service may not be available): %v", call.Err)
	} else {
		var notificationId uint32
		call.Store(&notificationId)
		fmt.Printf("✓ Notification sent (ID: %d)\n", notificationId)
	}

	fmt.Println()
}

// SystemBusExample demonstrates connecting to and using the system bus
func SystemBusExample() {
	fmt.Println("=== System Bus Example ===")
	
	// Connect to system bus
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Printf("Failed to connect to system bus: %v", err)
		return
	}
	defer conn.Close()

	fmt.Println("✓ Connected to system bus")

	// Example 1: List all names on system bus
	var names []string
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
	if err != nil {
		log.Printf("Failed to list names: %v", err)
		return
	}

	fmt.Printf("Found %d services on system bus\n", len(names))
	fmt.Println("Sample services:")
	count := 0
	for _, name := range names {
		if !isUniqueName(name) {
			fmt.Printf("  - %s\n", name)
			count++
			if count >= 5 {
				break
			}
		}
	}

	// Example 2: Query systemd (common system bus use case)
	obj := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
	var version string
	err = obj.Call("org.freedesktop.DBus.Properties.Get", 0,
		"org.freedesktop.systemd1.Manager",
		"Version").Store(&version)
	if err != nil {
		log.Printf("Note: Failed to get systemd version (may require permissions): %v", err)
	} else {
		fmt.Printf("✓ Systemd version: %s\n", version)
	}

	fmt.Println()
}

// ComparisonExample shows key differences in behavior
func ComparisonExample() {
	fmt.Println("=== Key Differences ===")
	
	fmt.Println("\n1. Scope:")
	fmt.Println("   Session Bus: Per-user session (isolated)")
	fmt.Println("   System Bus:  System-wide (shared across users)")
	
	fmt.Println("\n2. Lifetime:")
	fmt.Println("   Session Bus: Login → Logout")
	fmt.Println("   System Bus:  Boot → Shutdown")
	
	fmt.Println("\n3. Typical Services:")
	fmt.Println("   Session Bus: Notifications, media players, desktop apps")
	fmt.Println("   System Bus:  systemd, NetworkManager, UDisks2, BlueZ")
	
	fmt.Println("\n4. Security:")
	fmt.Println("   Session Bus: Permissive within user session")
	fmt.Println("   System Bus:  Strict policies, often requires privileges")
	
	fmt.Println("\n5. Use Cases:")
	fmt.Println("   Session Bus: User applications, desktop integration")
	fmt.Println("   System Bus:  System services, hardware management")
	
	// Check environment variables
	fmt.Println("\n6. Environment:")
	sessionAddr := os.Getenv("DBUS_SESSION_BUS_ADDRESS")
	if sessionAddr != "" {
		fmt.Printf("   Session bus address: %s\n", sessionAddr)
	} else {
		fmt.Println("   Session bus address: (not set, using default)")
	}
	fmt.Println("   System bus socket: /var/run/dbus/system_bus_socket")
	
	fmt.Println()
}

// isUniqueName checks if a D-Bus name is a unique connection name (starts with :)
func isUniqueName(name string) bool {
	return len(name) > 0 && name[0] == ':'
}

func main() {
	fmt.Println("GDBus Communication: Session Bus vs System Bus")
	fmt.Println("==============================================")
	fmt.Println()

	// Run session bus example
	SessionBusExample()

	// Run system bus example
	SystemBusExample()

	// Show comparison
	ComparisonExample()

	fmt.Println("Summary:")
	fmt.Println("--------")
	fmt.Println("Choose Session Bus for: User apps, desktop integration")
	fmt.Println("Choose System Bus for:  System services, hardware access")
}
