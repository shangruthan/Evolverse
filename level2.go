package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Projectile struct {
	X, Y   float64
	VX, VY float64
	Img    *ebiten.Image
}

type Level2 struct {
	playerX, playerY   float64
	playerVX, playerVY float64
	playerHealth       int
	onGround           bool
	isComplete         bool
	isGameOver         bool
	platforms          []image.Rectangle
	altar              image.Rectangle
	enemyX, enemyY     float64
	enemyShootTimer    int
	projectiles        []*Projectile
	playerImg          *ebiten.Image
	platformImg        *ebiten.Image
	altarImg           *ebiten.Image
	enemyImg           *ebiten.Image
	background         *ebiten.Image
}

const (
	gravity         = 0.6
	jumpStrength    = -12.0
	moveSpeed       = 4.0
	playerWidth     = 32
	playerHeight    = 32
	projectileSpeed = 3.5
)

func NewLevel2() *Level2 {
	player := ebiten.NewImage(playerWidth, playerHeight)
	player.Fill(color.RGBA{R: 210, G: 105, B: 30, A: 255})

	platform := ebiten.NewImage(1, 1)
	platform.Fill(color.RGBA{R: 34, G: 139, B: 34, A: 255})

	altar := ebiten.NewImage(48, 64)
	altar.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})

	enemy := ebiten.NewImage(32, 32)
	enemy.Fill(color.RGBA{R: 128, G: 0, B: 128, A: 255})

	bg := ebiten.NewImage(screenWidth, screenHeight)
	bg.Fill(color.RGBA{R: 173, G: 216, B: 230, A: 255})

	platforms := []image.Rectangle{
		image.Rect(0, 440, 640, 480), // Floor
		image.Rect(0, 340, 150, 360), // Start platform
		// --- THIS IS THE FIX ---
		// Moved the cover wall and the platform it's on closer to the start.
		image.Rect(210, 250, 230, 440), // Cover wall (Moved left)
		image.Rect(230, 250, 380, 270), // Mid platform (Moved left)
		image.Rect(480, 120, 640, 140), // Enemy platform
		image.Rect(0, 150, 150, 170),   // Altar platform
	}
	altarRect := image.Rect(60, 86, 108, 150)

	l2 := &Level2{
		playerX:         50,
		playerY:         300,
		playerHealth:    3,
		platforms:       platforms,
		altar:           altarRect,
		enemyX:          480,
		enemyY:          120 - 32,
		enemyShootTimer: 90,
		playerImg:       player,
		platformImg:     platform,
		altarImg:        altar,
		enemyImg:        enemy,
		background:      bg,
	}
	return l2
}

func (l *Level2) Reset() {
	newState := NewLevel2()
	*l = *newState
}

func (l *Level2) Update() {
	if l.isGameOver {
		l.Reset()
		return
	}

	// --- HORIZONTAL MOVEMENT & COLLISION ---
	l.playerVX = 0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		l.playerVX = -moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		l.playerVX = moveSpeed
	}
	l.playerX += l.playerVX
	playerRect := image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+playerWidth, int(l.playerY)+playerHeight)
	for _, p := range l.platforms {
		if playerRect.Overlaps(p) {
			l.playerX -= l.playerVX
			break
		}
	}

	// --- VERTICAL MOVEMENT & COLLISION ---
	l.playerVY += gravity
	l.playerY += l.playerVY
	l.onGround = false
	playerRect = image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+playerWidth, int(l.playerY)+playerHeight)
	for _, p := range l.platforms {
		if playerRect.Overlaps(p) {
			if l.playerVY > 0 {
				l.playerY = float64(p.Min.Y - playerHeight)
				l.playerVY = 0
				l.onGround = true
			} else if l.playerVY < 0 {
				l.playerY = float64(p.Max.Y)
				l.playerVY = 0
			}
		}
	}

	// --- JUMP INPUT ---
	if l.onGround && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		l.playerVY = jumpStrength
	}

	// --- ENEMY LOGIC ---
	l.enemyShootTimer--
	if l.enemyShootTimer <= 0 {
		projImg := ebiten.NewImage(10, 10)
		projImg.Fill(color.White)

		spawnX := l.enemyX + (playerWidth / 2) - 5
		spawnY := l.enemyY + (playerHeight / 2) - 5

		dirX := l.playerX - spawnX
		dirY := l.playerY - spawnY
		length := math.Sqrt(dirX*dirX + dirY*dirY)
		proj := &Projectile{
			X:   spawnX,
			Y:   spawnY,
			VX:  (dirX / length) * projectileSpeed,
			VY:  (dirY / length) * projectileSpeed,
			Img: projImg,
		}
		l.projectiles = append(l.projectiles, proj)
		l.enemyShootTimer = 120 // Reset timer (2 seconds)
	}

	// --- PROJECTILE LOGIC ---
	playerRect = image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+playerWidth, int(l.playerY)+playerHeight)
	for i := len(l.projectiles) - 1; i >= 0; i-- {
		p := l.projectiles[i]
		p.X += p.VX
		p.Y += p.VY

		projRect := p.Img.Bounds().Add(image.Pt(int(p.X), int(p.Y)))

		if projRect.Overlaps(playerRect) {
			l.playerHealth--
			l.projectiles = append(l.projectiles[:i], l.projectiles[i+1:]...)
			if l.playerHealth <= 0 {
				l.isGameOver = true
			}
			continue
		}

		for _, plat := range l.platforms {
			if projRect.Overlaps(plat) {
				l.projectiles = append(l.projectiles[:i], l.projectiles[i+1:]...)
				break
			}
		}
	}

	// --- OBJECTIVE CHECK ---
	if playerRect.Overlaps(l.altar) {
		l.isComplete = true
	}
}

func (l *Level2) Draw(screen *ebiten.Image) {
	screen.DrawImage(l.background, nil)
	for _, p := range l.platforms {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(p.Dx()), float64(p.Dy()))
		op.GeoM.Translate(float64(p.Min.X), float64(p.Min.Y))
		screen.DrawImage(l.platformImg, op)
	}
	altarOps := &ebiten.DrawImageOptions{}
	altarOps.GeoM.Translate(float64(l.altar.Min.X), float64(l.altar.Min.Y))
	screen.DrawImage(l.altarImg, altarOps)

	enemyOps := &ebiten.DrawImageOptions{}
	enemyOps.GeoM.Translate(l.enemyX, l.enemyY)
	screen.DrawImage(l.enemyImg, enemyOps)

	playerOps := &ebiten.DrawImageOptions{}
	playerOps.GeoM.Translate(l.playerX, l.playerY)
	screen.DrawImage(l.playerImg, playerOps)

	for _, p := range l.projectiles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.X, p.Y)
		screen.DrawImage(p.Img, op)
	}

	healthText := fmt.Sprintf("Health: %d", l.playerHealth)
	ebitenutil.DebugPrint(screen, healthText)
}

func (l *Level2) IsDone() bool {
	return l.isComplete
}
