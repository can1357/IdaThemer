package theme

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func toRGB(hex string) string {
	hex = strings.ToUpper(hex)
	switch len(hex) {
	case 4:
		// #RGB -> #RRGGBB
		return fmt.Sprintf("#%c%c%c%c%c%c", hex[1], hex[1], hex[2], hex[2], hex[3], hex[3])
	case 5:
		// #RGBA -> #RRGGBB
		return fmt.Sprintf("#%c%c%c%c%c%c%c%c", hex[1], hex[1], hex[2], hex[2], hex[3], hex[3], hex[4], hex[4])
	case 7:
		// #RRGGBB -> #RRGGBB
		return hex
	case 9:
		// #RRGGBBAA -> #RRGGBB
		return hex[:7]
	default:
		log.Fatalf("Invalid hex color: %s", hex)
		return ""
	}
}
func getComponents(hex string) (r, g, b int64) {
	if len(hex) != 7 && len(hex) != 9 {
		hex = toRGB(hex)
	}
	r, _ = strconv.ParseInt(hex[1:3], 16, 64)
	g, _ = strconv.ParseInt(hex[3:5], 16, 64)
	b, _ = strconv.ParseInt(hex[5:7], 16, 64)
	return
}

func getAlpha(hex string) (res string) {
	defer func() {
		if len(res) == 0 {
			res = "FF"
		}
		if res == "FF" {
			res = ""
		}
	}()

	hex = strings.ToUpper(hex)
	switch len(hex) {
	case 4:
		// #RGB -> FF
		break
	case 5:
		// #RGBA -> AA
		return hex[4:5] + hex[4:5]
	case 7:
		// #RRGGBB -> FF
		break
	case 9:
		// #RRGGBBAA -> AA
		return hex[7:9]
	default:
		return "FF"
	}
	return "FF"
}
func colorDistance(hex1, hex2 string) int64 {
	r1, g1, b1 := getComponents(hex1)
	r2, g2, b2 := getComponents(hex2)
	return (r1-r2)*(r1-r2) + (g1-g2)*(g1-g2) + (b1-b2)*(b1-b2)
}

type Instance struct {
	Data           Data
	Colors         map[string]string   // key -> hex
	InvertedColors map[string][]string // hex -> keys[]
}

func (p *Instance) addColor(key, hex string) {
	if len(hex) == 0 {
		return
	}
	hex = toRGB(hex)
	p.Colors[key] = hex
	p.InvertedColors[hex] = append(p.InvertedColors[hex], key)
}

func (p *Instance) remapRGB(hex string, other *Instance) string {
	if hex == "" {
		return ""
	}

	if v, ok := p.InvertedColors[hex]; ok {
		for _, k := range v {
			if v, ok := other.Colors[k]; ok {
				return v
			}
		}
	} else {
		// Find the closest color
		var closest string
		var minDist int64 = 1<<63 - 1
		for color, _ := range p.InvertedColors {
			dist := colorDistance(hex, color)
			if dist < minDist {
				minDist = dist
				closest = color
			}
		}
		if minDist > 5000 {
			return hex
		}
		return p.remapRGB(closest, other)
	}
	return hex
}

func (p *Instance) RemapCSS(css string, other *Instance) string {
	regexp := regexp.MustCompile(`#[0-9a-fA-F]{3,8};`)
	remapCache := make(map[string]string)
	css = regexp.ReplaceAllStringFunc(css, func(hex string) string {
		hex = hex[:len(hex)-1]
		alpha := getAlpha(hex)
		hex = toRGB(hex)

		if v, ok := remapCache[hex]; ok {
			return v + alpha + ";"
		}
		result := p.remapRGB(hex, other)
		remapCache[hex] = result
		if result != hex {
			fmt.Printf("Remapping %s -> %s (distance %d)\n", hex, result, colorDistance(hex, result))
		}
		return result + alpha + ";"
	})
	return css
}
