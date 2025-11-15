package main

import (
	"os"
	"os/signal"
	"syscall"

	"fmt"
	"log"
	"time"

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
	fmt.Println("Displaying rainbow animation...")
	fmt.Println("Press Ctrl+C to exit")

	// Color palette for the rainbow
	colors := []launchpad.Color{
		launchpad.ColorRed,
		launchpad.ColorAmber,
		launchpad.ColorYellow,
		launchpad.ColorGreen,
		launchpad.ColorGreen,
		launchpad.ColorYellow,
		launchpad.ColorAmber,
		launchpad.ColorRed,
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Animate the rainbow scrolling
	offset := 0
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			fmt.Println("\nReceived interrupt signal, cleaning up...")
			return
		case <-ticker.C:
		// Update each row with a different color
		for y := 0; y < 8; y++ {
			colorIndex := (y + offset) % len(colors)
			color := colors[colorIndex]

			for x := 0; x < 8; x++ {
				lp.SetLED(x, y, color, launchpad.BrightnessFull)
			}
		}

		// Update scene buttons to match the rainbow
		for y := 0; y < 8; y++ {
			colorIndex := (y + offset) % len(colors)
			color := colors[colorIndex]
			lp.SetSceneButton(y, color, launchpad.BrightnessMedium)
		}

		// Update top buttons
		for x := 0; x < 8; x++ {
			colorIndex := (x + offset) % len(colors)
			color := colors[colorIndex]
			lp.SetTopButton(x, color, launchpad.BrightnessMedium)
		}

		offset++
		if offset >= len(colors) {
			offset = 0
		}
		}
	}
}
