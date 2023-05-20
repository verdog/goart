package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
)

// TODO move CircularMask stuff to own file
type CircularMask struct {
	source image.Image
	center image.Point
	radius int
}

func (c *CircularMask) ColorModel() color.Model {
	return c.source.ColorModel()
}

func (c *CircularMask) Bounds() image.Rectangle {
	return image.Rect(c.center.X-c.radius, c.center.Y-c.radius, c.center.X+c.radius, c.center.Y+c.radius)
}

func (c *CircularMask) At(x, y int) color.Color {
	xx := float64(x - c.center.X)
	yy := float64(y - c.center.Y)
	rr := float64(c.radius)

	if xx*xx+yy*yy < rr*rr {
		return c.source.At(x, y)
	}
	return color.Alpha{0}
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("usage:", os.Args[0], "/path/to/image", "x", "y", "radius")
		return
	}

	// load input image
	img_path := os.Args[1]
	img, err := func(path string) (*image.NRGBA, error) {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			return nil, err
		}

		img, _, err := image.Decode(file)
		if err != nil {
			return nil, err
		}

		// input may not be NRGBA, so convert it
		bounds := img.Bounds()
		nrgba := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
		draw.Draw(nrgba, nrgba.Bounds(), img, bounds.Min, draw.Src)

		return nrgba, nil
	}(img_path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// load numbers
	x, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err, os.Args[2])
		return
	}

	y, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println(err, os.Args[3])
		return
	}

	r, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Println(err, os.Args[4])
		return
	}

	circle := &CircularMask{
		source: img,
		center: image.Point{X: x, Y: y},
		radius: r,
	}

	output, err := os.Create(img_path + ".masked")
	defer output.Close()
	if err != nil {
		fmt.Println("Failed to open file", img_path+".masked")
		fmt.Println(err)
		return
	}
	png.Encode(output, circle)

	fmt.Println("Saved", img_path+".masked")
}
