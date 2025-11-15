# Launchpad Go Library - API Documentation

## Package Overview

package launchpad // import "github.com/inegm/golp/pkg/launchpad"

Package launchpad provides a Go library for controlling the Novation Launchpad
Mini MIDI device.

The Launchpad Mini is an 8x8 grid controller with bi-color LEDs (red/green) and
additional scene and control buttons. This library provides a high-level API for
LED control, button event handling, and advanced features like double-buffering.

# Basic Usage

Create a Launchpad instance, open the connection, and control LEDs:

    package main

    import (
    	"log"
    	"github.com/inegm/golp/pkg/launchpad"
    )

    func main() {
    	lp := launchpad.New()
    	err := lp.Open()
    	if err != nil {
    		log.Fatal(err)
    	}
    	defer lp.Close()

    	// Set an LED
    	lp.SetLED(3, 4, launchpad.ColorRed, launchpad.BrightnessFull)
    }

# LED Control

The library provides several methods for controlling LEDs:

    // Individual LED
    lp.SetLED(x, y, launchpad.ColorGreen, launchpad.BrightnessFull)

    // Entire row
    lp.SetRow(3, launchpad.ColorAmber, launchpad.BrightnessMedium)

    // Entire column
    lp.SetColumn(5, launchpad.ColorRed, launchpad.BrightnessLow)

    // All grid LEDs
    lp.SetAllLEDs(launchpad.ColorYellow, launchpad.BrightnessFull)

    // Scene buttons (right column)
    lp.SetSceneButton(2, launchpad.ColorGreen, launchpad.BrightnessFull)

    // Top row buttons
    lp.SetTopButton(7, launchpad.ColorAmber, launchpad.BrightnessMedium)

    // Clear all LEDs
    lp.Clear()

# Button Events

Register callbacks to handle button presses and releases:

    lp.OnButton(func(event launchpad.ButtonEvent) {
    	if event.Pressed {
    		fmt.Printf("Button pressed: %v\n", event.Button)
    		lp.SetButtonLED(event.Button, launchpad.ColorGreen, launchpad.BrightnessFull)
    	} else {
    		fmt.Printf("Button released: %v\n", event.Button)
    		lp.SetButtonLED(event.Button, launchpad.ColorOff, launchpad.BrightnessOff)
    	}
    })

Alternatively, use channels for event handling:

    go func() {
    	for event := range lp.ButtonEvents() {
    		if event.Pressed {
    			fmt.Printf("Button: %v\n", event.Button)
    		}
    	}
    }()

# Double-Buffering

For smooth animations, use double-buffering to prepare the next frame while
displaying the current one:

    // Set up buffers
    lp.SetDisplayBuffer(launchpad.Buffer0)
    lp.SetUpdateBuffer(launchpad.Buffer1)

    // Update buffer 1 (invisible)
    for y := 0; y < 8; y++ {
    	for x := 0; x < 8; x++ {
    		lp.SetLED(x, y, launchpad.ColorRed, launchpad.BrightnessFull)
    	}
    }

    // Swap buffers for instant update
    lp.SwapBuffers()

# Colors and Brightness

Available colors:
  - ColorOff: LED off
  - ColorRed: Red
  - ColorGreen: Green
  - ColorAmber: Amber (red + green mix)
  - ColorYellow: Yellow (green with some red)

Brightness levels:
  - BrightnessOff: 0 (off)
  - BrightnessLow: 1
  - BrightnessMedium: 2
  - BrightnessFull: 3

# Advanced LED Control

For advanced users who need direct control over red and green components:

    state := launchpad.LEDState{
    	Red:   launchpad.BrightnessFull,
    	Green: launchpad.BrightnessMedium,
    	Flash: false,
    }
    lp.SetLEDState(x, y, state)

    // Or calculate raw velocity values
    velocity := launchpad.GetVelocityRGB(red, green, flash)

# Hardware Layout

