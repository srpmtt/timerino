package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"golang.org/x/term"
)

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func print(text string, height int, message string) {
	width, _, err := term.GetSize(0)
	if err != nil {
		width = 80
	}

	lines := strings.Split(text, "\n")
	verticalPadding := (height - len(lines) - 1) / 2

	for range verticalPadding {
		fmt.Println()
	}

	maxLineLength := 0
	for _, line := range lines {
		if len(line) > maxLineLength {
			maxLineLength = len(line)
		}
	}

	horizontalPadding := (width - maxLineLength) / 2

	for _, line := range lines {
		color.Green("%s%s\n", strings.Repeat(" ", horizontalPadding), line)
	}

	if message != "" {
		messageLines := strings.SplitSeq(message, "\n")
		for msgLine := range messageLines {
			msgPadding := (width - len(msgLine)) / 2
			fmt.Printf("%s%s\n", strings.Repeat(" ", msgPadding), msgLine)
		}
	}
}

func render(hours, minutes, seconds int, message string) {
	totalSeconds := hours*3600 + minutes*60 + seconds

	for totalSeconds >= 0 {
		clearScreen()
		h := totalSeconds / 3600
		m := (totalSeconds % 3600) / 60
		s := totalSeconds % 60

		countdownText := fmt.Sprintf("%02d:%02d:%02d", h, m, s)
		asciiArt := figure.NewFigure(countdownText, "doom", true).String()

		_, height, err := term.GetSize(0)
		if err != nil {
			height = 20
		}

		print(asciiArt, height, message)
		time.Sleep(1 * time.Second)
		totalSeconds--
	}

	clearScreen()
}

func main() {
	if len(os.Args) < 4 || len(os.Args) > 5 {
		fmt.Println("usage: timerino <hours> <minutes> <seconds> <message>")
		return
	}

	hours, errH := strconv.Atoi(os.Args[1])
	minutes, errM := strconv.Atoi(os.Args[2])
	seconds, errS := strconv.Atoi(os.Args[3])
	message := ""
	if len(os.Args) == 5 {
		message = os.Args[4]
	}

	if errH != nil || errM != nil || errS != nil {
		fmt.Println("invalid parameters: timerino <hours> <minutes> <seconds> <message>")
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		clearScreen()
		os.Exit(0)
	}()

	render(hours, minutes, seconds, message)
}
