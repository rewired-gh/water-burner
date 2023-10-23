package burner

import (
	"image"
	"math"
	"math/rand"
)

func calculateLocalMean(img *hiResImage, i, j, templateSize int, centerWeight float64) (mean hiResRGBA) {
	mean = hiResRGBA{0, 0, 0, 0}
	xStart, xEnd := max(0, i-templateSize), min(img.width-1, i+templateSize)
	yStart, yEnd := max(0, j-templateSize), min(img.height-1, j+templateSize)
	weightSum := 0.0

	for x := xStart; x <= xEnd; x += 1 {
		for y := yStart; y <= yEnd; y += 1 {
			pixel := img.at(x, y)
			weight := 1.0
			if x == i && y == j {
				pixel.scaleScalar(centerWeight)
				weight = centerWeight
			}
			mean.add(&pixel)
			weightSum += weight
		}
	}
	mean.scaleScalar(1 / weightSum)
	return
}

func calculateSimilarity(pixel1, pixel2 hiResRGBA, hSquared float64) (sim hiResRGBA) {
	return op(
		func(a, b float64) float64 {
			return math.Exp2(-math.Abs(a*a-b*b) / hSquared)
		},
		&pixel1,
		&pixel2,
	)
}

func fastNLM(img *hiResImage, h, centerWeight float64, searchSize, templateSize int) (outImg hiResImage) {
	hSquared := h * h
	outImg = makeBlankHiResImage(img.width, img.height)
	nMean := makeBlankHiResImage(img.width, img.height)

	for i := 0; i < img.width; i++ {
		for j := 0; j < img.height; j++ {
			nMean.set(i, j, calculateLocalMean(img, i, j, templateSize, centerWeight))
		}
	}

	for i := 0; i < img.width; i++ {
		for j := 0; j < img.height; j++ {
			pMean := nMean.at(i, j)

			weightSum := hiResRGBA{0, 0, 0, 0}
			weightedSum := hiResRGBA{0, 0, 0, 0}
			xStart, xEnd := max(0, i-searchSize), min(img.width-1, i+searchSize)
			yStart, yEnd := max(0, j-searchSize), min(img.height-1, j+searchSize)
			for x := xStart; x <= xEnd; x += 1 {
				for y := yStart; y <= yEnd; y += 1 {
					qMean := nMean.at(x, y)

					weight := calculateSimilarity(pMean, qMean, hSquared)
					weightSum.add(&weight)
					qMean.scale(&weight)
					weightedSum.add(&qMean)
				}
			}

			filteredPixel := op(
				func(a, b float64) float64 {
					return a / b
				},
				&weightedSum,
				&weightSum,
			)
			outImg.set(i, j, filteredPixel)
		}
	}
	return
}

func addNoise(img *hiResImage, strength, percent float64) {
	for i, pixel := range img.data {
		noiseRand := rand.Float64()
		if noiseRand > percent {
			continue
		}
		rNoise, gNoise, bNoise := rand.Intn(0xffff), rand.Intn(0xffff), rand.Intn(0xffff)
		noise := hiResRGBA{
			float64(rNoise),
			float64(gNoise),
			float64(bNoise),
			0,
		}
		a := pixel.A
		pixel.scaleScalar(1 - strength)
		noise.scaleScalar(strength)
		noise.add(&pixel)
		noise.A = a
		img.data[i] = noise
	}
}

func mapHighLevel(img *hiResImage, high, low float64) {
	area := high - low
	for i, pixel := range img.data {
		mapped := op(func(a, _ float64) float64 {
			if a > high {
				return 0xffff
			} else if a < low {
				return 0x0000
			}
			return (a - low) / area * 0xffff
		}, &pixel, &pixel)
		img.data[i] = mapped
	}
}

// Destroy visible and invisible watermarks in the image, along with parts of image itself
func BurnImage(img image.Image) image.Image {
	imgData := makeHiResImage(img)
	addNoise(&imgData, 0.85, 0.2)
	addNoise(&imgData, 0.15, 0.65)
	mapHighLevel(&imgData, 0xf300, 0x0cff)
	prod := fastNLM(&imgData, 3000, 6, 1, 1)
	return prod.toImage()
}
