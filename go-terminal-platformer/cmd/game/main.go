package main

import (
	"log"
	"os"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

func main() {
	// Initialize game
	g, err := engine.New()
	if err != nil {
		log.Fatalf("Failed to create game: %v", err)
	}

	// Initialize game systems
	if err := g.Initialize(); err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}
	defer g.Cleanup()

	// Run game loop
	if err := g.Run(); err != nil {
		log.Fatalf("Game error: %v", err)
	}

	os.Exit(0)
}
