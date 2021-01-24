package animation

import (
	"fmt"
)

const (
	ColorRed     = Color(0xFF0000)
	ColorGreen   = Color(0x00FF00)
	ColorBlue    = Color(0x0000FF)
	ColorMagenta = Color(0xFF00FF)
	ColorWhite   = Color(0xFFFFFF)
	ColorBlack   = Color(0x000000)
	ColorYellow  = Color(0xFFD23F)
)

// Color is an RGB representation of a color.
type Color uint32

// R is the red channel.
func (c Color) R() byte {
	return byte((c >> 16) & 0xFF)
}

// G is the green channel.
func (c Color) G() byte {
	return byte((c >> 8) & 0xFF)
}

// B is the blue channel.
func (c Color) B() byte {
	return byte(c & 0xFF)
}

func (c Color) String() string {
	return fmt.Sprintf("0x%x [R%x, G%x, B%x]", uint32(c), c.R(), c.G(), c.B())
}

// CreateColor creates a color from the components.
func CreateColor(r, g, b byte) Color {
	return Color(uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}

func (c Color) ARGB() uint32 {
	return 0xFF<<24 | uint32(c)
}

func (c Color) RGBA() uint32 {
	return uint32(c.R())<<24 | uint32(c.G())<<16 | uint32(c.B())<<8 | uint32(0xFF)
}

func (c Color) RGB() uint32 {
	return uint32(c)
}
