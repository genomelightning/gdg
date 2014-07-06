// Package Go DZI Generator is a pure Go implementation of `vips dzsave`.
package gdg

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/disintegration/imaging"
)

type ImageFormat string

const (
	JPEG ImageFormat = "jpeg"
	PNG  ImageFormat = "png"
)

// Option represents general DZI options.
type Option struct {
	DirPath       string
	Format        ImageFormat
	Overlap       uint
	TileSize      uint
	Width, Height uint
}

// GetMaxLevel computes and returns the maximum level of DZI files
// based on given width and height.
func GetMaxLevel(width, height uint) uint {
	return uint(math.Ceil(math.Log2(math.Max(float64(width), float64(height)))))
}

// GetLevelGrids returns columns and rows number of current level,
// tile size, width and height.
func GetLevelGrids(level, width, height, tileSize uint) (uint, uint) {
	return uint(math.Ceil(float64(width) / float64(tileSize))),
		uint(math.Ceil(float64(height) / float64(tileSize)))
}

// ComputeTileRect computes and returns corresponding rectangle
// based on given option, column and row.
func ComputeTileRect(opt *Option, col, row, maxCol, maxRow uint) (rect image.Rectangle) {
	if col == 0 {
		rect.Max.X = int(opt.TileSize + opt.Overlap)
	} else {
		rect.Min.X = int(col*opt.TileSize - opt.Overlap)
		if col == maxCol-1 {
			rect.Max.X = int(opt.Width)
		} else {
			rect.Max.X = int((col+1)*opt.TileSize + opt.Overlap)
		}
	}

	if row == 0 {
		rect.Max.Y = int(opt.TileSize + opt.Overlap)
	} else {
		rect.Min.Y = int(row*opt.TileSize - opt.Overlap)
		if row == maxRow-1 {
			rect.Max.Y = int(opt.Height)
		} else {
			rect.Max.Y = int((row+1)*opt.TileSize + opt.Overlap)
		}
	}

	return rect
}

// SaveTile saves tile to given path based on level, column and row.
func SaveTile(dirPath string, level, col, row uint, format ImageFormat, m *image.NRGBA, wg *sync.WaitGroup) error {
	imgPath := fmt.Sprintf("%s/%d/%d_%d.%s", dirPath, level, col, row, format)
	// fmt.Println(imgPath)
	os.MkdirAll(path.Dir(imgPath), os.ModePerm)
	fw, err := os.Create(imgPath)
	if err != nil {
		return err
	}
	defer fw.Close()

	switch format {
	case JPEG:
		err = jpeg.Encode(fw, m, &jpeg.Options{jpeg.DefaultQuality})
	case PNG:
		err = png.Encode(fw, m)
	}
	wg.Done()
	runtime.GC()
	return err
}

// Generate generates DZI files of given image and option.
// Width and height in option and image should be same.
func Generate(m *image.NRGBA, opt *Option) error {
	level := GetMaxLevel(opt.Width, opt.Height)
	wg := &sync.WaitGroup{}
	tm := m

	var col, row uint
	for ; level >= 0; level-- {
		cols, rows := GetLevelGrids(level, opt.Width, opt.Height, opt.TileSize)
		wg.Add(int(cols * rows))
		for col = 0; col < cols; col++ {
			for row = 0; row < rows; row++ {
				go SaveTile(opt.DirPath, level, col, row, opt.Format,
					imaging.Crop(tm, ComputeTileRect(opt, col, row, cols, rows)), wg)
				// if err := SaveTile(opt.DirPath, level, col, row, opt.Format,
				// 	imaging.Crop(tm, ComputeTileRect(opt, col, row, cols, rows))); err != nil {
				// 	return err
				// }
			}
		}

		opt.Width = uint(math.Ceil(float64(opt.Width) / 2))
		opt.Height = uint(math.Ceil(float64(opt.Height) / 2))
		tm = imaging.Thumbnail(tm, int(opt.Width), int(opt.Height), imaging.Box)
		runtime.GC()
		if level == 0 {
			break
		}
	}

	wg.Wait()
	return nil
}
