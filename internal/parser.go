package idathemer

import (
	"encoding/json"
	"os"

	"github.com/tidwall/jsonc"
)

type TokenColors struct {
	Name     string `json:"name,omitempty"`
	Scope    any    `json:"scope,omitempty"` // string or []string
	Settings struct {
		FontStyle  string `json:"fontStyle,omitempty"`
		Foreground string `json:"foreground,omitempty"`
	} `json:"settings,omitempty"`
}

func (tc TokenColors) GetScope() []string {
	switch v := tc.Scope.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	default:
		return nil
	}
}

type Data struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Colors      map[string]string `json:"colors"`
	TokenColors []TokenColors     `json:"tokenColors,omitempty"`
}

func Parse(data []byte) (*Instance, error) {
	data = jsonc.ToJSON(data)

	var d Data
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	parsed := &Instance{
		Data:           d,
		Colors:         make(map[string]RGBA),
		InvertedColors: make(map[RGB][]string),
	}
	for _, tc := range d.TokenColors {
		for _, scope := range tc.GetScope() {
			parsed.addColor(scope, tc.Settings.Foreground)
		}
	}
	for k, v := range d.Colors {
		parsed.addColor(k, v)
	}
	return parsed, nil
}
func ReadFile(path string) (*Instance, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}
