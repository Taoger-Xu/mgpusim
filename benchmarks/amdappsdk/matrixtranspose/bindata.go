// Code generated by "esc -o=bindata.go -pkg=matrixtranspose -private kernels.hsaco"; DO NOT EDIT.

package matrixtranspose

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/kernels.hsaco": {
		name:    "kernels.hsaco",
		local:   "kernels.hsaco",
		size:    9304,
		modtime: 1562383959,
		compressed: `
H4sIAAAAAAAC/+xaz2/b1h3/8r0nimJkiXG1NMu8lBMEewhiwXlRDS2XxY0XJ6iTOHWbH1sLg5EoiYtE
ChSVWcVGO4bjw2Bsw047FOg/sL9h9q7bpfUphxx66bXoYdtpmIZHPkqkZtbumqGDoy9gf8jv7+/3PUl8
fG/jJ8vXkSBcBU4YPgOBXSj+fSB48AMfL3i8MkhwFdIggwgAJKQ3ivtCFCXOF7hdHL2ZjWKQD7NLhO5H
8W84imE7liuonD+CbYhiYIe+pl1Q3zufO1VyDLtwfozufu5URfj6RIJ+IhgmHsIv01EkITuJx1+4teip
B2PzPW8++HwCyUFtAW/h1uLSynu+Lpsepzhfa1XrFXNWa1XZX6OjMagbH1aas/XaennuMvffSAPI3GZ2
dla+p9sdwzKvqAH9TL10UZ1TP5Df1m1Tb3aGEllVZ9XbWksfslRVbWmObay/a2tmp211dJnxVnutR1Yz
pDozonX1cXXG01zWzHpXqw9d3mnr5rVl9VpEOkjSS46qH3jSBbveiaQi+3BIklbXaXcdObh9t9fWIzoz
taalOaULMwOVVePDqIvyQLTQNOrmlUNF97RmV3/bMKuBeKlpPdKab3VrNd2OarEcAq3rl+nQe7Vqr7a1
in63qzWvDFwM5ZVKIPFpUa9p3aYTX7xhvrq1P2palcffqPZSfO2l+NoXe6bWMiqrDc3WqyuWYTrH7YGv
rA+DXZqP78+yVflG7fmFcd+oOo34BnUN03nZvXmr57Hi2/FeeEr8FzXd0I16wzlRRZnd1ppVW/Omc2dt
/UTVVretbnttfc2q1Tq6cwJL652A0l7Sb8INo1rVTf9r/Y7XlAfxCd6cL/3v4z/8luP/9FuIf9syv2Je
3Cwf8wd5nNVxs7pmVfUV22oPnle5FXvG1uz6ql5v6abj510OXC6xLw8uum6s61VfPsfFK7bxRHP0eIWo
c1554Py+9kSv2VYQVFUHc+12t7W6tPLO8NH6UnkouReR0Mtccktbv97UnPuW/djP2nNK35yXi8WiHLd+
EkJ/52JkJ5EEvkQUvNUdPrLQZ/BH+F2KrfWirXwQup6E16NrU0LEfr/f/3+sH2N5X/Eql/fZ0nwDdve2
AO2z9SkW5H2ms4l+6GbTxJUAnjI56qMdVk1KEJ6iJeICkn4vAH6KkKT0CZEBNu+gLHG3JKlABIluS1JB
hGcH6A2AJOx4mJbOUrSd+s1mWikB/OWTzQwCEeBT8TTDnQMxgUBCIk2dPUcxPDuYEBGIiQwVE2QeAXyK
sggQ85UE2BLFQhLJbhIpLgrioJyLIK8kkepiIpWyUno+K2UowF8/yUoIMPnoV/DzzR3mi0wgSCYIZJht
AiB17jzdyucLf4Z//ikjEIqShH5MxDyG3x5sEwQTZJqKIqGJlFgGeP8FAkAik00jmJiepmI6TaVspiyF
eYpC5dcmyzLnibkcVV4/U87gKXoaI7qNUCGN8jR1XqUsj8wUYjOSpvIFujU1VYDvA83ANG3D7h4IHz1P
IzZC778QAYSANyFgjyeFeFnOkwGEVCFHU7nzAx+KgL3xBvj13seA8sByAwQJRFwCFygIXzyXAU4hfLHM
rk97H5gvnk8CEIZZgOwG3t07hXPuBtrdA6y4G8LuHsGTLvObxWdcgF+++I4EkEDgBv1Q8BmXKBcovHax
nJmepgJWXOCyBJ500zjnJpDkEgAa2BOMyggr7ml8xhXxpDuBcy6zZXL/xQy393IHijAqY3zeTWLVzeC8
lwd479Ke7sGYxjSmMY1pTGMa05jGNKZXkoK95i8zPp7i92c5JjiW+D68DNH96r//q295+97Z6L5yIXt4
vGXDfKzbV9Tl5UW1VJyDSlMz6+oTf1NVLRfn1KPfV8hI/I/9fm9txhbDo/ZF03J0KFZ7ZqfXgmLd7BYb
WqcB/D/jOzYUHX3d8e60llGBYsVqtXTTgWKn13K0R1DsNDqO7V/5CCPbyLC2+PD2wq2b117ee5lkaBs/
bj8/vP8fpiQfSzQyvgGGx1cInVsI+Gz4/tHvW4F9ML4BTo2kJY3E/y73jUbmQ4DKiD0ZwTf4OQM0Mv8C
FA7pV5hmwmc+vuKcSJyDWW47UIs5v5EYqT8IM89dzo2EaXP7Qkz4AH8cHvsQ/YHb34Dh5xUfMn5L4dxD
9Bm3f3hE/+7G2P+InwM6d4T9vwMAAP//0BJj1VgkAAA=
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
