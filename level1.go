package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Seed the random number generator once
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Helper function to load images
func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Failed to load image %s: %v", path, err)
	}
	return img
}

const (
	// Character is now 64x64 pixels
	playerLogicWidth   = 80
	playerLogicHeight  = 80
	powerupLogicWidth  = 25
	powerupLogicHeight = 25
	// Bomb is now 20x40 pixels
	bombLogicWidth  = 40
	bombLogicHeight = 60
)

type Powerup struct {
	X, Y float64
}

type Bomb struct {
	X, Y float64
}

type Level1 struct {
	playerX           float64
	playerY           float64
	playerSpeed       float64
	playerHealth      int
	powerups          []*Powerup
	powerupsCollected int
	bombs             []*Bomb
	bombSpawnTimer    int
	isComplete        bool
	isGameOver        bool
	groundHeight      float64

	backgroundImg *ebiten.Image
	playerImg     *ebiten.Image
	bombImg       *ebiten.Image
	powerupImg    *ebiten.Image
}

func NewLevel1() *Level1 {
	bgImg := loadImage("assets/level1/level1_theme.png")
	playerImg := loadImage("assets/level1/level1_character.png")
	bombImg := loadImage("assets/level1/level1_bomb.png")
	powerupImg := loadImage("assets/level1/level1_powerup.png")

	groundH := 40.0

	l1 := &Level1{
		playerX:           screenWidth / 2,
		playerY:           screenHeight - groundH - playerLogicHeight,
		playerSpeed:       4.0,
		playerHealth:      3,
		powerupsCollected: 0,
		bombSpawnTimer:    120,
		isComplete:        false,
		isGameOver:        false,
		groundHeight:      groundH,
		backgroundImg:     bgImg,
		playerImg:         playerImg,
		bombImg:           bombImg,
		powerupImg:        powerupImg,
	}

	for i := 0; i < 5; i++ {
		pu := &Powerup{
			X: rand.Float64() * (screenWidth - powerupLogicWidth),
			Y: screenHeight - groundH - powerupLogicHeight,
		}
		l1.powerups = append(l1.powerups, pu)
	}
	return l1
}

func (l *Level1) Reset() {
	newState := NewLevel1()
	*l = *newState
}

func (l *Level1) Update() {
	if l.isGameOver {
		l.Reset()
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		l.playerX -= l.playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		l.playerX += l.playerSpeed
	}
	if l.playerX < 0 {
		l.playerX = 0
	}
	if l.playerX > screenWidth-playerLogicWidth {
		l.playerX = screenWidth - playerLogicWidth
	}

	l.bombSpawnTimer--
	if l.bombSpawnTimer <= 0 {
		newBomb := &Bomb{
			X: rand.Float64() * (screenWidth - bombLogicWidth),
			Y: -bombLogicHeight,
		}
		l.bombs = append(l.bombs, newBomb)
		l.bombSpawnTimer = rand.Intn(90) + 30
	}

	playerRect := image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+playerLogicWidth, int(l.playerY)+playerLogicHeight)

	for i := len(l.bombs) - 1; i >= 0; i-- {
		bomb := l.bombs[i]
		bomb.Y += 3.0
		if bomb.Y > screenHeight {
			l.bombs = append(l.bombs[:i], l.bombs[i+1:]...)
			continue
		}
		bombRect := image.Rect(int(bomb.X), int(bomb.Y), int(bomb.X)+bombLogicWidth, int(bomb.Y)+bombLogicHeight)
		if playerRect.Overlaps(bombRect) {
			l.playerHealth--
			l.bombs = append(l.bombs[:i], l.bombs[i+1:]...)
			if l.playerHealth <= 0 {
				l.isGameOver = true
			}
		}
	}

	for i := len(l.powerups) - 1; i >= 0; i-- {
		powerup := l.powerups[i]
		powerupRect := image.Rect(int(powerup.X), int(powerup.Y), int(powerup.X)+powerupLogicWidth, int(powerup.Y)+powerupLogicHeight)
		if playerRect.Overlaps(powerupRect) {
			l.powerupsCollected++
			l.powerups = append(l.powerups[:i], l.powerups[i+1:]...)
		}
	}

	if l.powerupsCollected >= 5 {
		l.isComplete = true
	}
}

func (l *Level1) Draw(screen *ebiten.Image) {
	bgOpts := &ebiten.DrawImageOptions{}
	bgW, bgH := l.backgroundImg.Size()
	bgOpts.GeoM.Scale(screenWidth/float64(bgW), screenHeight/float64(bgH))
	screen.DrawImage(l.backgroundImg, bgOpts)

	ground := ebiten.NewImage(screenWidth, int(l.groundHeight))
	ground.Fill(color.RGBA{R: 139, G: 69, B: 19, A: 255})
	groundOps := &ebiten.DrawImageOptions{}
	groundOps.GeoM.Translate(0, screenHeight-l.groundHeight)
	screen.DrawImage(ground, groundOps)

	playerOps := &ebiten.DrawImageOptions{}
	pW, pH := l.playerImg.Size()
	playerOps.GeoM.Scale(playerLogicWidth/float64(pW), playerLogicHeight/float64(pH))
	playerOps.GeoM.Translate(l.playerX, l.playerY)
	screen.DrawImage(l.playerImg, playerOps)

	for _, p := range l.powerups {
		puOps := &ebiten.DrawImageOptions{}
		puW, puH := l.powerupImg.Size()
		puOps.GeoM.Scale(powerupLogicWidth/float64(puW), powerupLogicHeight/float64(puH))
		puOps.GeoM.Translate(p.X, p.Y)
		screen.DrawImage(l.powerupImg, puOps)
	}

	for _, b := range l.bombs {
		bombOps := &ebiten.DrawImageOptions{}
		bW, bH := l.bombImg.Size()
		bombOps.GeoM.Scale(bombLogicWidth/float64(bW), bombLogicHeight/float64(bH))
		bombOps.GeoM.Translate(b.X, b.Y)
		screen.DrawImage(l.bombImg, bombOps)
	}

	healthText := fmt.Sprintf("Health: %d", l.playerHealth)
	powerupText := fmt.Sprintf("Powerups: %d/5", l.powerupsCollected)
	ebitenutil.DebugPrint(screen, healthText+"\n"+powerupText)
}

func (l *Level1) IsDone() bool {
	return l.isComplete
}
