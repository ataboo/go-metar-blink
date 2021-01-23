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

type ColorTheme struct {
	VFR   animation.Color `json:"vfr"`
	SVFR  animation.Color `json:"svfr"`
	IFR   animation.Color `json:"ifr"`
	LIFR  animation.Color `json:"lifr"`
	Error animation.Color `json:"error"`
}

type MetarAnimationFactory struct {
	theme *ColorTheme
}

func CreateMetarAnimationFactory(theme *ColorTheme) *MetarAnimationFactory {
	return &MetarAnimationFactory{
		theme: theme,
	}
}

func (f *MetarAnimationFactory) LoadingAnimation(channelCount int) animation.Animation {
	channels := make([]int, channelCount)
	for i := 0; i < channelCount; i++ {
		channels[i] = i
	}

	return animation.CreatePulseAnimation(time.Second*2, animation.ColorWhite, animation.ColorBlack, channels)
}

func (f *MetarAnimationFactory) ConditionsAnimation(stations map[string]*stationrepo.Station) animation.Animation {
	tracks := make([]*animation.Track, len(stations))
	for _, s := range stations {
		track, err := f.trackForConditions(s)
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

func (f *MetarAnimationFactory) trackForConditions(station *stationrepo.Station) (*animation.Track, error) {
	if station.FlightRules == common.FlightRuleError {
		return f.stationErrorTrack()
	}

	color := f.trackColorForFlightRules(station)

	if station.WindSpeedKts <= float64(MinBlinkingWindSpeed) {
		// TODO support single frame animation
		return animation.CreateTrack(2, false, []animation.KeyFrame{
			{0, color},
		})

	}

	frameCount := f.frameCountForWindSpeed(station.WindSpeedKts)
	return animation.CreateTrack(frameCount, true, []animation.KeyFrame{
		{5, color},
		{10, animation.ColorBlack},
		{15, animation.ColorBlack},
		{20, color},
		{frameCount - 1, color},
	})
}

// Two quick blinks in 1 second, off for 1 second
func (f *MetarAnimationFactory) stationErrorTrack() (*animation.Track, error) {
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

func (f *MetarAnimationFactory) trackColorForFlightRules(station *stationrepo.Station) animation.Color {
	switch station.FlightRules {
	case common.FlightRuleIFR:
		return f.theme.IFR
	case common.FlightRuleLIFR:
		return f.theme.LIFR
	case common.FlightRuleSVFR:
		return f.theme.SVFR
	case common.FlightRuleVFR:
		return f.theme.VFR
	case common.FlightRuleUnknown:
		return f.theme.Error
	case common.FlightRuleMVFR:
		return f.theme.SVFR
	default:
		common.LogWarn("flight rule '%s' has no color", station.FlightRules)
		return animation.ColorRed
	}
}

func (f *MetarAnimationFactory) frameCountForWindSpeed(windSpeedKts float64) int {
	if windSpeedKts > MaxPeriodWindSpeed {
		windSpeedKts = MaxPeriodWindSpeed
	}

	periodMultiplier := BasePeriodWindSpeed / windSpeedKts

	return int(MetarAnimationFPS * periodMultiplier)
}
