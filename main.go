package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"time"
)

var light = vector{-15, 15, 20}

var bgColor = color.RGBA{255, 255, 255, 255}

type sphere struct {
	Center vector
	Radius float64
	Color  color.Color
}

func avgColor(colors ...color.Color) color.Color {
	ret := color.RGBA{0, 0, 0, 255}
	var count, r, g, b uint32 = 0, 0, 0, 0
	for _, c := range colors {
		rr, gg, bb, _ := c.RGBA()
		r += rr
		g += gg
		b += bb
		count++
	}
	ret.R, ret.G, ret.B = uint8(r/count/0x101), uint8(g/count/0x101), uint8(b/count/0x101)
	return ret
}

func (s sphere) HitPoint(ray vector) (bool, vector) {
	// sphere center projection in ray
	p := ray.Dot(s.Center)
	if p < 0.0 {
		return false, vector{}
	}

	// distance between sphere center and ray
	d := math.Sqrt(s.Center.Len()*s.Center.Len() - p*p)
	if d > s.Radius {
		return false, vector{}
	}
	delta := math.Sqrt(s.Radius*s.Radius - d*d)
	return true, ray.Mul(p - delta)
}

func rayTrace(w, h, fov float64, s sphere) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	tanfov2 := math.Tan(((fov / 2.0) * math.Pi) / 180.0)

	for y := 0.0; y < h; y++ {
		for x := 0.0; x < w; x++ {
			xx := (2.0*(x+0.5) - w) * tanfov2 / h
			yy := (1 - 2.0*(y+0.5)/h) * tanfov2
			img.Set(int(x), int(y), getColor(xx, yy, s))
		}
	}
	return img
}

func getColor(x, y float64, s sphere) color.Color {
	hit, p := s.HitPoint(vector{x, y, -1}.Unit())
	if hit {
		col := color.RGBA{255, 255, 255, 255}
		cosa := p.Sub(s.Center).Unit().Dot(light.Sub(p)) / light.Sub(p).Len()
		factor := (1 + cosa) / 2
		factor = math.Exp2(10*factor) / 1024
		col.R = uint8(factor * float64(col.R))
		col.G = uint8(factor * float64(col.G))
		col.B = uint8(factor * float64(col.B))
		return avgColor(col, s.Color)
	}
	return bgColor
}

// SSAA/FSAA-style antialiasing
func antiAlias(count int, img *image.RGBA) *image.RGBA {
	img2 := img
	for i := 0; i < count; i++ {
		b := img.Bounds()
		img2 = image.NewRGBA(image.Rect(0, 0, img.Rect.Dx()/2, img.Rect.Dy()/2))
		for y, yy := b.Min.Y, 0; y < b.Max.Y; y, yy = y+2, yy+1 {
			for x, xx := b.Min.X, 0; x < b.Max.X; x, xx = x+2, xx+1 {
				c := avgColor(img.At(x, y), img.At(x+1, y), img.At(x, y+1), img.At(x+1, y+1))
				img2.Set(int(xx), int(yy), c)
			}
		}
		img = img2
	}
	return img2
}

func saveImage(img *image.RGBA, fname string) {
	f, err := os.Create(fname)
	if err != nil {
		log.Fatalf("unable to open file %s: %s\n", fname, err.Error())
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Fatalln("unable to encode to image file: ", err.Error())
	}
}

func main() {
	width, height, fov, aac := 640, 480, 60.0, 2
	pngFile := "out.png"

	s := sphere{
		Center: vector{0.0, 0.0, -9.0},
		Radius: 4.0,
		Color:  color.RGBA{240, 0, 0, 255},
	}

	w := math.Pow(2, float64(aac)) * float64(width)
	h := math.Pow(2, float64(aac)) * float64(height)

	t := time.Now()
	img := rayTrace(w, h, fov, s)
	fmt.Printf("raytrace: %.3f secs\n", time.Now().Sub(t).Seconds())

	t = time.Now()
	img2 := antiAlias(aac, img)
	fmt.Printf("antialiasing: %.3f secs\n", time.Now().Sub(t).Seconds())

	saveImage(img2, pngFile)
	fmt.Println("Image file:", pngFile)
}
