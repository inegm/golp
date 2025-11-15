package launchpad

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // auto-register rtmidi driver
)

// midiConnection wraps a MIDI input/output connection
type midiConnection struct {
	in  drivers.In
	out drivers.Out
}

// findLaunchpad searches for a connected Launchpad Mini device
// Returns the input and output ports if found
func findLaunchpad() (drivers.In, drivers.Out, error) {
	ins := midi.GetInPorts()
	outs := midi.GetOutPorts()

	// Look for Launchpad in input ports
	var inPort drivers.In
	for _, port := range ins {
		name := port.String()
		// Launchpad Mini typically shows up as "Launchpad Mini" or similar
		if containsLaunchpad(name) {
			inPort = port
			break
		}
	}

	if inPort == nil {
		return nil, nil, fmt.Errorf("launchpad input port not found")
	}

	// Look for matching output port
	var outPort drivers.Out
	for _, port := range outs {
		name := port.String()
		if containsLaunchpad(name) {
			outPort = port
			break
		}
	}

	if outPort == nil {
		return nil, nil, fmt.Errorf("launchpad output port not found")
	}

	return inPort, outPort, nil
}

// containsLaunchpad checks if a port name contains "launchpad"
func containsLaunchpad(name string) bool {
	// Convert to lowercase for case-insensitive matching
	lower := ""
	for _, r := range name {
		if r >= 'A' && r <= 'Z' {
			lower += string(r + 32)
		} else {
			lower += string(r)
		}
	}

	// Check for "launchpad" substring
	target := "launchpad"
	if len(lower) < len(target) {
		return false
	}

	for i := 0; i <= len(lower)-len(target); i++ {
		if lower[i:i+len(target)] == target {
			return true
		}
	}

	return false
}

// openMIDI opens MIDI input and output connections to the Launchpad
func openMIDI() (*midiConnection, error) {
	inPort, outPort, err := findLaunchpad()
	if err != nil {
		return nil, fmt.Errorf("failed to find launchpad: %w", err)
	}

	err = inPort.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open input port: %w", err)
	}

	err = outPort.Open()
	if err != nil {
		inPort.Close()
		return nil, fmt.Errorf("failed to open output port: %w", err)
	}

	return &midiConnection{
		in:  inPort,
		out: outPort,
	}, nil
}

// close closes the MIDI connection
func (mc *midiConnection) close() error {
	var inErr, outErr error
	if mc.in != nil {
		inErr = mc.in.Close()
	}
	if mc.out != nil {
		outErr = mc.out.Close()
	}

	if inErr != nil {
		return inErr
	}
	return outErr
}

// sendMessage sends a 3-byte MIDI message to the Launchpad
func (mc *midiConnection) sendMessage(status, data1, data2 byte) error {
	return mc.out.Send([]byte{status, data1, data2})
}

// sendNoteOn sends a note-on message (LED control)
func (mc *midiConnection) sendNoteOn(key, velocity byte) error {
	return mc.sendMessage(statusNoteOn, key, velocity)
}

// sendNoteOff sends a note-off message
func (mc *midiConnection) sendNoteOff(key byte) error {
	return mc.sendMessage(statusNoteOff, key, 0)
}

// sendControlChange sends a controller change message
func (mc *midiConnection) sendControlChange(controller, data byte) error {
	return mc.sendMessage(statusControlChange, controller, data)
}

// startListening starts listening for MIDI input messages
// Calls the handler function for each received message
// Returns a stop function that should be called to stop listening
func (mc *midiConnection) startListening(handler func([]byte)) (func(), error) {
	// Set up a listener that calls the handler for each message
	stop, err := midi.ListenTo(mc.in, func(msg midi.Message, timestampms int32) {
		// Message is already a []byte alias, pass it directly
		handler([]byte(msg))
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start MIDI listener: %w", err)
	}

	return stop, nil
}
