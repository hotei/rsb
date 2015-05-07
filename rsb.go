// rsb.go (c) 2015 David Rook

package rsb

// rsb6, rsb8, rsb9

import (
	"io"

	_ "github.com/hotei/rsb/rsb4"
	_ "github.com/hotei/rsb/rsb5"
	_ "github.com/hotei/rsb/rsb6"
	_ "github.com/hotei/rsb/rsb8"
	_ "github.com/hotei/rsb/rsb9"
	"github.com/hotei/rsb/rsbcomn"
)

func ReadRSB(r io.Reader) (img rsbcomn.RSBT, err error) {
	return rsbcomn.ReadRSB(r)
}
