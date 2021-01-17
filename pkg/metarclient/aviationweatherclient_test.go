package metarclient

import (
	"context"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/ataboo/go-metar-blink/pkg/common"
)

func TestAviationWeatherUrlBuilding(t *testing.T) {
	client, err := CreateMetarClient(&Settings{
		StationIDs: []string{"CYEG", "CYYC"},
		Strategy:   AviationWeatherMetarStrategy,
	})
	if err != nil {
		t.Error(err)
	}

	aviationClient := client.(*aviationWeatherClient)

	endPoint, err := aviationClient.buildQueryURL(false)
	if err != nil {
		t.Error(err)
	}

	stationString := endPoint.Query().Get("stationString")
	if stationString != "CYEG,CYYC" {
		t.Error("unnexpected stationString", stationString)
	}
}

func TestAviationWeatherParseResponse(t *testing.T) {
	client, err := CreateMetarClient(&Settings{
		StationIDs: []string{"CYEG", "CYYC"},
		Strategy:   AviationWeatherMetarStrategy,
	})
	if err != nil {
		t.Error(err)
	}

	aviationClient := client.(*aviationWeatherClient)

	exampleRaw, err := ioutil.ReadFile(path.Join(common.GetProjectRoot(), "resources/aviation-weather-example.xml"))
	if err != nil {
		t.Error(err)
	}

	response := http.Response{
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(strings.NewReader(string(exampleRaw))),
	}

	_, err = aviationClient.parseResponse(&response)
	if err == nil {
		t.Error("expected 404 error")
	}

	response.StatusCode = http.StatusOK

	reports, err := aviationClient.parseResponse(&response)
	if err != nil {
		t.Error("expected 404 error")
	}

	if len(reports) != 2 {
		t.Error("unexpected report count", len(reports))
	}

	if reports[0].FlightCategory != common.FlightRuleVFR {
		t.Error("unnexpected flight rules")
	}

	if reports[0].StationID != "CYEG" {
		t.Error("unnexpected station ID")
	}

	if reports[0].WindSpeedKts != 8 {
		t.Error("unnexpected wind speed")
	}
}

func TestAviationWeatherParseResponseWrongStations(t *testing.T) {
	_ = common.InitLoggersToTestWriter()

	client, err := CreateMetarClient(&Settings{
		StationIDs: []string{"CABC"},
		Strategy:   AviationWeatherMetarStrategy,
	})
	aviationClient := client.(*aviationWeatherClient)

	exampleRaw, err := ioutil.ReadFile(path.Join(common.GetProjectRoot(), "resources/aviation-weather-example.xml"))
	if err != nil {
		t.Error(err)
	}

	response := http.Response{
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(strings.NewReader(string(exampleRaw))),
	}

	_, err = aviationClient.parseResponse(&response)
	if err == nil {
		t.Error("expected 404 error")
	}

	response.StatusCode = http.StatusOK

	reports, err := aviationClient.parseResponse(&response)
	if err != nil {
		t.Error("expected 404 error")
	}

	if len(reports) != 1 {
		t.Error("unexpected report count", len(reports))
	}

	if reports[0].FlightCategory != common.FlightRuleError {
		t.Error("unnexpected flight rules")
	}

	if reports[0].StationID != "CABC" {
		t.Error("unnexpected station ID")
	}

	if reports[0].WindSpeedKts != 0 {
		t.Error("unnexpected wind speed")
	}
}

func TestClientFetchIntegrated(t *testing.T) {
	_ = common.InitLoggersToTestWriter()

	client := newAviationWeatherClient(&Settings{
		StationIDs: []string{"CYEG", "CYYC"},
		Strategy:   AviationWeatherMetarStrategy,
	}, "http://localhost:3000/go-metar-blink").(*aviationWeatherClient)

	server := http.Server{
		Addr:    ":3000",
		Handler: nil,
	}

	exampleRaw, err := ioutil.ReadFile(path.Join(common.GetProjectRoot(), "resources/aviation-weather-example.xml"))
	if err != nil {
		t.Error(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		expectedQuery, _ := client.buildQueryURL(false)
		if expectedQuery.RawQuery != r.URL.RawQuery {
			t.Error("Query mismatch: ", expectedQuery.RawQuery, r.URL.RawQuery)
		}

		w.Write(exampleRaw)
	})

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	doneServerChan := make(chan int, 0)

	client.Fetch(func(reports []*MetarReport, err error) {
		if err != nil {
			t.Error(err)
		}
		t.Logf("Got %d reports", len(reports))

		if len(reports) != 2 {
			t.Error("received unnexpected report count")
		}

		if reports[0].StationID != "CYEG" {
			t.Error("unnexpected first station id")
		}

		if reports[1].StationID != "CYYC" {
			t.Error("unnexpected second station id")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		err = server.Shutdown(ctx)
		if err != nil {
			t.Error(err)
		}

		doneServerChan <- 1
	})

	select {
	case <-doneServerChan:
		//done
	}
}
