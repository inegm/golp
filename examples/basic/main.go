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

	// Clear all LEDs
	fmt.Println("Clearing all LEDs...")
	err = lp.Clear()
	if err != nil {
		log.Fatalf("Failed to clear: %v", err)
	}

	// Light up some individual LEDs
	fmt.Println("Setting individual LEDs...")

	// Red corner (top-left)
	lp.SetLED(0, 0, launchpad.ColorRed, launchpad.BrightnessFull)

	// Green corner (top-right)
	lp.SetLED(7, 0, launchpad.ColorGreen, launchpad.BrightnessFull)

	// Amber corner (bottom-left)
	lp.SetLED(0, 7, launchpad.ColorAmber, launchpad.BrightnessFull)

	// Yellow corner (bottom-right)
	lp.SetLED(7, 7, launchpad.ColorYellow, launchpad.BrightnessFull)

	// Draw a cross in the middle
	fmt.Println("Drawing cross pattern...")
	for i := 0; i < 8; i++ {
		lp.SetLED(i, 3, launchpad.ColorGreen, launchpad.BrightnessLow)
		lp.SetLED(i, 4, launchpad.ColorGreen, launchpad.BrightnessLow)
		lp.SetLED(3, i, launchpad.ColorRed, launchpad.BrightnessLow)
		lp.SetLED(4, i, launchpad.ColorRed, launchpad.BrightnessLow)
	}

	// Light up scene buttons
	fmt.Println("Lighting scene buttons...")
	for i := 0; i < 8; i++ {
		lp.SetSceneButton(i, launchpad.ColorAmber, launchpad.BrightnessMedium)
	}

	// Light up top buttons
	fmt.Println("Lighting top buttons...")
	for i := 0; i < 8; i++ {
		lp.SetTopButton(i, launchpad.ColorRed, launchpad.BrightnessLow)
	}

	fmt.Println("Press Ctrl+C to exit")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-sigChan
	fmt.Println("\nReceived interrupt signal, cleaning up...")
}
