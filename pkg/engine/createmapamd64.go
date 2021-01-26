// +build amd64

package engine

import (
	"github.com/ataboo/go-metar-blink/pkg/logger"
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
	"github.com/ataboo/go-metar-blink/pkg/virtualmap"
)

func createMap(stations map[string]*stationrepo.Station, brightness byte) (MetarMap, error) {
	logger.LogInfo("Building virtual map on AMD64")
	return virtualmap.CreateVirtualMap(stations, brightness)
}
