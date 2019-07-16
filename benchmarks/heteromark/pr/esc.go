// Code generated by "esc -o esc.go -pkg=pagerank kernels.hsaco"; DO NOT EDIT.

package pagerank

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

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
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

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/kernels.hsaco": {
		name:    "kernels.hsaco",
		local:   "kernels.hsaco",
		size:    9640,
		modtime: 1563222183,
		compressed: `
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
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
