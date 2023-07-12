package main

import (
	"fmt"
	"github.com/gookit/color"
	"os"
)

// Log to send logs to system-wide configured log outputs. call as Log(Data String)
func Log(Data string) {
	TxtColor := color.FgGreen.Render
	fmt.Println(TxtColor(Data))
}

// Error send logs to system-wide configured log outputs. call as Error(Data String)
func Error(Data string) {
	TxtColor := color.FgRed.Render
	fmt.Println(TxtColor(Data))
}

// Warning send logs to system-wide configured log outputs. call as Warning(Data String)
func Warning(Data string) {
	TxtColor := color.FgYellow.Render
	fmt.Println(TxtColor(Data))
}

// Info send logs to system-wide configured log outputs. call as Info(Data String)
func Info(Data string) {
	TxtColor := color.FgBlue.Render
	fmt.Println(TxtColor(Data))
}

// Debug send logs to system-wide configured log outputs. call as Debug(Data String)
func Debug(Data string) {
	TxtColor := color.FgWhite.Render
	fmt.Println(TxtColor(Data))
}

// Panic send logs to system-wide configured log outputs and exit. call as Panic(Data String)
func Panic(Data string) {
	TxtColor := color.FgLightRed.Render
	fmt.Println(TxtColor(Data))
	os.Exit(-1)
}
