package rotary

type Encoder struct {
	// Pins.
	clk uint
	dt  uint
	sw  uint

	// Public.
	Direction int
	Count     int

	// Data.
	currentStateCLK bool
	lastStateCLK    bool
	lastButtonPress int64

	// Hardware dependencies.
	getPin func(uint) bool
	now    func() int64
	sleep  func(ms int)
}

// New creates a new Rotary Encoder.
// The clk and dt values are the pins that will be used, and should be set as GPIO inputs.
// The sw pin should be set to pull up.
func New(getPin func(uint) bool, now func() int64, sleep func(ms int), clk uint, dt uint, sw uint) *Encoder {
	return &Encoder{
		clk:          clk,
		dt:           dt,
		sw:           sw,
		getPin:       getPin,
		now:          now,
		sleep:        sleep,
		lastStateCLK: getPin(clk),
	}
}

func (e *Encoder) Update() bool {
	// Read the current state of CLK.
	e.currentStateCLK = e.getPin(e.clk)

	// If last and current state of CLK are different, then a pulse occurred.
	// React to only 1 state change to avoid double count.
	if e.currentStateCLK != e.lastStateCLK && e.currentStateCLK {
		// If the DT state is different than the CLK state then
		// the encoder is rotating CW so increment.
		if e.getPin(e.dt) != e.currentStateCLK {
			e.Count++
			e.Direction = +1
		} else {
			// Encoder is rotating CCW so decrement.
			e.Count--
			e.Direction = -1
		}
	}

	// Remember last CLK state.
	e.lastStateCLK = e.currentStateCLK

	// Read the button state.
	var pressed bool

	// If we detect LOW signal, button is pressed.
	if !e.getPin(e.sw) {
		// If 50ms have passed since last LOW pulse, it means that the
		// button has been pressed, released and pressed again
		now := e.now()
		pressed = (now - e.lastButtonPress) > 50000

		// Remember last button press event.
		e.lastButtonPress = now
	}

	// Put in a slight delay to help debounce the reading
	e.sleep(1)
	return pressed
}
