// +build arm

package engine

import (
	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/lightsmap"
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
)

func createMap(stations map[string]*stationrepo.Station) (MetarMap, error) {
	common.LogInfo("Building light map on arm")
	return lightsmap.CreateLightMap(stations)
}
