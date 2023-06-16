package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kalexmills/asebiten"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	_ "image/png"
	"log"
)

// NOTE: change path in the two places below to test another version.
const (
	path          = "v1.3-rc1-x64/animation_terminal.json"
	indicatorPath = "v1.3-rc1-x64/ui-indicators.json"
)

//go:embed v1.3-rc1-x64
var embedded embed.FS

const (
	screenWidth  = 320
	screenHeight = 240

	frameWidth  = 112
	frameHeight = 100
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
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	keys []ebiten.Key
}

var (
	ticks = 0
	frame = 0
)

func (g *Game) Update() error {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.Key1):
		anim.SetTag("ON")
	case inpututil.IsKeyJustPressed(ebiten.Key2):
		anim.SetTag("OFF")
	case inpututil.IsKeyJustPressed(ebiten.Key3):
		anim.SetTag("CRASH")
	case inpututil.IsKeyJustPressed(ebiten.Key4):
		anim.SetTag("")
	case inpututil.IsKeyJustPressed(ebiten.KeyP):
		anim.Toggle()
	case inpututil.IsKeyJustPressed(ebiten.KeyR):
		anim.Restart()
	}
	asebiten.Update()
	anim.Update()

	ticks++
	if ticks > 60 {
		frame++
		if err := indicator.SetFrame(frame); err != nil {
			frame = 1
			_ = indicator.SetFrame(frame)
		}
		ticks = 0
	}
	return nil
}

var (
	anim      *asebiten.Animation
	indicator *asebiten.Animation
)

const msg = "animations: select (1-5);\ntoggle pause: (P), restart: (R)"

func (g *Game) Draw(screen *ebiten.Image) {
	text.Draw(screen, msg, mplusFont, 10, 10, color.White)

	anim.DrawPackedTo(screen, func(opts *ebiten.DrawImageOptions) {
		opts.GeoM.Translate(screenWidth/2-frameWidth/2, screenHeight/2-frameHeight/2)
		//opts.GeoM.Scale(2, 4)
		//op.GeoM.Translate(screenWidth/2, screenHeight/2)
	})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(frameWidth)/2, float64(frameHeight))
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	indicator.DrawTo(screen, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	var err error

	anim, err = asebiten.LoadAnimation(embedded, path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(anim)

	indicator, err = asebiten.LoadAnimation(embedded, indicatorPath)
	if err != nil {
		log.Fatal(err)
	}
	indicator.Pause()
	_ = indicator.SetFrame(1)

	fmt.Println("loaded")

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Aesbiten Demo")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
