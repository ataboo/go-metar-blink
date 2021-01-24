package metarclient

import (
	"fmt"

	"github.com/ataboo/go-metar-blink/pkg/common"
)

type MetarStrategy string

type Settings struct {
	StationIDs []string
	Strategy   MetarStrategy
}

type MetarResponseHandler func(reports map[string]*MetarReport, err error)
type MetarPositionResponseHandler func(positions map[string]*MetarPosition, err error)

type MetarReport struct {
	Error           bool
	StationID       string
	ObservationTime string
	FlightRules     string
	WindSpeedKts    float64
	SkyCover        string
	CloudBaseFtAGL  int
	CloudType       string
}

type MetarPosition struct {
	Error     bool
	StationID string
	Latitude  float64
	Longitude float64
	Elevation float64
}

type MetarClient interface {
	GetReports() (reports map[string]*MetarReport, err error)
	GetStationPositions() (positions map[string]*MetarPosition, err error)
	Fetch(handler MetarResponseHandler)
	FetchStationPositions(handler MetarPositionResponseHandler)
}

func CreateMetarClient(settings *Settings) (MetarClient, error) {
	switch settings.Strategy {
	case common.AviationWeatherMetarStrategy:
		return newAviationWeatherClient(settings, AviationWeatherEndPoint), nil
	default:
		return nil, fmt.Errorf("configured metar strategy not supported")
	}
}
