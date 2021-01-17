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

const AviationWeatherEndPoint = "https://aviationweather.gov/adds/dataserver_current/httpparam"

type aviationWeatherMetar struct {
	StationID       string  `xml:"station_id"`
	ObservationTime string  `xml:"observation_time"`
	WindSpeedKts    float32 `xml:"wind_speed_kt"`
	FlightCategory  string  `xml:"flight_category"`
	Latitude        string  `xml:"latitude"`
	Longitude       string  `xml:"longitude"`
	Altitude        string  `xml:"altitude"`
}

type aviationWeatherData struct {
	Metars []*aviationWeatherMetar `xml:"data>METAR"`
}

type aviationWeatherClient struct {
	settings *Settings
	endPoint string
}

type rawMetarHandler func([]*aviationWeatherMetar, error)

func newAviationWeatherClient(settings *Settings, endPoint string) MetarClient {
	settings.StationIDs = sort.StringSlice(settings.StationIDs)

	return &aviationWeatherClient{
		settings: settings,
		endPoint: endPoint,
	}
}

func (c *aviationWeatherClient) Fetch(handler MetarResponseHandler) error {
	endPoint, err := c.buildQueryURL(false)
	if err != nil {
		return err
	}

	go c.fetchRoutine(func(awm []*aviationWeatherMetar, e error) {
		if err != nil {
			handler(nil, e)
			return
		}

		reports := make([]*MetarReport, len(awm))
		for i, a := range awm {
			reports[i] = &MetarReport{
				StationID:       a.StationID,
				ObservationTime: a.ObservationTime,
				FlightRules:     a.FlightCategory,
				WindSpeedKts:    a.WindSpeedKts,
			}
		}

		handler(reports, nil)
	}, endPoint)

	return nil
}

func (c *aviationWeatherClient) GetStationPositions(handler MetarPositionResponseHandler) error {
	endPoint, err := c.buildQueryURL(true)
	if err != nil {
		return err
	}

	go c.fetchRoutine(func(awm []*aviationWeatherMetar, e error) {
		if err != nil {
			handler(nil, e)
			return
		}

		reports := make([]*MetarPosition, len(awm))
		for i, a := range awm {
			reports[i] = &MetarPosition{
				StationID: a.StationID,
				Latitude:  a.Latitude,
				Longitude: a.Longitude,
			}
		}

		handler(reports, nil)
	}, endPoint)

	return nil
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
			"altitude",
		}
	} else {
		fields = []string{
			"station_id",
			"observation_time",
			"wind_speed_kt",
			"flight_category",
		}
	}

	q := u.Query()
	q.Set("dataSource", "metars")
	q.Set("requestType", "retrieve")
	q.Set("format", "xml")
	q.Set("stationString", strings.Join(c.settings.StationIDs, ","))
	q.Set("hoursBeforeNow", "3")
	q.Set("mostRecentForEachStation", "constraint")
	q.Set("fields", strings.Join(fields, ","))
	u.RawQuery = q.Encode()

	return u, nil
}

func (c *aviationWeatherClient) fetchRoutine(handler rawMetarHandler, endPoint *url.URL) {
	response, err := http.Get(endPoint.String())

	if err != nil {
		handler(nil, err)
		return
	}

	handler(c.parseResponse(response))
}

func (c *aviationWeatherClient) parseResponse(response *http.Response) ([]*aviationWeatherMetar, error) {
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

	for _, stationID := range c.settings.StationIDs {
		matched := false
		for rawIdx < len(data.Metars) {
			rawMetar := data.Metars[rawIdx]

			if rawMetar.StationID > stationID {
				break
			}

			rawIdx++

			if rawMetar.StationID == stationID {
				matched = true
				break
			} else {
				common.LogWarn("received data for unnexpected station: %s'\n'", rawMetar.StationID)

			}
		}

		if !matched {
			common.LogWarn("failed to receive data for station %s'\n'", stationID)
		}
	}

	return data.Metars, nil
}
