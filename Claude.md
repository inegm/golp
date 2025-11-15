# Launchpad Mini Go Library - Development Plan

## Project Overview

**Goal**: Create a Go library for controlling the Novation Launchpad Mini MIDI hardware device.

**Current Status**: Foundation phase - implementing core MIDI communication and device control.

**Documentation**: Complete programmer's reference available in `docs/launchpad-programmers-reference.md`

---

## Hardware Specifications

### Physical Layout
- **Total Buttons**: 80
  - 64 square grid buttons (8×8)
  - 8 round scene launch buttons (right side)
  - 8 round top buttons (Automap/Live control)

### LED System
- **Type**: Bi-colored LEDs (red + green elements)
- **Colors**: Red, Green, Amber (red+green), Yellow
- **Brightness Levels**: 4 levels (Off, Low, Medium, Full)

### MIDI Protocol
- **Communication**: MIDI note-on, note-off, controller change
- **Channel**: MIDI channel 1 (channel 3 for rapid updates)
- **Message Format**: Always 3 bytes
- **Rate Limit**: Maximum 400 messages/second
- **Update Time**: ~200ms for full surface (80 LEDs)

---

## Technical Specifications Summary

### Button Addressing

**X-Y Layout Mode (Default)**
- Formula: `Key = (16 × Row) + Column`
- Origin: Top-left corner (0,0)
- Range: 0-127
- Scene buttons: Column 8
- Top row: Controllers 104-111

**Drum Rack Layout Mode**
- 6 continuous musical octaves
- Optimized for drum programming

### LED Control

**Note-On Message Format**
```
90h, Key, Velocity (144, Key, Velocity)
```

**Velocity Byte Structure**
```
Bit 6:    Must be 0
Bit 5-4:  Green brightness (0-3)
Bit 3:    Clear bit (double-buffering)
Bit 2:    Copy bit (double-buffering)
Bit 1-0:  Red brightness (0-3)
```

**Velocity Calculation**
```
Hex:     Velocity = (10h × Green) + Red + Flags
Decimal: Velocity = (16 × Green) + Red + Flags
Flags:   12 (0Ch) for normal, 8 for flash, 0 for buffering
```

### Common Color Values

| Color | Brightness | Hex | Decimal |
|-------|-----------|-----|---------|
| Off | - | 0Ch | 12 |
| Red | Low | 0Dh | 13 |
| Red | Full | 0Fh | 15 |
| Amber | Low | 1Dh | 29 |
| Amber | Full | 3Fh | 63 |
| Yellow | Full | 3Eh | 62 |
| Green | Low | 1Ch | 28 |
| Green | Full | 3Ch | 60 |

### System Commands (Controller Change)

| Command | Hex | Decimal | Purpose |
|---------|-----|---------|---------|
| Reset | B0h, 00h, 00h | 176, 0, 0 | Reset device |
| X-Y Mode | B0h, 00h, 01h | 176, 0, 1 | Set X-Y layout |
| Drum Mode | B0h, 00h, 02h | 176, 0, 2 | Set drum layout |
| Test Low | B0h, 00h, 7Dh | 176, 0, 125 | All LEDs low |
| Test Med | B0h, 00h, 7Eh | 176, 0, 126 | All LEDs medium |
| Test Full | B0h, 00h, 7Fh | 176, 0, 127 | All LEDs full |

### Input Messages

**Grid Button**
- Press: `90h, Key, 7Fh` (144, Key, 127)
- Release: `90h, Key, 00h` (144, Key, 0)

**Top Row Button**
- Press: `B0h, Controller, 7Fh` (176, 104-111, 127)
- Release: `B0h, Controller, 00h` (176, 104-111, 0)

### Advanced Features

**Double-Buffering**
- Two independent LED buffers (0 and 1)
- Allows instant visual updates
- Command: `B0h, 00h, 20-3Dh` (176, 0, 32-61)
- Data formula: `(4 × Update) + Display + 32 + Flags`

