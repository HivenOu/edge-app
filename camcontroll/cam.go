package camcontroll

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

const PhotoCMD = "fswebcam"

// TakePhotograph 调用摄像头拍照
func TakePhotograph(photoName string) error {
	if photoName == "" {
		fmt.Println("photoName is nil")
		return nil
	}
	//realCMD := fmt.Sprintf(PhotoCMD, photoName)
	cmd := exec.Command(PhotoCMD, photoName)
	return cmd.Run()
}

// 获取照片
func GetPhotoByte(photoName string) ([]byte, error) {
	//base64.StdEncoding
	f, err := os.Open(photoName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil
	}
	return io.ReadAll(f)
}
