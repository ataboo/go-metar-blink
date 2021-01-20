package metaranimation

import (
	"time"

	"github.com/ataboo/go-metar-blink/pkg/animation"
	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
)

const (
	MetarAnimationFPS    = 50
	MinBlinkingWindSpeed = 5
	BasePeriodWindSpeed  = float64(40)
	MaxPeriodWindSpeed   = float64(80)
)

// var windPeriods = map[float64]int{
// 	10: 6 * MetarAnimationFPS,
// 	20: 4 * MetarAnimationFPS,
// 	30: 2 * MetarAnimationFPS,
// 	40: MetarAnimationFPS,
// 	50: MetarAnimationFPS / 2,
// }

func LoadingAnimation(channelCount int) animation.Animation {
	channels := make([]int, channelCount)
	for i := 0; i < channelCount; i++ {
		channels[i] = i
	}

	return animation.CreatePulseAnimation(time.Second*2, animation.ColorWhite, animation.ColorBlack, channels)
}

func ConditionsAnimation(stations map[string]*stationrepo.Station) animation.Animation {
	tracks := make([]*animation.Track, len(stations))
	for _, s := range stations {
		track, err := trackForConditions(s)
		if err != nil {
			common.LogError("failed to create animation track: %s", err)
			panic("aborting")
		}
		// TODO optimize group tracks by condition and windspeed range
		track.ChannelIDs = []int{s.Ordinal}
		tracks[s.Ordinal] = track
	}

	trackAnim := animation.CreateTrackAnimation(tracks, MetarAnimationFPS)

	return trackAnim
}

func trackForConditions(station *stationrepo.Station) (*animation.Track, error) {
	if station.FlightRules == common.FlightRuleError {
		return stationErrorTrack()
	}

	color := trackColorForFlightRules(station)

	if station.WindSpeedKts <= float64(MinBlinkingWindSpeed) {
		// TODO support single frame animation
		return animation.CreateTrack(2, false, []animation.KeyFrame{
			{0, color},
		})

	}

	frameCount := frameCountForWindSpeed(station.WindSpeedKts)
	return animation.CreateTrack(frameCount, true, []animation.KeyFrame{
		{5, color},
		{10, animation.ColorBlack},
		{15, animation.ColorBlack},
		{20, color},
		{frameCount - 1, color},
	})
}

// Two quick blinks in 1 second, off for 1 second
func stationErrorTrack() (*animation.Track, error) {
	t, err := animation.CreateTrack(100, true, []animation.KeyFrame{
		{0, animation.ColorBlack},
		{4, animation.ColorRed},
		{9, animation.ColorRed},
		{19, animation.ColorBlack},
		{24, animation.ColorBlack},
		{29, animation.ColorRed},
		{34, animation.ColorRed},
		{39, animation.ColorBlack},
		{99, animation.ColorBlack},
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}

func trackColorForFlightRules(station *stationrepo.Station) animation.Color {
	switch station.FlightRules {
	case common.FlightRuleIFR:
		return animation.ColorOrange
	case common.FlightRuleLIFR:
		return animation.ColorYellow
	case common.FlightRuleSVFR:
		return animation.ColorBlue
	case common.FlightRuleVFR:
		return animation.ColorGreen
	case common.FlightRuleUnknown:
		return animation.ColorRed
	case common.FlightRuleMVFR:
		return animation.ColorBlue
	default:
		common.LogWarn("flight rule '%s' has no color", station.FlightRules)
		return animation.ColorRed
	}
}

func frameCountForWindSpeed(windSpeedKts float64) int {
	if windSpeedKts > MaxPeriodWindSpeed {
		windSpeedKts = MaxPeriodWindSpeed
	}

	periodMultiplier := BasePeriodWindSpeed / windSpeedKts

	return int(MetarAnimationFPS * periodMultiplier)
}
