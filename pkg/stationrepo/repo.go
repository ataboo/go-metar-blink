package stationrepo

import (
	"strconv"
	"time"

	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/geo"
	"github.com/ataboo/go-metar-blink/pkg/metarclient"
)

type StationRepo struct {
	client      metarclient.MetarClient
	stations    []*Station
	coordinates map[string]*geo.Coordinate
	lastUpdate  time.Time
}

type Station struct {
	ID               string
	FlightConditions string
	WindSpeedKts     float64
	Coordinate       *geo.Coordinate
}

type FetchHandler func(stations []*Station, err error)

func (r *StationRepo) Fetch(handler FetchHandler) {

}

func (r *StationRepo) FetchFresh(handler FetchHandler) {

}

func (r *StationRepo) loadCoordinatesIfEmpty(next func(err error)) {
	if r.coordinates != nil {
		return
	}

	r.client.GetStationPositions(func(positions []*metarclient.MetarPosition, err error) {
		if err != nil {
			next(err)
			return
		}

		r.coordinates = make(map[string]*geo.Coordinate)
		for _, pos := range positions {
			latFl, err := strconv.ParseFloat(pos.Latitude, 64)
			if err != nil {
				common.LogError("failed to parse station %s latitude: '%s'", pos.StationID, pos.Latitude)
			}

			longFl, err := strconv.ParseFloat(pos.Longitude, 64)
			if err != nil {
				common.LogError("failed to parse station %s longitude: '%s'", pos.StationID, pos.Longitude)
			}

			altFl, err := strconv.ParseFloat(pos.Altitude, 64)
			if err != nil {
				common.LogError("failed to parse station %s altitude: '%s'", pos.StationID, pos.Altitude)
			}

			r.coordinates[pos.StationID] = &geo.Coordinate{
				Latitude:  latFl,
				Longitude: longFl,
				Altitude:  altFl,
			}
		}

		next()
	})
}
