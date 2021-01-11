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

type MetarResponseHandler func(summaries []*MetarSummary, err error)

type MetarSummary struct {
	StationID    string
	FlightRules  string
	WindSpeedKts float32
}

type MetarClient interface {
	Fetch(handler MetarResponseHandler) error
}

func InitMetarClient(settings *Settings) (MetarClient, error) {
	switch settings.Strategy {
	case AviationWeatherMetarStrategy:
		return newAviationWeatherClient(settings, AviationWeatherEndPoint), nil
	default:
		return nil, fmt.Errorf("configured metar strategy not supported")
	}
}
