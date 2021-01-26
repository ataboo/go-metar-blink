package stationrepo

import (
	"errors"
	"fmt"

	"github.com/ataboo/go-metar-blink/pkg/animation"
	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/geo"
	"github.com/ataboo/go-metar-blink/pkg/logger"
	"github.com/ataboo/go-metar-blink/pkg/metarclient"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

const (
	PositionCacheFileName = "station_positions.json"
)

type StationRepo struct {
	client      metarclient.MetarClient
	coordinates map[string]*geo.Coordinate
	config      *Config
}

type Station struct {
	ID           string
	Ordinal      int
	FlightRules  string
	WindSpeedKts float64
	Coordinate   *geo.Coordinate
	Color        animation.Color
}

type Config struct {
	StationIDs []string
}

func CreateStationRepo(client metarclient.MetarClient, config *Config) *StationRepo {
	return &StationRepo{
		client: client,
		config: config,
	}
}

func (r *StationRepo) LoadStations() (stations map[string]*Station, err error) {
	err = r.loadCoordinatesIfEmpty()
	if err != nil {
		return nil, err
	}

	stations = make(map[string]*Station, 0)
	idx := 0
	for _, id := range r.config.StationIDs {
		position, ok := r.coordinates[id]
		if !ok {
			return nil, fmt.Errorf("failed to get position matching station '%s'", id)
		}
		stations[id] = &Station{
			ID:           id,
			Ordinal:      idx,
			FlightRules:  common.FlightRuleError,
			WindSpeedKts: 0,
			Coordinate: &geo.Coordinate{
				Latitude:  position.Latitude,
				Longitude: position.Longitude,
				Altitude:  position.Altitude,
			},
		}
		idx++
	}

	return stations, nil
}

func (r *StationRepo) UpdateReports(stations map[string]*Station) error {
	logger.LogDebug("repo fetching fresh reports")
	reports, err := r.client.GetReports()
	if err != nil {
		return err
	}

	for _, s := range stations {
		r, ok := reports[s.ID]
		if !ok {
			r = &metarclient.MetarReport{
				Error:           true,
				StationID:       s.ID,
				ObservationTime: "",
				FlightRules:     common.FlightRuleError,
				WindSpeedKts:    0,
			}
		}

		s.FlightRules = r.FlightRules
		s.WindSpeedKts = r.WindSpeedKts
	}

	return nil
}

func (r *StationRepo) loadCoordinatesIfEmpty() error {
	if r.coordinates != nil {
		return nil
	}

	if err := r.loadCoordinatesFromCache(); err == nil {
		logger.LogInfo("successfully loaded cached station coordinates")
		return nil
	}

	err := r.loadCoordinatesFromClient()
	if err != nil {
		return err
	}

	bytes, err := json5.MarshalIndent(r.coordinates, "", "\t")
	if err != nil {
		logger.LogWarn("failed to marshal coordinates map")
		return nil
	}

	err = common.CacheToFile(PositionCacheFileName, bytes)
	if err != nil {
		logger.LogWarn("failed to save coordinates to cache")

	}

	logger.LogInfo("fetched and cached station coordinates")

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
		logger.LogError("failed to unmarshal cached positions")
		return err
	}

	success := true
	for _, stationID := range r.config.StationIDs {
		if _, ok := positionMap[stationID]; !ok {
			success = false
			logger.LogInfo("station '%s' not found in cached positions", stationID)
		}
	}

	if !success {
		return errors.New("failed to find a station in cached positions")
	}

	r.coordinates = positionMap

	return nil
}
