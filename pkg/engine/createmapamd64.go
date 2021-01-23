// +build amd64

package engine

import (
	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
	"github.com/ataboo/go-metar-blink/pkg/virtualmap"
)

func createMap(stations map[string]*stationrepo.Station) (MetarMap, error) {
	common.LogInfo("Building virtual map on AMD64")
	return virtualmap.CreateVirtualMap(stations)
}
