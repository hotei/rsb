// rsbcomn.go (c) 2015 David Rook

package rsbcomn

// common to rsb4, rsb5, rsb6, rsb8, rsb9

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type RSBT struct {
	FileType    uint32
	FileSubType uint32
	FileLen     uint32
	PixelCount  uint32
	PixelBits   uint32
	PixelBytes  uint32
	Width       uint32
	Height      uint32
	RedDepth    uint32
	GreenDepth  uint32
	BlueDepth   uint32
	AlphaDepth  uint32
	MagicStart  uint32
	RsbData     []byte
}

type OffsetRSB struct {
	width  int
	height int
	red    int
	green  int
	blue   int
	alpha  int
	magic  int
}

const (
	LittleEndian = true // all testing done to date is with LE
	BigEndian    = !LittleEndian
	MaxHeight    = 65537 // used by sanity check
	MaxWidth     = 65537 // used by sanity check
)

var (
	MagicBytes  = []byte{0, 0, 0, 0, '?', '?', '?', '?', '?', '?'}
	ErrBadMagic = errors.New("rsb: File Id byte is not a valid RSB magic number")
	ErrTooShort = errors.New("rsb: File is too short")
	ErrFmtUnk   = errors.New("rsb: Unknown format")
	ErrBadHdr   = errors.New("rsb: Cant understand header info")
	DepthList   = []byte{4, 5, 6, 8} // R,G,B legal bit depth values
)

var offSetMap map[int]OffsetRSB

var magicMap map[int]int

func init() {
	Verbose = true
	magicMap = make(map[int]int)

	offSetMap = make(map[int]OffsetRSB)
	offSetMap[4] = OffsetRSB{width: 4, height: 8, red: 12, green: 16, blue: 20, alpha: 24}
	magicMap[44444] = 28 //
	offSetMap[5] = OffsetRSB{width: 4, height: 8, red: 12, green: 16, blue: 20, alpha: 24}
	magicMap[58888] = 28 // ok
	offSetMap[6] = OffsetRSB{width: 4, height: 8, red: 12, green: 16, blue: 20, alpha: 24}
	magicMap[64444] = 28 // ok
	magicMap[65650] = 28 // ok
	magicMap[68880] = 28 + 16
	magicMap[68888] = 28
	offSetMap[8] = OffsetRSB{width: 4, height: 8, red: 19, green: 23, blue: 27, alpha: 31}
	magicMap[85650] = 35 // ok
	magicMap[88888] = 35 // ok
	offSetMap[9] = OffsetRSB{width: 4, height: 8, red: 19, green: 23, blue: 27, alpha: 31}
	magicMap[95650] = 43 // ok
	magicMap[98880] = 43 // ok
	magicMap[98888] = 43 // ok
}

// BUG(mdr): is return imgp  *RSBT better choice?
// ReadRSB returns a type that includes header and image raw bytes
func ReadRSB(r io.Reader) (img RSBT, err error) {
	Verbose = true
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Printf("ReadAll failed\n")
		return img, err
	}
	lenRSB := uint32(len(b))
	if lenRSB < 28 {
		return img, ErrTooShort // header is too short
	}
	Verbose.Printf("input length is %d bytes\n", lenRSB)
	// byte order is LittleEndian for some, BigEndian for others???
	// header is 28 bytes at front of file

	img.FileType = LSB32(b[0:4])
	fmt.Printf("File type = %v = %d\n", b[0:4], img.FileType)

	x, ok := offSetMap[int(img.FileType)]
	if !ok {
		log.Panicf("can't find map for file type %d\n", img.FileType)
	}

	img.FileLen = lenRSB
	img.Width = LSB32(b[x.width : x.width+4])
	img.Height = LSB32(b[x.height : x.height+4])
	img.RedDepth = LSB32(b[x.red : x.red+4])
	img.GreenDepth = LSB32(b[x.green : x.green+4])
	img.BlueDepth = LSB32(b[x.blue : x.blue+4])
	img.AlphaDepth = LSB32(b[x.alpha : x.alpha+4])
	img.FileSubType = img.RedDepth*1000 + img.GreenDepth*100 + img.BlueDepth*10 + img.AlphaDepth
	magicVal, ok := magicMap[int(img.FileType*10000+img.FileSubType)]
	if !ok {
		log.Panicf("unk fmt")
		return img, ErrFmtUnk
	}
	img.MagicStart = uint32(magicVal)
	img.PixelBits = img.RedDepth + img.GreenDepth + img.BlueDepth + img.AlphaDepth
	img.PixelBytes = img.PixelBits >> 3 //  same as /= 8
	img.PixelCount = img.Width * img.Height
	fmt.Printf("img = %v\n", img)
	if img.PixelBits <= 0 {
		log.Panicf("pixel bit count can't be zero")
	}
	if (img.PixelBits % 8) != 0 {
		log.Panicf("pixel bit count must be multiple of 8")
	}
	if img.FileLen < ((img.PixelCount * img.PixelBytes) + img.MagicStart) {
		return img, ErrTooShort
	}
	fmt.Printf("first 28 bytes: %v\n", b[:28])
	fmt.Printf("past magic    : %v\n", b[:img.MagicStart])
	fmt.Printf("first 70 bytes: %v\n", b[:70])

	img.RsbData = b

	//var pixelLength = img.RedDepth + img.BlueDepth + img.GreenDepth + img.AlphaDepth
	//var pixelCount = img.Width * img.Height

	extraBytes := img.FileLen - ((img.PixelCount * img.PixelBytes) + img.MagicStart)
	Verbose.Printf("There are %d extra bytes in file (over & above header)\n", extraBytes)
	Verbose.Printf("Exiting readRSB() normally\n")
	return img, nil
}

// convert from LITTLE endian four byte slice to uint32
// reverse function is LSBytesFromUint32
func LSB32(b []byte) uint32 {
	if len(b) != 4 {
		FatalError(fmt.Errorf("mdr: Slice must be exactly 4 bytes\n"))
	}
	var rc uint32
	rc = uint32(b[3])
	rc <<= 8
	rc |= uint32(b[2])
	rc <<= 8
	rc |= uint32(b[1])
	rc <<= 8
	rc |= uint32(b[0])
	return rc
}

// convert from BIG endian four byte slice to uint32
// reverse function is MSBytesFromUint32
func MSB32(b []byte) uint32 {
	if len(b) != 4 {
		FatalError(fmt.Errorf("mdr: Slice must be exactly 4 bytes\n"))
	}
	var rc uint32
	rc = uint32(b[0])
	rc <<= 8
	rc |= uint32(b[1])
	rc <<= 8
	rc |= uint32(b[2])
	rc <<= 8
	rc |= uint32(b[3])
	return rc
}

func GoodDepth(d uint32) bool {
	for _, depth := range DepthList {
		if d == uint32(depth) {
			return true
		}
	}
	return false
}

func FatalError(err error) {
	fmt.Printf("%v\n", err)
	os.Exit(1)
}
