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

const (
	l2_playerLogicWidth  = 60
	l2_playerLogicHeight = 60
	l2_enemyLogicWidth   = 80
	l2_enemyLogicHeight  = 80
	l2_projLogicWidth    = 40
	l2_projLogicHeight   = 40

	l2_gravity         = 0.6
	l2_jumpStrength    = -12.0
	l2_moveSpeed       = 4.0
	l2_projectileSpeed = 3.5
)

type Projectile struct {
	X, Y   float64
	VX, VY float64
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

	themeImg      *ebiten.Image
	playerImg     *ebiten.Image
	platformImg   *ebiten.Image
	altarImg      *ebiten.Image
	enemyImg      *ebiten.Image
	projectileImg *ebiten.Image
}

func NewLevel2() *Level2 {
	themeImg := loadImage("assets/level2/level2_theme.png")
	playerImg := loadImage("assets/level2/level2_character.png")
	enemyImg := loadImage("assets/level2/level2_enemy.png")
	altarImg := loadImage("assets/level2/level2_monolith.png")
	projectileImg := loadImage("assets/level2/level2_bomb.png")

	platformImg := ebiten.NewImage(1, 1)
	// --- THIS IS THE ONLY CHANGE ---
	// Changed the platform color from green to grey.
	platformImg.Fill(color.RGBA{R: 105, G: 105, B: 105, A: 255})

	platforms := []image.Rectangle{
		image.Rect(0, 440, 640, 480),   // Floor
		image.Rect(0, 340, 150, 360),   // Start platform
		image.Rect(210, 250, 230, 440), // Cover wall
		image.Rect(230, 250, 380, 270), // Mid platform
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
		enemyY:          120 - l2_enemyLogicHeight,
		enemyShootTimer: 90,
		themeImg:        themeImg,
		playerImg:       playerImg,
		platformImg:     platformImg,
		altarImg:        altarImg,
		enemyImg:        enemyImg,
		projectileImg:   projectileImg,
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

	l.playerVX = 0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		l.playerVX = -l2_moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		l.playerVX = l2_moveSpeed
	}
	l.playerX += l.playerVX
	playerRect := image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+l2_playerLogicWidth, int(l.playerY)+l2_playerLogicHeight)
	for _, p := range l.platforms {
		if playerRect.Overlaps(p) {
			l.playerX -= l.playerVX
			break
		}
	}

	l.playerVY += l2_gravity
	l.playerY += l.playerVY
	l.onGround = false
	playerRect = image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+l2_playerLogicWidth, int(l.playerY)+l2_playerLogicHeight)
	for _, p := range l.platforms {
		if playerRect.Overlaps(p) {
			if l.playerVY > 0 {
				l.playerY = float64(p.Min.Y - l2_playerLogicHeight)
				l.playerVY = 0
				l.onGround = true
			} else if l.playerVY < 0 {
				l.playerY = float64(p.Max.Y)
				l.playerVY = 0
			}
		}
	}

	if l.onGround && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		l.playerVY = l2_jumpStrength
	}

	l.enemyShootTimer--
	if l.enemyShootTimer <= 0 {
		spawnX := l.enemyX - l2_projLogicWidth
		spawnY := l.enemyY + (l2_enemyLogicHeight / 2) - (l2_projLogicHeight / 2)

		dirX := (l.playerX + l2_playerLogicWidth/2) - spawnX
		dirY := (l.playerY + l2_playerLogicHeight/2) - spawnY
		length := math.Sqrt(dirX*dirX + dirY*dirY)
		proj := &Projectile{
			X:  spawnX,
			Y:  spawnY,
			VX: (dirX / length) * l2_projectileSpeed,
			VY: (dirY / length) * l2_projectileSpeed,
		}
		l.projectiles = append(l.projectiles, proj)
		l.enemyShootTimer = 120
	}

	playerRect = image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+l2_playerLogicWidth, int(l.playerY)+l2_playerLogicHeight)
	for i := len(l.projectiles) - 1; i >= 0; i-- {
		p := l.projectiles[i]
		p.X += p.VX
		p.Y += p.VY
		projRect := image.Rect(int(p.X), int(p.Y), int(p.X)+l2_projLogicWidth, int(p.Y)+l2_projLogicHeight)
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

	if playerRect.Overlaps(l.altar) {
		l.isComplete = true
	}
}

func (l *Level2) Draw(screen *ebiten.Image) {
	bgOpts := &ebiten.DrawImageOptions{}
	bgW, bgH := l.themeImg.Size()
	bgOpts.GeoM.Scale(screenWidth/float64(bgW), screenHeight/float64(bgH))
	screen.DrawImage(l.themeImg, bgOpts)

	for _, p := range l.platforms {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(p.Dx()), float64(p.Dy()))
		op.GeoM.Translate(float64(p.Min.X), float64(p.Min.Y))
		screen.DrawImage(l.platformImg, op)
	}

	altarOps := &ebiten.DrawImageOptions{}
	aW, aH := l.altarImg.Size()
	altarOps.GeoM.Scale(float64(l.altar.Dx())/float64(aW), float64(l.altar.Dy())/float64(aH))
	altarOps.GeoM.Translate(float64(l.altar.Min.X), float64(l.altar.Min.Y))
	screen.DrawImage(l.altarImg, altarOps)

	enemyOps := &ebiten.DrawImageOptions{}
	eW, eH := l.enemyImg.Size()
	enemyOps.GeoM.Scale(l2_enemyLogicWidth/float64(eW), l2_enemyLogicHeight/float64(eH))
	enemyOps.GeoM.Translate(l.enemyX, l.enemyY)
	screen.DrawImage(l.enemyImg, enemyOps)

	playerOps := &ebiten.DrawImageOptions{}
	pW, pH := l.playerImg.Size()
	playerOps.GeoM.Scale(l2_playerLogicWidth/float64(pW), l2_playerLogicHeight/float64(pH))
	playerOps.GeoM.Translate(l.playerX, l.playerY)
	screen.DrawImage(l.playerImg, playerOps)

	for _, p := range l.projectiles {
		op := &ebiten.DrawImageOptions{}
		projW, projH := l.projectileImg.Size()
		op.GeoM.Scale(l2_projLogicWidth/float64(projW), l2_projLogicHeight/float64(projH))
		op.GeoM.Translate(p.X, p.Y)
		screen.DrawImage(l.projectileImg, op)
	}

	healthText := fmt.Sprintf("Health: %d", l.playerHealth)
	ebitenutil.DebugPrint(screen, healthText)
}

func (l *Level2) IsDone() bool {
	return l.isComplete
}
