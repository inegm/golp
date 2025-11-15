// Package launchpad provides a Go library for controlling the Novation Launchpad Mini MIDI device.
//
// The Launchpad Mini is an 8x8 grid controller with bi-color LEDs (red/green) and additional
// scene and control buttons. This library provides a high-level API for LED control, button
// event handling, and advanced features like double-buffering.
//
// Basic Usage:
//
//	lp := launchpad.New()
//	err := lp.Open()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer lp.Close()
//
//	// Set an LED
//	lp.SetLED(3, 4, launchpad.ColorRed, launchpad.BrightnessFull)
//
//	// Listen for button presses
//	lp.OnButton(func(event launchpad.ButtonEvent) {
//	    if event.Pressed {
//	        fmt.Printf("Button pressed: %v\n", event.Button)
//	    }
//	})
//
// Features:
//   - Simple LED control with intuitive color and brightness API
//   - Event-driven button input with callbacks or channels
//   - Double-buffering support for smooth animations
//   - Automatic MIDI message rate limiting (400 msg/sec)
//   - Type-safe button coordinates and colors
//   - Support for all 80 buttons (64 grid + 8 scene + 8 top)
//
// Hardware Specifications:
//   - 8x8 grid of bi-color LEDs (red/green/amber/yellow)
//   - 4 brightness levels per color
//   - 8 scene buttons (right column)
//   - 8 top row buttons (Automap/Live)
//   - MIDI communication over USB
package launchpad
