package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TimeState int

const (
	Past TimeState = iota
	Future
)

type Level4 struct {
	playerX, playerY float64
	isComplete       bool
	timeState        TimeState
	wallsPast        []image.Rectangle
	wallsFuture      []image.Rectangle
	core             image.Rectangle
	playerImg        *ebiten.Image
	wallImg          *ebiten.Image
	coreImg          *ebiten.Image
}

func NewLevel4() *Level4 {
	player := ebiten.NewImage(32, 32)
	player.Fill(color.RGBA{R: 200, G: 200, B: 255, A: 255})

	wall := ebiten.NewImage(1, 1)
	wall.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})

	core := ebiten.NewImage(32, 32)
	core.Fill(color.RGBA{R: 0, G: 191, B: 255, A: 255})

	return &Level4{
		playerX:   50,
		playerY:   240,
		timeState: Past,
		wallsPast: []image.Rectangle{
			image.Rect(200, 100, 220, 400),
		},
		wallsFuture: []image.Rectangle{
			image.Rect(400, 0, 420, 200),
			image.Rect(400, 300, 420, 480),
		},
		core:      image.Rect(580, 224, 612, 256),
		playerImg: player,
		wallImg:   wall,
		coreImg:   core,
	}
}

func (l *Level4) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		if l.timeState == Past {
			l.timeState = Future
		} else {
			l.timeState = Past
		}
	}

	prevX, prevY := l.playerX, l.playerY
	speed := 3.0

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		l.playerX -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		l.playerX += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		l.playerY -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		l.playerY += speed
	}

	playerRect := image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+32, int(l.playerY)+32)

	var currentWalls []image.Rectangle
	if l.timeState == Past {
		currentWalls = l.wallsPast
	} else {
		currentWalls = l.wallsFuture
	}

	for _, wall := range currentWalls {
		if playerRect.Overlaps(wall) {
			l.playerX, l.playerY = prevX, prevY
			break
		}
	}

	if playerRect.Overlaps(l.core) {
		l.isComplete = true
	}
}

func (l *Level4) Draw(screen *ebiten.Image) {
	bgColor := color.RGBA{R: 10, G: 0, B: 20, A: 255}
	if l.timeState == Future {
		bgColor = color.RGBA{R: 0, G: 20, B: 10, A: 255}
	}
	screen.Fill(bgColor)

	var currentWalls []image.Rectangle
	if l.timeState == Past {
		currentWalls = l.wallsPast
	} else {
		currentWalls = l.wallsFuture
	}

	for _, wall := range currentWalls {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(wall.Dx()), float64(wall.Dy()))
		op.GeoM.Translate(float64(wall.Min.X), float64(wall.Min.Y))
		screen.DrawImage(l.wallImg, op)
	}

	coreOps := &ebiten.DrawImageOptions{}
	coreOps.GeoM.Translate(float64(l.core.Min.X), float64(l.core.Min.Y))
	screen.DrawImage(l.coreImg, coreOps)

	playerOps := &ebiten.DrawImageOptions{}
	playerOps.GeoM.Translate(l.playerX, l.playerY)
	screen.DrawImage(l.playerImg, playerOps)
}

func (l *Level4) IsDone() bool {
	return l.isComplete
}
