package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/inegm/golp/pkg/launchpad"
)

const (
	gridSize = 8
)

// Grid represents the game state
type Grid [gridSize][gridSize]bool

// NewRandomGrid creates a grid with random live cells
func NewRandomGrid(density float64) Grid {
	var grid Grid
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			grid[y][x] = rand.Float64() < density
		}
	}
	return grid
}

// CountNeighbors counts the number of live neighbors for a cell
func (g *Grid) CountNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue // Skip the cell itself
			}

			// Wrap around edges (toroidal topology)
			nx := (x + dx + gridSize) % gridSize
			ny := (y + dy + gridSize) % gridSize

			if g[ny][nx] {
				count++
			}
		}
	}
	return count
}

// Step advances the game by one generation
func (g *Grid) Step() Grid {
	var next Grid

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			neighbors := g.CountNeighbors(x, y)
			alive := g[y][x]

			// Conway's rules:
			// 1. Any live cell with 2-3 neighbors survives
			// 2. Any dead cell with exactly 3 neighbors becomes alive
			// 3. All other cells die or stay dead
			if alive {
				next[y][x] = neighbors == 2 || neighbors == 3
			} else {
				next[y][x] = neighbors == 3
			}
		}
	}

	return next
}

// IsEmpty checks if the grid has no live cells
func (g *Grid) IsEmpty() bool {
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			if g[y][x] {
				return false
			}
		}
	}
	return true
}

// CountLiveCells returns the number of live cells
func (g *Grid) CountLiveCells() int {
	count := 0
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			if g[y][x] {
				count++
			}
		}
	}
	return count
}

// RenderToLaunchpad displays the grid on the Launchpad
func RenderToLaunchpad(lp *launchpad.Launchpad, grid Grid, generation int) {
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			if grid[y][x] {
				// Live cells are green
				lp.SetLED(x, y, launchpad.ColorGreen, launchpad.BrightnessFull)
			} else {
				// Dead cells are off
				lp.SetLED(x, y, launchpad.ColorOff, launchpad.BrightnessOff)
			}
		}
	}

	// Show generation count on scene buttons (binary representation)
	for i := 0; i < 8; i++ {
		if generation&(1<<i) != 0 {
			lp.SetSceneButton(i, launchpad.ColorAmber, launchpad.BrightnessLow)
		} else {
			lp.SetSceneButton(i, launchpad.ColorOff, launchpad.BrightnessOff)
		}
	}

	// Show population on top buttons (scaled)
	liveCells := grid.CountLiveCells()
	populationBar := (liveCells * 8) / 64 // Scale to 0-8
	for i := 0; i < 8; i++ {
		if i < populationBar {
			lp.SetTopButton(i, launchpad.ColorRed, launchpad.BrightnessMedium)
		} else {
			lp.SetTopButton(i, launchpad.ColorOff, launchpad.BrightnessOff)
		}
	}
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Conway's Game of Life on Launchpad Mini")
	fmt.Println("=======================================")

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

	// Set up double-buffering for smooth animation
	lp.SetDisplayBuffer(launchpad.Buffer0)
	lp.SetUpdateBuffer(launchpad.Buffer1)

	// Create initial random grid (30% density)
	grid := NewRandomGrid(0.3)
	generation := 0

	fmt.Println("Starting simulation...")
	fmt.Println("- Live cells: Green")
	fmt.Println("- Scene buttons: Generation count (binary)")
	fmt.Println("- Top buttons: Population bar")
	fmt.Println("Press Ctrl+C to exit")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Game loop
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	// Track if simulation has stagnated
	lastLiveCells := -1
	stagnantCount := 0

	for {
		select {
		case <-sigChan:
			fmt.Println("\nReceived interrupt signal, cleaning up...")
			return

		case <-ticker.C:
			// Render current state
			RenderToLaunchpad(lp, grid, generation)

			// Swap buffers for instant update
			lp.SwapBuffers()

			// Check for empty or stagnant grid
			liveCells := grid.CountLiveCells()
			if grid.IsEmpty() {
				fmt.Printf("\nGeneration %d: All cells died. Restarting...\n", generation)
				grid = NewRandomGrid(0.3)
				generation = 0
				stagnantCount = 0
				continue
			}

			// Check for stagnation (population not changing)
			if liveCells == lastLiveCells {
				stagnantCount++
				if stagnantCount > 20 {
					fmt.Printf("\nGeneration %d: Simulation stagnant. Restarting...\n", generation)
					grid = NewRandomGrid(0.3)
					generation = 0
					stagnantCount = 0
					continue
				}
			} else {
				stagnantCount = 0
			}
			lastLiveCells = liveCells

			// Advance to next generation
			grid = grid.Step()
			generation++

			// Print status every 10 generations
			if generation%10 == 0 {
				fmt.Printf("Generation %d: %d live cells\n", generation, liveCells)
			}
		}
	}
}
