// rsb6.go (c) 2015 David Rook

// for decoder to work like other image formats the package must expose
// func Decode(r io.Reader) (img image.Image, err error) {}
// and
// func DecodeConfig(r io.Reader) (config image.Config, err error) {}
package rsb6

import (
	"fmt"
	"image"
	"io"
	"log"
	"os"

	"github.com/hotei/rsb/rsbcomn"
)

const (
	MagicNumber = 6
)

func init() {
	Verbose = true
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
	// fmt.Printf("DecodeConfig: %d colorsPallete, width %d, height %d\n")
	return image.Config{ColorModel: nil, Width: int(bf.Width), Height: int(bf.Height)}, nil
}

func Decode(r io.Reader) (pic image.Image, err error) {
	var inputType uint32
	img, err := rsbcomn.ReadRSB(r)
	if err != nil {
		return pic, err
	} else {
		Verbose.Printf("rsbcomn.ReadRSB succeeded\n")
	}

	Verbose.Printf("FileType[%d] R Depth[%d] G Depth[%d] B Depth[%d] Alpha Depth [%d]\n",
		img.FileType, img.RedDepth, img.GreenDepth, img.BlueDepth, img.AlphaDepth)
	inputType = img.RedDepth*1000 + img.GreenDepth*100 + img.BlueDepth*10 + img.AlphaDepth
	Verbose.Printf("Input type is %d\n", inputType)
	rgba := image.NewRGBA(image.Rect(0, 0, int(img.Width), int(img.Height)))
	ib := img.RsbData
	p := 0
	n := 0
	pixelCount := int(img.PixelCount)
	pixelBytes := int(img.PixelBytes)

	if inputType == 4444 {
		Verbose.Printf("Converting type 6-4444\n")
		if false {
			os.Exit(0)
		}

		magicStartPlace := int(img.MagicStart)
		magicEndPlace := int(img.PixelBytes*img.PixelCount) + magicStartPlace
		for i := magicStartPlace; i < magicEndPlace; i += pixelBytes {
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

	if inputType == 5650 { // working on village-map.rsb
		Verbose.Printf("Converting type 6-5650\n")
		if false {
			os.Exit(0)
		}
		var v uint16
		magicStartPlace := int(img.MagicStart)
		magicEndPlace := int(img.PixelBytes*img.PixelCount) + magicStartPlace
		for i := magicStartPlace; i < magicEndPlace; i += pixelBytes {
			if n >= pixelCount {
				fmt.Printf("early loop break because expected pixelCount exceeded\n")
				break
			}
			if rsbcomn.LittleEndian {
				v = (uint16(ib[i+1]) << 8) + uint16(ib[i])
			} else {
				v = (uint16(ib[i]) << 8) + uint16(ib[i+1])
			}
			rgba.Pix[p] = uint8((v & 0xf800) >> 8)
			p++
			rgba.Pix[p] = uint8((v & 0x07E0) >> 3)
			p++
			rgba.Pix[p] = uint8((v & 0x001f) << 3)
			p++
			rgba.Pix[p] = 0xff // implied by format type
			p++
			n++
		}
		return rgba, nil
	}

	if inputType == 8880 {
		Verbose.Printf("Converting filetype 6-8880\n")
		if false {
			os.Exit(0)
		}
		magicStartPlace := int(img.MagicStart)
		magicEndPlace := int(img.PixelBytes*img.PixelCount) + magicStartPlace
		for i := magicStartPlace; i < magicEndPlace; i += pixelBytes {
			if n >= pixelCount {
				fmt.Printf("early loop break because pixelCount exceeded\n")
				break
			}
			rgba.Pix[p] = ib[i+2]
			p++
			rgba.Pix[p] = ib[i]
			p++
			rgba.Pix[p] = ib[i+1]
			p++
			rgba.Pix[p] = 0xff // alpha value is implied by format type
			p++
			n++
		}
		return rgba, nil
	}

	if inputType == 8888 {
		Verbose.Printf("Converting filetype 6-8888\n")
		if false {
			os.Exit(0)
		}
		magicStartPlace := int(img.MagicStart)
		magicEndPlace := int(img.PixelBytes*img.PixelCount) + magicStartPlace
		for i := magicStartPlace; i < magicEndPlace; i += pixelBytes {
			if n >= int(pixelCount) {
				fmt.Printf("early loop break because pixelCount exceeded\n")
				break
			}
			rgba.Pix[p] = ib[i+3]
			p++
			rgba.Pix[p] = ib[i+2]
			p++
			rgba.Pix[p] = ib[i+1]
			p++
			rgba.Pix[p] = ib[i]
			p++
			n++
		}
		return rgba, nil
	}
	return nil, rsbcomn.ErrFmtUnk

}
