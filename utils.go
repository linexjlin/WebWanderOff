package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"strings"
)

// 解压缩 Gzip
func UncompressGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	uncompressed, err := ioutil.ReadAll(reader)
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
