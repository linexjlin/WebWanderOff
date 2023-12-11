package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 解压缩 Gzip
func UncompressGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	uncompressed, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return uncompressed, nil
}

func isTextMimeType(mimeType string) bool {
	textTypes := []string{
		"text/",
		"application/json",
		"application/javascript",
		"application/xml",
		"application/xhtml+xml",
	}

	for _, prefix := range textTypes {
		if strings.HasPrefix(mimeType, prefix) {
			return true
		}
	}
	return false
}

func findFavicon(dir string) string {
	pattern := filepath.Join(dir, "favicon.*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return ""
	}

	for _, match := range matches {
		if strings.HasSuffix(match, ".ico") || strings.HasSuffix(match, ".png") || strings.HasSuffix(match, ".jpg") || strings.HasSuffix(match, ".jpeg") || strings.HasSuffix(match, ".svg") {

			data, err := os.ReadFile(match)
			if err != nil {
				continue
			}

			base64String := base64.StdEncoding.EncodeToString(data)
			return base64String
		}
	}
	return ""
}
