package metarclient

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/ataboo/go-metar-blink/pkg/common"
)

var _ MetarClient = (*aviationWeatherClient)(nil)

type aviationWeatherMetar struct {
	StationID       string  `xml:"station_id"`
	ObservationTime string  `xml:"observation_time"`
	WindSpeedKt     float32 `xml:"wind_speed_kt"`
	FlightCategory  string  `xml:"flight_category"`
}

type aviationWeatherData struct {
	Metars []*aviationWeatherMetar `xml:"data>METAR"`
}

type aviationWeatherClient struct {
	settings Settings
}

func newAviationWeatherClient(settings Settings) MetarClient {
	settings.StationIDs = sort.StringSlice(settings.StationIDs)

	return &aviationWeatherClient{
		settings: settings,
	}
}

func (c *aviationWeatherClient) Fetch(handler MetarResponseHandler) error {
	endPoint, err := c.buildQueryURL()
	if err != nil {
		return err
	}

	go c.fetchRoutine(handler, endPoint)

	return nil
}

func (c *aviationWeatherClient) buildQueryURL() (*url.URL, error) {
	u, err := url.Parse("https://aviationweather.gov/adds/dataserver_current/httpparam")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("dataSource", "metars")
	q.Set("requestType", "retrieve")
	q.Set("format", "xml")
	q.Set("stationString", strings.Join(c.settings.StationIDs, ","))
	q.Set("hoursBeforeNow", "3")
	q.Set("mostRecentForEachStation", "constraint")

	u.RawQuery = q.Encode()

	return u, nil
}

func (c *aviationWeatherClient) fetchRoutine(handler MetarResponseHandler, endPoint *url.URL) {
	response, err := http.Get(endPoint.String())

	if err != nil {
		handler(nil, err)
		return
	}

	handler(c.parseResponse(response))
}

func (c *aviationWeatherClient) parseResponse(response *http.Response) ([]*MetarSummary, error) {
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response %d", response.StatusCode)
	}

	data := aviationWeatherData{}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(responseBytes, &data)
	if err != nil {
		return nil, err
	}

	sort.Slice(data.Metars, func(i int, j int) bool {
		return data.Metars[i].StationID < data.Metars[j].StationID
	})

	rawIdx := 0

	summaries := make([]*MetarSummary, len(c.settings.StationIDs))
	for i, stationID := range c.settings.StationIDs {
		summaries[i] = &MetarSummary{
			StationID: stationID,
		}

		var matchingMetar *aviationWeatherMetar
		for rawIdx < len(data.Metars) {
			rawMetar := data.Metars[rawIdx]

			if rawMetar.StationID > stationID {
				break
			}

			rawIdx++

			if rawMetar.StationID == stationID {
				matchingMetar = rawMetar
				break
			} else {
				common.LogWarn("received data for unnexpected station: %s'\n'", rawMetar.StationID)
			}
		}

		if matchingMetar == nil {
			common.LogWarn("failed to receive data for station %s'\n'", stationID)

			matchingMetar = &aviationWeatherMetar{
				StationID:      stationID,
				WindSpeedKt:    0,
				FlightCategory: common.FlightRuleError,
			}
		}

		summaries[i] = &MetarSummary{
			StationID:    stationID,
			FlightRules:  matchingMetar.FlightCategory,
			WindSpeedKts: matchingMetar.WindSpeedKt,
		}
	}

	return summaries, nil
}
