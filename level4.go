package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	l4_playerLogicWidth    = 48
	l4_playerLogicHeight   = 48
	l4_keyItemLogicWidth   = 40
	l4_keyItemLogicHeight  = 40
	l4_dropZoneLogicWidth  = 64
	l4_dropZoneLogicHeight = 64
)

type TimeState int

const (
	Past TimeState = iota
	Future
)

type Level4 struct {
	playerX, playerY    float64
	playerHasItem       bool
	playerWasInMonolith bool
	isComplete          bool
	timeState           TimeState
	timeMonolith        image.Rectangle
	wallsPast           []image.Rectangle
	wallsFuture         []image.Rectangle
	keyItemX, keyItemY  float64
	dropZone            image.Rectangle

	theme1Img       *ebiten.Image
	theme2Img       *ebiten.Image
	playerImg       *ebiten.Image
	wallImg         *ebiten.Image
	keyItemImg      *ebiten.Image
	dropZoneImg     *ebiten.Image
	timeMonolithImg *ebiten.Image
}

func NewLevel4() *Level4 {
	theme1Img := loadImage("assets/level4/level4_theme1.png")
	theme2Img := loadImage("assets/level4/level4_theme2.png")
	playerImg := loadImage("assets/level4/level4_character.png")
	keyItemImg := loadImage("assets/level4/level4_altar.png")
	dropZoneImg := loadImage("assets/level4/level4_chest.png")
	// --- NEW ---
	// Load the rocket image for the time-shifting monolith
	rocketImg := loadImage("assets/level4/level4_rocket.png")

	wall := ebiten.NewImage(1, 1)
	wall.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	// The old purple rectangle is no longer needed

	return &Level4{
		playerX:   50,
		playerY:   240,
		timeState: Past,
		wallsPast: []image.Rectangle{
			image.Rect(150, 0, 170, 200),
			image.Rect(150, 300, 170, 480),
			image.Rect(450, 200, 470, 280),
		},
		wallsFuture: []image.Rectangle{
			image.Rect(300, 100, 320, 400),
		},
		keyItemX:     580,
		keyItemY:     100,
		dropZone:     image.Rect(560, 380, 560+l4_dropZoneLogicWidth, 380+l4_dropZoneLogicHeight),
		timeMonolith: image.Rect(screenWidth/2-12, screenHeight/2-32, screenWidth/2+12, screenHeight/2+32),
		playerImg:    playerImg,
		wallImg:      wall,
		keyItemImg:   keyItemImg,
		dropZoneImg:  dropZoneImg,
		// --- CHANGED ---
		// Assign the loaded rocket image
		timeMonolithImg: rocketImg,
		theme1Img:       theme1Img,
		theme2Img:       theme2Img,
	}
}

func (l *Level4) Update() {
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

	playerRect := image.Rect(int(l.playerX), int(l.playerY), int(l.playerX)+l4_playerLogicWidth, int(l.playerY)+l4_playerLogicHeight)
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

	isInsideNow := playerRect.Overlaps(l.timeMonolith)
	if isInsideNow && !l.playerWasInMonolith {
		if l.timeState == Past {
			l.timeState = Future
		} else {
			l.timeState = Past
		}
	}
	l.playerWasInMonolith = isInsideNow

	if !l.playerHasItem {
		if l.timeState == Past {
			itemRect := image.Rect(int(l.keyItemX), int(l.keyItemY), int(l.keyItemX)+l4_keyItemLogicWidth, int(l.keyItemY)+l4_keyItemLogicHeight)
			if playerRect.Overlaps(itemRect) {
				l.playerHasItem = true
			}
		}
	} else {
		l.keyItemX = l.playerX + (l4_playerLogicWidth / 2) - (l4_keyItemLogicWidth / 2)
		l.keyItemY = l.playerY - l4_keyItemLogicHeight

		if l.timeState == Future {
			if playerRect.Overlaps(l.dropZone) {
				l.isComplete = true
			}
		}
	}
}

func (l *Level4) Draw(screen *ebiten.Image) {
	bgOpts := &ebiten.DrawImageOptions{}
	if l.timeState == Past {
		bgW, bgH := l.theme1Img.Size()
		bgOpts.GeoM.Scale(screenWidth/float64(bgW), screenHeight/float64(bgH))
		screen.DrawImage(l.theme1Img, bgOpts)
	} else {
		bgW, bgH := l.theme2Img.Size()
		bgOpts.GeoM.Scale(screenWidth/float64(bgW), screenHeight/float64(bgH))
		screen.DrawImage(l.theme2Img, bgOpts)
	}

	// --- CHANGED ---
	// Draw the rocket image scaled to the monolith's logical size
	monolithOps := &ebiten.DrawImageOptions{}
	mW, mH := l.timeMonolithImg.Size()
	monolithOps.GeoM.Scale(float64(l.timeMonolith.Dx())/float64(mW), float64(l.timeMonolith.Dy())/float64(mH))
	monolithOps.GeoM.Translate(float64(l.timeMonolith.Min.X), float64(l.timeMonolith.Min.Y))
	screen.DrawImage(l.timeMonolithImg, monolithOps)

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

	if l.timeState == Future && !l.isComplete {
		dropZoneOps := &ebiten.DrawImageOptions{}
		dzW, dzH := l.dropZoneImg.Size()
		dropZoneOps.GeoM.Scale(l4_dropZoneLogicWidth/float64(dzW), l4_dropZoneLogicHeight/float64(dzH))
		dropZoneOps.GeoM.Translate(float64(l.dropZone.Min.X), float64(l.dropZone.Min.Y))
		screen.DrawImage(l.dropZoneImg, dropZoneOps)
	}

	if l.timeState == Past && !l.playerHasItem {
		itemOps := &ebiten.DrawImageOptions{}
		itW, itH := l.keyItemImg.Size()
		itemOps.GeoM.Scale(l4_keyItemLogicWidth/float64(itW), l4_keyItemLogicHeight/float64(itH))
		itemOps.GeoM.Translate(l.keyItemX, l.keyItemY)
		screen.DrawImage(l.keyItemImg, itemOps)
	}

	playerOps := &ebiten.DrawImageOptions{}
	pW, pH := l.playerImg.Size()
	playerOps.GeoM.Scale(l4_playerLogicWidth/float64(pW), l4_playerLogicHeight/float64(pH))
	playerOps.GeoM.Translate(l.playerX, l.playerY)
	screen.DrawImage(l.playerImg, playerOps)

	if l.playerHasItem {
		itemOps := &ebiten.DrawImageOptions{}
		itW, itH := l.keyItemImg.Size()
		itemOps.GeoM.Scale(l4_keyItemLogicWidth/float64(itW), l4_keyItemLogicHeight/float64(itH))
		itemOps.GeoM.Translate(l.keyItemX, l.keyItemY)
		screen.DrawImage(l.keyItemImg, itemOps)
	}

	var objectiveText string
	if !l.playerHasItem {
		objectiveText = "Find the altar in the Past."
	} else {
		objectiveText = "Take the altar to the chest in the Future."
	}
	ebitenutil.DebugPrint(screen, objectiveText)
}

func (l *Level4) IsDone() bool {
	return l.isComplete
}
