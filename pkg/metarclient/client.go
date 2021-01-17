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
	StationID       string
	ObservationTime string
	FlightRules     string
	WindSpeedKts    float32
}

type MetarPosition struct {
	StationID string
	Latitude  string
	Longitude string
	Altitude  string
}

type MetarClient interface {
	Fetch(handler MetarResponseHandler) error
	GetStationPositions(handler MetarPositionResponseHandler) error
}

func CreateMetarClient(settings *Settings) (MetarClient, error) {
	switch settings.Strategy {
	case AviationWeatherMetarStrategy:
		return newAviationWeatherClient(settings, AviationWeatherEndPoint), nil
	default:
		return nil, fmt.Errorf("configured metar strategy not supported")
	}
}
