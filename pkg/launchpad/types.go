package launchpad

import "fmt"

// Color represents the color of an LED on the Launchpad
type Color int

const (
	ColorOff    Color = iota // LED off
	ColorRed                 // Red LED
	ColorGreen               // Green LED
	ColorAmber               // Amber (red + green mix)
	ColorYellow              // Yellow (red + green, more green)
)

// String returns the string representation of a Color
func (c Color) String() string {
	switch c {
	case ColorOff:
		return "Off"
	case ColorRed:
		return "Red"
	case ColorGreen:
		return "Green"
	case ColorAmber:
		return "Amber"
	case ColorYellow:
		return "Yellow"
	default:
		return fmt.Sprintf("Color(%d)", c)
	}
}

// Brightness represents the brightness level of an LED
type Brightness int

// String returns the string representation of a Brightness
func (b Brightness) String() string {
	switch b {
	case BrightnessOff:
		return "Off"
	case BrightnessLow:
		return "Low"
	case BrightnessMedium:
		return "Medium"
	case BrightnessFull:
		return "Full"
	default:
		return fmt.Sprintf("Brightness(%d)", b)
	}
}

// Valid returns true if the brightness is valid (0-3)
func (b Brightness) Valid() bool {
	return b >= 0 && b <= 3
}

// MappingMode represents the button layout mode
type MappingMode int

const (
	MappingXY   MappingMode = iota // X-Y layout mode (default)
	MappingDrum                    // Drum rack layout mode
)

// String returns the string representation of a MappingMode
func (m MappingMode) String() string {
	switch m {
	case MappingXY:
		return "X-Y"
	case MappingDrum:
		return "Drum"
	default:
		return fmt.Sprintf("MappingMode(%d)", m)
	}
}

// Button represents a button on the Launchpad
type Button struct {
	X       int        // Column position (0-7 for grid, 8 for scene buttons)
	Y       int        // Row position (0-7 for grid, -1 for top buttons)
	IsScene bool       // True if this is a scene button (right column)
	IsTop   bool       // True if this is a top row button
}

// NewGridButton creates a Button for a grid position
func NewGridButton(x, y int) Button {
	return Button{
		X:       x,
		Y:       y,
		IsScene: false,
		IsTop:   false,
	}
}

// NewSceneButton creates a Button for a scene button
func NewSceneButton(y int) Button {
	return Button{
		X:       8,
		Y:       y,
		IsScene: true,
		IsTop:   false,
	}
}

// NewTopButton creates a Button for a top row button
func NewTopButton(x int) Button {
	return Button{
		X:       x,
		Y:       -1,
		IsScene: false,
		IsTop:   true,
	}
}

// String returns the string representation of a Button
func (b Button) String() string {
	if b.IsTop {
		return fmt.Sprintf("Top[%d]", b.X)
	}
	if b.IsScene {
		return fmt.Sprintf("Scene[%d]", b.Y)
	}
	return fmt.Sprintf("Grid[%d,%d]", b.X, b.Y)
}

// Valid returns true if the button coordinates are valid
func (b Button) Valid() bool {
	if b.IsTop {
		return b.X >= 0 && b.X < TopButtons
	}
	if b.IsScene {
		return b.Y >= 0 && b.Y < SceneButtons
	}
	return b.X >= 0 && b.X < GridWidth && b.Y >= 0 && b.Y < GridHeight
}

// MIDIKey returns the MIDI key number for this button in X-Y mode
func (b Button) MIDIKey() int {
	if b.IsTop {
		return -1 // Top buttons use controller change, not note
	}
	// Formula: Key = (16 × Row) + Column
	return (16 * b.Y) + b.X
}

// MIDIController returns the MIDI controller number for top buttons
func (b Button) MIDIController() int {
	if !b.IsTop {
		return -1
	}
	return controllerTopButton0 + b.X
}

// ButtonEvent represents a button press or release event
type ButtonEvent struct {
	Button  Button // The button that triggered the event
	Pressed bool   // True if button was pressed, false if released
}

// String returns the string representation of a ButtonEvent
func (e ButtonEvent) String() string {
	action := "released"
	if e.Pressed {
		action = "pressed"
	}
	return fmt.Sprintf("%s %s", e.Button, action)
}

// BufferID represents a double-buffer identifier (0 or 1)
type BufferID int

const (
	Buffer0 BufferID = 0
	Buffer1 BufferID = 1
)

// Valid returns true if the buffer ID is valid (0 or 1)
func (b BufferID) Valid() bool {
	return b == Buffer0 || b == Buffer1
}

// String returns the string representation of a BufferID
func (b BufferID) String() string {
	return fmt.Sprintf("Buffer%d", b)
}

// LEDState represents the state of a single LED
type LEDState struct {
	Red   Brightness // Red component brightness (0-3)
	Green Brightness // Green component brightness (0-3)
	Flash bool       // True if LED should flash
}

// Velocity calculates the MIDI velocity byte for this LED state
func (s LEDState) Velocity() byte {
	flags := velocityFlagsNormal
	if s.Flash {
		flags = velocityFlagsFlash
	}
	// Formula: velocity = (16 × green) + red + flags
	return byte((16 * int(s.Green)) + int(s.Red) + flags)
}

// NewLEDState creates an LEDState from a color and brightness
func NewLEDState(color Color, brightness Brightness) LEDState {
	state := LEDState{Flash: false}

	if brightness == BrightnessOff {
		return state // Both red and green are 0
	}

	switch color {
	case ColorOff:
		// Both 0 (already set)
	case ColorRed:
		state.Red = brightness
	case ColorGreen:
		state.Green = brightness
	case ColorAmber:
		state.Red = brightness
		state.Green = brightness
	case ColorYellow:
		// Yellow is green=full, red=medium (or close to it)
		state.Green = brightness
		if brightness > BrightnessLow {
			state.Red = brightness - 1
		}
	}

	return state
}

// String returns the string representation of an LEDState
func (s LEDState) String() string {
	flash := ""
	if s.Flash {
		flash = " (flash)"
	}
	return fmt.Sprintf("R:%d G:%d%s", s.Red, s.Green, flash)
}
