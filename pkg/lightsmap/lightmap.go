// +build arm amd64

package lightsmap

import (
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

type LightMap struct {
	stations map[string]*stationrepo.Station
	device   *ws2811.WS2811
}

func CreateLightMap(stations map[string]*stationrepo.Station, brightness byte) (lMap *LightMap, err error) {
	lMap = &LightMap{
		stations: stations,
	}

	options := ws2811.DefaultOptions
	options.Channels[0].Brightness = int(brightness)
	options.Channels[0].LedCount = len(stations)

	dev, err := ws2811.MakeWS2811(&options)
	if err != nil {
		return nil, err
	}

	lMap.device = dev

	if err := dev.Init(); err != nil {
		return nil, err
	}

	return lMap, nil
}

func (l *LightMap) Update() error {

	for _, s := range l.stations {
		l.device.Leds(0)[s.Ordinal] = s.Color.RGB()
	}

	return l.device.Render()
}

func (l *LightMap) Dispose() {
	l.device.Fini()
}
