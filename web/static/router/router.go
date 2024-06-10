package router

import (
	"net/http"
	"strings"
	"sync"
)

type Router struct {
	s                  sync.Mutex
	hardcodedResources map[string]*Resource
	// locations are paths to the directories, from which the files will be served.
	locations []*directory

	// defaultResourceName points to the default file to be served in case Resource resolution fails.
	defaultResourceName string
}

func New() *Router {
	return &Router{
		hardcodedResources:  map[string]*Resource{},
		defaultResourceName: "index.html",
	}
}

func (l *Router) HttpHandle(w http.ResponseWriter, r *http.Request, tryDefault bool) bool {
	var path = strings.TrimPrefix(r.URL.Path, "/")
	var res = l.FindResource(path, tryDefault)
	if res != nil {
		w.Header().Set("Content-Type", res.mime)
		_, _ = w.Write(res.data)
		return true
	}
	return false
}

func (l *Router) RegisterDefaultResource(name string) {
	l.defaultResourceName = name
}

func (l *Router) RegisterStatic(name, filePath string) error {
	l.s.Lock()
	defer l.s.Unlock()

	r, e := newResource(filePath)
	if e != nil {
		return e
	}
	l.hardcodedResources[name] = r
	return nil
}

func (l *Router) RegisterDirectory(filePath string) {
	l.s.Lock()
	defer l.s.Unlock()
	l.locations = append(l.locations, newDirectory(filePath))
}

func (l *Router) FindResource(name string, tryDefault bool) *Resource {
	// check if there is a hardcoded resource available at that address
	l.s.Lock()
	r, h := l.hardcodedResources[name]
	locations := l.locations
	l.s.Unlock()
	if h {
		_ = r.KeepUpToDate()
		return r
	}

	// walk through all directories and try and find the resource.
	for _, dir := range locations {
		if r = dir.Resource(name); r != nil {
			return r
		}
	}
	if tryDefault {
		return l.FindResource(l.defaultResourceName, false)
	}
	return nil
}
