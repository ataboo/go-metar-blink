package metarclient

import "fmt"

const (
	AviationWeatherMetarStrategy = "AviationWeather"
)

type MetarStrategy string

type Settings struct {
	StationIDs []string
	Strategy   MetarStrategy
}

type MetarResponseHandler func(reports []*MetarReport, err error)
type MetarPositionResponseHandler func(positions []*MetarPosition, err error)

type MetarReport struct {
	Error           bool
	StationID       string
	ObservationTime string
	FlightRules     string
	WindSpeedKts    float64
}

type MetarPosition struct {
	Error     bool
	StationID string
	Latitude  float64
	Longitude float64
	Elevation float64
}

type MetarClient interface {
	GetReports() (reports []*MetarReport, err error)
	GetStationPositions() (positions []*MetarPosition, err error)
	Fetch(handler MetarResponseHandler)
	FetchStationPositions(handler MetarPositionResponseHandler)
}

func CreateMetarClient(settings *Settings) (MetarClient, error) {
	switch settings.Strategy {
	case AviationWeatherMetarStrategy:
		return newAviationWeatherClient(settings, AviationWeatherEndPoint), nil
	default:
		return nil, fmt.Errorf("configured metar strategy not supported")
	}
}
