package main

import (
	"fmt"
	"github.com/fatih/color"
)

func main() {

	fmt.Println("")
	fmt.Println("Regular           Bold")
	fmt.Println("----------------------")
	color.New(color.FgRed).Printf("red\t")
	color.New(color.BgRed).Print("         ")
	color.New(color.FgRed, color.Bold).Println(" red")

	color.New(color.FgGreen).Printf("green\t")
	color.New(color.BgGreen).Print("         ")
	color.New(color.FgGreen, color.Bold).Println(" green")

	color.New(color.FgYellow).Printf("yellow\t")
	color.New(color.BgYellow).Print("         ")
	color.New(color.FgYellow, color.Bold).Println(" yellow")

	color.New(color.FgBlue).Printf("blue\t")
	color.New(color.BgBlue).Print("         ")
	color.New(color.FgBlue, color.Bold).Println(" blue")

	color.New(color.FgMagenta).Printf("magenta\t")
	color.New(color.BgMagenta).Print("         ")
	color.New(color.FgMagenta, color.Bold).Println(" magenta")

	color.New(color.FgCyan).Printf("cyan\t")
	color.New(color.BgCyan).Print("         ")
	color.New(color.FgCyan, color.Bold).Println(" cyan")

	color.New(color.FgWhite).Printf("white\t")
	color.New(color.BgWhite).Print("         ")
	color.New(color.FgWhite, color.Bold).Println(" white")

	fmt.Println("")
	fmt.Println("Hi                Bold")
	fmt.Println("----------------------")
	color.New(color.FgHiRed).Printf("red\t")
	color.New(color.BgHiRed).Print("         ")
	color.New(color.FgHiRed, color.Bold).Println(" red")

	color.New(color.FgHiGreen).Printf("green\t")
	color.New(color.BgHiGreen).Print("         ")
	color.New(color.FgHiGreen, color.Bold).Println(" green")

	color.New(color.FgHiYellow).Printf("yellow\t")
	color.New(color.BgHiYellow).Print("         ")
	color.New(color.FgHiYellow, color.Bold).Println(" yellow")

	color.New(color.FgHiBlue).Printf("blue\t")
	color.New(color.BgHiBlue).Print("         ")
	color.New(color.FgHiBlue, color.Bold).Println(" blue")

	color.New(color.FgHiMagenta).Printf("magenta\t")
	color.New(color.BgHiMagenta).Print("         ")
	color.New(color.FgHiMagenta, color.Bold).Println(" magenta")

	color.New(color.FgHiCyan).Printf("cyan\t")
	color.New(color.BgHiCyan).Print("         ")
	color.New(color.FgHiCyan, color.Bold).Println(" cyan")

	color.New(color.FgHiWhite).Printf("white\t")
	color.New(color.BgHiWhite).Print("         ")
	color.New(color.FgHiWhite, color.Bold).Println(" white")

}
