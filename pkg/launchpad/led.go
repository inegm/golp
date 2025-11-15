package launchpad

import "fmt"

// SetLED sets the color and brightness of a single LED
func (lp *Launchpad) SetLED(x, y int, color Color, brightness Brightness) error {
	btn := NewGridButton(x, y)
	return lp.SetButtonLED(btn, color, brightness)
}

// SetButtonLED sets the color and brightness of an LED for any button type
func (lp *Launchpad) SetButtonLED(btn Button, color Color, brightness Brightness) error {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.midi == nil {
		return fmt.Errorf("launchpad not open")
	}

	if !btn.Valid() {
		return fmt.Errorf("invalid button: %v", btn)
	}

	if !brightness.Valid() {
		return fmt.Errorf("invalid brightness: %v", brightness)
	}

	// Create LED state
	state := NewLEDState(color, brightness)
	velocity := state.Velocity()

	// Send appropriate MIDI message based on button type
	if btn.IsTop {
		// Top row uses controller change
		controller := byte(btn.MIDIController())
		return lp.sendControlChange(controller, velocity)
	} else {
		// Grid and scene buttons use note-on
		key := byte(btn.MIDIKey())
		return lp.sendNoteOn(key, velocity)
	}
}

// SetLEDState sets the LED state using a custom LEDState (for advanced control)
func (lp *Launchpad) SetLEDState(x, y int, state LEDState) error {
	btn := NewGridButton(x, y)
	return lp.SetButtonLEDState(btn, state)
}

// SetButtonLEDState sets the LED state for any button type using a custom LEDState
func (lp *Launchpad) SetButtonLEDState(btn Button, state LEDState) error {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.midi == nil {
		return fmt.Errorf("launchpad not open")
	}

	if !btn.Valid() {
		return fmt.Errorf("invalid button: %v", btn)
	}

	velocity := state.Velocity()

	if btn.IsTop {
		controller := byte(btn.MIDIController())
		return lp.sendControlChange(controller, velocity)
	} else {
		key := byte(btn.MIDIKey())
		return lp.sendNoteOn(key, velocity)
	}
}

// Clear turns off all LEDs
func (lp *Launchpad) Clear() error {
	// Turning off all LEDs by setting them to off
	// We could iterate through all buttons, but Reset() also clears LEDs
	// For a more targeted clear, we iterate through all positions

	// Clear grid
	for y := 0; y < GridHeight; y++ {
		for x := 0; x < GridWidth; x++ {
			err := lp.SetLED(x, y, ColorOff, BrightnessOff)
			if err != nil {
				return err
			}
		}
	}

	// Clear scene buttons
	for y := 0; y < SceneButtons; y++ {
		btn := NewSceneButton(y)
		err := lp.SetButtonLED(btn, ColorOff, BrightnessOff)
		if err != nil {
			return err
		}
	}

	// Clear top buttons
	for x := 0; x < TopButtons; x++ {
		btn := NewTopButton(x)
		err := lp.SetButtonLED(btn, ColorOff, BrightnessOff)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetAllLEDs sets all grid LEDs to the same color and brightness
func (lp *Launchpad) SetAllLEDs(color Color, brightness Brightness) error {
	for y := 0; y < GridHeight; y++ {
		for x := 0; x < GridWidth; x++ {
			err := lp.SetLED(x, y, color, brightness)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SetRow sets all LEDs in a row to the same color and brightness
func (lp *Launchpad) SetRow(y int, color Color, brightness Brightness) error {
	if y < 0 || y >= GridHeight {
		return fmt.Errorf("invalid row: %d", y)
	}

	for x := 0; x < GridWidth; x++ {
		err := lp.SetLED(x, y, color, brightness)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetColumn sets all LEDs in a column to the same color and brightness
func (lp *Launchpad) SetColumn(x int, color Color, brightness Brightness) error {
	if x < 0 || x >= GridWidth {
		return fmt.Errorf("invalid column: %d", x)
	}

	for y := 0; y < GridHeight; y++ {
		err := lp.SetLED(x, y, color, brightness)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetSceneButton sets the LED for a scene button (right column)
func (lp *Launchpad) SetSceneButton(y int, color Color, brightness Brightness) error {
	if y < 0 || y >= SceneButtons {
		return fmt.Errorf("invalid scene button index: %d", y)
	}

	btn := NewSceneButton(y)
	return lp.SetButtonLED(btn, color, brightness)
}

// SetTopButton sets the LED for a top row button
func (lp *Launchpad) SetTopButton(x int, color Color, brightness Brightness) error {
	if x < 0 || x >= TopButtons {
		return fmt.Errorf("invalid top button index: %d", x)
	}

	btn := NewTopButton(x)
	return lp.SetButtonLED(btn, color, brightness)
}

// SetAllTopButtons sets all top row LEDs to the same color and brightness
func (lp *Launchpad) SetAllTopButtons(color Color, brightness Brightness) error {
	for x := 0; x < TopButtons; x++ {
		err := lp.SetTopButton(x, color, brightness)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetAllSceneButtons sets all scene button LEDs to the same color and brightness
func (lp *Launchpad) SetAllSceneButtons(color Color, brightness Brightness) error {
	for y := 0; y < SceneButtons; y++ {
		err := lp.SetSceneButton(y, color, brightness)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetVelocity calculates the MIDI velocity byte for a color and brightness
// This is a helper function for advanced users who want to work with raw velocity values
func GetVelocity(color Color, brightness Brightness) byte {
	state := NewLEDState(color, brightness)
	return state.Velocity()
}

// GetVelocityRGB calculates the MIDI velocity byte from red and green brightness values
// This is a helper function for advanced users who want direct RGB control
func GetVelocityRGB(red, green Brightness, flash bool) byte {
	state := LEDState{
		Red:   red,
		Green: green,
		Flash: flash,
	}
	return state.Velocity()
}
