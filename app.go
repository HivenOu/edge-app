package main

import (
	"crypto/tls"
	"crypto/x509"
	"edge-app/camcontroll"
	"edge-app/domain"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	subClient := InitMqttClient(onSubConnectionLost)
	//pubClient := InitMqttClient(onPubConnectionLost)

	wait := sync.WaitGroup{}
	wait.Add(1)
	//go func() {
	//	for {
	//		time.Sleep(1 * time.Second)
	//		pubClient.Publish("topic", 0, false, "hello world")
	//	}
	//}()

	//subClient.Subscribe("#", 0, onReceived)
	subClient.Subscribe("/cam_control", 0, camControl)
	//subClient.Subscribe("", 0, onReceived)

	wait.Wait()
}

func InitMqttClient(onConnectionLost MQTT.ConnectionLostHandler) MQTT.Client {
	pool := x509.NewCertPool()
	cert, err := tls.LoadX509KeyPair("/home/manvw/ief/edge_cert/i1m4AbeL7I_private_cert.crt", "/home/manvw/ief/edge_cert/i1m4AbeL7I_private_cert.key")
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{cert},
		// 单向认证，client不校验服务端证书
		InsecureSkipVerify: true,
	}
	// 使用tls或者ssl协议，连接8883端口
	opts := MQTT.NewClientOptions().AddBroker("tls://127.0.0.1:8883").SetClientID(fmt.Sprintf("%f", rand.Float64()))
	opts.SetTLSConfig(tlsConfig)
	opts.OnConnect = onConnect
	opts.AutoReconnect = false
	// 回调函数，客户端与服务端断连后立刻被触发
	opts.OnConnectionLost = onConnectionLost
	client := MQTT.NewClient(opts)
	loopConnect(client)
	return client
}

func camControl(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Receive topic: %s,  payload: %s \n", message.Topic(), string(message.Payload()))
	if string(message.Payload()) != "" {
		msg := new(domain.MQTTMsg)
		if err := json.Unmarshal(message.Payload(), msg); err != nil {
			fmt.Printf(err.Error())
			return
		}
		if msg.Action == "photo" {
			err := camcontroll.TakePhotograph("./images/image.jpg")
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
			return
		}
		if msg.Action == "get_photo" {
			b, err := camcontroll.GetPhotoByte("./images/image.jpg")
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
			s := base64.StdEncoding.EncodeToString(b)
			fmt.Println(s)
			return
		}
		fmt.Printf("invalid action : %s", msg.Action)
	}
}

func onReceived(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Receive topic: %s,  payload: %s \n", message.Topic(), string(message.Payload()))

}

// sub客户端与服务端断连后，触发重连机制
func onSubConnectionLost(client MQTT.Client, err error) {
	fmt.Println("on sub connect lost, try to reconnect")
	loopConnect(client)
	client.Subscribe("topic", 0, onReceived)
}

// pub客户端与服务端断连后，触发重连机制
func onPubConnectionLost(client MQTT.Client, err error) {
	fmt.Println("on pub connect lost, try to reconnect")
	loopConnect(client)
}

func onConnect(client MQTT.Client) {
	fmt.Println("on connect")
}

func loopConnect(client MQTT.Client) {
	for {
		token := client.Connect()
		if rs, err := CheckClientToken(token); !rs {
			fmt.Printf("connect error: %s\n", err.Error())
		} else {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func CheckClientToken(token MQTT.Token) (bool, error) {
	if token.Wait() && token.Error() != nil {
		return false, token.Error()
	}
	return true, nil
}
