package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

type LoadingScreen struct {
	startTime time.Time
	image     *ebiten.Image
	duration  time.Duration
}

func NewLoadingScreen(imagePath string) *LoadingScreen {
	img, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		log.Printf("Failed to load loading image %s: %v", imagePath, err)
	}

	return &LoadingScreen{
		startTime: time.Now(),
		image:     img,
		duration:  3 * time.Second,
	}
}

func (ls *LoadingScreen) Update() {
	// Nothing to update for loading screen
}

func (ls *LoadingScreen) Draw(screen *ebiten.Image) {
	if ls.image != nil {
		// Get image dimensions
		bounds := ls.image.Bounds()
		imgWidth := float64(bounds.Dx())
		imgHeight := float64(bounds.Dy())

		// Calculate scale to fit screen
		scaleX := float64(screenWidth) / imgWidth
		scaleY := float64(screenHeight) / imgHeight

		// Use the smaller scale to maintain aspect ratio and fit entirely
		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}

		// Calculate position to center the scaled image
		scaledWidth := imgWidth * scale
		scaledHeight := imgHeight * scale
		x := (float64(screenWidth) - scaledWidth) / 2
		y := (float64(screenHeight) - scaledHeight) / 2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x, y)
		screen.DrawImage(ls.image, op)
	}
}

func (ls *LoadingScreen) IsDone() bool {
	return time.Since(ls.startTime) >= ls.duration
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
			log.Println("Loading Level 1 -> Level 2...")
			g.currentLevel = NewLoadingScreen("assets/loading_assets/loading1_2.png")
		case 3:
			log.Println("Starting Level 2: Dawn of Intelligence")
			g.currentLevel = NewLevel2()
		case 4:
			log.Println("Loading Level 2 -> Level 3...")
			g.currentLevel = NewLoadingScreen("assets/loading_assets/loading2_3.png")
		case 5:
			log.Println("Starting Level 3: Age of Empires")
			g.currentLevel = NewLevel3()
		case 6:
			log.Println("Loading Level 3 -> Level 4...")
			g.currentLevel = NewLoadingScreen("assets/loading_assets/loading3_4.png")
		case 7:
			log.Println("Starting Level 4: Space-Time Frontier")
			g.currentLevel = NewLevel4()
		case 8:
			log.Println("Game Complete!")
			g.currentLevel = NewLoadingScreen("assets/loading_assets/game_complete.png")
		default:
			log.Println("Restarting game...")
			g.levelIndex = 0 // Reset to start over
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