The Launchpad Mini consists of:
  - 64 grid buttons (8×8, coordinates 0-7 for x and y)
  - 8 scene buttons (right column, y coordinates 0-7)
  - 8 top buttons (top row, x coordinates 0-7)

Grid coordinate system:

        0   1   2   3   4   5   6   7   [Scene]
    0  [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ]   [8]
    1  [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ]   [9]
    2  [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ]   [10]
    ...

# System Commands

    // Reset device to defaults
    lp.Reset()

    // Set mapping mode (X-Y or Drum rack layout)
    lp.SetMappingMode(launchpad.MappingXY)

    // Test all LEDs at specified brightness
    lp.TestLEDs(launchpad.BrightnessFull)

# Error Handling

Most methods return an error which should be checked:

    if err := lp.SetLED(x, y, color, brightness); err != nil {
    	log.Printf("Failed to set LED: %v", err)
    }

# Cleanup

Always close the connection when done to reset the device and free resources:

    defer lp.Close()

The Close method automatically:
  - Resets the device (turns off all LEDs)
  - Stops MIDI listeners
  - Closes MIDI connections
  - Closes the MIDI driver

# Performance

The library automatically handles MIDI message rate limiting (400 messages per
second) to prevent overwhelming the device. Messages are queued and sent at the
appropriate rate.

# Thread Safety

All public methods are thread-safe and can be called from multiple goroutines.
The library uses mutexes internally to protect shared state.

# Examples

See the examples directory for complete working examples:
  - examples/basic: Basic LED control
  - examples/buttons: Button event handling
  - examples/rainbow: Animated rainbow effect
  - examples/animation: Double-buffered bouncing ball
  - examples/gameoflife: Conway's Game of Life implementation

const GridWidth = 8 ...
const VelocityOff = 0x0C ...
const VelocityRedFullFlash = 0x0B ...
const MaxMessagesPerSecond = 400
func GetVelocity(color Color, brightness Brightness) byte
func GetVelocityRGB(red, green Brightness, flash bool) byte
type Brightness int
    const BrightnessOff Brightness = 0 ...
type BufferID int
    const Buffer0 BufferID = 0 ...
type Button struct{ ... }
    func NewGridButton(x, y int) Button
    func NewSceneButton(y int) Button
    func NewTopButton(x int) Button
type ButtonEvent struct{ ... }
type ButtonHandler func(ButtonEvent)
type Color int
    const ColorOff Color = iota ...
type LEDState struct{ ... }
    func NewLEDState(color Color, brightness Brightness) LEDState
type Launchpad struct{ ... }
    func New() *Launchpad
type MappingMode int
    const MappingXY MappingMode = iota ...

---

## Main Type: Launchpad

package launchpad // import "github.com/inegm/golp/pkg/launchpad"

type Launchpad struct {
	// Has unexported fields.
}
    Launchpad represents a connection to a Launchpad Mini device

func New() *Launchpad
func (lp *Launchpad) ButtonEvents() <-chan ButtonEvent
func (lp *Launchpad) Clear() error
func (lp *Launchpad) Close() error
func (lp *Launchpad) CopyBuffer() error
func (lp *Launchpad) EnableFlash(enabled bool) error
func (lp *Launchpad) GetDisplayBuffer() BufferID
func (lp *Launchpad) GetMappingMode() MappingMode
func (lp *Launchpad) GetUpdateBuffer() BufferID
func (lp *Launchpad) IsFlashEnabled() bool
func (lp *Launchpad) OnButton(handler ButtonHandler)
func (lp *Launchpad) Open() error
func (lp *Launchpad) Reset() error
func (lp *Launchpad) SetAllLEDs(color Color, brightness Brightness) error
func (lp *Launchpad) SetAllSceneButtons(color Color, brightness Brightness) error
func (lp *Launchpad) SetAllTopButtons(color Color, brightness Brightness) error
func (lp *Launchpad) SetButtonLED(btn Button, color Color, brightness Brightness) error
func (lp *Launchpad) SetButtonLEDState(btn Button, state LEDState) error
func (lp *Launchpad) SetColumn(x int, color Color, brightness Brightness) error
func (lp *Launchpad) SetDisplayBuffer(buffer BufferID) error
func (lp *Launchpad) SetLED(x, y int, color Color, brightness Brightness) error
func (lp *Launchpad) SetLEDState(x, y int, state LEDState) error
func (lp *Launchpad) SetMappingMode(mode MappingMode) error
func (lp *Launchpad) SetRow(y int, color Color, brightness Brightness) error
func (lp *Launchpad) SetSceneButton(y int, color Color, brightness Brightness) error
func (lp *Launchpad) SetTopButton(x int, color Color, brightness Brightness) error
func (lp *Launchpad) SetUpdateBuffer(buffer BufferID) error
func (lp *Launchpad) SwapBuffers() error
func (lp *Launchpad) TestLEDs(brightness Brightness) error

