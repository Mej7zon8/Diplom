package router

import (
	"io"
	"os"
	"path/filepath"
	"sync"
)

type directory struct {
	s        sync.Mutex
	filePath string

	cached map[string]*Resource
}

func newDirectory(filePath string) *directory {
	return &directory{filePath: filePath, cached: map[string]*Resource{}}
}

func (d *directory) Resource(name string) *Resource {
	d.s.Lock()
	defer d.s.Unlock()

	var path = filepath.Join(d.filePath, name)
	cached, h := d.cached[path]
	if h {
		if cached.KeepUpToDate() != nil {
			delete(d.cached, path)
			return nil
		}
		return cached
	}
	r, e := newResource(path)
	if e != nil {
		return nil
	}
	d.cached[path] = r
	return r
}

type Resource struct {
	info os.FileInfo
	path string
	mime string
	data []byte
}

func newResource(path string) (*Resource, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	info, e := f.Stat()
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(f)
	if e != nil {
		return nil, e
	}
	return &Resource{
		info: info,
		path: path,
		mime: TypeByExtension(filepath.Ext(path)),
		data: data,
	}, nil
}

func (r *Resource) KeepUpToDate() error {
	info, e := os.Stat(r.path)
	if e != nil {
		return e
	}
	if info.ModTime().Equal(r.info.ModTime()) {
		return nil
	}
	r.data, e = os.ReadFile(r.path)
	return e
}
