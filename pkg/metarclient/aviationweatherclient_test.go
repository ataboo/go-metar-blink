package metarclient

import (
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/ataboo/go-metar-blink/pkg/common"
)

func TestAviationWeatherUrlBuilding(t *testing.T) {
	aviationClient := newAviationWeatherClient(Settings{
		StationIDs: []string{"CYEG", "CYYC"},
	}).(*aviationWeatherClient)

	endPoint, err := aviationClient.buildQueryURL()
	if err != nil {
		t.Error(err)
	}

	stationString := endPoint.Query().Get("stationString")
	if stationString != "CYEG,CYYC" {
		t.Error("unnexpected stationString", stationString)
	}
}

func TestAviationWeatherParseResponse(t *testing.T) {
	aviationClient := newAviationWeatherClient(Settings{
		StationIDs: []string{"CYEG", "CYYC"},
	}).(*aviationWeatherClient)

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

	summaries, err := aviationClient.parseResponse(&response)
	if err != nil {
		t.Error("expected 404 error")
	}

	if len(summaries) != 2 {
		t.Error("unexpected summary count", len(summaries))
	}

	if summaries[0].FlightRules != common.FlightRuleVFR {
		t.Error("unnexpected flight rules")
	}

	if summaries[0].StationID != "CYEG" {
		t.Error("unnexpected station ID")
	}

	if summaries[0].WindSpeedKts != 8 {
		t.Error("unnexpected wind speed")
	}
}

func TestAviationWeatherParseResponseWrongStations(t *testing.T) {
	aviationClient := newAviationWeatherClient(Settings{
		StationIDs: []string{"CABC"},
	}).(*aviationWeatherClient)

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

	summaries, err := aviationClient.parseResponse(&response)
	if err != nil {
		t.Error("expected 404 error")
	}

	if len(summaries) != 1 {
		t.Error("unexpected summary count", len(summaries))
	}

	if summaries[0].FlightRules != common.FlightRuleError {
		t.Error("unnexpected flight rules")
	}

	if summaries[0].StationID != "CABC" {
		t.Error("unnexpected station ID")
	}

	if summaries[0].WindSpeedKts != 0 {
		t.Error("unnexpected wind speed")
	}
}

/*
Stations      []string
EndPoint      string
Strategy      MetarStrategy
WindySpeedKts float32

https://aviationweather.gov/adds/dataserver_current/httpparam
?dataSource=metars
&requestType=retrieve
&format=xml
&stationString=PHNL,KDE,CYYC,CYBR
&hoursBeforeNow=2
&mostRecentForEachStation=constraint
*/
