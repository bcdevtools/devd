package utils

import (
	"fmt"
	"net/http"
	"strconv"
)

func TryGetDownloadFileSizeInBytes(url string) (contentLength int64, descContentLength string, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}

	contentLengthVal := res.Header.Get("Content-Length")
	contentLength, err = strconv.ParseInt(contentLengthVal, 10, 64)
	if err != nil {
		return 0, "", err
	}

	if contentLength < 1024 {
		return contentLength, fmt.Sprintf("%s bytes", contentLengthVal), nil
	}

	if contentLength < 1024*1024 {
		return contentLength, fmt.Sprintf("%.2f KB", float64(contentLength)/1024), nil
	}

	if contentLength < 1024*1024*1024 {
		return contentLength, fmt.Sprintf("%.2f MB", float64(contentLength)/(1024*1024)), nil
	}

	return contentLength, fmt.Sprintf("%.2f GB", float64(contentLength)/(1024*1024*1024)), nil
}
