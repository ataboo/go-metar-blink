package common

import "github.com/ataboo/go-metar-blink/pkg/animation"

type ColorThemeStrings struct {
	VFR        string `json:"vfr"`
	SVFR       string `json:"svfr"`
	IFR        string `json:"ifr"`
	LIFR       string `json:"lifr"`
	Error      string `json:"error"`
	Brightness string `json:"brightness"`
}

type ColorTheme struct {
	VFR        animation.Color
	SVFR       animation.Color
	IFR        animation.Color
	LIFR       animation.Color
	Error      animation.Color
	Brightness byte
}

func (t *ColorThemeStrings) ParseColors(errors map[string]string) *ColorTheme {
	brightness, err := ParseByteHexString(t.Brightness)
	if err != nil {
		errors["Color.Brightness"] = "Expecting byte hex string 0x00 - 0xFF"
	}

	return &ColorTheme{
		VFR:        t.parseColor(errors, t.VFR, "Color.VFR"),
		SVFR:       t.parseColor(errors, t.SVFR, "Color.SVFR"),
		IFR:        t.parseColor(errors, t.IFR, "Color.IFR"),
		LIFR:       t.parseColor(errors, t.LIFR, "Color.LIFR"),
		Error:      t.parseColor(errors, t.Error, "Color.Error"),
		Brightness: brightness,
	}
}

func (t *ColorThemeStrings) parseColor(errors map[string]string, colorStr string, fieldName string) animation.Color {
	color, err := ParseColorHexString(colorStr)
	if err != nil || color > 0xFFFFFF {
		errors[fieldName] = "Expecting RGB uint32 hex string 0x0 - 0xFFFFFF"
	}

	return color
}
