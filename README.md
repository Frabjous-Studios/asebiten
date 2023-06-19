# asebiten

[![Go Report Card](https://goreportcard.com/badge/github.com/Frabjous-Studios/asebiten)](https://goreportcard.com/report/github.com/kalexmills/asebiten)
![GitHub](https://img.shields.io/github/license/Frabjous-Studios/asebiten)

Load exported Aseprite animations and use them in Ebitengine games.

Inspired by [ebiten-aseprite](https://pkg.go.dev/github.com/tducasse/ebiten-aseprite), with a simplified interface.

Tested with Aesprite v1.3-rc1.

## Usage:

- Export a SpriteSheet from Aseprite, using the 'Array' option from the dropdown, not 'Hash'.

- Use as in the following example
```go
package main

import (
	"embed"
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/Frabjous-Studios/asebiten"
    "github.com/Frabjous-Studios/asebiten/models/asepritev3"
	_ "image/png"
)


// example requires 'img.json' and 'img.png' from Aesprite output in the imgs directory
//go:embed imgs
var embedded embed.FS

var anim *asebiten.Animation

func init() {
	var err error
	// image is loaded relative to the JSON file; according to the path specified in the json output.
	anim, err = asepritev3.LoadAnimation(embedded, "imgs/img.json") 
	if err != nil {
		log.Fatal(err)
	}
}

const (
	width = 320
	height = 240
)

type Game struct {}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(width/2, height/2)
	anim.DrawTo(screen, op)
}

func (g *Game) Update() error {
	asebiten.Update() // call Update() once in each frame prior to calling Update() on any animations.
	anim.Update()     // call Update() once on each animation to update it based on the current frame.
	return nil
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Aesbiten Demo")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
```



### Loading SpriteSheets

SpriteSheets can be loaded by referencing the JSON file output by Aseprite. The corresponding image is loaded relative
to the location of this JSON file.

```go
anim, err = asepritev3.LoadAnimation(embedded, "imgs/img.json") 
if err != nil {
	panic(err)
}
```

To use the loaded animation, update the `Update()`  and `Draw()` funcs in your game as follows.

```go
func (g *Game) Update() error {
  asebiten.Update() // call once per update, prior to updating individual animations
  
  anim.Update() // call once per animation to make progress
  
  // ...
}

func (g *Game) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	anim.DrawTo(screen, opts)
}
```

### Tags

All animation tags set in Aseprite are loaded and made available for use in the loaded Animation. `SetTag` sets the
currently running animation to the provided tag and restarts the animation from the first frame.

```go
anim.SetTag("running")
```

### Callbacks

All animations loop by default. If this is undesirable, each animation can be modified via callbacks which are run at
the end of the provided tag. This can be used to pause a running animation, change the current tag, play a sound, or do
anything else. 

```go
anim.SetCallback("falling", func(a *Animation) {
	a.SetTag("idle")
	game.PlaySound("oof.mp3")
})
```

Take care not to put too much logic in callbacks as it can become hard to manage.

### Cloning Animations

To reuse the images loaded as part of the original sprite sheet, just `Clone()` an existing animation as follows.
Each cloned animation tracks its state separately. 

```go
ball, err = asepritev3.LoadAnimation(embedded, "imgs/ball.json") 
if err != nil {
	panic(err)
}
ball1 := ball.Clone() 
ball2 := ball.Clone()
ball3 := ball.Clone()
```

All callbacks which are set on the source animation will be copied as well. To create a lightweight animation which
doesn't receive a copy of all its callbacks, use `NewFlyweightAnimation`, which only inherits the SpriteSheet of its
argument.

```go
ball, err = asepritev3.LoadAnimation(embedded, "imgs/ball.json")
if err != nil {
	panic(err)
}
balls := make([]Animation, 100)
for i := 0; i < 100; i++ {
	balls[i] = NewFlyweightAnimation(ball) // lightweight copy; no callbacks transferred
}
```

## Comparison

There are other ways to import Aesprite files and use them in Go. This library aims to provide an opinionated and flexible option for those
using Ebitengine. You might be interested in these projects:

* [SolarLune/goaseprite](https://pkg.go.dev/github.com/solarlune/goaseprite): low-dependency library for making sense of Aseprite export formats.
  It doesn't integrate with Ebitengine, so if all you want to do is read Aseprite output formats, this could be useful.
* [tducasse/ebiten-aesprite](https://pkg.go.dev/github.com/tducasse/ebiten-aseprite): easy-to-use but hardcoded for 60 Hz, and doesn't support arbitrary
  matrix transformations.
