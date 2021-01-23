package engine

import (
	"sync"
	"time"

	"github.com/ataboo/go-metar-blink/pkg/animation"
	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/metaranimation"
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
)

type MetarMap interface {
	Update() error
	Dispose()
}

type Engine struct {
	repo         *stationrepo.StationRepo
	stations     map[string]*stationrepo.Station
	frameTicker  *time.Ticker
	fetchTicker  *time.Ticker
	lastFrame    time.Time
	animation    animation.Animation
	quitChan     chan int
	metarMap     MetarMap
	updatePeriod time.Duration
	fps          int
	lock         sync.Mutex
	colorMap     map[int]animation.Color
	doneSubs     []chan int
	animFactory  *metaranimation.MetarAnimationFactory
}

func CreateEngine(repo *stationrepo.StationRepo, settings *common.AppSettings) (*Engine, error) {
	stations, err := repo.LoadStations()
	if err != nil {
		return nil, err
	}

	theme := metaranimation.ColorTheme{
		Error: settings.Colors.ParseColorHexString(settings.Colors.Error),
		IFR:   settings.Colors.ParseColorHexString(settings.Colors.IFR),
		LIFR:  settings.Colors.ParseColorHexString(settings.Colors.LIFR),
		VFR:   settings.Colors.ParseColorHexString(settings.Colors.VFR),
		SVFR:  settings.Colors.ParseColorHexString(settings.Colors.SVFR),
	}

	e := &Engine{
		repo:         repo,
		stations:     stations,
		quitChan:     make(chan int),
		updatePeriod: time.Duration(settings.UpdatePeriodMins) * time.Minute,
		fps:          50,
		lock:         sync.Mutex{},
		colorMap:     make(map[int]animation.Color),
		doneSubs:     make([]chan int, 0),
		animFactory:  metaranimation.CreateMetarAnimationFactory(&theme),
	}

	mMap, err := createMap(stations)
	if err != nil {
		common.LogError("failed to init map: %s", err)
		return nil, err
	}

	e.metarMap = mMap

	return e, nil
}

func (e *Engine) Start() error {
	common.LogInfo("started loading animation")
	e.animation = e.animFactory.LoadingAnimation(len(e.stations))
	e.animation.Start()

	e.fetchTicker = time.NewTicker(e.updatePeriod)
	e.frameTicker = time.NewTicker(time.Second / time.Duration(e.fps))

	go e.fetchRoutine()
	go e.mainLoop()

	return nil
}

func (e *Engine) DoneSubscribe() chan int {
	newChan := make(chan int)
	e.doneSubs = append(e.doneSubs, newChan)

	return newChan
}

func (e *Engine) Dispose() {
	e.metarMap.Dispose()
}

func (e *Engine) mainLoop() {
	running := true
	for running {
		select {
		case currentTime := <-e.frameTicker.C:
			e.lock.Lock()
			running = e.updateFrame(currentTime)
			e.lock.Unlock()

		case <-e.fetchTicker.C:
			go e.fetchRoutine()
		case <-e.quitChan:
			running = false
			break
		}
	}

	for _, c := range e.doneSubs {
		doneChan := c
		go func() {
			select {
			case doneChan <- 0:
				break
			case <-time.After(time.Second):
				common.LogWarn("timed out sending done to sub")
				break
			}
		}()
	}
}

func (e *Engine) updateFrame(currentTime time.Time) bool {
	e.animation.Step(e.colorMap)

	for _, s := range e.stations {
		s.Color = e.colorMap[s.Ordinal]
	}
	e.lastFrame = currentTime

	err := e.metarMap.Update()
	if err != nil {
		if _, ok := err.(*common.MapQuitError); ok {
			common.LogInfo("map has quit")
			return false
		}

		common.LogError("failed to update vmap %s", err)
	}

	return true
}

func (e *Engine) fetchRoutine() {
	e.repo.UpdateReports(e.stations)
	e.lock.Lock()
	e.animation = e.animFactory.ConditionsAnimation(e.stations)
	e.animation.Start()
	common.LogInfo("updated conditions animation")
	e.lock.Unlock()
}
