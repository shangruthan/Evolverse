// package main

// import (
// 	"fmt"
// 	"image"
// 	"image/color"
// 	"math/rand"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
// 	"github.com/hajimehoshi/ebiten/v2/inpututil"
// )

// type Powerup struct {
// 	X, Y float64
// 	Img  *ebiten.Image
// }

// type Bomb struct {
// 	X, Y float64
// 	Img  *ebiten.Image
// }

// type Level1 struct {
// 	playerX           float64
// 	playerY           float64
// 	playerSpeed       float64
// 	playerHealth      int
// 	playerImg         *ebiten.Image
// 	powerups          []*Powerup
// 	powerupsCollected int
// 	bombs             []*Bomb
// 	bombSpawnTimer    int
// 	isComplete        bool
// 	background        *ebiten.Image
// 	groundHeight      float64
// 	isGameOver        bool
// }

// func NewLevel1() *Level1 {
// 	player := ebiten.NewImage(32, 32)
// 	player.Fill(color.RGBA{R: 50, G: 205, B: 50, A: 255})

// 	bg := ebiten.NewImage(screenWidth, screenHeight)
// 	bg.Fill(color.RGBA{R: 135, G: 206, B: 235, A: 255})

// 	groundH := 40.0

// 	l1 := &Level1{
// 		playerX:           screenWidth / 2,
// 		playerY:           screenHeight - groundH - 32,
// 		playerSpeed:       4.0,
// 		playerHealth:      3,
// 		playerImg:         player,
// 		powerupsCollected: 0,
// 		bombSpawnTimer:    120, // Spawn a bomb every 2 seconds
// 		isComplete:        false,
// 		isGameOver:        false,
// 		background:        bg,
// 		groundHeight:      groundH,
// 	}

// 	// Create 5 powerups at random positions
// 	powerupImg := ebiten.NewImage(16, 16)
// 	powerupImg.Fill(color.RGBA{R: 255, G: 255, B: 0, A: 255})
// 	for i := 0; i < 5; i++ {
// 		pu := &Powerup{
// 			X:   rand.Float64() * (screenWidth - 16),
// 			Y:   screenHeight - groundH - 16,
// 			Img: powerupImg,
// 		}
// 		l1.powerups = append(l1.powerups, pu)
// 	}

// 	return l1
// }

// func (l *Level1) Reset() {
// 	// By calling NewLevel1(), we get a fresh state
// 	newState := NewLevel1()
// 	*l = *newState
// }

// func (l *Level1) Update() {
// 	if l.isGameOver {
// 		// You could add a timer or key press to restart
// 		l.Reset()
// 		return
// 	}

// 	// Player Movement
// 	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
// 		l.playerX -= l.playerSpeed
// 	}
// 	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
// 		l.playerX += l.playerSpeed
// 	}
// 	if l.playerX < 0 {
// 		l.playerX = 0
// 	}
// 	if l.playerX > screenWidth-32 {
// 		l.playerX = screenWidth - 32
// 	}

// 	// Bomb Spawning
// 	l.bombSpawnTimer--
// 	if l.bombSpawnTimer <= 0 {
// 		bombImg := ebiten.NewImage(10, 20)
// 		bombImg.Fill(color.Black)
// 		newBomb := &Bomb{
// 			X:   rand.Float64() * (screenWidth - 10),
// 			Y:   -20, // Start just above the screen
// 			Img: bombImg,
// 		}
// 		l.bombs = append(l.bombs, newBomb)
// 		l.bombSpawnTimer = rand.Intn(90) + 30 // Spawn next bomb in 0.5 to 1.5 seconds
// 	}

// 	playerRect := ebiten.NewImage(32, 32).Bounds()
// 	playerRect = playerRect.Add(image.Pt(int(l.playerX), int(l.playerY)))

// 	// Update Bombs and Check Collision
// 	for i := len(l.bombs) - 1; i >= 0; i-- {
// 		bomb := l.bombs[i]
// 		bomb.Y += 3.0 // Bomb fall speed

// 		// Remove bomb if it's off-screen
// 		if bomb.Y > screenHeight {
// 			l.bombs = append(l.bombs[:i], l.bombs[i+1:]...)
// 			continue
// 		}

