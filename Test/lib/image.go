package lib

import (
	"bufio"
	"bytes"
	"golang.org/x/image/draw"
	"image"
	_ "image/jpeg"
	"image/png"
	"math"
	"sync"
)

const (
	ResizeWidth  = 1400
	ResizeHeight = 1050

	MiddleRate float64 = 0.9
	LeftRate   float64 = 0.85
	RightRate  float64 = 0.85

	CropLocationMiddle = "m"
	CropLocationLeft   = "l"
	CropLocationRight  = "r"
)

type CropLocation string

type ImgLoc struct {
	X      int
	Y      int
	Width  int
	Height int
}

// ImgCompress 图片压缩缩放
func ImgCompress(data []byte) []byte {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return data
	}
	x, y := rect(src.Bounds().Max.X, src.Bounds().Max.Y)
	img := image.NewRGBA(image.Rectangle{
		Min: src.Bounds().Min,
		Max: image.Point{
			X: x,
			Y: y,
		},
	})

	draw.CatmullRom.Scale(img, img.Bounds(), src, src.Bounds(), draw.Src, nil)

	buf := bytes.Buffer{}
	if err = png.Encode(&buf, img); err != nil {
		return data
	}
	return buf.Bytes()
}

// ImgExpand 图片扩展
func ImgExpand(img []byte, locations []CropLocation) ([][]byte, error) {
	var list [][]byte

	var wg sync.WaitGroup
	wg.Add(len(locations))

	for _, v := range locations {
		go func(v CropLocation, img []byte) {
			defer wg.Done()
			i, _, err := image.Decode(bytes.NewReader(img))
			if err != nil {
				return
			}
			loc := location(v)
			var _img image.Image
			switch i.(type) {
			case *image.RGBA:
				_img = i.(*image.RGBA).SubImage(image.Rect(
					loc.X, loc.Y,
					loc.Width, loc.Height,
				))
			case *image.NRGBA:
				_img = i.(*image.NRGBA).SubImage(image.Rect(
					loc.X, loc.Y,
					loc.Width, loc.Height,
				))
			}
			var res bytes.Buffer
			b := bufio.NewWriter(&res)
			if err = png.Encode(b, _img); err != nil {
				return
			}
			if err = b.Flush(); err != nil {
				return
			}

			list = append(list, res.Bytes())
		}(v, img)
	}

	wg.Wait()

	return list, nil
}

func rect(x, y int) (int, int) {
	if x >= y && x > ImageMaxSize {
		rate := float64(ImageMaxSize) / float64(x)
		y = int(math.Floor(float64(y) * rate))
		x = ImageMaxSize
	} else if y >= x && y > ImageMaxSize {
		rate := float64(ImageMaxSize) / float64(y)
		x = int(math.Floor(float64(x) * rate))
		y = ImageMaxSize
	}
	return x, y
}

func location(crop CropLocation) ImgLoc {
	var rate float64
	var xMultiplier int
	switch crop {
	case CropLocationLeft:
		rate = LeftRate
		xMultiplier = 0
	case CropLocationRight:
		rate = RightRate
		xMultiplier = 2
	default:
		rate = MiddleRate
		xMultiplier = 1
	}

	width := int(ResizeWidth * rate)
	height := int(ResizeHeight * rate)

	x := (ResizeWidth - width) / 2
	x = x * xMultiplier

	y := (ResizeHeight - height) / 2

	return ImgLoc{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}
