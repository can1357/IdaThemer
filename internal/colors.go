package idathemer

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

type RGB = colorful.Color

type RGBA struct {
	RGB
	Alpha float64
}

func prettyColor(c RGB) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex())).Render(c.Hex())
}

func (c RGBA) ToCssRgb() string {
	return c.RGB.Hex()
}
func (c RGBA) ToCssRgba() string {
	if c.Alpha > (254.0 / 255.0) {
		return c.RGB.Hex()
	} else {
		return fmt.Sprintf("%s%02X", c.RGB.Hex(), int(c.Alpha*255))
	}
}
func (c RGBA) Distance(other RGBA) float64 {
	return c.RGB.DistanceHSLuv(other.RGB)
}

func u4parse(b byte) int64 {
	switch b {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return int64(b - '0')
	case 'A', 'B', 'C', 'D', 'E', 'F':
		return int64(b - 'A' + 10)
	default:
		return int64(b - 'a' + 10)
	}
}
func u8parse(b []byte) int64 {
	return (u4parse(b[0]) << 4) + u4parse(b[1])
}

func u8vparse(hex string) (r, g, b, a int64) {
	r = u8parse([]byte(hex[1:3]))
	g = u8parse([]byte(hex[3:5]))
	b = u8parse([]byte(hex[5:7]))
	if len(hex) == 9 {
		a = u8parse([]byte(hex[7:9]))
	} else {
		a = 0xFF
	}
	return
}
func u4vparse(hex string) (r, g, b, a int64) {
	r = u4parse(hex[1]) * 0x11
	g = u4parse(hex[2]) * 0x11
	b = u4parse(hex[3]) * 0x11
	if len(hex) == 5 {
		a = u4parse(hex[4]) * 0x11
	} else {
		a = 0xFF
	}
	return
}

func NewCssColor(hex string) (RGBA, error) {
	var r, g, b, a int64
	switch len(hex) {
	case 4, 5:
		// #RGB, #RGBA
		r, g, b, a = u4vparse(hex)
	case 7, 9:
		// #RRGGBB, #RRGGBBAA
		r, g, b, a = u8vparse(hex)
	default:
		return RGBA{}, fmt.Errorf("invalid hex color: %s", hex)
	}
	return RGBA{
		RGB:   RGB{R: float64(r) / 255, G: float64(g) / 255, B: float64(b) / 255},
		Alpha: float64(a) / 255,
	}, nil
}
