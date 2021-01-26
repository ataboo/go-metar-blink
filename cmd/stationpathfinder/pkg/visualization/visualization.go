package visualization

import (
	"errors"
	"fmt"
	"math"
	"path"

	"github.com/ataboo/go-metar-blink/cmd/stationpathfinder/pkg/wirepath"
	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/logger"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type MapQuitError struct{}

func (e *MapQuitError) Error() string { return "SDL Map has quit" }

type TextTexture struct {
	Texture *sdl.Texture
	Surface *sdl.Surface
}

func (t TextTexture) Dispose() {
	t.Texture.Destroy()
	t.Surface.Free()
}

type PathFindingVisualization struct {
	screenPositions map[string]*sdl.Point
	running         bool
	window          *sdl.Window
	windowSurface   *sdl.Surface
	renderedIDs     map[string]*TextTexture
	renderer        *sdl.Renderer
	font            *ttf.Font
}

const (
	PaddingPx     = 80
	ImageWidthPx  = 1920
	ImageHeightPx = 1080
)

func (v PathFindingVisualization) Update(pathfinder wirepath.PathFinder, bestPath []string, bestScore float64) error {
	if !v.running {
		return &MapQuitError{}
	}

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			v.running = false
			return &MapQuitError{}
		}
	}

	path := pathfinder.GetBestPath()
	stats := pathfinder.Stats()

	v.renderer.Clear()
	v.renderer.SetDrawColor(0x55, 0x55, 0x55, 0xFF)
	v.renderer.FillRect(&sdl.Rect{0, 0, v.windowSurface.W, v.windowSurface.H})

	statsText, err := v.renderStats(stats, bestScore)
	if err != nil {
		return err
	}
	defer statsText.Dispose()

	statsRect := statsText.Surface.ClipRect
	v.renderer.Copy(statsText.Texture, &statsRect, &sdl.Rect{X: PaddingPx, Y: PaddingPx, W: statsRect.W, H: statsRect.H})

	if bestPath != nil {
		lastBestID := bestPath[len(bestPath)-1]
		for _, id := range bestPath {
			thisPos := v.screenPositions[id]
			lastPos := v.screenPositions[lastBestID]
			v.renderer.SetDrawColor(0, 255, 0, 255)

			dy := thisPos.Y - lastPos.Y
			dx := thisPos.X - lastPos.X

			var xOffset int32
			var yOffset int32
			if math.Abs(float64(dy)) > math.Abs(float64(dx)) {
				xOffset = 4
				yOffset = 0
			} else {
				yOffset = 4
				xOffset = 0
			}

			v.renderer.DrawLine(thisPos.X+xOffset, thisPos.Y+yOffset, lastPos.X+xOffset, lastPos.Y+yOffset)
			v.renderer.DrawLine(thisPos.X-xOffset, thisPos.Y-yOffset, lastPos.X-xOffset, lastPos.Y-yOffset)
			lastBestID = id
		}
	}

	var lastID string = path[len(path)-1]
	for _, id := range path {
		thisPos := v.screenPositions[id]
		idText := v.renderedIDs[id]

		lastPos := v.screenPositions[lastID]
		v.renderer.SetDrawColor(255, 0, 0, 255)
		v.renderer.DrawLine(thisPos.X, thisPos.Y, lastPos.X, lastPos.Y)
		lastID = id

		v.renderer.SetDrawColor(255, 255, 255, 255)
		v.renderer.FillRect(&sdl.Rect{
			X: thisPos.X - 1,
			Y: thisPos.Y - 1,
			W: 3,
			H: 3,
		})

		err := v.renderer.Copy(idText.Texture, &idText.Surface.ClipRect, &sdl.Rect{
			X: thisPos.X - idText.Surface.ClipRect.W/2,
			Y: thisPos.Y - idText.Surface.ClipRect.W/2 - 16,
			W: idText.Surface.W,
			H: idText.Surface.H,
		})
		if err != nil {
			logger.LogError(err.Error())
		}
	}

	v.renderer.Present()

	return nil
}

func CreatePathFindingVisualization(pathFinder wirepath.PathFinder) (*PathFindingVisualization, error) {
	positions := pathFinder.GetPositions()

	if len(positions) == 0 {
		return nil, errors.New("need at least one station position")
	}

	screenPosMap := make(map[string]*sdl.Point)
	for _, p := range positions {
		screenPosMap[p.Name] = &sdl.Point{
			X: int32(p.X),
			Y: int32(p.Y),
		}
	}

	vMap := &PathFindingVisualization{
		screenPositions: screenPosMap,
		running:         true,
	}

	if err := ttf.Init(); err != nil {
		return nil, err
	}
	font, err := ttf.OpenFont(path.Join(common.GetResourcesRoot(), "dev", "meslo_powerline.ttf"), 16)
	if err != nil {
		return nil, err
	}
	vMap.font = font

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, err
	}

	wnd, err := sdl.CreateWindow(
		"Go Metar Blink - Pathfinder",
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

	renderer, err := sdl.CreateRenderer(wnd, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}
	vMap.renderer = renderer

	wndSurface, err := wnd.GetSurface()
	if err != nil {
		return nil, err
	}
	vMap.windowSurface = wndSurface

	if err := vMap.renderIds(); err != nil {
		return nil, err
	}

	return vMap, nil
}

func (m *PathFindingVisualization) renderIds() error {
	m.renderedIDs = make(map[string]*TextTexture, len(m.screenPositions))
	for id := range m.screenPositions {
		surface, err := m.font.RenderUTF8Solid(id, sdl.Color{255, 255, 255, 255})
		if err != nil {
			return err
		}

		texture, err := m.renderer.CreateTextureFromSurface(surface)
		if err != nil {
			return err
		}

		m.renderedIDs[id] = &TextTexture{
			Texture: texture,
			Surface: surface,
		}
	}

	return nil
}

func (v *PathFindingVisualization) renderStats(stats *wirepath.PathfinderStats, bestScore float64) (*TextTexture, error) {

	bestScore = math.Min(bestScore, stats.ShortestPath)

	statsStr := fmt.Sprintf("Shortest:  %f\nCurrent:   %f\nRuntime:   %s\nGenerated: %d", bestScore, stats.ShortestPath, stats.RunTime, stats.PathsGenerated)

	surface, err := v.font.RenderUTF8BlendedWrapped(statsStr, sdl.Color{255, 255, 255, 255}, 800)
	if err != nil {
		return nil, err
	}

	texture, err := v.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}

	return &TextTexture{
		Texture: texture,
		Surface: surface,
	}, nil
}

func (m *PathFindingVisualization) Dispose() {
	for _, i := range m.renderedIDs {
		i.Dispose()
	}

	m.font.Close()
	ttf.Quit()
	sdl.Quit()
	m.window.Destroy()
	m.renderer.Destroy()
}
