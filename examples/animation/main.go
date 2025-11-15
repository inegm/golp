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
	fmt.Println("Displaying double-buffered animation...")
	fmt.Println("Press Ctrl+C to exit")

	// Clear both buffers
	lp.Clear()

	// Set up double-buffering
	// Display buffer 0, update buffer 1
	lp.SetDisplayBuffer(launchpad.Buffer0)
	lp.SetUpdateBuffer(launchpad.Buffer1)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Animation: bouncing ball
	x, y := 0, 0
	dx, dy := 1, 1

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			fmt.Println("\nReceived interrupt signal, cleaning up...")
			return
		case <-ticker.C:
		// Clear the update buffer (turn off all LEDs)
		for row := 0; row < 8; row++ {
			for col := 0; col < 8; col++ {
				lp.SetLED(col, row, launchpad.ColorOff, launchpad.BrightnessOff)
			}
		}

		// Draw the ball at new position
		lp.SetLED(x, y, launchpad.ColorRed, launchpad.BrightnessFull)

		// Draw a trail
		if x > 0 {
			lp.SetLED(x-1, y, launchpad.ColorAmber, launchpad.BrightnessLow)
		}
		if x < 7 {
			lp.SetLED(x+1, y, launchpad.ColorAmber, launchpad.BrightnessLow)
		}
		if y > 0 {
			lp.SetLED(x, y-1, launchpad.ColorAmber, launchpad.BrightnessLow)
		}
		if y < 7 {
			lp.SetLED(x, y+1, launchpad.ColorAmber, launchpad.BrightnessLow)
		}

		// Swap buffers for instant update
		lp.SwapBuffers()

		// Update position
		x += dx
		y += dy

		// Bounce off walls
		if x <= 0 || x >= 7 {
			dx = -dx
		}
		if y <= 0 || y >= 7 {
			dy = -dy
		}

		// Keep in bounds
		if x < 0 {
			x = 0
		}
		if x > 7 {
			x = 7
		}
		if y < 0 {
			y = 0
		}
		if y > 7 {
			y = 7
		}
		}
	}
}
