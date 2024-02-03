package main

import (
	"flag"
	"fmt"
	"github.com/bgrewell/commander/internal/assistants"
	"github.com/fatih/color"
	"os"
	"strings"
)

var (
	version string = "0.0.1"
	build   string = "debug"
	rev     string = "debug"
	branch  string = "debug"
)

func PrintUsageLine(parameter string, defaultValue interface{}, description string, units string, extra string) {
	yellow := color.New(color.FgHiYellow)
	cyan := color.New(color.FgHiCyan)
	red := color.New(color.FgHiRed)
	yellow.Printf("    %-22s", parameter)
	cyan.Printf("  %-14v", defaultValue)
	yellow.Printf("  %-36s", description)
	cyan.Printf("  %-10s", units)
	red.Printf("  %s\n", extra)
}

func Usage() (usage func()) {
	return func() {
		white := color.New(color.FgWhite)
		boldWhite := color.New(color.FgWhite, color.Bold)
		boldGreen := color.New(color.FgGreen, color.Bold)
		usageLineFormat := "    %-22s  %-14v  %s\n"
		//ruleLineFormat := "    %-22s  %-14v  %-36s  %s\n"
		boldGreen.Printf("[+] Commander :: Version %v :: Build %v :: Rev %v :: Branch %v\n", version, build, rev, branch)
		boldWhite.Print("Usage: ")
		fmt.Printf("commander <flags> [question]\n")
		boldGreen.Print("  General Options:\n")
		white.Printf(usageLineFormat, "Parameter", "Default", "Description")
		PrintUsageLine("--h[elp]", false, "show this help output", "[flag]", "")
		PrintUsageLine("--json", false, "output machine readable json", "[flag]", "not implemented")
		PrintUsageLine("--explain", false, "provide an explanation of the output", "[flag]", "")
	}
}

func main() {

	yellow := color.New(color.FgHiYellow)
	cyan := color.New(color.FgHiCyan)
	red := color.New(color.FgHiRed)

	var explain = flag.Bool("explain", false, "")
	flag.Usage = Usage()
	flag.Parse()
	args := flag.Args()

	question := strings.Join(args, " ")

	if question == "" {
		fmt.Println("Error: You need to ask a question")
		os.Exit(-1)
	}

	assistant, err := assistants.NewOpenAIAssistant()
	if err != nil {
		panic(err)
	}

	response, err := assistant.Query(question)
	if err != nil {
		panic(err)
	}

	cyan.Printf("Command:\n")
	fmt.Println(response[0])

	if *explain {
		explanation, err := assistant.Explain(response[0])
		if err != nil {
			panic(err)
		}
		lines := strings.Split(explanation[0], "\n")
		red.Printf("\nExplanation:\n")
		for _, line := range lines {
			yellow.Printf("  %s\n", line)
		}
	}
}
