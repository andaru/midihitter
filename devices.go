package main

import (
	"fmt"
	"strings"

	"github.com/rakyll/portmidi"
)

func deviceID(idint int) portmidi.DeviceId {
	return portmidi.DeviceId(idint)
}

type devices struct {
	info map[portmidi.DeviceId]*portmidi.DeviceInfo
}

func newDevices() (devs devices) {
	portmidi.Initialize()
	devs = devices{
		info: map[portmidi.DeviceId]*portmidi.DeviceInfo{},
	}

	for id := 0; id < portmidi.CountDevices(); id++ {
		devID := deviceID(id)
		devs.info[devID] = portmidi.GetDeviceInfo(devID)
	}
	return
}

func (d devices) Input(name string) (id portmidi.DeviceId, err error) {
	var device *portmidi.DeviceInfo
	for id, device = range d.info {
		if device.Name == name && device.IsInputAvailable {
			return
		}
	}
	// Try a prefix match if we didn't get an exact match
	for id, device = range d.info {
		if strings.HasPrefix(device.Name, name) && device.IsInputAvailable {
			return
		}
	}
	err = fmt.Errorf("no input found for: %s", name)
	return
}

func (d devices) Output(name string) (id portmidi.DeviceId, err error) {
	var device *portmidi.DeviceInfo
	for id, device = range d.info {
		if device.Name == name && device.IsOutputAvailable {
			return
		}
	}
	// Try a prefix match if we didn't get an exact match
	for id, device = range d.info {
		if strings.HasPrefix(device.Name, name) && device.IsOutputAvailable {
			return
		}
	}
	err = fmt.Errorf("no output found for: %s", name)
	return
}

func (d devices) DeviceNames() (names []string) {
	namemap := map[string]struct{}{}
	for _, device := range d.info {
		namemap[device.Name] = struct{}{}
	}
	names = make([]string, len(namemap))
	i := 0
	for name, _ := range namemap {
		names[i] = name
		i++
	}
	return
}
