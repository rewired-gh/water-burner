package burner

import (
	"image"
	"image/color"
	"image/draw"
)

type hiResRGBA struct {
	R float64
	G float64
	B float64
	A float64
}

func (rgba *hiResRGBA) add(other *hiResRGBA) {
	rgba.R += other.R
	rgba.G += other.G
	rgba.B += other.B
	rgba.A += other.A
}

func (rgba *hiResRGBA) scale(other *hiResRGBA) {
	rgba.R *= other.R
	rgba.G *= other.G
	rgba.B *= other.B
	rgba.A *= other.A
}

func (rgba *hiResRGBA) scaleScalar(other float64) {
	rgba.R *= other
	rgba.G *= other
	rgba.B *= other
	rgba.A *= other
}

func op(fn func(float64, float64) float64, lhs, rhs *hiResRGBA) hiResRGBA {
	return hiResRGBA{
		fn(lhs.R, rhs.R),
		fn(lhs.G, rhs.G),
		fn(lhs.B, rhs.B),
		fn(lhs.A, rhs.A),
	}
}

type hiResImage struct {
	data   []hiResRGBA
	width  int
	height int
}

func makeHiResImage(img image.Image) (outImg hiResImage) {
	bounds := img.Bounds()
	outImg.width = bounds.Max.X - bounds.Min.X
	outImg.height = bounds.Max.Y - bounds.Min.Y
	outImg.data = make([]hiResRGBA, outImg.width*outImg.height)
	for i := bounds.Min.X; i < bounds.Max.X; i++ {
		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
			r, g, b, a := img.At(i, j).RGBA()
			outImg.data[j*outImg.width+i] = hiResRGBA{
				float64(r),
				float64(g),
				float64(b),
				float64(a),
			}
		}
	}
	return
}

func makeBlankHiResImage(width, height int) (outImg hiResImage) {
	outImg.width = width
	outImg.height = height
	outImg.data = make([]hiResRGBA, width*height)
	return
}

func (img hiResImage) toImage() (outImg draw.Image) {
	outImg = draw.Image(image.NewRGBA(image.Rect(0, 0, img.width, img.height)))
	for i := 0; i < img.width; i++ {
		for j := 0; j < img.height; j++ {
			pixel := img.at(i, j)
			outImg.Set(i, j, color.RGBA{
				uint8(uint32(pixel.R) >> 8),
				uint8(uint32(pixel.G) >> 8),
				uint8(uint32(pixel.B) >> 8),
				uint8(uint32(pixel.A) >> 8),
			})
		}
	}
	return
}

func (img hiResImage) set(x, y int, pixel hiResRGBA) {
	img.data[y*img.width+x] = pixel
}

func (img hiResImage) at(x, y int) hiResRGBA {
	return img.data[y*img.width+x]
}
