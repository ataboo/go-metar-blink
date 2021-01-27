package metaranimation

import (
	"fmt"
	"strings"

	"github.com/ataboo/go-metar-blink/pkg/animation"
	"github.com/ataboo/go-metar-blink/pkg/logger"
)

type morseChar string

const (
	FramesPerUnit  = 15
	DitLength      = 1 * FramesPerUnit
	DahLength      = 3 * FramesPerUnit
	IntraCharSpace = 1 * FramesPerUnit
	InterCharSpace = 3 * FramesPerUnit
	InterWordSpace = 7 * FramesPerUnit
)

func CreateMorseAnimation(str string) (*animation.Track, error) {
	off := animation.ColorBlack
	on := animation.ColorGreen

	position := InterWordSpace
	keyFrames := []animation.KeyFrame{
		{0, off},
		{position, off},
	}

	morseChars := stringToMorseChars(str)

	for _, m := range morseChars {
		if m == "" {
			position += InterWordSpace
			keyFrames = append(keyFrames, animation.KeyFrame{position, off})
			continue
		}

		for _, c := range m {
			switch c {
			case '*':
				position++
				keyFrames = append(keyFrames, animation.KeyFrame{position, on})
				position += DitLength
				keyFrames = append(keyFrames, animation.KeyFrame{position, on})
				position++
				keyFrames = append(keyFrames, animation.KeyFrame{position, off})
				position += IntraCharSpace
				keyFrames = append(keyFrames, animation.KeyFrame{position, off})
			case '-':
				position++
				keyFrames = append(keyFrames, animation.KeyFrame{position, on})
				position += DahLength
				keyFrames = append(keyFrames, animation.KeyFrame{position, on})
				position++
				keyFrames = append(keyFrames, animation.KeyFrame{position, off})
				position += IntraCharSpace
				keyFrames = append(keyFrames, animation.KeyFrame{position, off})
			default:
				return nil, fmt.Errorf("unsuported morse char: %s", m)
			}
		}

		position += InterCharSpace - IntraCharSpace
		keyFrames = append(keyFrames, animation.KeyFrame{position, off})
	}

	position += InterWordSpace * 4
	keyFrames = append(keyFrames, animation.KeyFrame{position, off})

	return animation.CreateTrack(position+1, true, keyFrames)
}

func stringToMorseChars(str string) []morseChar {
	trimmed := strings.ToLower(strings.ReplaceAll(str, " ", ""))

	morseChars := make([]morseChar, len(trimmed))
	for i, c := range trimmed {
		morseChars[i] = charToMorseString(c)
	}

	return morseChars
}

func charToMorseString(c rune) morseChar {
	switch c {
	case 'a':
		return "*-"
	case 'b':
		return "-***"
	case 'c':
		return "-*-*"
	case 'd':
		return "-**"
	case 'e':
		return "*"
	case 'f':
		return "**-*"
	case 'g':
		return "--*"
	case 'h':
		return "****"
	case 'i':
		return "**"
	case 'j':
		return "*---"
	case 'k':
		return "-*-"
	case 'l':
		return "*-**"
	case 'm':
		return "--"
	case 'n':
		return "-*"
	case 'o':
		return "---"
	case 'p':
		return "*--*"
	case 'q':
		return "--*-"
	case 'r':
		return "*-*"
	case 's':
		return "***"
	case 't':
		return "-"
	case 'u':
		return "**-"
	case 'v':
		return "***-"
	case 'w':
		return "*--"
	case 'x':
		return "-**-"
	case 'y':
		return "-*--"
	case 'z':
		return "--**"
	case '.':
		return "*-*-*-"
	case ',':
		return "--**--"
	case '?':
		return "**--**"
	case '/':
		return "-**-*"
	case '@':
		return "*--*-*"
	case '1':
		return "*----"
	case '2':
		return "**---"
	case '3':
		return "***--"
	case '4':
		return "****-"
	case '5':
		return "*****"
	case '6':
		return "-****"
	case '7':
		return "--***"
	case '8':
		return "---**"
	case '9':
		return "----*"
	case '0':
		return "-----"
	case ' ':
		return ""
	default:
		logger.LogError("morse char not supported '%s'", c)
		return ""
	}
}
