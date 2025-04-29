package model

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GetFont(numFonts int) string {
	rand.Seed(time.Now().UnixNano())

	fonts := []string{
		"Georgia, Times New Roman, Times, serif",
		"Arial, Helvetica, sans-serif",
		"Times New Roman, Times, serif",
		"sans-serif",
		"serif",
		"monospace",
		"Inconsolata, monospace",
		"loveletter",
		"spacedock",
		"neuropol",
		"pixelpoiiz",
		"messypup",
	}

	arrLen := len(fonts)
	if numFonts > arrLen || numFonts <= 0 {
		numFonts = arrLen
	}

	randomFonts := make([]string, numFonts)
	perm := rand.Perm(arrLen)
	for i := 0; i < numFonts; i++ {
		randomFonts[i] = fonts[perm[i]]
	}

	return randomFonts[rand.Intn(numFonts)]
}

func GetRandomSizeValue(min, max uint) uint {
	if min > max {
		min, max = max, min
	}

	if max-min > (1<<32)-1 {
		// Handle the case where the range is too large for uint
		panic("Invalid range: difference between max and min exceeds uint max value")
	}

	rand.Seed(time.Now().UnixNano())
	return uint(rand.Intn(int(max-min+1))) + min
}

func GetRandomValue(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func GetZIndex() string {
	return fmt.Sprintf("-%d", GetRandomValue(1, 10))
}

func GetOpacity() string {
	return fmt.Sprintf("0.%d", GetRandomValue(5, 9))
}

func GetRandomImageLeft() string {
	return fmt.Sprintf("%dpx", GetRandomValue(8, 600))
}

func GetRandomImageTop() string {
	return fmt.Sprintf("%dpx", GetRandomValue(55, 600))
}

func getRandomFontSize() string {
	return fmt.Sprintf("%dpx", GetRandomValue(15, 25))
}

func GetRandomPercentageValue() string {
	return fmt.Sprintf("%d%%", GetRandomValue(0, 65))
}

func getRandomFontWeight() string {
	fontWeights := []string{"normal", "bold"}
	return fontWeights[rand.Intn(2)]
}

func RandomOrientation() string {
    orientation := []string{"writing-mode: vertical-lr; text-orientation: mixed;", "writing-mode: vertical-rl; text-orientation: upright;", ""}
	return orientation[rand.Intn(3)]
}

func GenerateImageValues() string {
	zIndex := GetZIndex()
	opacity := GetOpacity()
	topValue := GetRandomImageTop()
	leftValue := GetRandomPercentageValue()
	return fmt.Sprintf("style='z-index: %s; opacity: %s; top:%s; left: %s;'", zIndex, opacity, topValue, leftValue)
}

func GenerateRandomStyles() string {
	fontSize := getRandomFontSize()
	topValue := GetRandomPercentageValue()
	leftValue := GetRandomPercentageValue()
	fontWeight := getRandomFontWeight()
    //orientation := randomOrientation()

	return fmt.Sprintf("font-size: %s; left: %s; top:%s; font-weight: %s", fontSize, leftValue, topValue, fontWeight)
}

func ColorCalibrator() string {
	rand.Seed(time.Now().UnixNano())

	hue := 185
	saturation := GetRandomSizeValue(0, 100)
	lightness := 50

	return fmt.Sprintf("hsl(%d, %d%%, %d%%)", hue, saturation, lightness)
}

func Corrupt(str string) string {
	corruptions := []map[string]string{
		{"u": "ü"},
		{"e": "è"},
		{"e": "ë"},
		{"a": "@"},
		{"u": "ù"},
		{"a": "à"},
		{"o": "ò"},
		{"s": "$"},
		{"i": "ï"},
		{"y": "ÿ"},
		{"i": "î"},
		{"a": "á"},
		{"a": "ã"},
		{"e": "ê"},
		{"i": "ï"},
		{"o": "ô"},
		{"o": "ø"},
		{"i": "1"},
	}

	rand.Seed(time.Now().UnixNano())

	if rand.Intn(2) == 1 {
		n := rand.Intn(2) + 1
		indices := rand.Perm(len(corruptions))[:n]

		for _, index := range indices {
			corruption := corruptions[index]
			for k, v := range corruption {
				str = strings.ReplaceAll(str, k, v)
			}
		}

		return str
	}

	return str
}

func GenerateHead() string {
	return `<meta http-equiv="content-type" content="text/html; charset=utf-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta property="og:image" content="/static/navi.jpg">
    	<link rel="stylesheet" type="text/css" href="static/main.css">
    	<link rel="stylesheet" type="text/css" href="static/fonts.css">
    	<link rel="icon" href="/favicon.ico">
    	<meta name="description" content="Echo into the wired.">
        <meta name="keywords" content="anonymous,chan,textboard,text board,text-board,message board, messageboard, message-board">
	`
}

func GenerateFoot() string {
	return `<footer><img src="static/navi.jpg" alt="laincore"></footer>`
}
