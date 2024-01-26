package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	idathemer "github.com/can1357/IdaThemer/internal"
)

//go:embed rsrc/theme.css
var refThemeCSS string

//go:embed rsrc/theme.json
var refThemeJSON []byte

func createTheme(refJSON []byte, refCSS string, targetJSONPath string, outputDir string, nameAlt string) error {
	// Parse the source theme
	src, err := idathemer.Parse(refJSON)
	if err != nil {
		return err
	}

	// Parse the target theme
	dst, err := idathemer.ReadFile(targetJSONPath)
	if err != nil {
		return err
	}

	// Remap the colors
	newCSS := src.RemapCSS(refCSS, dst)
	if dst.Data.Type == "light" {
		newCSS = strings.ReplaceAll(newCSS, `@importtheme "dark";`, `@importtheme "_base";`)
	}

	// Create the output directory
	dst.Data.Name = strings.ReplaceAll(dst.Data.Name, "/", "_")
	dst.Data.Name = strings.ReplaceAll(dst.Data.Name, "\\", "_")
	if dst.Data.Name == "" {
		dst.Data.Name = nameAlt
		if dst.Data.Name == "" {
			dst.Data.Name = "theme"
		}
	}

	outputDir = filepath.Join(outputDir, dst.Data.Name)
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	// Write the new theme.css
	err = os.WriteFile(filepath.Join(outputDir, "theme.css"), []byte(newCSS), 0644)
	if err != nil {
		return err
	}
	return nil
}

type PackageJSON struct {
	Contributes struct {
		Themes []struct {
			Label string `json:"label"`
			Path  string `json:"path"`
		} `json:"themes"`
	} `json:"contributes"`
}

func createThemesFromExtension(refJSON []byte, refCSS string, extensionPath string, extensionId string, outputDir string) error {
	packageJSONPath := filepath.Join(extensionPath, "package.json")
	if _, err := os.Stat(packageJSONPath); err != nil {
		return err
	}

	var packageJSON PackageJSON
	packageJSONData, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(packageJSONData, &packageJSON)
	if err != nil {
		return err
	}

	for _, theme := range packageJSON.Contributes.Themes {
		abs := filepath.Join(extensionPath, theme.Path)
		err = createTheme(refJSON, refCSS, abs, outputDir, extensionId)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Not enough arguments")
		fmt.Printf("Usage: %s [theme.json] [ida themes dir]\n", os.Args[0])
		fmt.Printf("       %s *            [ida themes dir]\n", os.Args[0])
		os.Exit(1)
	}
	targetTheme := os.Args[1]
	outputDir := os.Args[2]

	var err error
	if targetTheme == "*" {
		home := os.Getenv("VSCODE_DATA")
		if home == "" {
			home = os.Getenv("HOME")
			if home == "" {
				home, err = os.UserHomeDir()
			}
		}
		if err == nil {
			vscodeExtDir := filepath.Join(home, ".vscode/extensions")
			var folders []os.DirEntry
			folders, err = os.ReadDir(vscodeExtDir)
			if err == nil {
				for _, folder := range folders {
					err := createThemesFromExtension(refThemeJSON, refThemeCSS, filepath.Join(vscodeExtDir, folder.Name()), folder.Name(), outputDir)
					fmt.Printf("Created themes from %s | %v\n", folder.Name(), err)
				}
			}
		}
	} else {
		err = createTheme(refThemeJSON, refThemeCSS, targetTheme, outputDir, "")
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
