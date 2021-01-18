package virtualmap

import (
	"errors"
	"path"

	"github.com/ataboo/go-metar-blink/pkg/common"
	"github.com/ataboo/go-metar-blink/pkg/geo"
	"github.com/ataboo/go-metar-blink/pkg/stationrepo"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	PaddingPx     = 20
	ImageWidthPx  = 1280
	ImageHeightPx = 720
)

func ShowMap(stations []*stationrepo.Station) error {

	if len(stations) == 0 {
		return errors.New("need atleast one station")
	}

	sdl.Init(sdl.INIT_VIDEO)

	coordinates := make([]*geo.Coordinate, len(stations))
	for i, s := range stations {
		coordinates[i] = s.Coordinate
	}

	renderSpec := CreateRenderSpec(coordinates, ImageWidthPx, ImageHeightPx, PaddingPx)

	wnd, cv, err := sdlcanvas.CreateWindow(ImageWidthPx, ImageHeightPx, "Hello")
	if err != nil {
		return err
	}
	defer wnd.Destroy()

	wndSurface, err := wnd.Window.GetSurface()
	if err != nil {
		return err
	}

	var solid *sdl.Surface
	font, err := ttf.OpenFont(path.Join(common.GetResourcesRoot(), "meslo_powerline.ttf"), 16)
	if err != nil {
		return err
	}
	defer font.Close()

	wnd.MainLoop(func() {
		w, h := float64(cv.Width()), float64(cv.Height())
		cv.SetFillStyle("#000")
		cv.FillRect(0, 0, w, h)

		for _, s := range stations {
			if solid, err = font.RenderUTF8Solid(s.ID, sdl.Color{255, 0, 0, 255}); err != nil {
				common.LogError("failed to render font: %s", err)
				continue
			}

			x, y := renderSpec.ProjectCoordinate(s.Coordinate)

			solid.Blit(nil, wndSurface, &sdl.Rect{
				X: int32(x),
				Y: int32(y),
				W: 0,
				H: 0,
			})

		}

		// for r := 0.0; r < math.Pi*2; r += math.Pi * 0.1 {
		// 	cv.SetFillStyle(int(r*10), int(r*20), int(r*40))
		// 	cv.BeginPath()
		// 	cv.MoveTo(w*0.5, h*0.5)
		// 	cv.Arc(w*0.5, h*0.5, math.Min(w, h)*0.4, r, r+0.1*math.Pi, false)
		// 	cv.ClosePath()
		// 	cv.Fill()
		// }

		// if step > 100 {
		// 	cv.SetStrokeStyle("#222")
		// } else {
		// 	cv.SetStrokeStyle("#FFF")
		// }

		// if step > 200 {
		// 	step = 0
		// }

		// cv.SetLineWidth(10)
		// cv.BeginPath()
		// cv.Arc(w*0.5, h*0.5, math.Min(w, h)*0.4, 0, math.Pi*2, false)
		// cv.Stroke()

		// step++
	})

	return nil
}
