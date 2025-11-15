package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/inegm/golp/pkg/launchpad"
)

func main() {
	// Create a new Launchpad instance
	lp := launchpad.New()

	// Open connection to the device
	fmt.Println("Opening Launchpad...")
	err := lp.Open()
	if err != nil {
		log.Fatalf("Failed to open Launchpad: %v", err)
	}
	defer lp.Close()

	fmt.Println("Launchpad connected!")
	fmt.Println("Press buttons on the Launchpad...")
	fmt.Println("Press Ctrl+C to exit")

	// Clear all LEDs
	lp.Clear()

	// Register button event handler using callback
	lp.OnButton(func(event launchpad.ButtonEvent) {
		if event.Pressed {
			// Light up the button when pressed
			fmt.Printf("Button pressed: %v\n", event.Button)
			lp.SetButtonLED(event.Button, launchpad.ColorGreen, launchpad.BrightnessFull)
		} else {
			// Turn off the button when released
			fmt.Printf("Button released: %v\n", event.Button)
			lp.SetButtonLED(event.Button, launchpad.ColorOff, launchpad.BrightnessOff)
		}
	})

	// Alternative: Use channel-based event handling
	// Uncomment this block to use channels instead of callbacks
	/*
	go func() {
		for event := range lp.ButtonEvents() {
			if event.Pressed {
				fmt.Printf("Channel: Button pressed: %v\n", event.Button)
			}
		}
	}()
	*/

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-sigChan
	fmt.Println("\nReceived interrupt signal, cleaning up...")
}
