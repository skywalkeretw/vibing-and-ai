package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// Initialize tcell screen
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}

	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}
	defer screen.Fini()

	// Set up screen
	screen.Clear()
	screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))

	// Display welcome message
	message := "Go Terminal Platformer - Press ESC to exit"
	row := 0
	for _, r := range message {
		screen.SetContent(row, 0, r, nil, tcell.StyleDefault)
		row++
	}

	status := "Game starting..."
	row = 0
	for _, r := range status {
		screen.SetContent(row, 2, r, nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
		row++
	}

	screen.Show()

	// Simple event loop
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}
