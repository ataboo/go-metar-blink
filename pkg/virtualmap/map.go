package virtualmap

import (
	"errors"
	"path"

	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/geo"
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	PaddingPx     = 80
	ImageWidthPx  = 1600
	ImageHeightPx = 900
)

type MapQuitError struct{}

func (e *MapQuitError) Error() string { return "this virtual map is no longer running" }

type VirtualMap struct {
	stations         map[string]*stationrepo.Station
	renderedIDs      map[string]*sdl.Surface
	stationIDs       []string
	window           *sdl.Window
	windowSurface    *sdl.Surface
	stationScreenPos map[string]*sdl.Point
	running          bool
}

func CreateVirtualMap(stations map[string]*stationrepo.Station) (vMap *VirtualMap, err error) {
	if len(stations) == 0 {
		return nil, errors.New("need at least one station")
	}

	vMap = &VirtualMap{
		stations: stations,
		running:  true,
	}

	if err := ttf.Init(); err != nil {
		return nil, err
	}

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, err
	}

	wnd, err := sdl.CreateWindow(
		"Go Metar Blink",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		ImageWidthPx,
		ImageHeightPx,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return nil, err
	}
	vMap.window = wnd

	wndSurface, err := wnd.GetSurface()
	if err != nil {
		return nil, err
	}
	vMap.windowSurface = wndSurface

	if err := vMap.renderIds(); err != nil {
		return nil, err
	}

	coordinates := make([]*geo.Coordinate, len(stations))
	idx := 0
	for _, s := range stations {
		coordinates[idx] = s.Coordinate
		idx++
	}
	renderSpec := CreateRenderSpec(coordinates, ImageWidthPx, ImageHeightPx, PaddingPx)

	vMap.stationScreenPos = make(map[string]*sdl.Point)
	vMap.stationIDs = make([]string, len(stations))
	for _, s := range stations {
		x, y := renderSpec.ProjectCoordinate(s.Coordinate)
		vMap.stationScreenPos[s.ID] = &sdl.Point{
			X: int32(x),
			Y: int32(y),
		}
		vMap.stationIDs[s.Ordinal] = s.ID
	}

	return vMap, nil
}

func (m *VirtualMap) renderIds() error {
	font, err := ttf.OpenFont(path.Join(common.GetResourcesRoot(), "meslo_powerline.ttf"), 16)
	defer font.Close()

	m.renderedIDs = make(map[string]*sdl.Surface, len(m.stations))
	for _, s := range m.stations {
		m.renderedIDs[s.ID], err = font.RenderUTF8Solid(s.ID, sdl.Color{255, 255, 255, 255})
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *VirtualMap) Close() {
	for _, i := range m.renderedIDs {
		i.Free()
	}

	ttf.Quit()
	sdl.Quit()
	m.window.Destroy()
}

func (m *VirtualMap) Update() error {
	if !m.running {
		return &MapQuitError{}
	}

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			m.running = false
			return &MapQuitError{}
		}
	}

	m.windowSurface.FillRect(&sdl.Rect{0, 0, m.windowSurface.W, m.windowSurface.H}, 0xFF555555)

	for _, stationID := range m.stationIDs {
		station := m.stations[stationID]
		idSolid := m.renderedIDs[stationID]
		screenPos := m.stationScreenPos[stationID]

		m.windowSurface.FillRect(&sdl.Rect{
			X: screenPos.X - idSolid.W/2,
			Y: screenPos.Y - idSolid.H/2,
			W: idSolid.W,
			H: idSolid.H,
		}, station.Color.ARGB())

		err := idSolid.Blit(&idSolid.ClipRect, m.windowSurface, &sdl.Rect{
			X: screenPos.X - idSolid.W/2,
			Y: screenPos.Y - idSolid.H/2,
			W: idSolid.W,
			H: idSolid.H,
		})

		if err != nil {
			common.LogError(err.Error())
		}
	}

	return m.window.UpdateSurface()
}
