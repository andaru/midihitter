package main

func setupWatcher(w watcher) {
	// button/controller events
	comBtnDown := Event(MIDI_CONTROL, 49, 127)
	leftBtn := Event(MIDI_CONTROL, 41, 127)
	rightBtn := Event(MIDI_CONTROL, 42, 127)
	playBtn := Event(MIDI_CONTROL, 39, 127)
	stopBtn := Event(MIDI_CONTROL, 40, 127)
	blueKnobRel := Event(MIDI_CONTROL, 1)
	greenKnobRel := Event(MIDI_CONTROL, 2)
	whiteKnobRel := Event(MIDI_CONTROL, 3)
	orangeKnobRel := Event(MIDI_CONTROL, 4)
	helpBtn := Event(MIDI_CONTROL, 5, 127)
	blueKnobClick := Event(MIDI_CONTROL, 64, 127)
	greenKnobClick := Event(MIDI_CONTROL, 65, 127)
	micBtn := Event(MIDI_CONTROL, 48, 127)
	tempoBtn := Event(MIDI_CONTROL, 6, 127)

	f1 := Event(MIDI_NOTE_ON, 53)
	g1 := Event(MIDI_NOTE_ON, 55)

	// fullscreen toggle
	w.AddEvent(comBtnDown, rule{Key: "backslash"})

	// scrub left
	w.AddEvent(leftBtn, rule{Key: "Left"})
	// scrub right
	w.AddEvent(rightBtn, rule{Key: "Right"})
	// play
	w.AddEvent(playBtn, rule{Key: "space"})
	// stop
	w.AddEvent(stopBtn, rule{Key: "x"})
	// blue knob in relative mode controls volume up/down
	w.AddEvent(blueKnobRel, rule{Data2: Exact(127), Key: "plus"})
	w.AddEvent(blueKnobRel, rule{Data2: Exact(1), Key: "minus"})
	// green knob moves up and down
	w.AddEvent(greenKnobRel, rule{Data2: Exact(127), Key: "Up"})
	w.AddEvent(greenKnobRel, rule{Data2: Exact(1), Key: "Down"})
	// white knob moves left and right also
	w.AddEvent(whiteKnobRel, rule{Data2: Exact(127), Key: "Left"})
	w.AddEvent(whiteKnobRel, rule{Data2: Exact(1), Key: "Right"})
	// blue and green knob click hits enter
	w.AddEvent(blueKnobClick, rule{Key: "Return"})
	w.AddEvent(greenKnobClick, rule{Key: "Return"})
	// help/popup key gives us the menu
	w.AddEvent(helpBtn, rule{Key: "Tab"})
	// tempo button hits backspace, to return through menus
	w.AddEvent(tempoBtn, rule{Key: "BackSpace"})
	// mic button brings up the system menu
	w.AddEvent(micBtn, rule{Key: "s"})

	// ff/rw
	w.AddEvent(orangeKnobRel, rule{Data2: Exact(127), Key: "r"})
	w.AddEvent(orangeKnobRel, rule{Data2: Exact(1), Key: "f"})

	w.AddEvent(f1, rule{Key: "Ctrl+b"})
	w.AddEvent(g1, rule{Key: "Ctrl+f"})
}
