package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
	IsDone() bool
}

type Game struct {
	currentLevel Scene
	levelIndex   int
}

func (g *Game) Update() error {
	if g.currentLevel == nil || g.currentLevel.IsDone() {
		g.levelIndex++
		switch g.levelIndex {
		case 1:
			log.Println("Starting Level 1: Primordial Path")
			g.currentLevel = NewLevel1()
		case 2:
			log.Println("Starting Level 2: Dawn of Intelligence")
			g.currentLevel = NewLevel2()
		case 3:
			log.Println("Starting Level 3: Age of Empires")
			g.currentLevel = NewLevel3()
		case 4:
			log.Println("Starting Level 4: Space-Time Frontier")
			g.currentLevel = NewLevel4()
		default:
			log.Println("Game Complete! Restarting...")
			g.levelIndex = 1
			g.currentLevel = NewLevel1()
		}
	}
	g.currentLevel.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentLevel != nil {
		g.currentLevel.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		levelIndex: 0,
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Evolverse: The Dimensional Journey")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
