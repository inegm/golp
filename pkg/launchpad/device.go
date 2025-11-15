package launchpad

import (
	"fmt"
	"sync"
	"time"

	"gitlab.com/gomidi/midi/v2"
)

// Launchpad represents a connection to a Launchpad Mini device
type Launchpad struct {
	midi *midiConnection
	mu   sync.Mutex

	// State
	mappingMode   MappingMode
	displayBuffer BufferID
	updateBuffer  BufferID
	flashEnabled  bool

	// Message rate limiting
	msgQueue      chan message
	stopQueue     chan struct{}
	queueRunning  bool

	// Event handling
	buttonHandlers []ButtonHandler
	eventChan      chan ButtonEvent
	stopListener   chan struct{}
	listenerStop   func() // Function to stop MIDI listener
}

// message represents a queued MIDI message
type message struct {
	status byte
	data1  byte
	data2  byte
}

// ButtonHandler is a function that handles button events
type ButtonHandler func(ButtonEvent)

// New creates a new Launchpad instance but does not connect to the device
func New() *Launchpad {
	return &Launchpad{
		mappingMode:   MappingXY,
		displayBuffer: Buffer0,
		updateBuffer:  Buffer0,
		flashEnabled:  false,
		msgQueue:      make(chan message, 100), // Buffer up to 100 messages
		stopQueue:     make(chan struct{}),
		stopListener:  make(chan struct{}),
		eventChan:     make(chan ButtonEvent, 50), // Buffer up to 50 events
	}
}

// Open opens a connection to the Launchpad device
func (lp *Launchpad) Open() error {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.midi != nil {
		return fmt.Errorf("launchpad already open")
	}

	// Open MIDI connection
	midi, err := openMIDI()
	if err != nil {
		return err
	}

	lp.midi = midi

	// Start message queue processor
	go lp.processMessageQueue()
	lp.queueRunning = true

	// Start input listener
	stopFunc, err := lp.midi.startListening(lp.handleIncomingMessage)
	if err != nil {
		lp.midi.close()
		lp.midi = nil
		return fmt.Errorf("failed to start listener: %w", err)
	}
	lp.listenerStop = stopFunc

	// Reset the device to a known state (without locking - we already have the lock)
	err = lp.sendControlChange(controllerSystem, systemReset)
	if err != nil {
		lp.midi.close()
		lp.midi = nil
		return fmt.Errorf("failed to reset device: %w", err)
	}

	// Update internal state to match reset
	lp.mappingMode = MappingXY
	lp.displayBuffer = Buffer0
	lp.updateBuffer = Buffer0
	lp.flashEnabled = false

	return nil
}

// Close closes the connection to the Launchpad
// Resets the device (turns off all LEDs) before closing
func (lp *Launchpad) Close() error {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.midi == nil {
		return nil // Already closed
	}

	// Reset the device to turn off all LEDs (ignore errors - we're closing anyway)
	lp.sendControlChange(controllerSystem, systemReset)

	// Give the reset command time to be sent
	time.Sleep(50 * time.Millisecond)

	// Stop MIDI listener
	if lp.listenerStop != nil {
		lp.listenerStop()
		lp.listenerStop = nil
	}

	// Stop message queue
	if lp.queueRunning {
		close(lp.stopQueue)
		lp.queueRunning = false
	}

	// Stop listener channel
	close(lp.stopListener)

	// Close MIDI connection
	err := lp.midi.close()
	lp.midi = nil

	// Close the MIDI driver
	midi.CloseDriver()

	return err
}

// Reset resets the Launchpad to default state
// Turns off all LEDs, resets mapping mode, buffers, and duty cycle
func (lp *Launchpad) Reset() error {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.midi == nil {
		return fmt.Errorf("launchpad not open")
	}

	// Send reset command
	err := lp.sendControlChange(controllerSystem, systemReset)
	if err != nil {
		return fmt.Errorf("failed to reset: %w", err)
	}

	// Update internal state
	lp.mappingMode = MappingXY
	lp.displayBuffer = Buffer0
	lp.updateBuffer = Buffer0
	lp.flashEnabled = false

	return nil
}

