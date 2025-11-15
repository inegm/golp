package launchpad

import "fmt"

// SetDisplayBuffer sets which buffer is displayed
// The display buffer is what's currently visible on the Launchpad
func (lp *Launchpad) SetDisplayBuffer(buffer BufferID) error {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.midi == nil {
		return fmt.Errorf("launchpad not open")
	}

	if !buffer.Valid() {
		return fmt.Errorf("invalid buffer ID: %v", buffer)
	}

	// Calculate buffer command data byte
	// Formula: data = (4 × update) + display + 32 + flags
	updateBit := int(lp.updateBuffer)
	displayBit := int(buffer)
	flags := 0
	if lp.flashEnabled {
		flags = bufferFlagFlash
	}

	data := byte((4 * updateBit) + displayBit + bufferBase + flags)

	err := lp.sendControlChange(controllerSystem, data)
	if err != nil {
		return fmt.Errorf("failed to set display buffer: %w", err)
	}

	lp.displayBuffer = buffer
	return nil
}

// SetUpdateBuffer sets which buffer receives LED updates
// The update buffer is where new LED states are written
func (lp *Launchpad) SetUpdateBuffer(buffer BufferID) error {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.midi == nil {
		return fmt.Errorf("launchpad not open")
	}

	if !buffer.Valid() {
		return fmt.Errorf("invalid buffer ID: %v", buffer)
	}

	// Calculate buffer command data byte
	updateBit := int(buffer)
	displayBit := int(lp.displayBuffer)
	flags := 0
	if lp.flashEnabled {
		flags = bufferFlagFlash
	}

	data := byte((4 * updateBit) + displayBit + bufferBase + flags)

	err := lp.sendControlChange(controllerSystem, data)
	if err != nil {
		return fmt.Errorf("failed to set update buffer: %w", err)
	}

	lp.updateBuffer = buffer
	return nil
}

// SwapBuffers swaps the display and update buffers
// This is useful for double-buffering: update one buffer while displaying the other,
// then swap for instant visual update
func (lp *Launchpad) SwapBuffers() error {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.midi == nil {
		return fmt.Errorf("launchpad not open")
	}

	// Swap the buffers
	newDisplay := lp.updateBuffer
	newUpdate := lp.displayBuffer

	// Calculate command
	updateBit := int(newUpdate)
	displayBit := int(newDisplay)
	flags := 0
	if lp.flashEnabled {
		flags = bufferFlagFlash
	}

	data := byte((4 * updateBit) + displayBit + bufferBase + flags)

	err := lp.sendControlChange(controllerSystem, data)
	if err != nil {
		return fmt.Errorf("failed to swap buffers: %w", err)
	}

	lp.displayBuffer = newDisplay
	lp.updateBuffer = newUpdate
	return nil
}

// CopyBuffer copies the update buffer to the display buffer
// This makes both buffers show the same content
func (lp *Launchpad) CopyBuffer() error {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.midi == nil {
		return fmt.Errorf("launchpad not open")
	}

	// Calculate command with copy flag
	updateBit := int(lp.updateBuffer)
	displayBit := int(lp.displayBuffer)
	flags := bufferFlagCopy
	if lp.flashEnabled {
		flags |= bufferFlagFlash
	}

	data := byte((4 * updateBit) + displayBit + bufferBase + flags)

	err := lp.sendControlChange(controllerSystem, data)
	if err != nil {
		return fmt.Errorf("failed to copy buffer: %w", err)
	}

	return nil
}

// EnableFlash enables automatic LED flashing
// LEDs marked with the flash flag will automatically flash
func (lp *Launchpad) EnableFlash(enabled bool) error {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.midi == nil {
		return fmt.Errorf("launchpad not open")
	}

	// Calculate command with flash flag
	updateBit := int(lp.updateBuffer)
	displayBit := int(lp.displayBuffer)
	flags := 0
	if enabled {
		flags = bufferFlagFlash
	}

	data := byte((4 * updateBit) + displayBit + bufferBase + flags)

	err := lp.sendControlChange(controllerSystem, data)
	if err != nil {
		return fmt.Errorf("failed to set flash mode: %w", err)
	}

	lp.flashEnabled = enabled
	return nil
}

// GetDisplayBuffer returns the current display buffer ID
func (lp *Launchpad) GetDisplayBuffer() BufferID {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	return lp.displayBuffer
}

// GetUpdateBuffer returns the current update buffer ID
func (lp *Launchpad) GetUpdateBuffer() BufferID {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	return lp.updateBuffer
}

// IsFlashEnabled returns whether auto-flash mode is enabled
func (lp *Launchpad) IsFlashEnabled() bool {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	return lp.flashEnabled
}
