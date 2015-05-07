// rsb9.go (c) 2015 David Rook

// for decoder to work like other image formats the package must expose
// func Decode(r io.Reader) (img image.Image, err error) {}
// and
// func DecodeConfig(r io.Reader) (config image.Config, err error) {}
package rsb9

import (
	"fmt"
	"image"
	"io"
	"log"

	"github.com/hotei/rsb/rsbcomn"
)

const (
	MagicNumber = 9
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
	fmt.Printf("DecodeConfig: %d colorsPallete, width %d, height %d\n")
	//return image.Config{ColorModel: bf.AColors, Width: int(bf.Width), Height: int(bf.Height)}, nil
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
	inputType = img.FileSubType
	Verbose.Printf("Input type is 9-%d\n", inputType)
	rgba := image.NewRGBA(image.Rect(0, 0, int(img.Width), int(img.Height)))
	ib := img.RsbData
	p := 0
	n := 0
	pixelCount := int(img.PixelCount)
	pixelBytes := int(img.PixelBytes)

	if inputType == 5650 { // working
		Verbose.Printf("rsb decoding type 9-5650\n")
		magicStartPlace := int(img.MagicStart)
		magicEndPlace := int(img.PixelBytes*img.PixelCount) + magicStartPlace
		for i := magicStartPlace; i < magicEndPlace; i += pixelBytes {
			if n >= pixelCount {
				fmt.Printf("early loop break because expected pixelCount exceeded\n")
				break
			}
			var v uint16
			if rsbcomn.LittleEndian {
				v = (uint16(ib[i+1]) << 8) + uint16(ib[i])
			} else {
				v = (uint16(ib[i]) << 8) + uint16(ib[i+1])
			}
			r := (v & 0xf800) >> 8 // RED 5 bits
			g := (v & 0x07E0) >> 3 // GREEN 6 bits
			b := (v & 0x001f) << 3 // BLUE 5 bits
			rgba.Pix[p] = uint8(r)
			p++
			rgba.Pix[p] = uint8(g)
			p++
			rgba.Pix[p] = uint8(b)
			p++
			rgba.Pix[p] = 0xff // alpha value is implied by format type
			p++
			n++
		}
		return rgba, nil
	}

	if inputType == 8880 { // working
		Verbose.Printf("rsb decoding type 9-8880\n")
		magicStartPlace := int(img.MagicStart)
		magicEndPlace := int(img.PixelBytes*img.PixelCount) + magicStartPlace
		for i := magicStartPlace; i < magicEndPlace; i += pixelBytes {
			if n >= pixelCount {
				fmt.Printf("early loop break because pixelCount exceeded\n")
				break
			}
			rgba.Pix[p] = ib[i+0] // r
			p++
			rgba.Pix[p] = ib[i+1] // g
			p++
			rgba.Pix[p] = ib[i+2] // b
			p++
			rgba.Pix[p] = 0xff // alpha value is implied by format type
			p++
			n++
		}
		return rgba, nil
	}

	if inputType == 8888 { // working for 98888gun.rsb
		Verbose.Printf("rsb decoding type 9-8888 \n")
		magicStartPlace := int(img.MagicStart)
		magicEndPlace := int(img.PixelBytes*img.PixelCount) + magicStartPlace
		for i := magicStartPlace; i < magicEndPlace; i += pixelBytes {
			if n >= pixelCount {
				fmt.Printf("early loop break because pixelCount exceeded\n")
				break
			}
			rgba.Pix[p] = ib[i+1]
			p++
			rgba.Pix[p] = ib[i+2]
			p++
			rgba.Pix[p] = ib[i+3]
			p++
			rgba.Pix[p] = ib[i+0] // right
			p++
			n++
		}
		return rgba, nil
	}

	return nil, rsbcomn.ErrFmtUnk
}