// SetMappingMode sets the button layout mapping mode
func (lp *Launchpad) SetMappingMode(mode MappingMode) error {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.midi == nil {
		return fmt.Errorf("launchpad not open")
	}

	var data byte
	switch mode {
	case MappingXY:
		data = systemLayoutXY
	case MappingDrum:
		data = systemLayoutDrum
	default:
		return fmt.Errorf("invalid mapping mode: %v", mode)
	}

	err := lp.sendControlChange(controllerSystem, data)
	if err != nil {
		return fmt.Errorf("failed to set mapping mode: %w", err)
	}

	lp.mappingMode = mode
	return nil
}

// GetMappingMode returns the current mapping mode
func (lp *Launchpad) GetMappingMode() MappingMode {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	return lp.mappingMode
}

// TestLEDs turns on all LEDs at the specified brightness for testing
// This also resets all other device state
func (lp *Launchpad) TestLEDs(brightness Brightness) error {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.midi == nil {
		return fmt.Errorf("launchpad not open")
	}

	var data byte
	switch brightness {
	case BrightnessLow:
		data = systemTestLow
	case BrightnessMedium:
		data = systemTestMedium
	case BrightnessFull:
		data = systemTestFull
	default:
		return fmt.Errorf("invalid brightness for test mode: %v", brightness)
	}

	return lp.sendControlChange(controllerSystem, data)
}

// processMessageQueue processes queued messages with rate limiting
func (lp *Launchpad) processMessageQueue() {
	ticker := time.NewTicker(time.Second / MaxMessagesPerSecond)
	defer ticker.Stop()

	for {
		select {
		case <-lp.stopQueue:
			return
		case msg := <-lp.msgQueue:
			<-ticker.C // Wait for rate limit
			if lp.midi != nil {
				lp.midi.sendMessage(msg.status, msg.data1, msg.data2)
			}
		}
	}
}

// queueMessage adds a message to the send queue
func (lp *Launchpad) queueMessage(status, data1, data2 byte) {
	lp.msgQueue <- message{status, data1, data2}
}

// sendNoteOn sends a note-on message (bypassing queue for immediate sending)
func (lp *Launchpad) sendNoteOn(key, velocity byte) error {
	if lp.midi == nil {
		return fmt.Errorf("launchpad not open")
	}
	return lp.midi.sendNoteOn(key, velocity)
}

// sendControlChange sends a controller change message (bypassing queue)
func (lp *Launchpad) sendControlChange(controller, data byte) error {
	if lp.midi == nil {
		return fmt.Errorf("launchpad not open")
	}
	return lp.midi.sendControlChange(controller, data)
}

// handleIncomingMessage processes incoming MIDI messages
func (lp *Launchpad) handleIncomingMessage(msg []byte) {
	if len(msg) < 3 {
		return // Invalid message
	}

	status := msg[0]
	data1 := msg[1]
	data2 := msg[2]

	var btn Button
	var pressed bool

	switch status {
	case statusNoteOn:
		// Grid or scene button
		key := int(data1)
		x := key % 16
		y := key / 16

		if x == 8 {
			// Scene button
			btn = NewSceneButton(y)
		} else {
			// Grid button
			btn = NewGridButton(x, y)
		}

		pressed = data2 == velocityPressed

	case statusControlChange:
		// Top row button
		controller := int(data1)
		if controller >= controllerTopButton0 && controller <= controllerTopButton7 {
			x := controller - controllerTopButton0
			btn = NewTopButton(x)
			pressed = data2 == velocityPressed
		} else {
			return // Not a button event
		}

	default:
		return // Unknown message type
	}

	// Create event
	event := ButtonEvent{
		Button:  btn,
		Pressed: pressed,
	}

	// Send to event channel
	select {
	case lp.eventChan <- event:
	default:
		// Channel full, drop event
	}

	// Call registered handlers
	lp.mu.Lock()
	handlers := make([]ButtonHandler, len(lp.buttonHandlers))
	copy(handlers, lp.buttonHandlers)
	lp.mu.Unlock()

	for _, handler := range handlers {
		handler(event)
	}
}

// OnButton registers a handler for button events
func (lp *Launchpad) OnButton(handler ButtonHandler) {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	lp.buttonHandlers = append(lp.buttonHandlers, handler)
}

// ButtonEvents returns a channel that receives button events
func (lp *Launchpad) ButtonEvents() <-chan ButtonEvent {
	return lp.eventChan
}