// 		bombRect := bomb.Img.Bounds().Add(image.Pt(int(bomb.X), int(bomb.Y)))
// 		if playerRect.Overlaps(bombRect) {
// 			l.playerHealth--
// 			l.bombs = append(l.bombs[:i], l.bombs[i+1:]...) // Remove bomb on hit
// 			if l.playerHealth <= 0 {
// 				l.isGameOver = true
// 			}
// 		}
// 	}

// 	// Check Powerup Collection
// 	for i := len(l.powerups) - 1; i >= 0; i-- {
// 		powerup := l.powerups[i]
// 		powerupRect := powerup.Img.Bounds().Add(image.Pt(int(powerup.X), int(powerup.Y)))
// 		if playerRect.Overlaps(powerupRect) {
// 			l.powerupsCollected++
// 			l.powerups = append(l.powerups[:i], l.powerups[i+1:]...) // Remove collected powerup
// 		}
// 	}

// 	// Win Condition
// 	if l.powerupsCollected >= 5 {
// 		l.isComplete = true
// 	}
// }

// func (l *Level1) Draw(screen *ebiten.Image) {
// 	screen.DrawImage(l.background, nil)

// 	ground := ebiten.NewImage(screenWidth, int(l.groundHeight))
// 	ground.Fill(color.RGBA{R: 139, G: 69, B: 19, A: 255})
// 	groundOps := &ebiten.DrawImageOptions{}
// 	groundOps.GeoM.Translate(0, screenHeight-l.groundHeight)
// 	screen.DrawImage(ground, groundOps)

// 	playerOps := &ebiten.DrawImageOptions{}
// 	playerOps.GeoM.Translate(l.playerX, l.playerY)
// 	screen.DrawImage(l.playerImg, playerOps)

// 	for _, p := range l.powerups {
// 		puOps := &ebiten.DrawImageOptions{}
// 		puOps.GeoM.Translate(p.X, p.Y)
// 		screen.DrawImage(p.Img, puOps)
// 	}

// 	for _, b := range l.bombs {
// 		bombOps := &ebiten.DrawImageOptions{}
// 		bombOps.GeoM.Translate(b.X, b.Y)
// 		screen.DrawImage(b.Img, bombOps)
// 	}

// 	// Draw HUD
// 	healthText := fmt.Sprintf("Health: %d", l.playerHealth)
// 	powerupText := fmt.Sprintf("Powerups: %d/5", l.powerupsCollected)
// 	ebitenutil.DebugPrint(screen, healthText+"\n"+powerupText)
// }

// func (l *Level1) IsDone() bool {
// 	return l.isComplete
// }

package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Helper function to load images
func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Failed to load image %s: %v", path, err)
	}
	return img
}

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

	// Image assets
	backgroundImg *ebiten.Image
	playerImg     *ebiten.Image
	bombImg       *ebiten.Image
	powerupImg    *ebiten.Image
}

