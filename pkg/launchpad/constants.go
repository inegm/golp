package launchpad

// MIDI message status bytes
const (
	statusNoteOff         = 0x80 // 128
	statusNoteOn          = 0x90 // 144
	statusControlChange   = 0xB0 // 176
	statusNoteOnChannel3  = 0x92 // 146 - for rapid LED updates
)

// MIDI controller numbers for top row buttons (Automap/Live)
const (
	controllerTopButton0 = 104 // 0x68 - leftmost
	controllerTopButton1 = 105 // 0x69
	controllerTopButton2 = 106 // 0x6A
	controllerTopButton3 = 107 // 0x6B
	controllerTopButton4 = 108 // 0x6C
	controllerTopButton5 = 109 // 0x6D
	controllerTopButton6 = 110 // 0x6E
	controllerTopButton7 = 111 // 0x6F - rightmost
)

// System command controller (always 0 for system commands)
const (
	controllerSystem = 0
)

// System command data values
const (
	systemReset         = 0    // 0x00 - Reset to defaults
	systemLayoutXY      = 1    // 0x01 - X-Y layout mode
	systemLayoutDrum    = 2    // 0x02 - Drum rack layout mode
	systemTestLow       = 125  // 0x7D - All LEDs on low brightness
	systemTestMedium    = 126  // 0x7E - All LEDs on medium brightness
	systemTestFull      = 127  // 0x7F - All LEDs on full brightness
)

// Double-buffering base values (add to 32)
const (
	bufferBase         = 32   // Base value for buffer commands
	bufferFlagCopy     = 16   // Copy update buffer to display buffer
	bufferFlagFlash    = 8    // Enable auto-flash mode
)

// Velocity byte bit positions and masks
const (
	velocityFlagsClear = 0x08 // Clear bit (bit 3)
	velocityFlagsCopy  = 0x04 // Copy bit (bit 2)
	velocityFlagsNormal = velocityFlagsCopy | velocityFlagsClear // 0x0C = 12
	velocityFlagsFlash  = velocityFlagsClear // 0x08 = 8
)

// Grid dimensions
const (
	GridWidth  = 8 // Number of columns in main grid
	GridHeight = 8 // Number of rows in main grid
	SceneButtons = 8 // Number of scene buttons (right column)
	TopButtons   = 8 // Number of top row buttons
)

// Button velocity values
const (
	velocityPressed  = 127 // 0x7F - Button pressed
	velocityReleased = 0   // 0x00 - Button released
)

// MIDI rate limiting
const (
	MaxMessagesPerSecond = 400 // Maximum MIDI messages per second
)

// Pre-calculated velocity values for common colors (normal mode, flags = 12)
// Formula: velocity = (16 × green) + red + flags
const (
	// Off
	VelocityOff = 0x0C // 12

	// Red (red=3, green=0)
	VelocityRedLow    = 0x0D // 13 - red=1
	VelocityRedMedium = 0x0E // 14 - red=2
	VelocityRedFull   = 0x0F // 15 - red=3

	// Amber (red+green mix)
	VelocityAmberLow    = 0x1D // 29 - red=1, green=1
	VelocityAmberMedium = 0x2E // 46 - red=2, green=2
	VelocityAmberFull   = 0x3F // 63 - red=3, green=3

	// Yellow (green=3, red=2)
	VelocityYellowFull = 0x3E // 62 - red=2, green=3

	// Green (red=0, green=3)
	VelocityGreenLow    = 0x1C // 28 - green=1
	VelocityGreenMedium = 0x2C // 44 - green=2
	VelocityGreenFull   = 0x3C // 60 - green=3
)

// Pre-calculated velocity values for flashing colors (flags = 8)
const (
	VelocityRedFullFlash    = 0x0B // 11
	VelocityAmberFullFlash  = 0x3B // 59
	VelocityYellowFullFlash = 0x3A // 58
	VelocityGreenFullFlash  = 0x38 // 56
)

// Brightness constants for LED brightness levels (0-3)
const (
	BrightnessOff    Brightness = 0 // LED off
	BrightnessLow    Brightness = 1 // Low brightness
	BrightnessMedium Brightness = 2 // Medium brightness
	BrightnessFull   Brightness = 3 // Full brightness
)
