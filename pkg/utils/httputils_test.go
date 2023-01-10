package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	t.Run("test time format", func(t *testing.T) {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	})
}

func TestExecuteRequest(t *testing.T) {
	//定义test参数，包括：名称 入仓 想要的返回值
	type args struct {
		host, path, method, rawQuery, filename string
		body                                   []byte
	}
	//获取文件的byte
	file, err := os.Open("/Users/mac/appData/images/xxx.jpg")
	if err != nil {
		fmt.Println(err.Error())
		t.Skip()
		return
	}
	b, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	tests := []struct {
		name        string
		args        args
		want        string
		returnError error
	}{
		{
			name: "success",
			args: args{
				"127.0.0.1:8080", "/dis/ief-images", http.MethodPost, "", "xxxx.jpg", b,
			},
			returnError: nil,
		},
		{
			name: "fail",
			args: args{
				"192.168.0.104:8080", "/dis/ief-images", http.MethodPost, "", "xxxx.jpg", b,
			},
			returnError: nil,
		},
	}

	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			if _, err := ExecuteRequest(te.args.host, te.args.path, te.args.method, te.args.rawQuery, te.args.filename, te.args.body); err != te.returnError {
				t.Errorf("DealResourceVersion() = %v, want %v", err, te.returnError)
			}
		})
	}
}
