package tools

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
)

var fontSize float64 = 26

func AddDateToThumbnail(r io.Reader, label string, x int, y int) (io.Reader, error) {
	img, err := png.Decode(r)
	if err != nil {
		log.Printf("Error converting Reader image to Image %v", err)
		return nil, err
	}
	// Read the font data.
	fontBytes, err := ioutil.ReadFile("Fonts/luximr.ttf")
	if err != nil {
		log.Printf("Error reading font file %v", err)
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Printf("Error parsing font %v", err)
		return nil, err
	}
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img.(*image.RGBA))
	c.SetSrc(image.Black)

	pt := freetype.Pt(70, 10+int(c.PointToFixed(fontSize)>>6))
	_, err = c.DrawString(label, pt)
	if err != nil {
		log.Printf("Error drawing text on iamge %v", err)
		return nil, err
	}

	tempFilename := fmt.Sprintf("Images/Temp/csm_%s.png", label)

	file, err := os.Create(tempFilename)
	if err != nil {
		log.Printf("Error creating temporary png file %v", err)
		return nil, err
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		log.Printf("Error encoding to temporary png file %v", err)
		return nil, err
	}

	png, err := os.Open(tempFilename)
	if err != nil {
		log.Printf("Error openingtemporary png %v", err)
		return nil, err
	}
	return png, nil
}
