package common

const (
	FlightRuleVFR                = "VFR"
	FlightRuleSVFR               = "SVFR"
	FlightRuleMVFR               = "MVFR"
	FlightRuleIFR                = "IFR"
	FlightRuleLIFR               = "LIFR"
	FlightRuleUnknown            = "Unknown"
	FlightRuleError              = "Error"
	AviationWeatherMetarStrategy = "AviationWeather"
)

type MapQuitError struct{}

func (e *MapQuitError) Error() string { return "this map is no longer running" }
