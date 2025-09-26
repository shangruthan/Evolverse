package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	l3_playerLogicWidth  = 60
	l3_playerLogicHeight = 60
	l3_enemyLogicWidth   = 55
	l3_enemyLogicHeight  = 55
	l3_allyLogicWidth    = 40
	l3_allyLogicHeight   = 40
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

	themeImg        *ebiten.Image
	playerImg       *ebiten.Image
	playerAttackImg *ebiten.Image
	enemyImg        *ebiten.Image
	allyImg         *ebiten.Image
	obstacleImg     *ebiten.Image
	exitZoneImg     *ebiten.Image
}

func NewLevel3() *Level3 {
	themeImg := loadImage("assets/level3/level3_theme.png")
	playerImg := loadImage("assets/level3/level3_character.png")
	enemyImg := loadImage("assets/level3/level3_enemy.png")
	allyImg := loadImage("assets/level3/level3_ally.png")
	monolithImg := loadImage("assets/level3/level3_monolith.png")

	playerAttack := ebiten.NewImage(48, 48)
	playerAttack.Fill(color.RGBA{R: 255, G: 255, B: 150, A: 128})
	obstacle := ebiten.NewImage(1, 1)
	// --- THIS IS THE ONLY CHANGE ---
	// Changed the obstacle color from grey to black.
	obstacle.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 255})

	return &Level3{
		playerX:         50,
		playerY:         400,
		playerImg:       playerImg,
		playerAttackImg: playerAttack,
		allyX:           550,
		allyY:           60,
		allyImg:         allyImg,
		enemyX:          300,
		enemyY:          100,
		enemySpeed:      1.5,
		isEnemyAlive:    true,
		enemyImg:        enemyImg,
		obstacles: []image.Rectangle{
			image.Rect(150, 0, 180, 300),
			image.Rect(150, 300, 450, 330),
			image.Rect(420, 150, 450, 300),
		},
		obstacleImg: obstacle,
		exitZone:    image.Rect(500, 400, 620, 460),
		exitZoneImg: monolithImg,
		themeImg:    themeImg,
	}
}

func (l *Level3) Update() {
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
	playerRect := image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+l3_playerLogicWidth, int(l.playerY)+l3_playerLogicHeight)
	for _, obs := range l.obstacles {
		if playerRect.Overlaps(obs) {
			l.playerX, l.playerY = prevX, prevY
			break
		}
	}

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

	if l.isEnemyAlive {
		prevEnemyX, prevEnemyY := l.enemyX, l.enemyY
		dirX := l.playerX - l.enemyX
		dirY := l.playerY - l.enemyY
		length := math.Sqrt(dirX*dirX + dirY*dirY)
		if length > 1 {
			l.enemyX += (dirX / length) * l.enemySpeed
			l.enemyY += (dirY / length) * l.enemySpeed
		}
		enemyRect := image.Rect(int(l.enemyX), int(l.enemyY), int(l.enemyX)+l3_enemyLogicWidth, int(l.enemyY)+l3_enemyLogicHeight)
		for _, obs := range l.obstacles {
			if enemyRect.Overlaps(obs) {
				l.enemyX, l.enemyY = prevEnemyX, prevEnemyY
				break
			}
		}

		if l.isAttacking {
			distX := (l.playerX + l3_playerLogicWidth/2) - (l.enemyX + l3_enemyLogicWidth/2)
			distY := (l.playerY + l3_playerLogicHeight/2) - (l.enemyY + l3_enemyLogicHeight/2)
			if math.Sqrt(distX*distX+distY*distY) < 40 {
				l.isEnemyAlive = false
			}
		}
	} else if !l.isAllyFollowing {
		distX := l.playerX - l.allyX
		distY := l.playerY - l.allyY
		if math.Sqrt(distX*distX+distY*distY) < 40 {
			l.isAllyFollowing = true
		}
	}

	if l.isAllyFollowing {
		distX := l.playerX - l.allyX
		distY := l.playerY - l.allyY
		dist := math.Sqrt(distX*distX + distY*distY)
		if dist > 50 {
			prevAllyX, prevAllyY := l.allyX, l.allyY
			l.allyX += (distX / dist) * (playerSpeed * 0.9)
			l.allyY += (distY / dist) * (playerSpeed * 0.9)
			allyRect := image.Rect(int(l.allyX), int(l.allyY), int(l.allyX)+l3_allyLogicWidth, int(l.allyY)+l3_allyLogicHeight)
			for _, obs := range l.obstacles {
				if allyRect.Overlaps(obs) {
					l.allyX, l.allyY = prevAllyX, prevAllyY
					break
				}
			}
		}
		playerRect := image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+l3_playerLogicWidth, int(l.playerY)+l3_playerLogicHeight)
		if playerRect.Overlaps(l.exitZone) {
			l.isComplete = true
		}
	}
}

func (l *Level3) Draw(screen *ebiten.Image) {
	bgOpts := &ebiten.DrawImageOptions{}
	bgW, bgH := l.themeImg.Size()
	bgOpts.GeoM.Scale(screenWidth/float64(bgW), screenHeight/float64(bgH))
	screen.DrawImage(l.themeImg, bgOpts)

	op := &ebiten.DrawImageOptions{}
	eW, eH := l.exitZoneImg.Size()
	op.GeoM.Scale(float64(l.exitZone.Dx())/float64(eW), float64(l.exitZone.Dy())/float64(eH))
	op.GeoM.Translate(float64(l.exitZone.Min.X), float64(l.exitZone.Min.Y))
	screen.DrawImage(l.exitZoneImg, op)

	for _, obs := range l.obstacles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(obs.Dx()), float64(obs.Dy()))
		op.GeoM.Translate(float64(obs.Min.X), float64(obs.Min.Y))
		screen.DrawImage(l.obstacleImg, op)
	}

	allyOps := &ebiten.DrawImageOptions{}
	aW, aH := l.allyImg.Size()
	allyOps.GeoM.Scale(l3_allyLogicWidth/float64(aW), l3_allyLogicHeight/float64(aH))
	allyOps.GeoM.Translate(l.allyX, l.allyY)
	screen.DrawImage(l.allyImg, allyOps)

	if l.isEnemyAlive {
		enemyOps := &ebiten.DrawImageOptions{}
		eW, eH := l.enemyImg.Size()
		enemyOps.GeoM.Scale(l3_enemyLogicWidth/float64(eW), l3_enemyLogicHeight/float64(eH))
		enemyOps.GeoM.Translate(l.enemyX, l.enemyY)
		screen.DrawImage(l.enemyImg, enemyOps)
	}

	playerOps := &ebiten.DrawImageOptions{}
	pW, pH := l.playerImg.Size()
	playerOps.GeoM.Scale(l3_playerLogicWidth/float64(pW), l3_playerLogicHeight/float64(pH))
	playerOps.GeoM.Translate(l.playerX, l.playerY)
	screen.DrawImage(l.playerImg, playerOps)

	if l.isAttacking {
		attackOps := &ebiten.DrawImageOptions{}
		attackOps.GeoM.Translate(l.playerX-8, l.playerY-8)
		screen.DrawImage(l.playerAttackImg, attackOps)
	}

	objectiveText := "Defeat the enemy!"
	if !l.isEnemyAlive && !l.isAllyFollowing {
		objectiveText = "Rescue the ally!"
	} else if l.isAllyFollowing {
		objectiveText = "Escort the ally to the monolith!"
	}
	ebitenutil.DebugPrint(screen, objectiveText)
}

func (l *Level3) IsDone() bool {
	return l.isComplete
}
