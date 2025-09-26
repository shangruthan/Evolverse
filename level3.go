package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Level3 struct {
	playerX, playerY float64
	allyX, allyY     float64
	isAllyFollowing  bool
	enemyX, enemyY   float64
	enemySpeed       float64
	isEnemyAlive     bool
	isAttacking      bool
	attackTimer      int
	attackCooldown   int
	isComplete       bool
	exitZone         image.Rectangle
	obstacles        []image.Rectangle
	playerImg        *ebiten.Image
	playerAttackImg  *ebiten.Image
	enemyImg         *ebiten.Image
	allyImg          *ebiten.Image
	obstacleImg      *ebiten.Image
	exitZoneImg      *ebiten.Image
	background       *ebiten.Image
}

func NewLevel3() *Level3 {
	player := ebiten.NewImage(32, 32)
	player.Fill(color.White)
	playerAttack := ebiten.NewImage(48, 48)
	playerAttack.Fill(color.RGBA{R: 255, G: 255, B: 150, A: 128})
	enemy := ebiten.NewImage(32, 32)
	enemy.Fill(color.RGBA{R: 220, G: 20, B: 60, A: 255})
	ally := ebiten.NewImage(32, 32)
	ally.Fill(color.RGBA{R: 100, G: 149, B: 237, A: 255})
	obstacle := ebiten.NewImage(1, 1)
	obstacle.Fill(color.RGBA{R: 80, G: 80, B: 80, A: 255})
	exitZone := ebiten.NewImage(1, 1)
	exitZone.Fill(color.RGBA{R: 0, G: 255, B: 0, A: 100})
	bg := ebiten.NewImage(screenWidth, screenHeight)
	bg.Fill(color.RGBA{R: 128, G: 128, B: 128, A: 255})

	return &Level3{
		playerX:         50,
		playerY:         400,
		playerImg:       player,
		playerAttackImg: playerAttack,
		allyX:           550,
		allyY:           60,
		allyImg:         ally,
		enemyX:          300,
		enemyY:          100,
		enemySpeed:      1.5,
		isEnemyAlive:    true,
		enemyImg:        enemy,
		obstacles: []image.Rectangle{
			image.Rect(150, 0, 180, 300),
			image.Rect(150, 300, 450, 330),
			image.Rect(420, 150, 450, 300),
		},
		obstacleImg: obstacle,
		exitZone:    image.Rect(500, 400, 620, 460),
		exitZoneImg: exitZone,
		background:  bg,
	}
}

func (l *Level3) Update() {
	// --- PLAYER MOVEMENT ---
	prevX, prevY := l.playerX, l.playerY
	playerSpeed := 3.0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		l.playerX -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		l.playerX += playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		l.playerY -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		l.playerY += playerSpeed
	}
	playerRect := image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+32, int(l.playerY)+32)
	for _, obs := range l.obstacles {
		if playerRect.Overlaps(obs) {
			l.playerX, l.playerY = prevX, prevY
			break
		}
	}

	// --- PLAYER ATTACK ---
	if l.attackCooldown > 0 {
		l.attackCooldown--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && l.attackCooldown <= 0 {
		l.isAttacking = true
		l.attackTimer = 15
		l.attackCooldown = 60
	}
	if l.isAttacking {
		l.attackTimer--
		if l.attackTimer <= 0 {
			l.isAttacking = false
		}
	}

	// --- ENEMY LOGIC ---
	if l.isEnemyAlive {
		dirX := l.playerX - l.enemyX
		dirY := l.playerY - l.enemyY
		length := math.Sqrt(dirX*dirX + dirY*dirY)
		if length > 1 { // Simple chase
			l.enemyX += (dirX / length) * l.enemySpeed
			l.enemyY += (dirY / length) * l.enemySpeed
		}
		if l.isAttacking { // Check if gets hit by attack
			distX := (l.playerX + 16) - (l.enemyX + 16)
			distY := (l.playerY + 16) - (l.enemyY + 16)
			if math.Sqrt(distX*distX+distY*distY) < 40 {
				l.isEnemyAlive = false
			}
		}
	} else if !l.isAllyFollowing { // --- ALLY RECRUIT LOGIC ---
		distX := l.playerX - l.allyX
		distY := l.playerY - l.allyY
		if math.Sqrt(distX*distX+distY*distY) < 40 {
			l.isAllyFollowing = true
		}
	}

	// --- ALLY FOLLOWING LOGIC ---
	if l.isAllyFollowing {
		distX := l.playerX - l.allyX
		distY := l.playerY - l.allyY
		dist := math.Sqrt(distX*distX + distY*distY)
		if dist > 50 { // Follow if too far away
			prevAllyX, prevAllyY := l.allyX, l.allyY
			l.allyX += (distX / dist) * (playerSpeed * 0.9) // Slightly slower than player
			l.allyY += (distY / dist) * (playerSpeed * 0.9)
			allyRect := image.Rect(int(l.allyX), int(l.allyY), int(l.allyX)+32, int(l.allyY)+32)
			for _, obs := range l.obstacles {
				if allyRect.Overlaps(obs) {
					l.allyX, l.allyY = prevAllyX, prevAllyY
					break
				}
			}
		}
		// --- WIN CONDITION ---
		playerRect := image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+32, int(l.playerY)+32)
		if playerRect.Overlaps(l.exitZone) {
			l.isComplete = true
		}
	}
}

func (l *Level3) Draw(screen *ebiten.Image) {
	screen.DrawImage(l.background, nil)
	// Draw Exit Zone
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(l.exitZone.Dx()), float64(l.exitZone.Dy()))
	op.GeoM.Translate(float64(l.exitZone.Min.X), float64(l.exitZone.Min.Y))
	screen.DrawImage(l.exitZoneImg, op)
	// Draw Obstacles
	for _, obs := range l.obstacles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(obs.Dx()), float64(obs.Dy()))
		op.GeoM.Translate(float64(obs.Min.X), float64(obs.Min.Y))
		screen.DrawImage(l.obstacleImg, op)
	}
	// Draw Ally
	allyOps := &ebiten.DrawImageOptions{}
	allyOps.GeoM.Translate(l.allyX, l.allyY)
	screen.DrawImage(l.allyImg, allyOps)
	// Draw Enemy
	if l.isEnemyAlive {
		enemyOps := &ebiten.DrawImageOptions{}
		enemyOps.GeoM.Translate(l.enemyX, l.enemyY)
		screen.DrawImage(l.enemyImg, enemyOps)
	}
	// Draw Player
	playerOps := &ebiten.DrawImageOptions{}
	playerOps.GeoM.Translate(l.playerX, l.playerY)
	screen.DrawImage(l.playerImg, playerOps)
	// Draw Attack Visual
	if l.isAttacking {
		attackOps := &ebiten.DrawImageOptions{}
		attackOps.GeoM.Translate(l.playerX-8, l.playerY-8)
		screen.DrawImage(l.playerAttackImg, attackOps)
	}
	// Draw HUD
	objectiveText := "Defeat the red enemy!"
	if !l.isEnemyAlive && !l.isAllyFollowing {
		objectiveText = "Rescue the blue ally!"
	} else if l.isAllyFollowing {
		objectiveText = "Escort the ally to the green exit zone!"
	}
	ebitenutil.DebugPrint(screen, objectiveText)
}

func (l *Level3) IsDone() bool {
	return l.isComplete
}
