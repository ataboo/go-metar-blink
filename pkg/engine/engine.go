package engine

import (
	"sync"
	"time"

	"github.com/ataboo/go-metar-blink/pkg/animation"
	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/metaranimation"
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
	"github.com/ataboo/go-metar-blink/pkg/virtualmap"
)

type Engine struct {
	repo         *stationrepo.StationRepo
	stations     map[string]*stationrepo.Station
	frameTicker  *time.Ticker
	fetchTicker  *time.Ticker
	lastFrame    time.Time
	animation    animation.Animation
	quitChan     chan int
	vMap         *virtualmap.VirtualMap
	updatePeriod time.Duration
	fps          int
	lock         sync.Mutex
	colorMap     map[int]animation.Color
	doneSubs     []chan int
}

func CreateEngine(repo *stationrepo.StationRepo, settings *common.AppSettings) (*Engine, error) {
	stations, err := repo.LoadStations()
	if err != nil {
		return nil, err
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
	}

	vMap, err := virtualmap.CreateVirtualMap(stations)
	if err != nil {
		common.LogError("failed to init virtual map: %s", err)
		return nil, err
	}
	e.vMap = vMap

	return e, nil
}

func (e *Engine) Start() error {
	e.animation = metaranimation.LoadingAnimation(len(e.stations))
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

	err := e.vMap.Update()
	if err != nil {
		if _, ok := err.(*virtualmap.MapQuitError); ok {
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
	e.animation = metaranimation.ConditionsAnimation(e.stations)
	e.animation.Start()
	e.lock.Unlock()
}
