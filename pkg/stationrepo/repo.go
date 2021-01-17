package stationrepo

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/geo"
	"github.com/ataboo/go-metar-blink/pkg/metarclient"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

const (
	PositionCacheFileName = "station_positions.json"
)

type StationRepo struct {
	client      metarclient.MetarClient
	coordinates map[string]*geo.Coordinate
	reports     []*metarclient.MetarReport
	lastUpdate  time.Time
}

type Station struct {
	ID           string
	FlightRules  string
	WindSpeedKts float64
	Coordinate   *geo.Coordinate
}

func CreateStationRepo(client metarclient.MetarClient) *StationRepo {
	return &StationRepo{
		client: client,
	}
}

func (r *StationRepo) GetStations() (stations []*Station, err error) {
	err = r.loadCoordinatesIfEmpty()
	if err != nil {
		return nil, err
	}

	if r.reports == nil || r.lastUpdate.Add(time.Minute*time.Duration(common.GetAppSettings().UpdatePeriodMins)).Before(time.Now()) {
		common.LogInfo("repo fetching fresh reports")
		reports, err := r.client.GetReports()
		if err != nil {
			return nil, err
		}

		r.reports = reports
	}

	stations = make([]*Station, len(r.reports))
	for i, report := range r.reports {
		position, ok := r.coordinates[report.StationID]
		if !ok {
			return nil, fmt.Errorf("failed to get position matching station '%s'", report.StationID)
		}
		stations[i] = &Station{
			ID:           report.StationID,
			FlightRules:  report.FlightRules,
			WindSpeedKts: report.WindSpeedKts,
			Coordinate: &geo.Coordinate{
				Latitude:  position.Latitude,
				Longitude: position.Longitude,
				Altitude:  position.Altitude,
			},
		}
	}

	return stations, nil
}

func (r *StationRepo) loadCoordinatesIfEmpty() error {
	if r.coordinates != nil {
		return nil
	}

	if err := r.loadCoordinatesFromCache(); err == nil {
		return nil
	}

	err := r.loadCoordinatesFromClient()
	if err != nil {
		return err
	}

	bytes, err := json5.MarshalIndent(r.coordinates, "", "\t")
	if err != nil {
		common.LogWarn("failed to marshal coordinates map")
		return nil
	}

	err = common.CacheToFile(PositionCacheFileName, bytes)
	if err != nil {
		common.LogWarn("failed to save coordinates to cache")
	}

	return nil
}

func (r *StationRepo) loadCoordinatesFromClient() error {
	positions, err := r.client.GetStationPositions()
	if err != nil {
		return err
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

	return nil
}

func (r *StationRepo) loadCoordinatesFromCache() error {
	bytes, err := common.LoadCachedFile(PositionCacheFileName)
	if err != nil {
		return err
	}

	positionMap := map[string]*geo.Coordinate{}

	err = json5.Unmarshal(bytes, &positionMap)
	if err != nil {
		common.LogError("failed to unmarshal cached positions")
		return err
	}

	success := true
	for _, stationID := range common.GetAppSettings().StationIDs {
		if _, ok := positionMap[stationID]; !ok {
			success = false
			common.LogInfo("station '%s' not found in cached positions", stationID)
		}
	}

	if !success {
		return errors.New("failed to find a station in cached positions")
	}

	r.coordinates = positionMap

	return nil
}
