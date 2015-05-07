// rsb4.go (c) 2015 David Rook

// for rsb decoder to work like other image formats the package must expose
// func Decode(r io.Reader) (img image.Image, err error) {}
// and
// func DecodeConfig(r io.Reader) (config image.Config, err error) {}
package rsb4

import (
	"fmt"
	"github.com/hotei/rsb/rsbcomn"
	"image"
	"io"
	"log"
	"os"
)

const (
	MagicNumber = 4
)

func init() {
	rsbcomn.MagicBytes[0] = MagicNumber
	rsbRegisterValue := string(rsbcomn.MagicBytes)

	if len(rsbRegisterValue) != 10 {
		log.Panic("here,len=", len(rsbRegisterValue))
	}
	image.RegisterFormat("rsb", rsbRegisterValue, Decode, DecodeConfig)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Verbose.Printf("init done\n")
}

func DecodeConfig(r io.Reader) (config image.Config, err error) {
	bf, err := rsbcomn.ReadRSB(r) // not efficient but simple wins
	fmt.Printf("DecodeConfig: %d colorsPallete, width %d, height %d\n")
	return image.Config{ColorModel: nil, Width: int(bf.Width), Height: int(bf.Height)}, nil
}

func Decode(r io.Reader) (pic image.Image, err error) {
	var inputType uint32
	img, err := rsbcomn.ReadRSB(r)
	if err != nil {
		Verbose.Printf("rsbcomn.ReadRSB: failed with %v\n", err)
		panic("here")
	} else {
		Verbose.Printf("rsbcomn.ReadRSB succeeded\n")
	}
	if img.FileType != MagicNumber { // can't happen?
		return pic, rsbcomn.ErrBadMagic // return blank image
	}
	fmt.Printf("FileType[%d] R Depth[%d] G Depth[%d] B Depth[%d] Alpha Depth [%d]\n",
		img.FileType, img.RedDepth, img.GreenDepth, img.BlueDepth, img.AlphaDepth)
	inputType = img.RedDepth*1000 + img.GreenDepth*100 + img.BlueDepth*10 + img.AlphaDepth

	Verbose.Printf("Input type is %d\n", inputType)

	if inputType == 4444 {
		rgba := image.NewRGBA(image.Rect(0, 0, int(img.Width), int(img.Height)))
		Verbose.Printf("Converting filetype rsb4/4444 \n")
		if false {
			os.Exit(0)
		}
		ib := img.RsbData // raw rsb data plus some extra we're guessing about
		p := 0
		pixelCount := int(img.PixelCount)
		pixelBytes := int(img.PixelBytes)
		n := 0 // since we cant trust bufLen test to exit on time
		magicStartPlace := int(img.MagicStart)
		magicEndPlace := int(img.PixelBytes*img.PixelCount) + magicStartPlace
		for i := magicStartPlace; i < magicEndPlace; i += pixelBytes {
			// BUG(mdr): 4-4444 add Endian test?
			if n >= pixelCount {
				fmt.Printf("early loop break because pixelCount exceeded\n")
				break
			}
			rgba.Pix[p] = (ib[i] & 0x0f) << 4
			p++
			rgba.Pix[p] = ib[i] & 0xf0
			p++
			rgba.Pix[p] = (ib[i+1] & 0x0f) << 4
			p++
			rgba.Pix[p] = ib[i+1] & 0xf0
			p++
			n++
		}
		return rgba, nil
	}
	return nil, rsbcomn.ErrFmtUnk
}