**Rapid Update Mode**
- MIDI Channel 3: `92h, Vel1, Vel2, 92h, Vel3, Vel4...`
- Updates 2 LEDs per message
- Doubles update speed (80 LEDs in 40 messages)

**Auto-Flash Mode**
- Enable: `B0h, 00h, 28h` (176, 0, 40)
- LEDs marked with flash bit automatically flash

---

## Implementation Plan

### Phase 1: Foundation ✓

1. **Claude.md** - This documentation file
2. **MIDI Dependency** - Add `gitlab.com/gomidi/midi`
3. **Package Structure** - Organize into logical files

### Phase 2: Core Implementation

4. **constants.go** - MIDI message types, color constants, controller numbers
5. **types.go** - Core types:
   - `Color` enum (Red, Green, Amber, Yellow, Off)
   - `Brightness` enum (Off, Low, Medium, Full)
   - `Button` struct (X, Y coordinates)
   - `ButtonEvent` (button, pressed/released)
   - `MappingMode` enum (XY, DrumRack)

6. **midi.go** - Low-level MIDI wrapper:
   - Device discovery
   - Send/receive MIDI messages
   - Connection management

7. **device.go** - Main Launchpad device:
   - `NewLaunchpad()` constructor
   - `Open()` / `Close()` connection
   - `Reset()` system reset
   - `SetMappingMode()` layout selection
   - Message queue and rate limiting

8. **led.go** - LED control API:
   - `SetLED(x, y, color, brightness)` individual LED
   - `SetAllLEDs(color, brightness)` bulk operation
   - `Clear()` turn off all LEDs
   - `Test(brightness)` test mode
   - Velocity calculation helpers

9. **input.go** - Event handling:
   - Goroutine for MIDI listening
   - `OnButtonPress(callback)` registration
   - `OnButtonRelease(callback)` registration
   - `ButtonEvents()` channel-based API
   - Parse incoming MIDI to events

10. **buffer.go** - Double-buffering:
    - `SetDisplayBuffer(0/1)` select visible buffer
    - `SetUpdateBuffer(0/1)` select write buffer
    - `SwapBuffers()` instant switch
    - `EnableFlash(bool)` auto-flash mode
    - `CopyBuffer()` duplicate buffer contents

### Phase 3: Advanced Features

11. **Rapid Update Mode**
    - Detect when rapid mode beneficial
    - Implement channel 3 messaging
    - Batch LED updates for performance

12. **Rate Limiting**
    - Message queue with 400/sec limit
    - Automatic throttling
    - Priority queue (system commands > LED updates)

13. **Mapping Modes**
    - Translate X-Y to drum rack coordinates
    - Support both modes transparently

### Phase 4: Developer Experience

14. **Examples**
    - `examples/basic/` - Simple LED control
    - `examples/rainbow/` - Animated rainbow grid
    - `examples/buttons/` - Button event handling
    - `examples/animation/` - Double-buffer animation

15. **README.md**
    - Installation instructions
    - Quick start guide
    - API reference
    - Code examples

16. **Tests**
    - Unit tests with mock MIDI
    - Integration tests (optional, requires hardware)
    - Benchmarks for performance

---

## API Design Principles

### 1. Simplicity First
```go
// Simple case: just set an LED
lp.SetLED(3, 4, launchpad.Red, launchpad.Full)

// Advanced case: double-buffering
lp.SetUpdateBuffer(1)
lp.SetLED(0, 0, launchpad.Green, launchpad.Full)
lp.SwapBuffers()
```

### 2. Event-Driven Input
```go
// Callback style
lp.OnButtonPress(func(btn Button) {
    fmt.Printf("Pressed: %d, %d\n", btn.X, btn.Y)
})

// Channel style
for event := range lp.ButtonEvents() {
    if event.Pressed {
        // Handle press
    }
}
```

### 3. Safe Defaults
- X-Y mapping mode by default
- Automatic rate limiting
- Clear/copy bits set for normal use
- Error handling with meaningful messages

### 4. Type Safety
```go
// Use enums instead of magic numbers
lp.SetLED(x, y, launchpad.Amber, launchpad.Medium)

// Not: lp.SetLED(x, y, 29)
```

