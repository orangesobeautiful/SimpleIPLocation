/* 標準函示庫 http.FileServer 的修改版 */

package httpfs

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var compressExtList = []string{"js", "css", "html", "ico"}

const indexPage = "/index.html"

// FileHandler 用來處理靜態檔案的 http handler
type FileHandler struct {
	spaMode      bool
	acptEncoding []string
	root         http.FileSystem
}

// NewFileServer create new File Server
func NewFileServer(rootFS http.FileSystem, spaMode bool) *FileHandler {
	return &FileHandler{
		spaMode:      spaMode,
		acptEncoding: []string{"br", "gzip", "deflate"},
		root:         rootFS,
	}
}

func (fHandler *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

// fOpen 套用壓縮格式開啟檔案
func (fHandler *FileHandler) open(r *http.Request, name string) (http.File, string, error) {
	var encodingDir string
	var encoding string
	if IsNeedCompress(name) {
		encoding = ParseAcceptHeader(r.Header.Get("Accept-Encoding"), fHandler.acptEncoding)
	}
	if encoding == "" {
		encodingDir = "original"
	} else {
		encodingDir = encoding
	}

	f, err := fHandler.root.Open(path.Join(encodingDir, name))
	if err != nil {
		return nil, "", err
	}
	return f, encoding, nil
}

// serveContent 回應客戶端要求的檔案並寫入相應的編碼格式
func (fHandler *FileHandler) serveContent(
	w http.ResponseWriter, req *http.Request, encoding, name string, modtime time.Time, content io.ReadSeeker) {
	if encoding != "" {
		w.Header().Set("content-encoding", encoding)
		w.Header().Set("vary", "Accept-Encoding")
	}
	http.ServeContent(w, req, name, modtime, content)
}

func (fHandler *FileHandler) openErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if fHandler.spaMode && (errors.Is(err, fs.ErrNotExist) || errors.Is(err, fs.ErrPermission)) {
		var f http.File
		var encoding string
		f, encoding, err = fHandler.open(r, indexPage)
		if err == nil {
			defer f.Close()
			var d fs.FileInfo
			d, err = f.Stat()
			if err == nil {
				fHandler.serveContent(w, r, encoding, d.Name(), d.ModTime(), f)
				return
			}
		}
	}

	var msg string
	var code int
	if errors.Is(err, fs.ErrNotExist) {
		msg, code = "Not found", http.StatusNotFound
	} else if errors.Is(err, fs.ErrPermission) {
		msg, code = "Forbidden", http.StatusForbidden
	} else {
		msg, code = "Internal Server Error", http.StatusInternalServerError
	}

	http.Error(w, msg, code)
}

func IsNeedCompress(name string) bool {
	ext := path.Ext(name)
	if len(ext) < 2 {
		return false
	}
	ext = ext[1:]

	for idx := range compressExtList {
		if ext == compressExtList[idx] {
			return true
		}
	}
	return false
}

// name is '/'-separated, not filepath.Separator.
func (fHandler *FileHandler) serveFile(w http.ResponseWriter, r *http.Request, name string) {
	// redirect .../index.html to .../
	// can't use Redirect() because that would make the path absolute,
	// which would be a problem running under StripPrefix
	if strings.HasSuffix(r.URL.Path, indexPage) {
		localRedirect(w, r, "./")
		return
	}

	f, encoding, err := fHandler.open(r, name)
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
		var idxEncoding string
		ff, idxEncoding, err = fHandler.open(r, index)
		if err == nil {
			defer ff.Close()
			var dd fs.FileInfo
			dd, err = ff.Stat()
			if err == nil {
				d = dd
				f = ff
				encoding = idxEncoding
			}
		} else {
			// 請求的 URL 為目錄，禁止訪問
			err = os.ErrPermission
			fHandler.openErrorHandler(w, r, err)
			return
		}
	}

	fHandler.serveContent(w, r, encoding, d.Name(), d.ModTime(), f)
}
