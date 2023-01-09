package utils

import (
	"bytes"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
)

func ExecuteRequest(host, path, method, rawQuery, filename string, body []byte) ([]byte, error) {
	urlInfo := &url.URL{
		Scheme: "http",
		Path:   path,
		Host:   host,
	}
	if rawQuery != "" {
		urlInfo.RawQuery = rawQuery
	}
	contentType := "application/json"
	var reader io.Reader

	//文件上传时使用其他的reader
	if filename != "" && len(body) != 0 {
		var rb = &bytes.Buffer{} // 创建一个buffer
		w := multipart.NewWriter(rb)
		fw, err := w.CreateFormFile("file", filename) // 自定义文件名，发送文件流
		if err != nil {
			log.Fatalln(err)
		}
		_, err = fw.Write(body)
		if err != nil {
			return nil, err
		}
		err = w.Close()
		if err != nil {
			return nil, err
		}
		contentType = w.FormDataContentType()
		reader = rb
	} else if len(body) != 0 {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, path, reader)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	//http.File()
	header := http.Header{"content-type": []string{contentType}}
	req.Header = header
	req.URL = urlInfo
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		err = errors.New(string(data))
		logrus.Error(err)
		return nil, err
	}
	return data, err
}

func test(filename string) error {
	var rb = &bytes.Buffer{} // 创建一个buffer
	w := multipart.NewWriter(rb)
	fw, err := w.CreateFormFile("file", filename) // 自定义文件名，发送文件流
	if err != nil {
		log.Fatalln(err)
	}
	fi, err := os.Open("data.go")
	if err != nil {
		log.Fatalln(err)
	}
	_, _ = io.Copy(fw, fi) // 把文件内容，复制到fw中

	//_ = w.WriteField("test3", "test-v") // 上传其他字段
	//defer w.Close() // 很重要，一定要关闭写入，不然服务端会报EOF错误，而且度不到数据
	c := &http.Client{}
	req, err := http.NewRequest("POST", "", rb)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := c.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}
