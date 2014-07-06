package main

import (
	"fmt"
	"image"
	"log"
	"os"
	// "time"
	"runtime"

	"github.com/genomelightning/gdg"
)

func main() {
	fr, err := os.Open("all-in-one-863*4000.png")
	if err != nil {
		log.Fatalf("Fail to open image file: %v", err)
	}
	defer fr.Close()

	src, _, err := image.Decode(fr)
	if err != nil {
		log.Fatalf("Fail to decode image file: %v", err)
	}

	fmt.Println("Start generating...")

	runtime.GOMAXPROCS(runtime.NumCPU())

	gdg.Generate((*image.NRGBA)(src.(*image.RGBA)), &gdg.Option{
		DirPath:  "gdg",
		Format:   gdg.PNG,
		Overlap:  1,
		TileSize: 256,
		Width:    uint(src.Bounds().Max.X),
		Height:   uint(src.Bounds().Max.Y),
	})
}
