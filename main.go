package main

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/can1357/IdaThemer/theme"
)

//go:embed rsrc/*.svg rsrc/*.png
var resources embed.FS

//go:embed rsrc/theme.css
var themeCSS string

//go:embed rsrc/theme.json
var themeJSON []byte

func main() {
	// Parse arguments [theme.json] [output dir]
	if len(os.Args) != 3 {
		fmt.Printf("Not enough arguments")
		fmt.Printf("Usage: %s [theme.json] [output dir]\n", os.Args[0])
		os.Exit(1)
	}
	targetTheme := os.Args[1]
	outputDir := os.Args[2]

	// Parse the source theme
	src, err := theme.Parse(themeJSON)
	if err != nil {
		panic(err)
	}

	// Parse the target theme
	dst, err := theme.ReadFile(targetTheme)
	if err != nil {
		panic(err)
	}

	// Remap the colors
	newCSS := src.RemapCSS(themeCSS, dst)
	if dst.Data.Type == "light" {
		newCSS = strings.ReplaceAll(newCSS, `@importtheme "dark";`, "")
	}

	// Create the output directory
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(err)
	}

	// Copy all files from resources to output directory
	files, err := resources.ReadDir("rsrc")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		data, err := resources.ReadFile("rsrc/" + file.Name())
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(outputDir+"/"+file.Name(), data, 0644)
		if err != nil {
			panic(err)
		}
	}

	// Write the new theme.css
	err = os.WriteFile(outputDir+"/theme.css", []byte(newCSS), 0644)
	if err != nil {
		panic(err)
	}
}
