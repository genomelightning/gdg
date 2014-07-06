package gdg

import (
	"fmt"
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetMaxLevel(t *testing.T) {
	type Val struct {
		width, height uint
		level         uint
	}
	vals := []Val{
		{1598, 318, 11},
		{3198, 3198, 12},
		{4224, 3168, 13},
		{127998, 27614, 17},
	}
	Convey("Get maximum level", t, func() {
		for _, v := range vals {
			So(GetMaxLevel(v.width, v.height), ShouldEqual, v.level)
		}
	})
}

// :TODO
func TestGetLevelGrids(t *testing.T) {
	type Val struct {
		width, height   uint
		level, tileSize uint
		cols, rows      uint
	}
	vals := []Val{
		{4224, 3186, 13, 256, 17, 13},
		{2112, 1584, 12, 256, 9, 7},
		{1056, 792, 11, 256, 5, 4},
		{528, 396, 10, 256, 3, 2},
		{127998, 27614, 17, 256, 500, 108},
		{63999, 13807, 16, 256, 250, 54},
		// {32000, 6904, 15, 256, 126, 27},
	}
	Convey("Get level grids", t, func() {
		for _, v := range vals {
			cols, rows := GetLevelGrids(v.level, v.width, v.height, v.tileSize)
			So(cols, ShouldEqual, v.cols)
			So(rows, ShouldEqual, v.rows)
		}

		fmt.Println(uint(math.Ceil(float64(63999) / 2)))
	})
}
