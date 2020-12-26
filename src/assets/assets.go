package assets

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

const (
	ContentEncoding = "Content-Encoding"
	ContentLength   = "Content-Length"
	ContentType     = "Content-Type"
	IfNoneMatch     = "If-None-Match"
	ETag            = "ETag"
	CacheControl    = "Cache-Control"

	GZip      = "gzip"
	Identity  = "identity"
	NoCache   = "public, no-cache"
	Immutable = "public, immutable, max-age=2592000"
)

func NewAssetsHandler(assetsDir string) gin.HandlerFunc {
	assetsContext := newAssetsContext(assetsDir)
	return func(c *gin.Context) {
		assetsContext.Handle(c)
	}
}

type Context struct {
	AssetsDir  string
	AssetsInfo map[string]*AssetInfo
}

type AssetInfo struct {
	FilePath        string
	ContentType     string
	ContentLength   string
	ContentEncoding string
	ETag            string
}

func newAssetsContext(assetsDir string) *Context {
	return &Context{
		AssetsDir:  assetsDir,
		AssetsInfo: map[string]*AssetInfo{},
	}
}

func (context *Context) Handle(c *gin.Context) {
	info := context.GetAssetInfo(c.Param("asset"))
	if info == nil {
		c.String(http.StatusNotFound, fmt.Sprintf("unknown asset: %q", c.Param("asset")))
		return
	}

	cacheControlValue := Immutable
	if info.ETag != c.Request.URL.Query().Get("v") {
		cacheControlValue = NoCache
	}

	w := c.Writer
	h := w.Header()

	if info.ETag == c.Request.Header.Get(IfNoneMatch) {
		h.Set(ETag, info.ETag)
		h.Set(CacheControl, cacheControlValue)
		w.WriteHeader(http.StatusOK)
		return
	}

	h.Set(ContentType, info.ContentType)
	h.Set(ContentLength, info.ContentLength)
	h.Set(ContentEncoding, info.ContentEncoding)
	h.Set(ETag, info.ETag)
	h.Set(CacheControl, cacheControlValue)
	w.WriteHeader(http.StatusOK)
	c.File(info.FilePath)
}

func (context *Context) GetAssetInfo(assetName string) *AssetInfo {
	contentType := MimeTypes[path.Ext(assetName)]
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return &AssetInfo{
		FilePath:        path.Join(context.AssetsDir, assetName),
		ContentType:     contentType,
		ContentLength:   "33",
		ContentEncoding: Identity,
		ETag:            "abcde",
	}
}
