package main

// #cgo LDFLAGS: -lX11 -lXtst
// #include <X11/Xlib.h>
// #include <X11/keysym.h>
// #include <X11/keysymdef.h>
// #include <X11/extensions/XTest.h>
import "C"

import (
	"fmt"
	"os"

	"github.com/rakyll/portmidi"
	"strconv"
	"strings"
)

// The watcher considers portmidi events in terms of considering
// whether a particular MIDI event should cause a keyboard event to
// occur.

// MIDIMsgType is a MIDI status message code. For some reason it comes
// from the golang portmidi bindings in an event as an int64
type MIDIMsgType int64

const (
	MIDI_CLOCK       MIDIMsgType = 0xf8
	MIDI_ACTIVE      MIDIMsgType = 0xfe
	MIDI_STATUS_MASK MIDIMsgType = 0x80
	MIDI_SYSEX       MIDIMsgType = 0xf0
	MIDI_EOX         MIDIMsgType = 0xf7
	MIDI_START       MIDIMsgType = 0xFA
	MIDI_STOP        MIDIMsgType = 0xFC
	MIDI_CONTINUE    MIDIMsgType = 0xFB
	MIDI_F9          MIDIMsgType = 0xF9
	MIDI_FD          MIDIMsgType = 0xFD
	MIDI_RESET       MIDIMsgType = 0xFF
	MIDI_NOTE_ON     MIDIMsgType = 0x90
	MIDI_NOTE_OFF    MIDIMsgType = 0x80
	MIDI_CHANNEL_AT  MIDIMsgType = 0xD0
	MIDI_POLY_AT     MIDIMsgType = 0xA0
	MIDI_PROGRAM     MIDIMsgType = 0xC0
	MIDI_CONTROL     MIDIMsgType = 0xB0
	MIDI_PITCHBEND   MIDIMsgType = 0xE0
	MIDI_MTC         MIDIMsgType = 0xF1
	MIDI_SONGPOS     MIDIMsgType = 0xF2
	MIDI_SONGSEL     MIDIMsgType = 0xF3
	MIDI_TUNE        MIDIMsgType = 0xF6
)

func (mt MIDIMsgType) String() string {
	switch mt {
	case MIDI_ACTIVE:
		return "ACTIVE"
	case MIDI_CHANNEL_AT:
		return "CHANNEL_AT"
	case MIDI_CLOCK:
		return "CLOCK"
	case MIDI_CONTINUE:
		return "CONTINUE"
	case MIDI_CONTROL:
		return "CONTROL"
	case MIDI_EOX:
		return "EOX"
	case MIDI_F9:
		return "F9"
	case MIDI_FD:
		return "FD"
	case MIDI_MTC:
		return "MTC"
	case MIDI_NOTE_OFF:
		return "NOTE_OFF"
	case MIDI_NOTE_ON:
		return "NOTE_ON"
	case MIDI_PITCHBEND:
		return "PITCHBEND"
	case MIDI_POLY_AT:
		return "POLY_AT"
	case MIDI_PROGRAM:
		return "PROGRAM"
	case MIDI_RESET:
		return "RESET"
	case MIDI_SONGPOS:
		return "SONGPOS"
	case MIDI_SONGSEL:
		return "SONGSEL"
	case MIDI_START:
		return "START"
	case MIDI_STOP:
		return "STOP"
	case MIDI_SYSEX:
		return "SYSEX"
	case MIDI_TUNE:
		return "TUNE"
	default:
		return strconv.Itoa(int(mt))
	}
}

// ValueMatch is the value match condition for this rule in the
// watcher. It describes how the watched value should be changing in
// order to match the rule.
type ValueMatch int

const (
	ValueExact ValueMatch = iota
	ValueIncreasing
	ValueDecreasing
	ValuePassesIncreasing
	ValuePassesDecreasing
)

type key struct {
	Status MIDIMsgType
	Data1  int64
	Data2  int64
}

