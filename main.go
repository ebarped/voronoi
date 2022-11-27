package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

const (
	imagePath   = "voronoi.png"
	imageWidth  = 800
	imageHeight = 600
	seedCount   = 20 // number of points in the image with an associated color
	seedRadius  = 3  // number of pixels of the radius
)

type Circle struct {
	p image.Point
	r int
}

// IsInside returns true if point p is inside the circle
func (c Circle) IsInside(p image.Point) bool {
	dx := float64(c.p.X - p.X)
	dy := float64(c.p.Y - p.Y)
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist <= float64(c.r) {
		return true
	}
	return false
}

// Dist returns the distance between the point p and the circle
func (c Circle) Dist(p image.Point) float64 {
	// distance to the center - radius

	dx := float64(c.p.X - p.X)
	dy := float64(c.p.Y - p.Y)
	dist := math.Sqrt(dx*dx+dy*dy) - float64(c.r)
	return dist
}

// Draw draws a circle in the img with the specified color
func (c Circle) Draw(img image.RGBA, color color.Color) {
	pUpLeft := image.Point{
		c.p.X - c.r,
		c.p.Y - c.r,
	}
	pDownRight := image.Point{
		c.p.X + c.r,
		c.p.Y + c.r,
	}

	// we create a rectagle around the central point of the
	// circle so we only check these pixels instead
	// checking the pixels of all the image
	outerRectangle := image.Rectangle{
		pUpLeft,
		pDownRight,
	}

	for x := outerRectangle.Min.X; x <= outerRectangle.Max.X; x++ {
		for y := outerRectangle.Min.Y; y <= outerRectangle.Max.Y; y++ {
			pixel := image.Point{x, y}
			if c.IsInside(pixel) {
				img.Set(pixel.X, pixel.Y, color)
			}
		}

	}

}

type Seed struct {
	Circle
	color color.Color
}

var (
	COLOR_BACKGROUND = color.RGBA{18, 18, 18, 255}
	COLOR_SEED       = color.Black
)

var colorList = []color.RGBA{
	{0, 255, 255, 255},   // "Aqua"
	{0, 0, 255, 255},     // "Blue"
	{255, 0, 255, 255},   // "Fuchsia"
	{0, 255, 0, 255},     // "Lime"
	{255, 255, 0, 255},   // "Yellow"
	{0, 128, 128, 255},   // "Teal"
	{192, 192, 192, 255}, // "Silver"
	{128, 128, 128, 255}, // "Gray"
	{0, 128, 0, 255},     // "Green"
	{0, 0, 128, 255},     // "Navy"
	{128, 0, 0, 255},     // "Maroon"
	{128, 128, 0, 255},   // "Olive"
	{128, 0, 128, 255},   // "Purple"
	{255, 0, 0, 255},     // "Red"
	{255, 255, 255, 255}, // "White"
}

func fillImage(img image.RGBA, c color.Color) {
	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			img.Set(x, y, c)
		}
	}
}

func getClosestSeed(p image.Point, seeds []Seed) Seed {
	currentDist := math.MaxFloat64
	closestSeed := seeds[0]
	for _, seed := range seeds {
		dist := seed.Dist(p)
		if dist < currentDist {
			closestSeed = seed
			currentDist = dist
		}
	}
	return closestSeed
}

func main() {
	// define the size of the image
	upLeft := image.Point{0, 0}
	lowRight := image.Point{imageWidth, imageHeight}

	// create new png image
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// fill image with background color
	fillImage(*img, COLOR_BACKGROUND)

	// generate the seed points
	randomSeed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randomSeed)
	seeds := make([]Seed, seedCount)
	for i := 0; i < seedCount; i++ {
		c := Circle{
			image.Point{
				random.Intn(imageWidth + 1),
				random.Intn(imageHeight + 1),
			},
			seedRadius,
		}
		seeds[i] = Seed{c, colorList[i%len(colorList)]}
	}

	// for each pixel in the image, get the closest seed
	// and set its color
	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			p := image.Point{x, y}
			seed := getClosestSeed(p, seeds)
			img.Set(x, y, seed.color)
		}
	}

	// draw the seeds
	for _, seed := range seeds {
		seed.Circle.Draw(*img, COLOR_SEED)
	}

	// save image as png
	f, _ := os.Create(imagePath)
	png.Encode(f, img)
}
