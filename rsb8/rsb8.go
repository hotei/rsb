// rsb8.go (c) 2015 David Rook

// Decode() is working for type 8/5650
// still testing type 8/8888
package rsb8

import (
	"fmt"
	"image"
	"io"
	"log"

	"github.com/hotei/rsb/rsbcomn"
)

const (
	MagicNumber = 8
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
	Verbose.Printf("DecodeConfig: %d colorsPallete, width %d, height %d\n")
	return image.Config{ColorModel: nil, Width: int(bf.Width), Height: int(bf.Height)}, nil
}

func Decode(r io.Reader) (pic image.Image, err error) {
	var inputType uint32
	img, err := rsbcomn.ReadRSB(r)
	if err != nil {
		Verbose.Printf("rsbcomn.ReadRSB: failed with %v\n", err)
		return pic, err
	} else {
		Verbose.Printf("rsbcomn.ReadRSB succeeded\n")
	}
	Verbose = true
	Verbose.Printf("FileType[%d] R Depth[%d] G Depth[%d] B Depth[%d] Alpha Depth [%d]\n",
		img.FileType, img.RedDepth, img.GreenDepth, img.BlueDepth, img.AlphaDepth)
	inputType = img.RedDepth*1000 + img.GreenDepth*100 + img.BlueDepth*10 + img.AlphaDepth
	Verbose.Printf("Input type is 8-%d\n", inputType)
	rgba := image.NewRGBA(image.Rect(0, 0, int(img.Width), int(img.Height)))
	ib := img.RsbData // give input a shorter handle
	p := 0
	n := 0
	pixelCount := int(img.PixelCount)
	pixelBytes := int(img.PixelBytes)

	if inputType == 5650 {
		Verbose.Printf("rsb decoding type 8-5650\n")
		magicStartPlace := int(img.MagicStart)
		magicEndPlace := int(img.PixelBytes*img.PixelCount) + magicStartPlace
		var v uint16
		for i := magicStartPlace; i < magicEndPlace; i += pixelBytes {
			if n >= pixelCount {
				fmt.Printf("warning:-> early loop break because expected pixelCount exceeded\n")
				break
			}
			if rsbcomn.LittleEndian {
				v = (uint16(ib[i+1]) << 8) | uint16(ib[i])
			} else {
				v = (uint16(ib[i]) << 8) | uint16(ib[i+1])
			}
			r := (v & 0xf800) >> 8 // RED 5 bits
			g := (v & 0x07E0) >> 3 // GREEN 6 bits
			b := (v & 0x001f) << 3 // BLUE 5 bits
			if false {             // sanity check during debug phase
				if r > 0xff {
					fmt.Printf("%d %d: r[%d] g[%d] b[%d]\n", n, i, r, g, b)
					log.Panicf("bad computer\n")
				}
				if g > 0xff {
					fmt.Printf("%d %d: r[%d] g[%d] b[%d]\n", n, i, r, g, b)
					log.Panicf("bad computer\n")
				}
				if b > 0xff {
					fmt.Printf("%d %d: r[%d] g[%d] b[%d]\n", n, i, r, g, b)
					log.Panicf("bad computer\n")
				}
				fmt.Printf("%d %d: r[%d] g[%d] b[%d]\n", n, i, r, g, b)
			}
			rgba.Pix[p] = uint8(r)
			p++
			rgba.Pix[p] = uint8(g)
			p++
			rgba.Pix[p] = uint8(b)
			p++
			rgba.Pix[p] = 0xff
			p++
			n++
		}
		return rgba, nil
	}

	if inputType == 8888 {
		Verbose.Printf("rsb decoding type 8-8888\n")
		magicStartPlace := int(img.MagicStart)
		magicEndPlace := int(img.PixelBytes*img.PixelCount) + magicStartPlace
		for i := magicStartPlace; i < magicEndPlace; i += pixelBytes {
			if n >= pixelCount {
				fmt.Printf("warning:-> early loop break because expected pixelCount exceeded\n")
				break
			}
			// 1230 is close - works for b&w at least
			rgba.Pix[p] = ib[i+1]
			p++
			rgba.Pix[p] = ib[i+2]
			p++
			rgba.Pix[p] = ib[i+3]
			p++
			rgba.Pix[p] = ib[i+0]
			p++
			n++
		}
		return rgba, nil
	}
	return nil, rsbcomn.ErrFmtUnk
}
