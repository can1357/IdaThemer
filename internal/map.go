package idathemer

import (
	"fmt"
	"math"
	"regexp"

	"github.com/lucasb-eyer/go-colorful"
)

type Instance struct {
	Data           Data
	Colors         map[string]RGBA  // key -> hex
	InvertedColors map[RGB][]string // hex -> keys[]
}

func (p *Instance) addColor(key, hex string) {
	col, err := NewCssColor(hex)
	if err != nil {
		return
	}
	p.Colors[key] = col
	p.InvertedColors[col.RGB] = append(p.InvertedColors[col.RGB], key)
}

func threeWayAdjustment(src, dst, remapped RGB) RGB {
	_, s1, l1 := src.Hsl()
	_, s2, l2 := dst.Hsl()
	h, s, l := remapped.Hsl()

	if s < s1 {
		s1, s2 = s2, s1
	}
	if l < l1 {
		l1, l2 = l2, l1
	}
	s += (s2 - s1) * 0.25
	l += (l2 - l1) * 0.25
	return colorful.Hsl(h, max(0, min(1, s)), max(0, min(1, l)))
}

func (p *Instance) remapColorRGB(col RGB, other *Instance, internal bool) RGB {
	if keys, ok := p.InvertedColors[col]; ok {

		options := make(map[RGB][]string)
		for _, key := range keys {
			if color, ok := other.Colors[key]; ok {
				options[color.RGB] = append(options[color.RGB], key)
			}
		}
		if len(options) != 0 {
			var best RGB
			var bestKeys []string
			fmt.Printf("Options for %s:\n", prettyColor(col))
			for k, v := range options {
				fmt.Printf("  %s -> %v\n", prettyColor(k), v)
				if len(v) > len(bestKeys) {
					bestKeys = v
					best = k
				}
			}
			fmt.Printf("  Remapping as %s\n", prettyColor(best))
			if !internal {
				fmt.Println()
			}
			return best
		}

	} else {
		fmt.Printf("Color %s not found\n", prettyColor(col))

		// Find the closest color
		var closest RGB
		var minDist float64 = math.Inf(1)
		for color := range p.InvertedColors {
			dist := RGBA{color, 1}.Distance(RGBA{col, 1})
			if dist < minDist {
				minDist = dist
				closest = color
			}
		}
		fmt.Printf("Closest color is %s (distance %f)\n", prettyColor(closest), minDist)
		remapped := p.remapColorRGB(closest, other, true)
		remapped = threeWayAdjustment(closest, col, remapped)
		fmt.Printf("Adjusted to %s\n\n", prettyColor(remapped))
		return remapped
	}
	return col
}

var cssHexRegex = regexp.MustCompile(`#[0-9a-fA-F]{3,8};`)

func (p *Instance) RemapCSS(css string, other *Instance) string {
	remapCache := make(map[RGB]RGB)
	css = cssHexRegex.ReplaceAllStringFunc(css, func(hex string) string {
		col, err := NewCssColor(hex[:len(hex)-1])
		if err != nil {
			return hex
		}
		prev := col.RGB
		if v, ok := remapCache[prev]; ok {
			return RGBA{v, col.Alpha}.ToCssRgba() + ";"
		}

		col.RGB = p.remapColorRGB(col.RGB, other, false)
		remapCache[prev] = col.RGB
		//fmt.Printf("Remapping %s -> %s\n", hex, col.ToCssRgba())
		return col.ToCssRgba() + ";"
	})
	return css
}
