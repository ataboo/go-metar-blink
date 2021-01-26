// +build arm

package engine

import (
	"github.com/ataboo/go-metar-blink/pkg/lightsmap"
	"github.com/ataboo/go-metar-blink/pkg/logger"
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
)

func createMap(stations map[string]*stationrepo.Station, brightness byte) (MetarMap, error) {
	logger.LogInfo("Building light map on arm")
	return lightsmap.CreateLightMap(stations, brightness)
}