func Event(status MIDIMsgType, data ...int64) key {
	k := key{Status: status, Data1: -1, Data2: -1}
	if len(data) == 0 {
		return k
	} else if len(data) >= 1 {
		k.Data1 = data[0]
	}
	if len(data) == 2 {
		k.Data2 = data[1]
	}
	return k
}

func newKey(e *portmidi.Event) key {
	return key{
		Status: MIDIMsgType(e.Status),
		Data1:  e.Data1,
		Data2:  e.Data2,
	}
}

type match struct {
	Type    ValueMatch
	Value   int64
	Passing int64
}

func Exact(v int64) *match {
	return &match{
		Type:  ValueExact,
		Value: v,
	}
}

func Increasing() *match {
	return &match{
		Type:  ValueIncreasing,
		Value: -1,
	}
}

func Decreasing() *match {
	return &match{
		Type:  ValueDecreasing,
		Value: -1,
	}
}

func IncreasingPast(v int64) *match {
	return &match{
		Type:    ValueIncreasing,
		Passing: v,
		Value:   -1,
	}
}

func DecreasingPast(v int64) *match {
	return &match{
		Type:    ValueDecreasing,
		Passing: v,
		Value:   -1,
	}
}

type rule struct {
	Data1 *match
	Data2 *match
	Key   string // the X11 keyboard code to send as an action
}

type rules struct {
	d map[MIDIMsgType]map[int64]map[int64][]rule
	h map[MIDIMsgType]map[int64]map[int64]*portmidi.Event
}

func newRules() rules {
	return rules{
		d: map[MIDIMsgType]map[int64]map[int64][]rule{},
		h: map[MIDIMsgType]map[int64]map[int64]*portmidi.Event{},
	}
}

type watcher struct {
	midiport portmidi.DeviceId
	history  map[key]*portmidi.Event
	r        rules
}

func newWatcher(port portmidi.DeviceId) watcher {
	return watcher{
		midiport: port,
		history:  map[key]*portmidi.Event{},
		r:        newRules(),
	}
}

func (w *watcher) AddEvent(k key, r rule) {
	_, ok := w.r.d[k.Status]
	if !ok {
		w.r.d[k.Status] = map[int64]map[int64][]rule{}
	}
	_, ok = w.r.d[k.Status][k.Data1]
	if !ok {
		w.r.d[k.Status][k.Data1] = map[int64][]rule{}
	}
	_, ok = w.r.d[k.Status][k.Data1][k.Data2]
	if !ok {
		w.r.d[k.Status][k.Data1][k.Data2] = []rule{}
	}
	w.r.d[k.Status][k.Data1][k.Data2] = append(w.r.d[k.Status][k.Data1][k.Data2], r)
}

func (w *watcher) addHistory(k key, e *portmidi.Event) {
	_, ok := w.r.h[k.Status]
	if !ok {
		w.r.h[k.Status] = map[int64]map[int64]*portmidi.Event{}
	}
	_, ok = w.r.h[k.Status][k.Data1]
	if !ok {
		w.r.h[k.Status][k.Data1] = map[int64]*portmidi.Event{}
	}
	w.r.h[k.Status][k.Data1][k.Data2] = e
}

func (w *watcher) getRules(k key) (rs []rule, err error) {
	_, ok := w.r.d[k.Status]
	rs = []rule{}
	if !ok {
		err = fmt.Errorf("no rule for key: %#v", k)
		return
	}
	root, ok := w.r.d[k.Status][k.Data1]
	if !ok {
		root, ok = w.r.d[k.Status][-1]
		if !ok {
			err = fmt.Errorf("no rule for key: %#v", k)
			return
		}
	}
	rs, ok = root[k.Data2]
	if ok {
		return
	}

	// we may have a -1 entry in there, meaning
	// a wildcard
	rs, ok = root[-1]
	if ok {
		return
	}
	err = fmt.Errorf("no rule for key: %#v", k)
	return
}

