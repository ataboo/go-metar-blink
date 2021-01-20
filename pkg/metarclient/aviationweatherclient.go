package metarclient

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/ataboo/go-metar-blink/pkg/common"
)

var _ MetarClient = (*aviationWeatherClient)(nil)

const AviationWeatherEndPoint = "https://aviationweather.gov/adds/dataserver_current/httpparam"

//https://aviationweather.gov/docs/dataserver/schema/metar1_2.xsd
type aviationWeatherMetar struct {
	Error           bool
	StationID       string  `xml:"station_id"`
	ObservationTime string  `xml:"observation_time"`
	WindSpeedKts    float64 `xml:"wind_speed_kt"`
	FlightCategory  string  `xml:"flight_category"`
	Latitude        float64 `xml:"latitude"`
	Longitude       float64 `xml:"longitude"`
	Elevation       float64 `xml:"elevation_m"`
	VisibilitySM    float64 `xml:"visibility_statute_mi"`
	SkyCover        string  `xml:"sky_condition>sky_cover"`
	CloudBaseFtAGL  int     `xml:"sky_condition>cloud_base_ft_agl"`
	CloudType       string  `xml:"sky_condition>cloud_type"`
}

type aviationWeatherData struct {
	Metars []*aviationWeatherMetar `xml:"data>METAR"`
	Errors []string                `xml:"errors>error"`
}

type aviationWeatherClient struct {
	settings *Settings
	endPoint string
}

type rawMetarHandler func([]*aviationWeatherMetar, error)

func newAviationWeatherClient(settings *Settings, endPoint string) MetarClient {
	sort.Strings(settings.StationIDs)

	return &aviationWeatherClient{
		settings: settings,
		endPoint: endPoint,
	}
}

func (c *aviationWeatherClient) GetReports() (reports map[string]*MetarReport, err error) {
	endPoint, err := c.buildQueryURL(false)
	if err != nil {
		return nil, err
	}

	awm, err := c.getRawMetarData(endPoint)
	if err != nil {
		return nil, err
	}

	reports = make(map[string]*MetarReport, len(awm))
	for _, a := range awm {
		reports[a.StationID] = &MetarReport{
			Error:           a.Error,
			StationID:       a.StationID,
			ObservationTime: a.ObservationTime,
			FlightRules:     a.FlightCategory,
			WindSpeedKts:    a.WindSpeedKts,
			SkyCover:        a.SkyCover,
			CloudBaseFtAGL:  a.CloudBaseFtAGL,
			CloudType:       a.CloudType,
		}
	}

	return reports, nil
}

func (c *aviationWeatherClient) GetStationPositions() (positions map[string]*MetarPosition, err error) {
	endPoint, err := c.buildQueryURL(true)
	if err != nil {
		return nil, err
	}

	awm, err := c.getRawMetarData(endPoint)
	if err != nil {
		return nil, err
	}

	reports := make(map[string]*MetarPosition, len(awm))
	for id, a := range awm {
		reports[id] = &MetarPosition{
			Error:     a.Error,
			StationID: a.StationID,
			Latitude:  a.Latitude,
			Longitude: a.Longitude,
			Elevation: a.Elevation,
		}
	}

	return reports, nil
}

func (c *aviationWeatherClient) Fetch(handler MetarResponseHandler) {
	go func() {
		reports, err := c.GetReports()
		handler(reports, err)
	}()
}

func (c *aviationWeatherClient) FetchStationPositions(handler MetarPositionResponseHandler) {
	go func() {
		positions, err := c.GetStationPositions()
		handler(positions, err)
	}()
}

func (c *aviationWeatherClient) buildQueryURL(getPosition bool) (*url.URL, error) {
	u, err := url.Parse(c.endPoint)
	if err != nil {
		return nil, err
	}

	var fields []string

	if getPosition {
		fields = []string{
			"station_id",
			"latitude",
			"longitude",
			"elevation_m",
		}
	} else {
		fields = []string{
			"station_id",
			"observation_time",
			"wind_speed_kt",
			"flight_category",
			"visibility_statute_mi",
			"sky_cover",
			"cloud_base_ft_agl",
		}
	}

	q := u.Query()
	q.Set("dataSource", "metars")
	q.Set("requestType", "retrieve")
	q.Set("format", "xml")
	q.Set("stationString", strings.Join(c.settings.StationIDs, ","))
	q.Set("hoursBeforeNow", "6")
	q.Set("mostRecentForEachStation", "constraint")
	q.Set("fields", strings.Join(fields, ","))
	u.RawQuery = q.Encode()

	return u, nil
}

func (c *aviationWeatherClient) getRawMetarData(endPoint *url.URL) (map[string]*aviationWeatherMetar, error) {
	response, err := http.Get(endPoint.String())
	if err != nil {
		return nil, err
	}

	return c.parseResponse(response)
}

func (c *aviationWeatherClient) parseResponse(response *http.Response) (map[string]*aviationWeatherMetar, error) {
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response %d", response.StatusCode)
	}

	data := aviationWeatherData{}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	common.CacheToFile("last_aviation_weather_response.xml", responseBytes)

	err = xml.Unmarshal(responseBytes, &data)
	if err != nil {
		common.LogError("failed to parse aviation weather response")
		return nil, err
	}

	if len(data.Errors) > 0 {
		common.LogError("errors from aviation weather: %s", strings.Join(data.Errors, ", "))
		return nil, errors.New("received errors from aviation weather")
	}

	inputMap := make(map[string]*aviationWeatherMetar)
	for _, m := range data.Metars {
		inputMap[m.StationID] = m

		if m.FlightCategory == "" {
			common.LogWarn("station %s has no flight rule", m.StationID)
			m.FlightCategory = common.FlightRuleUnknown
		}
	}

	outputMap := make(map[string]*aviationWeatherMetar)

	for _, stationID := range c.settings.StationIDs {
		if m, ok := inputMap[stationID]; ok {
			outputMap[stationID] = m
		} else {
			outputMap[stationID] = &aviationWeatherMetar{
				Error:           true,
				StationID:       stationID,
				ObservationTime: "",
				FlightCategory:  common.FlightRuleError,
				WindSpeedKts:    0,
				Latitude:        0,
				Longitude:       0,
				Elevation:       0,
			}

			common.LogWarn("failed to receive data for station '%s'", stationID)
		}
	}

	return outputMap, nil
}
