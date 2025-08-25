package console

import (
	"github.com/fatih/color"
)

type Console struct{}

type styleAttrs []color.Attribute

var styles = map[string]styleAttrs{
	"Title":  {color.FgCyan, color.Italic, color.Bold},
	"Hdr":    {color.FgBlue, color.Italic, color.Bold},
	"Ok":     {color.FgHiGreen, color.Bold},
	"Err":    {color.FgHiRed, color.Bold},
	"Warn":   {color.FgHiYellow, color.Bold},
	"Info":   {color.FgHiBlue, color.Bold},
	"Prompt": {color.FgMagenta},
	"Input":  {color.FgHiMagenta, color.Bold},
	"Usr":    {color.FgHiCyan, color.BlinkSlow},
	"AI":     {color.FgHiBlue},
	"Sys":    {color.FgHiYellow, color.Bold},
}

func (c *Console) styledPrint(style string, format string, a ...any) {
	attr := styles[style]
	color.New(attr...).Printf(format, a...)
}

func (c *Console) Title(format string, a ...any)  { c.styledPrint("Title", format, a...) }
func (c *Console) Hdr(format string, a ...any)    { c.styledPrint("Hdr", format, a...) }
func (c *Console) Ok(format string, a ...any)     { c.styledPrint("Ok", format, a...) }
func (c *Console) Err(format string, a ...any)    { c.styledPrint("Err", format, a...) }
func (c *Console) Warn(format string, a ...any)   { c.styledPrint("Warn", format, a...) }
func (c *Console) Info(format string, a ...any)   { c.styledPrint("Info", format, a...) }
func (c *Console) Prompt(format string, a ...any) { c.styledPrint("Prompt", format, a...) }
func (c *Console) Input(format string, a ...any)  { c.styledPrint("Input", format, a...) }
func (c *Console) Usr(format string, a ...any)    { c.styledPrint("Usr", format, a...) }
func (c *Console) AI(format string, a ...any)     { c.styledPrint("AI", format, a...) }
func (c *Console) Sys(format string, a ...any)    { c.styledPrint("Sys", format, a...) }

func (c *Console) Write(a ...any) {
	color.New(color.FgHiWhite).Println(a...)
}

func ConsoleInstance() *Console {
	return &Console{}
}
