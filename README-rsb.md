<center>
# Red Storm Bitmap - rsb package
</center>


## OVERVIEW

This package decodes rsb image files to RGBA image format.  
The rsbWeb example uses it to display them in a web page.  
The showRSB example is a trivial X-window viewer for rsb images.

## Why?

Ghost Recon is a game (c) 2001 Red Storm Entertainment (RSE) and
published by Ubisoft. The game textures
are stored in a proprietary format with extension of ".rsb".
RSB format was never intended for use outside the game.
It isn't compressed and doesn't have "standard" magic numbers that
distinguish the more universally used formats
like gif, jpg and png.
<p>
Since the game is 14 years old the obvious question is WHY?  Mainly
because the game designers allowed users to make their own missions, staged
on the maps provided on the original disks or maps created by people with
programs like 3DMax (later Autodesk Maya). Designing scenarios was half the
fun of the game and RSE made it relatively easy. User groups created websites
to distribute their work and the game stayed active far longer than expected.
<p>
RSE provided a plugin for
PhotoShop so that you can create or edit textures in rsb format.  However,
with so many textures used in the game it becomes difficult to track them. I
wanted to be able to scan a directory (of faces for example) and pick out the
one that I wanted for the scenario.  Typically this might mean loading 50 or more
images into PhotoShop - a slow process.  With rsbWeb (one of the included examples)
I can easily create a web page that has all the faces shown along with the name of the
(sub)directory where I can find it.  <p>
As another example, of all the textures on a disk only a small
percent of them are useful in a user-created scenario.  I can scan a large 
volume of rsbs and copy the useful items into one directory that I then pick from.
It makes a very tedious operation fairly painless.



### Installation

If you have a working go installation on a Unix-like OS:

> ```go get github.com/hotei/rsb```

Will copy github.com/hotei/program to the first entry of your $GOPATH

or if go is not installed yet :

> ```cd DestinationDirectory```

> ```git clone https://github.com/hotei/rsb.git```

Comments can be sent to <hotei1352@gmail.com>.  Issues - should there be any - can be 
registered at _github.com/hotei/rsb_.

### Features

* Registers 5 versions of rsb as denoted by 4/5/6/8/9 as first byte of file.
* Decodes to image.image with same size as the rsb.

### Limitations

* <font color="red">Not all formats have been tested.</font>  There is a list in rsbcomn.go
that shows the formats at their testing/ok status.  If you have a file that doesn't
decode properly please file an issue and make a copy of the file available by a 
link if possible.  If it's a file on the Ghost Recon CD just provide the file name.
* <font color="red">Only tested with LittleEndian hardware(Intel/AMD 64 bit) so far - and only a few
feeble attempts to anticipate problems with BE.  </font>

### Usage

Typical usage is in a program that uses images in the go "image" format. This
package registers an _rsb_ image type that can then be used like other built-in
types (gif,jpg,png etc).  The following code is typical, but certainly not the
only way it can be used:

``` go

	import _ "github.com/hotei/rsb"

	ior,err := os.Open("somefile.rsb")
	img, _, err := image.Decode(ior)
```

* Examples
  * showRSB - X-11 viewer for rsb files
  * rsbWeb - walks a directory and displays all rsb files found on one page

### BUGS

* None known as of 2015-05-07

### To-Do

* Essential:
  * TBD
* Nice:
  * TBD
* Nice but no immediate need:
  * TBD

### Change Log

* 2015-05-07 initial version working for all available formats
* 2015-05-01 Started
  * compiled with go 1.4.2
 
### Resources

* I can find no resources on the www that document the rsb format. So I resorted to
[od] [4] along with the Photoshop Plugin _rsb8_ provided on the Ghost Recon disk.
* [go language reference] [1] 
* [go standard library package docs] [2]
* Source for [package rsb] [3]
* [Ghost Recon]  [5] wikipedia page

[1]: http://golang.org/ref/spec/ "go reference spec"
[2]: http://golang.org/pkg/ "go stdlib package docs"
[3]: http://github.com/hotei/rsb "github.com/hotei/rsb"
[4]: http://www.freebsd.org/cgi/man.cgi?query=od&apropos=0&sektion=0&manpath=FreeBSD+10.1-RELEASE&arch=default&format=html "od"
[5]: http://en.wikipedia.org/wiki/Tom_Clancy%27s_Ghost_Recon "Ghost Recon"

Comments can be sent to <hotei1352@gmail.com> or to user "hotei" at github.com.
License is BSD-two-clause, in file "LICENSE"

License
-------
The 'program' go package/program is distributed under the Simplified BSD License:

> Copyright (c) 2015 David Rook. All rights reserved.
> 
> Redistribution and use in source and binary forms, with or without modification, are
> permitted provided that the following conditions are met:
> 
>    1. Redistributions of source code must retain the above copyright notice, this list of
>       conditions and the following disclaimer.
> 
>    2. Redistributions in binary form must reproduce the above copyright notice, this list
>       of conditions and the following disclaimer in the documentation and/or other materials
>       provided with the distribution.
> 
> THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDER ``AS IS'' AND ANY EXPRESS OR IMPLIED
> WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
> FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> OR
> CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
> CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
> SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
> ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
> NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
> ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

Documentation (c) 2015 David Rook 

// EOF README.md  (this is a markdown document and tested OK with blackfriday)