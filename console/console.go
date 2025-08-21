package console

import (
	"github.com/fatih/color"
)

type Console struct{}

func ConsoleInstance() *Console {
    return &Console{}
}

// Title
func (c *Console) Title(format string, a ...any) {
	color.New(color.FgCyan, color.Italic, color.Bold).Printf(format, a...)
}

// Header
func (c *Console) Hdr(format string, a ...any) {
	color.New(color.FgBlue, color.Italic, color.Bold).Printf(format, a...)
}

// OK or Success
func (c *Console) Ok(format string, a ...any) {
	color.New(color.FgHiGreen, color.Bold).Printf(format, a...)
}

// Error
func (c *Console) Err(format string, a ...any) {
	color.New(color.FgHiRed, color.Bold).Printf(format, a...)
}

// Warning
func (c *Console) Warn(format string, a ...any) {
	color.New(color.FgHiYellow, color.Bold).Printf(format, a...)
}

// Info
func (c *Console) Info(format string, a ...any) {
	color.New(color.FgHiBlue, color.Bold).Printf(format, a...)
}

// Prompt
func (c *Console) Prompt(format string, a ...any) {
    color.New(color.FgMagenta).Printf(format, a...)
}

// Input
func (c *Console) Input(format string, a ...any) {
    color.New(color.FgHiMagenta, color.Bold).Printf(format, a...)
}

// User messages
func (c *Console) Usr(format string, a ...any) {
    color.New(color.FgHiCyan, color.BlinkSlow).Printf(format, a...)
}

// AI messages
func (c *Console) AI(format string, a ...any) {
    color.New(color.FgHiBlue).Printf(format, a...)
}

// System messages
func (c *Console) Sys(format string, a ...any) {
    color.New(color.FgHiYellow, color.Bold).Printf(format, a...)
}

// Print white text with newline
func (c *Console) Write(a ...any) {
	color.New(color.FgHiWhite).Println(a...)
}

//