// showRSB.go (c) 2015 David Rook

// NOTE:
// * created windows might be too small to hold decorations for close.
//     if so use key comb ALT+F4 to close the active window
// * This is true bare bones viewer.  You can't adjust anything, including size.
// * the test files have been created with a one pixel black bar at the right margin
//     if the bar isn't there then the decode failed.  Not obvious in these
//     small windows, but clear enough using rsbWeb program since thumbs are
//     larger than original.
// * not all file formats are represented here since PhotoShop plugin can't create
//     all of them

// usage: 
//  go build
//  cd ok
//  ../showRSB

package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path"

	// Alien
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"

	_ "github.com/hotei/rsb"
)

func showOne(fname string) {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	Verbose.Printf("testprog: working on %s\n", fname)
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	ximg := xgraphics.NewConvert(X, img)
	ximg.XShowExtra(fname, true)

	xevent.Main(X)
	xevent.Quit(X)
}

func main() {
	Verbose = true
	fmt.Printf("Starting testprog\n")

	dirlist, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Printf("error %v\n", err)
		os.Exit(1)
	}
	for _, fi := range dirlist { // fileinfo items
		fn := fi.Name()
		fmt.Printf("%v\n", fn)
	ext := path.Ext(fn)
		if  ext == ".rsb" ||ext ==  ".RSB" {
			fmt.Printf("\n\n\n Found rsb file called %s\n", fn)
			showOne(fn)
		}
	}
}