func NewLevel1() *Level1 {
	// Load all assets from the "assets" folder
	bgImg := loadImage("assets/level1_theme.png")
	playerImg := loadImage("assets/level1_character.png")
	bombImg := loadImage("assets/level1_bomb.png")
	powerupImg := loadImage("assets/level1_powerup.png")

	groundH := 40.0
	playerW, playerH := playerImg.Size()

	l1 := &Level1{
		playerX:           screenWidth / 2,
		playerY:           screenHeight - groundH - float64(playerH),
		playerSpeed:       4.0,
		playerHealth:      3,
		powerupsCollected: 0,
		bombSpawnTimer:    120, // Spawn a bomb every 2 seconds
		isComplete:        false,
		isGameOver:        false,
		groundHeight:      groundH,
		// Assign loaded images
		backgroundImg: bgImg,
		playerImg:     playerImg,
		bombImg:       bombImg,
		powerupImg:    powerupImg,
	}

	// Create 5 powerups at random positions
	powerupW, _ := l1.powerupImg.Size()
	for i := 0; i < 5; i++ {
		pu := &Powerup{
			X: rand.Float64() * (screenWidth - float64(powerupW)),
			Y: screenHeight - groundH - float64(l1.powerupImg.Bounds().Dy()),
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

	playerW, _ := l.playerImg.Size()

	// Player Movement
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		l.playerX -= l.playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		l.playerX += l.playerSpeed
	}
	if l.playerX < 0 {
		l.playerX = 0
	}
	if l.playerX > screenWidth-float64(playerW) {
		l.playerX = screenWidth - float64(playerW)
	}

	// Bomb Spawning
	l.bombSpawnTimer--
	if l.bombSpawnTimer <= 0 {
		bombW, _ := l.bombImg.Size()
		newBomb := &Bomb{
			X: rand.Float64() * (screenWidth - float64(bombW)),
			Y: -float64(l.bombImg.Bounds().Dy()), // Start just above the screen
		}
		l.bombs = append(l.bombs, newBomb)
		l.bombSpawnTimer = rand.Intn(90) + 30 // Spawn next bomb in 0.5 to 1.5 seconds
	}

	playerRect := l.playerImg.Bounds().Add(image.Pt(int(l.playerX), int(l.playerY)))

	// Update Bombs and Check Collision
	for i := len(l.bombs) - 1; i >= 0; i-- {
		bomb := l.bombs[i]
		bomb.Y += 3.0 // Bomb fall speed

		if bomb.Y > screenHeight {
			l.bombs = append(l.bombs[:i], l.bombs[i+1:]...)
			continue
		}

		bombRect := l.bombImg.Bounds().Add(image.Pt(int(bomb.X), int(bomb.Y)))
		if playerRect.Overlaps(bombRect) {
			l.playerHealth--
			l.bombs = append(l.bombs[:i], l.bombs[i+1:]...) // Remove bomb on hit
			if l.playerHealth <= 0 {
				l.isGameOver = true
			}
		}
	}

	// Check Powerup Collection
	for i := len(l.powerups) - 1; i >= 0; i-- {
		powerup := l.powerups[i]
		powerupRect := l.powerupImg.Bounds().Add(image.Pt(int(powerup.X), int(powerup.Y)))
		if playerRect.Overlaps(powerupRect) {
			l.powerupsCollected++
			l.powerups = append(l.powerups[:i], l.powerups[i+1:]...) // Remove collected powerup
		}
	}

	// Win Condition
	if l.powerupsCollected >= 5 {
		l.isComplete = true
	}
}

func (l *Level1) Draw(screen *ebiten.Image) {
	// Draw Background
	bgOpts := &ebiten.DrawImageOptions{}
	bgW, bgH := l.backgroundImg.Size()
	bgOpts.GeoM.Scale(screenWidth/float64(bgW), screenHeight/float64(bgH))
	screen.DrawImage(l.backgroundImg, bgOpts)

	// Draw Ground (a simple colored bar, can be replaced with an asset)
	ground := ebiten.NewImage(screenWidth, int(l.groundHeight))
	ground.Fill(color.RGBA{R: 139, G: 69, B: 19, A: 255})
	groundOps := &ebiten.DrawImageOptions{}
	groundOps.GeoM.Translate(0, screenHeight-l.groundHeight)
	screen.DrawImage(ground, groundOps)

	// Draw Player
	playerOps := &ebiten.DrawImageOptions{}
	playerOps.GeoM.Translate(l.playerX, l.playerY)
	screen.DrawImage(l.playerImg, playerOps)

	// Draw Powerups
	for _, p := range l.powerups {
		puOps := &ebiten.DrawImageOptions{}
		puOps.GeoM.Translate(p.X, p.Y)
		screen.DrawImage(l.powerupImg, puOps)
	}

	// Draw Bombs
	for _, b := range l.bombs {
		bombOps := &ebiten.DrawImageOptions{}
		bombOps.GeoM.Translate(b.X, b.Y)
		screen.DrawImage(l.bombImg, bombOps)
	}

	// Draw HUD
	healthText := fmt.Sprintf("Health: %d", l.playerHealth)
	powerupText := fmt.Sprintf("Powerups: %d/5", l.powerupsCollected)
	ebitenutil.DebugPrint(screen, healthText+"\n"+powerupText)
}

func (l *Level1) IsDone() bool {
	return l.isComplete
}
