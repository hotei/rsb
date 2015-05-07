// rsbWeb.go (c) 2015 David Rook

// 2015-05-06 added more flavors (4&5) of rsb
// 2015-05-05 working on first cut

// runs webserver on /www/rsb directory on your local disk, serving at localhost
package main

// working
import (
	// go 1.X only below here
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/hotei/rsb"
	"github.com/hotei/rsb/rsbcomn"
)

const (
	hostIPstr  = "127.0.0.1"
	portNum    = 8284
	serverRoot = "/www/"

	rsbURL = "/rsb/"
)

var (
	portNumString = fmt.Sprintf(":%d", portNum)
	listenOnPort  = hostIPstr + portNumString
	g_fileNames   []string
	myRsbDir      = []byte{}
)

func checkRsbName(pathname string, info os.FileInfo, err error) error {
	fmt.Printf("checking if rsb %s\n", pathname)
	if info == nil {
		fmt.Printf("WARNING --->  no stat info: %s\n", pathname)
		os.Exit(1)
	}
	if info.IsDir() {
		// dont create line for dirs
	} else { // regular file
		//fmt.Printf("found %s %s\n", pathname, filepath.Ext(pathname))
		ext := filepath.Ext(pathname)
		if  ext == ".rsb" || ext == ".RSB" {
			g_fileNames = append(g_fileNames, pathname)
		}
	}
	return nil
}

// create html to show thumbnail and text link to original image
func makeRsbLine(s string) []byte {
	//return []byte(fmt.Sprintf("<a href=\"%s\">View %s</a><br>\n",s,s))
	workDir := serverRoot + rsbURL[1:]
	s = s[len(workDir):]
	fname := serverRoot + rsbURL + s
	fmt.Printf("makeRsbLine: using %s\n", fname)
	f, err := os.Open(fname)
	defer f.Close()
	if err != nil {
		return []byte(fmt.Sprintf("makeRsbfailed for %s<br>\n", fname))
	}
	r, err := rsbcomn.ReadRSB(f)
	if err != nil {
		return []byte(fmt.Sprintf("makeRsbfailed for %s<br>\n", fname))
	}
	ftype := fmt.Sprintf("%d-%d", r.FileType, r.FileSubType)
	x := []byte(fmt.Sprintf("<img src=\"%s\" height=100 width=150> <a href= \"%s\"> %s %s</a><br>\n", rsbURL+s, rsbURL+s, s, ftype))
	fmt.Printf("%s\n", x)
	return x
}

func init() {
	fmt.Printf("Starting init()\n")

	pathName := serverRoot + rsbURL[1:]
	g_fileNames = make([]string, 0, 20)
	myRsbDir = []byte(`<html><!-- comment --><head><title>Test the rsb package</title></head><body>click on image to see in original size<br>`) // {}
	stats, err := os.Stat(pathName)
	if err != nil {
		fmt.Printf("Can't get fileinfo for %s\n", pathName)
		os.Exit(1)
	}
	if stats.IsDir() {
		filepath.Walk(pathName, checkRsbName)
	} else {
		fmt.Printf("this argument must be a directory (but %s isn't)\n", pathName)
		os.Exit(-1)
	}
	//fmt.Printf("g_fileNames = %v\n", g_fileNames)
	for _, val := range g_fileNames {
		fmt.Printf("%v\n", val)
		line := makeRsbLine(val)
		myRsbDir = append(myRsbDir, line...)
	}
	t := []byte(`</body></html>`)
	myRsbDir = append(myRsbDir, t...)
}

/*
// just hand the raw file over to the browser
func rawWriteOut(fileName string, w http.ResponseWriter) {
	rawBuf, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("rawWriteOut: cant open file %s\n", fileName)
		return
	}
	w.Write(rawBuf)
}
*/

func rsbWriteOut(imageName string, w http.ResponseWriter) {
	fmt.Printf("rsbWriteOut: imageName = %s\n", imageName)
	bf, err := os.Open(imageName)
	if err != nil {
		fmt.Printf("rsbWriteOut: cant open rsb %s\n", imageName)
		return
	}
	img, _, err := image.Decode(bf)
	if err != nil {
		fmt.Printf("rsbWriteOut: rsb decode failed for %s\n", imageName)
		w.Write([]byte(fmt.Sprintf("Decode failed for %s\n", imageName)))
		return
	}
	b := make([]byte, 0, 10000)	// hint for buffer size
	wo := bytes.NewBuffer(b) 
	err = png.Encode(wo, img) // encode it as something the browser understands
	if err != nil {
		fmt.Printf("rsbWriteOut: png encode failed for %s\n", imageName)
		return
	}
	w.Write(wo.Bytes())
}

func rsbHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("rsbHandler: r.URL.Path %q\n", r.URL.Path)
	if r.URL.Path == rsbURL {
		w.Write(myRsbDir)
		return
	}
	fileName := serverRoot + r.URL.Path[1:]
	fmt.Printf("rsbHandler: fname = %s\n", fileName)
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == ".rsb" || ext == ".RSB" {
		rsbWriteOut(fileName, w)
		return
	}
}

func main() {
	http.HandleFunc(rsbURL, rsbHandler)

	log.Printf("load rsb urls with %s%s\n", listenOnPort, rsbURL)
	log.Printf("rsbWeb is ready to serve at %s\n", listenOnPort)
	err := http.ListenAndServe(listenOnPort, nil)
	if err != nil {
		log.Printf("rsbWeb: error running webserver %v", err)
	}
}
