package main

import (
	"edge-app/pkg/utils"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	//打开文件获得[]byte
	file, err := os.Open("/Users/mac/appData/images/image2022-12-28_10-4-34.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	b, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = utils.ExecuteRequest("192.168.0.104:8080", "/dis/ief-images", http.MethodPost, "", "xxxx.jpg", b)
	if err != nil {
		fmt.Println(err.Error())
	}
}
