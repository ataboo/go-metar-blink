package common

import "testing"

func TestColorThemeParse(t *testing.T) {
	colorStr := &ColorThemeStrings{
		VFR:        "0x123456",
		SVFR:       "0x234567",
		IFR:        "0x345678",
		LIFR:       "0x456789",
		Error:      "0x567890",
		Brightness: "0x7C",
	}

	errors := make(map[string]string)
	parsed := colorStr.ParseColors(errors)

	if len(errors) > 0 {
		for field, err := range errors {
			t.Errorf("%s|%s", field, err)
		}
	}

	if parsed.VFR != 0x123456 {
		t.Error("unnexpected values")
	}

	if parsed.SVFR != 0x234567 {
		t.Error("unnexpected values")
	}

	if parsed.IFR != 0x345678 {
		t.Error("unnexpected values")
	}

	if parsed.LIFR != 0x456789 {
		t.Error("unnexpected values")
	}

	if parsed.Error != 0x567890 {
		t.Error("unnexpected values")
	}

	if parsed.Brightness != 0x7c {
		t.Error("unnexpected values")
	}
}

func TestColorParseErrors(t *testing.T) {
	colorStr := &ColorThemeStrings{
		VFR:        "nothex",
		SVFR:       "-1",
		IFR:        "",
		LIFR:       "",
		Error:      "",
		Brightness: "",
	}

	errors := make(map[string]string)
	colorStr.ParseColors(errors)

	if len(errors) != 6 {
		t.Error("expected 6 field errors")
	}
}
