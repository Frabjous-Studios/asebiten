package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kalexmills/asebiten"
	asepritev3 "github.com/kalexmills/asebiten/models/v3"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	_ "image/png"
	"log"
)

// NOTE: change path in the two places below to test another version.
const path = "v1.3-rc1-x64/anim-test.json"

//go:embed v1.3-rc1-x64
var embedded embed.FS

const (
	screenWidth  = 320
	screenHeight = 240

	frameWidth  = 16
	frameHeight = 16
)

var (
	mplusFont font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 36
	mplusFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull, // Use quantization to save glyph cache images.
	})
}

type Game struct {
	keys []ebiten.Key
}

func (g *Game) Update() error {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.Key3):
		anim.SetTag("Tag")
	case inpututil.IsKeyJustPressed(ebiten.Key2):
		anim.SetTag("Tag2")
	case inpututil.IsKeyJustPressed(ebiten.Key1):
		anim.SetTag("")
	case inpututil.IsKeyJustPressed(ebiten.KeyP):
		anim.Toggle()
	case inpututil.IsKeyJustPressed(ebiten.KeyR):
		anim.Restart()
	}
	asebiten.Update()
	anim.Update()
	return nil
}

var anim *asebiten.Animation

const msg = "animations: press (1), (2), or (3);\ntoggle pause: (P), restart: (R)"

func (g *Game) Draw(screen *ebiten.Image) {
	text.Draw(screen, msg, mplusFont, 10, 10, color.White)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	anim.DrawTo(screen, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	var err error

	anim, err = asepritev3.LoadAnimation(embedded, path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("loaded")

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Aesbiten Demo")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
