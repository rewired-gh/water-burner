package burner_test

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/rewired-gh/water-burner/pkg/burner"
)

func TestProcessImage(t *testing.T) {
	fileStrs := []string{"../assets/monet_encoded_freq", "../assets/monet_encoded_gan"}

	for _, fileStr := range fileStrs {
		imgFile, err := os.Open(fmt.Sprintf("%s.png", fileStr))
		if err != nil {
			t.Error(err)
		}
		defer imgFile.Close()

		img, err := png.Decode(imgFile)
		if err != nil {
			t.Error(err)
		}

		outImg := burner.BurnImage(img)

		outputFile, err := os.Create(fmt.Sprintf("%s_burnt.png", fileStr))
		if err != nil {
			t.Error(err)
		}
		defer outputFile.Close()

		err = png.Encode(outputFile, outImg)
		if err != nil {
			t.Error(err)
		}
	}
}