### Methods

#### New
```go
```

#### Open
```go
func (lp *Launchpad) Open() error
    Open opens a connection to the Launchpad device

```

#### Close
```go
func (lp *Launchpad) Close() error
    Close closes the connection to the Launchpad Resets the device (turns off
    all LEDs) before closing

```

#### Reset
```go
func (lp *Launchpad) Reset() error
    Reset resets the Launchpad to default state Turns off all LEDs, resets
    mapping mode, buffers, and duty cycle

```

#### SetLED
```go
func (lp *Launchpad) SetLED(x, y int, color Color, brightness Brightness) error
    SetLED sets the color and brightness of a single LED

```

#### SetRow
```go
func (lp *Launchpad) SetRow(y int, color Color, brightness Brightness) error
    SetRow sets all LEDs in a row to the same color and brightness

```

#### SetColumn
```go
func (lp *Launchpad) SetColumn(x int, color Color, brightness Brightness) error
    SetColumn sets all LEDs in a column to the same color and brightness

```

#### SetAllLEDs
```go
func (lp *Launchpad) SetAllLEDs(color Color, brightness Brightness) error
    SetAllLEDs sets all grid LEDs to the same color and brightness

```

#### Clear
```go
func (lp *Launchpad) Clear() error
    Clear turns off all LEDs

```

#### SetSceneButton
```go
func (lp *Launchpad) SetSceneButton(y int, color Color, brightness Brightness) error
    SetSceneButton sets the LED for a scene button (right column)

```

#### SetTopButton
```go
func (lp *Launchpad) SetTopButton(x int, color Color, brightness Brightness) error
    SetTopButton sets the LED for a top row button

```

#### OnButton
```go
func (lp *Launchpad) OnButton(handler ButtonHandler)
    OnButton registers a handler for button events

```

#### ButtonEvents
```go
func (lp *Launchpad) ButtonEvents() <-chan ButtonEvent
    ButtonEvents returns a channel that receives button events

```

#### SetDisplayBuffer
```go
func (lp *Launchpad) SetDisplayBuffer(buffer BufferID) error
    SetDisplayBuffer sets which buffer is displayed The display buffer is what's
    currently visible on the Launchpad

```

#### SetUpdateBuffer
```go
func (lp *Launchpad) SetUpdateBuffer(buffer BufferID) error
    SetUpdateBuffer sets which buffer receives LED updates The update buffer is
    where new LED states are written

```

#### SwapBuffers
```go
func (lp *Launchpad) SwapBuffers() error
    SwapBuffers swaps the display and update buffers This is useful for
    double-buffering: update one buffer while displaying the other, then swap
    for instant visual update

```

#### SetMappingMode
```go
func (lp *Launchpad) SetMappingMode(mode MappingMode) error
    SetMappingMode sets the button layout mapping mode

```

#### TestLEDs
```go
func (lp *Launchpad) TestLEDs(brightness Brightness) error
    TestLEDs turns on all LEDs at the specified brightness for testing This also
    resets all other device state

```