func (w *watcher) getPrev(k key) (ev *portmidi.Event, err error) {
	_, ok := w.r.h[k.Status]
	err = fmt.Errorf("no rule for key: %#v", k)
	if !ok {
		return
	}
	root, ok := w.r.h[k.Status][k.Data1]
	if !ok {
		return
	}
	ev, ok = root[k.Data2]
	if ok {
		return ev, nil
	}

	// we may have a -1 entry in there, meaning
	// a wildcard
	ev, ok = root[-1]
	if ok {
		return ev, nil
	}
	return
}

func ck(mtype ValueMatch, cur, prev int64, m *match) bool {
	if mtype == ValueExact {
		return m.Value == cur
	} else if mtype == ValueIncreasing {
		return cur > prev
	} else if mtype == ValueDecreasing {
		return cur < prev
	} else if mtype == ValuePassesIncreasing {
		return cur > m.Passing && prev <= m.Passing
	} else if mtype == ValuePassesDecreasing {
		return cur < m.Passing && prev >= m.Passing
	}
	panic("should never get here")
}

var (
	x11display = C.CString(os.Getenv("DISPLAY"))
	dpy        = C.XOpenDisplay(x11display)
)

func sendKey(disp *C.struct__XDisplay, strname string) {
	if disp == nil {
		fmt.Println("no X11 display available")
		return
	}

	var mod _Ctype_KeyCode

	if strings.HasPrefix(strname, "Ctrl+") {
		strname = strname[len("Ctrl+"):]
		modkey := "Control_L"
		mod = C.XKeysymToKeycode(disp, C.XStringToKeysym(C.CString(modkey)))
		if debug {
			fmt.Printf("send key %T (Control) code=%v\n", mod, mod)
		}
		C.XTestFakeKeyEvent(disp, C.uint(mod), 1, 0)
		C.XFlush(disp)
	}
	chstr := C.CString(strname)
	keysym := C.XStringToKeysym(chstr)
	keycode := C.XKeysymToKeycode(disp, keysym)
	if debug {
		fmt.Printf("send key (%s) sym=%v code=%v\n", strname, keysym, keycode)
	}
	C.XTestFakeKeyEvent(disp, C.uint(keycode), 1, 0)
	C.XFlush(disp)
	C.XTestFakeKeyEvent(disp, C.uint(keycode), 0, 0)
	C.XFlush(disp)
	if debug {
		fmt.Printf("hit key %v (keysym %v keycode %v)\n", strname, keysym, keycode)
	}
	fmt.Printf("mod is == %v\n", mod)
	if mod != 0 {
		C.XTestFakeKeyEvent(disp, C.uint(mod), 0, 0)
		C.XFlush(disp)
	}
}

func eventMatches(cur, prev *portmidi.Event, r rule) bool {
	var match bool
	var pv int64 = -1
	if prev != nil {
		pv = prev.Data1
	}
	if r.Data1 != nil {
		match = ck(r.Data1.Type, cur.Data1, pv, r.Data1)
	} else {
		match = true
	}
	if !match {
		return false
	}
	if r.Data2 == nil {
		return true
	}
	if prev == nil && r.Data2.Type != ValueExact {
		return true
	}
	if prev == nil {
		return ck(r.Data2.Type, cur.Data2, -1, r.Data2)
	}
	return ck(r.Data2.Type, cur.Data2, prev.Data2, r.Data2)
}

func (w *watcher) HandleEvent(e *portmidi.Event) {
	var prev *portmidi.Event
	var matched bool
	var key = newKey(e)
	rules, err := w.getRules(key)
	if err != nil {
		goto addHistory
	}
	prev, err = w.getPrev(key)
	for _, rule := range rules {
		if eventMatches(e, prev, rule) {
			sendKey(dpy, rule.Key)
			matched = true
		}
	}
addHistory:
	if !matched && debug {
		fmt.Printf("No handler for MIDI %#v\n", e)
	}
	w.addHistory(key, e)
}

func (w *watcher) Run() (err error) {
	midiIn, err := portmidi.NewInputStream(w.midiport, 1024)
	if err != nil {
		return
	}
	events := midiIn.Listen()
	for e := range events {
		w.HandleEvent(&e)
	}
	return nil
}
