/* 標準函示庫 http.FileServer 的修改版 */

package httpfs

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

const indexPage = "/index.html"

type fileHandler struct {
	spaMode bool
	root    http.FileSystem
}

// NewFileServer create new File Server
func NewFileServer(rootFS http.FileSystem, spaMode bool) *fileHandler {
	return &fileHandler{
		spaMode: spaMode,
		root:    rootFS,
	}
}

func (fHandler *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	fHandler.serveFile(w, r, path.Clean(upath))
}

// localRedirect gives a Moved Permanently response.
// It does not convert relative paths to absolute paths like Redirect does.
func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}

func (fHandler *fileHandler) openErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if fHandler.spaMode && (errors.Is(err, fs.ErrNotExist) || errors.Is(err, fs.ErrPermission)) {
		var f http.File
		f, err = fHandler.root.Open(indexPage)
		if err == nil {
			defer f.Close()
			var d fs.FileInfo
			d, err = f.Stat()
			if err == nil {
				http.ServeContent(w, r, d.Name(), d.ModTime(), f)
				return
			}
		}
	}

	var msg string
	var code int
	if errors.Is(err, fs.ErrNotExist) {
		msg, code = "not found", http.StatusNotFound
	} else if errors.Is(err, fs.ErrPermission) {
		msg, code = "Forbidden", http.StatusForbidden
	} else {
		msg, code = "Internal Server Error", http.StatusInternalServerError
	}

	http.Error(w, msg, code)
}

// name is '/'-separated, not filepath.Separator.
func (fHandler *fileHandler) serveFile(w http.ResponseWriter, r *http.Request, name string) {
	// redirect .../index.html to .../
	// can't use Redirect() because that would make the path absolute,
	// which would be a problem running under StripPrefix
	if strings.HasSuffix(r.URL.Path, indexPage) {
		localRedirect(w, r, "./")
		return
	}

	f, err := fHandler.root.Open(name)
	if err != nil {
		fHandler.openErrorHandler(w, r, err)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		fHandler.openErrorHandler(w, r, err)
		return
	}

	if !fHandler.spaMode {
		// redirect to canonical path: / at end of directory url
		// r.URL.Path always begins with /
		url := r.URL.Path
		if d.IsDir() {
			if url[len(url)-1] != '/' {
				localRedirect(w, r, path.Base(url)+"/")
				return
			}
		} else {
			if url[len(url)-1] == '/' {
				localRedirect(w, r, "../"+path.Base(url))
				return
			}
		}
	}

	if d.IsDir() {
		// use contents of index.html for directory, if present
		index := strings.TrimSuffix(name, "/") + indexPage
		var ff http.File
		ff, err = fHandler.root.Open(index)
		if err == nil {
			defer ff.Close()
			var dd fs.FileInfo
			dd, err = ff.Stat()
			if err == nil {
				d = dd
				f = ff
			}
		} else {
			// 請求的 URL 為目錄，禁止訪問
			err = os.ErrPermission
			fHandler.openErrorHandler(w, r, err)
			return
		}
	}

	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}
