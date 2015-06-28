package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rakyll/portmidi"
)

var (
	debug bool
)

func main() {
	var (
		flagPort    = flag.String("port", "", "MIDI port name (or prefix) to watch")
		flagList    = flag.Bool("l", false, "List the found input ports")
		flagVerbose = flag.Bool("v", false, "Display unmapped MIDI events and be more verbose")

		err  error
		port portmidi.DeviceId
	)

	flag.Parse()
	if *flagVerbose {
		debug = true
	}
	devs := newDevices()

	if *flagList {
		for _, name := range devs.DeviceNames() {
			fmt.Println(name)
		}
		return
	}

	if *flagPort != "" {
		port, err = devs.Input(*flagPort)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}
	} else {
		port = portmidi.GetDefaultInputDeviceId()
	}

	if debug {
		fmt.Printf("monitoring MIDI port: %s\n", devs.info[port].Name)
	}
	watcher := newWatcher(port)
	watcher.Run()
}