### 5. Performance Aware
- Batch updates when possible
- Automatic rapid mode selection
- Non-blocking event handling
- Efficient message queuing

---

## Package Structure

```
github.com/inegm/golp/
├── go.mod
├── go.sum
├── Claude.md                    # This file
├── README.md                    # User documentation
├── docs/
│   └── launchpad-programmers-reference.md
├── pkg/launchpad/
│   ├── launchpad.go            # Package docs + re-exports
│   ├── constants.go            # MIDI constants, colors
│   ├── types.go                # Core type definitions
│   ├── midi.go                 # Low-level MIDI I/O
│   ├── device.go               # Main Launchpad device
│   ├── led.go                  # LED control functions
│   ├── input.go                # Button event handling
│   ├── buffer.go               # Double-buffering
│   └── queue.go                # Message rate limiting
├── examples/
│   ├── basic/main.go
│   ├── rainbow/main.go
│   ├── buttons/main.go
│   └── animation/main.go
└── pkg/launchpad/launchpad_test.go
```

---

## Critical Implementation Notes

### Rate Limiting
- **Hard Limit**: 400 messages/second
- **Implementation**: Token bucket or sliding window
- **Behavior**: Queue messages, send at controlled rate
- **Priority**: System commands before LED updates

### Message Format
- **No Running Status**: Always send full 3-byte messages
- **Status Byte**: First byte (80h, 90h, B0h, 92h)
- **Data Bytes**: Always < 128 (bit 7 = 0)

### Buffer Management
- **Default State**: Both buffers = 0, no flash
- **Copy/Clear Bits**: Should be set (12) for normal LED updates
- **Flash Bit**: Only set when flash desired (8 instead of 12)

### Concurrency
- **MIDI I/O**: Single goroutine for sending
- **Event Listener**: Separate goroutine for receiving
- **Thread Safety**: Mutex for device state
- **Channels**: Buffered for event distribution

### Error Handling
- **Device Not Found**: Clear error message
- **Connection Lost**: Graceful degradation
- **Invalid Coordinates**: Bounds checking
- **MIDI Errors**: Wrap with context

---

## Dependencies

### External Libraries
- **MIDI**: `gitlab.com/gomidi/midi/v2` (pure Go, cross-platform)
  - Alternative: `github.com/rakyll/portmidi` (CGo, requires native library)

### Standard Library
- `fmt` - Formatting and errors
- `sync` - Mutexes and coordination
- `time` - Rate limiting timing
- `context` - Cancellation for goroutines

---

## Testing Strategy

### Unit Tests
- Color/brightness calculations
- Coordinate transformations
- Velocity byte encoding
- Message queue behavior

### Mock Testing
- Mock MIDI device for integration tests
- Test full API without hardware
- Verify message sequences

### Hardware Tests (Optional)
- Require actual Launchpad
- Manual verification of visual output
- Performance benchmarks

---

## Future Enhancements

### Beyond MVP
- **Launchpad Pro Support** - Extended grid, RGB LEDs
- **Launchpad S Support** - Additional models
- **Animation Helpers** - Built-in effects library
- **Pattern Recording** - Record/playback LED patterns
- **Performance Mode** - Zero-allocation hot path

### Community
- **Examples Gallery** - Community-contributed examples
- **Application Showcase** - Projects using the library
- **Video Tutorials** - Visual learning resources

---

## Development Status

- [x] Documentation review
- [x] Architecture planning
- [ ] Core implementation
- [ ] Advanced features
- [ ] Examples and documentation
- [ ] Testing and refinement

---

## Resources

- **Programmer's Reference**: `docs/launchpad-programmers-reference.md`
- **MIDI Specification**: https://www.midi.org/specifications
- **Go MIDI Library**: https://gitlab.com/gomidi/midi
- **Novation Support**: https://novationmusic.com/support

---

*Last Updated: 2025-11-15*
*Project: github.com/inegm/golp*
*Go Version: 1.24.4*
