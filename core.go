package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type CacheSystem struct {
	ListenAddr       string
	ThirdPartyPrefix string
	DefaultServer    string
	DefaultScheme    string
	CacheRoot        string
	OfflineDomains   []string
}

func (c *CacheSystem) Listen() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", c.cacheProxyHandler)
	log.Println("Listening on:", c.ListenAddr)
	log.Fatal(http.ListenAndServe(c.ListenAddr, mux))
}

func (c *CacheSystem) cacheProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Remove '/' from the request URI to get the actual URL
	targetURL := r.RequestURI[1:]
	urlPath := r.URL.Path

	log.Println("path", r.URL.Path)
	log.Println("targetURL", targetURL)

	cachePath := ""
	if strings.HasPrefix(targetURL, "https/") {
		cachePath = c.CacheRoot + "/" + strings.Replace(urlPath, "https/", "", 1)
		targetURL = strings.Replace(targetURL, "https/", "https://", 1)
	} else if strings.HasPrefix(targetURL, "http/") {
		cachePath = c.CacheRoot + "/" + strings.Replace(urlPath, "https/", "", 1)
		targetURL = strings.Replace(targetURL, "http/", "http://", 1)
	} else {
		cachePath = c.CacheRoot + "/" + c.DefaultServer + "/" + urlPath
		targetURL = c.DefaultScheme + "://" + c.DefaultServer + "/" + targetURL
	}

	// Determine the cache path

	if strings.HasSuffix(cachePath, "/") {
		cachePath = cachePath + "index"
	}
	log.Println("try cache:", cachePath)
	cacheDir := path.Dir(cachePath)

	// Check if the resource is already cached
	if _, err := os.Stat(cachePath); err == nil {
		// Serve the resource from the cache
		log.Println("hit cache", cachePath)
		c.serveFile(w, r, cachePath)
		return
	}

	log.Println("Download static data from:", targetURL)
	// The resource is not cached, make a request to the target URL
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, targetURL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Forward headers and body if method is POST
	//req.Header = r.Header
	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	isGzip := resp.Header.Get("Content-Encoding") == "gzip"

	// Read the response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 如果经过 Gzip 压缩，则解压缩
	if isGzip {
		data, err = UncompressGzip(data)
		if err != nil {
			fmt.Println("解压缩错误:", err)
			return
		}
	}

	// Create the cache directory if it does not exist
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response to the cache
	if err := os.WriteFile(cachePath, data, 0666); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.serveFile(w, r, cachePath)
}

func (c *CacheSystem) serveFile(w http.ResponseWriter, r *http.Request, name string) {
	contentType := "application/octet-stream" // Default MIME type
	// Check if the resource is already cached
	if fileInfo, err := os.Stat(name); err == nil {
		//print the content-type of cachePath
		if mimeType := mime.TypeByExtension(path.Ext(fileInfo.Name())); mimeType != "" {
			contentType = mimeType
		} else if mimeType, err := mimetype.DetectFile(name); err == nil {
			contentType = mimeType.String()
		}

		log.Println("Detected Content-Type:", contentType)
		//https://github.com/ffmpegwasm/ffmpeg.wasm/issues/263
		//SharedArrayBuffer is only available to pages that are cross-origin isolated. So you need to host your own server with Cross-Origin-Embedder-Policy: require-corp and Cross-Origin-Opener-Policy: same-origin headers to use ffmpeg.wasm.
		if strings.Contains(contentType, "text/html") {
			//set cross-origin
			w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
			w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		}

		w.Header().Set("Content-Type", contentType)

		if isTextMimeType(contentType) {
			log.Println("It is text, replace")
			w.Write(c.replaceUrlInText(name))
			return
		}

		// Serve the resource from the cache
		http.ServeFile(w, r, name)
		return
	}
}

func (c *CacheSystem) replaceUrlInText(file string) []byte {
	// readAll data from file
	// replace https://cdn.jsdelivr.net with http://127.0.0.1:8099/http/cdn.jsdelivr.net
	//replaceRules =
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	for _, domain := range c.OfflineDomains {
		target := fmt.Sprintf("http://%s/%s", c.ListenAddr, strings.Replace(domain, "://", "/", 1))
		data = bytes.ReplaceAll(data, []byte(domain), []byte(target))
	}

	return data
}
