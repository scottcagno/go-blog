package templates

import (
	"bytes"
	"fmt"
	"github.com/scottcagno/go-blog/pkg/logging"
	"html/template"
	"net/http"
	"sync"
)

const (
	_ = iota
	JAR
	JS
	OCTET_STREAM
	OGG
	PDF
	XHTML
	FLASH
	JSON
	XML
	ZIP
	WWW_FORM_URLENCODED
	_
	AUDIO_MPEG
	AUDIO_WMA
	AUDIO_REALAUDIO
	AUDIO_XWAV
	_
	IMAGE_GIF
	IMAGE_JPEG
	IMAGE_PNG
	IMAGE_TIFF
	IMAGE_XICON
	IMAGE_SVG
	_
	MULTI_MIXED
	MULTI_ALT
	MULTI_RELATED
	MULTI_FORMDATA
	_
	TEXT_CSS
	TEXT_CSV
	TEXT_HTML
	TEXT_JAVASCRIPT
	TEXT_PLAIN
	TEXT_XML
)

var contentType = map[int]string{
	// application types
	JAR:                 "application/java-archive",
	JS:                  "application/javascript",
	OCTET_STREAM:        "application/octet-stream",
	OGG:                 "application/ogg",
	PDF:                 "application/pdf",
	XHTML:               "application/xhtml+xml",
	FLASH:               "application/x-shockwave-flash",
	JSON:                "application/json",
	XML:                 "application/xml",
	ZIP:                 "application/zip",
	WWW_FORM_URLENCODED: "application/x-www-form-urlencoded",
	// audio types
	AUDIO_MPEG:      "audio/mpeg",
	AUDIO_WMA:       "audio/x-ms-wma",
	AUDIO_REALAUDIO: "audio/vnd.rn-realaudio",
	AUDIO_XWAV:      "audio/x-wav",
	// image types
	IMAGE_GIF:   "image/gif",
	IMAGE_JPEG:  "image/jpeg",
	IMAGE_PNG:   "image/png",
	IMAGE_TIFF:  "image/tiff",
	IMAGE_XICON: "image/x-icon",
	IMAGE_SVG:   "image/svg+xml",
	// multipart types
	MULTI_MIXED:    "multipart/mixed",
	MULTI_ALT:      "multipart/alternative",
	MULTI_RELATED:  "multipart/related",
	MULTI_FORMDATA: "multipart/form-data",
	// text types
	TEXT_CSS:        "text/css",
	TEXT_CSV:        "text/csv",
	TEXT_HTML:       "text/html",
	TEXT_JAVASCRIPT: "text/javascript",
	TEXT_PLAIN:      "text/plain",
	TEXT_XML:        "text/xml",
}

var sp = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

var fm = template.FuncMap{}

type TemplateCache struct {
	cache *template.Template
	logr  *logging.Logger
	buff  sync.Pool
	ct    map[int]string
	sync.RWMutex
}

func NewTemplateCache(pattern string, logger *logging.Logger) *TemplateCache {
	//example pattern: "web/templates/*.html"
	return &TemplateCache{
		cache: template.Must(template.New("*").Funcs(fm).ParseGlob(pattern)),
		logr:  logger,
		buff: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		ct: contentType,
	}
}

func (t *TemplateCache) Render(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	t.Lock()
	defer t.Unlock()
	buffer := t.buff.Get().(*bytes.Buffer)
	buffer.Reset()
	err := t.cache.ExecuteTemplate(buffer, tmpl, map[string]interface{}{"title": tmpl, "data": data})
	if err != nil {
		t.buff.Put(buffer)
		t.logr.Error.Printf("Error while executing template (%s): %v\n", tmpl, err)
		http.Redirect(w, r, "/error/404", http.StatusTemporaryRedirect)
		return
	}
	_, err = buffer.WriteTo(w)
	if err != nil {
		t.logr.Error.Printf("Error while writing (Render) to ResponseWriter: %v\n", err)
	}
	t.buff.Put(buffer)
	return
}

func (t *TemplateCache) Raw(w http.ResponseWriter, format string, data ...interface{}) {
	t.Lock()
	defer t.Unlock()
	_, err := fmt.Fprintf(w, format, data...)
	if err != nil {
		t.logr.Error.Printf("Error while writing (Raw) to ResponseWriter: %v\n", err)
		return
	}
	return
}

func (t *TemplateCache) ContentType(w http.ResponseWriter, content int) {
	t.Lock()
	defer t.Unlock()
	ct, ok := contentType[content]
	if !ok {
		t.logr.Error.Printf("Error, incompatible content type!\n")
		return
	}
	w.Header().Set("Content-Type", ct)
	return
}
