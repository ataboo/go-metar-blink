package stationrepo

import (
	"errors"
	"fmt"
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
	config      *Config
}

type Station struct {
	ID           string
	FlightRules  string
	WindSpeedKts float64
	Coordinate   *geo.Coordinate
}

type Config struct {
	UpdatePeriod time.Duration
}

func CreateStationRepo(client metarclient.MetarClient, config *Config) *StationRepo {
	return &StationRepo{
		client: client,
		config: config,
	}
}

func (r *StationRepo) GetStations() (stations []*Station, err error) {
	err = r.loadCoordinatesIfEmpty()
	if err != nil {
		return nil, err
	}

	if r.reports == nil || r.lastUpdate.Add(r.config.UpdatePeriod).Before(time.Now()) {
		common.LogInfo("repo fetching fresh reports")
		reports, err := r.client.GetReports()
		if err != nil {
			return nil, err
		}

		r.reports = reports
	}

	stations = make([]*Station, 0, len(stations))
	for _, report := range r.reports {
		if report.Error {
			continue
		}

		position, ok := r.coordinates[report.StationID]
		if !ok {
			return nil, fmt.Errorf("failed to get position matching station '%s'", report.StationID)
		}
		stations = append(stations, &Station{
			ID:           report.StationID,
			FlightRules:  report.FlightRules,
			WindSpeedKts: report.WindSpeedKts,
			Coordinate: &geo.Coordinate{
				Latitude:  position.Latitude,
				Longitude: position.Longitude,
				Altitude:  position.Altitude,
			},
		})
	}

	return stations, nil
}

func (r *StationRepo) loadCoordinatesIfEmpty() error {
	if r.coordinates != nil {
		return nil
	}

	if err := r.loadCoordinatesFromCache(); err == nil {
		common.LogInfo("successfully loaded cached station coordinates")
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

	common.LogInfo("fetched and cached station coordinates")

	return nil
}

func (r *StationRepo) loadCoordinatesFromClient() error {
	positions, err := r.client.GetStationPositions()
	if err != nil {
		return err
	}

	r.coordinates = make(map[string]*geo.Coordinate)
	for _, pos := range positions {
		r.coordinates[pos.StationID] = &geo.Coordinate{
			Latitude:  pos.Latitude,
			Longitude: pos.Longitude,
			Altitude:  pos.Elevation,
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
