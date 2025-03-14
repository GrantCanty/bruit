package ws_client

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	logging "github.com/sacOO7/go-logger"
)

type Empty struct {
}

var logger = logging.GetLogger(reflect.TypeOf(Empty{}).PkgPath()).SetLevel(logging.OFF)

func (socket Socket) EnableLogging() {
	logger.SetLevel(logging.TRACE)
}

func (socket Socket) GetLogger() logging.Logger {
	return logger
}

type Socket struct {
	Conn              *websocket.Conn
	WebsocketDialer   *websocket.Dialer
	Url               string
	ConnectionOptions ConnectionOptions
	RequestHeader     http.Header
	OnConnected       func(socket Socket)
	OnTextMessage     func(message string, socket Socket)
	OnBinaryMessage   func(data []byte, socket Socket)
	OnConnectError    func(err error, socket Socket)
	OnDisconnected    func(err error, socket Socket)
	OnPingReceived    func(data string, socket Socket)
	OnPongReceived    func(data string, socket Socket)
	IsConnected       bool
	isConnectedMu     sync.RWMutex
	Timeout           time.Duration
	sendMu            *sync.Mutex // Prevent "concurrent write to websocket connection"
	receiveMu         *sync.Mutex
	privateConn       bool
	publicConn        bool
	bookConn          bool
}

type ConnectionOptions struct {
	UseCompression bool
	UseSSL         bool
	Proxy          func(*http.Request) (*url.URL, error)
	Subprotocols   []string
}

// todo Yet to be done
type ReconnectionOptions struct {
}

func New(url string) Socket {
	return Socket{
		Url:           url,
		RequestHeader: http.Header{},
		ConnectionOptions: ConnectionOptions{
			UseCompression: false,
			UseSSL:         true,
		},
		WebsocketDialer: &websocket.Dialer{},
		Timeout:         0,
		sendMu:          &sync.Mutex{},
		receiveMu:       &sync.Mutex{},
	}
}

func (socket *Socket) setConnectionOptions() {
	socket.WebsocketDialer.EnableCompression = socket.ConnectionOptions.UseCompression
	socket.WebsocketDialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: socket.ConnectionOptions.UseSSL}
	socket.WebsocketDialer.Proxy = socket.ConnectionOptions.Proxy
	socket.WebsocketDialer.Subprotocols = socket.ConnectionOptions.Subprotocols
}

func (socket *Socket) Connect() {
	var err error
	var resp *http.Response
	socket.setConnectionOptions()

	socket.Conn, resp, err = socket.WebsocketDialer.Dial(socket.Url, socket.RequestHeader)

	if err != nil {
		logger.Error.Println("Error while connecting to server ", err)
		if resp != nil {
			logger.Error.Println("HTTP Response ", resp.StatusCode, " status: ", resp.Status)
		}
		socket.SetIsConnected(false)
		if socket.OnConnectError != nil {
			socket.OnConnectError(err, *socket)
		}
		return
	}

	logger.Info.Println("Connected to server")

	if socket.OnConnected != nil {
		socket.SetIsConnected(true)
		socket.OnConnected(*socket)
	}

	defaultPingHandler := socket.Conn.PingHandler()
	socket.Conn.SetPingHandler(func(appData string) error {
		logger.Trace.Println("Received PING from server")
		if socket.OnPingReceived != nil {
			socket.OnPingReceived(appData, *socket)
		}
		return defaultPingHandler(appData)
	})

	defaultPongHandler := socket.Conn.PongHandler()
	socket.Conn.SetPongHandler(func(appData string) error {
		logger.Trace.Println("Received PONG from server")
		if socket.OnPongReceived != nil {
			socket.OnPongReceived(appData, *socket)
		}
		return defaultPongHandler(appData)
	})

	defaultCloseHandler := socket.Conn.CloseHandler()
	socket.Conn.SetCloseHandler(func(code int, text string) error {
		result := defaultCloseHandler(code, text)
		logger.Warning.Println("Disconnected from server ", result)
		if socket.OnDisconnected != nil {
			socket.SetIsConnected(false)
			socket.OnDisconnected(errors.New(text), *socket)
		}
		return result
	})

	go func() {
		for {
			socket.receiveMu.Lock()
			if socket.Timeout != 0 {
				socket.Conn.SetReadDeadline(time.Now().Add(socket.Timeout))
			}
			messageType, message, err := socket.Conn.ReadMessage()
			socket.receiveMu.Unlock()
			if err != nil {
				logger.Error.Println("read:", err)
				if socket.OnDisconnected != nil {
					socket.SetIsConnected(false)
					socket.OnDisconnected(err, *socket)
				}
				return
			}
			logger.Info.Println("recv: ", message)

			switch messageType {
			case websocket.TextMessage:
				socket.receiveMu.Lock()
				if socket.OnTextMessage != nil {
					socket.OnTextMessage(string(message), *socket)
				}
				socket.receiveMu.Unlock()
			case websocket.BinaryMessage:
				socket.receiveMu.Lock()
				if socket.OnBinaryMessage != nil {
					socket.OnBinaryMessage(message, *socket)
				}
				socket.receiveMu.Unlock()
			}
		}
	}()
}

func (socket *Socket) SendText(message string) {
	err := socket.send(websocket.TextMessage, []byte(message))
	if err != nil {
		logger.Error.Println("write:", err)
		return
	}
}

func (socket *Socket) SendBinary(data []byte) {
	err := socket.send(websocket.BinaryMessage, data)
	if err != nil {
		logger.Error.Println("write:", err)
		return
	}
}

func (socket *Socket) send(messageType int, data []byte) error {
	socket.sendMu.Lock()
	err := socket.Conn.WriteMessage(messageType, data)
	socket.sendMu.Unlock()
	return err
}

func (socket *Socket) Close() {
	err := socket.send(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		logger.Error.Println("write close:", err)
	}
	socket.Conn.Close()
	if socket.OnDisconnected != nil {
		socket.SetIsConnected(false)
		socket.OnDisconnected(err, *socket)
	}
}

func ReceiveLocker(socket *Socket) {
	socket.receiveMu.Lock()
}

func ReceiveUnlocker(socket *Socket) {
	socket.receiveMu.Unlock()
}

func SendLocker(socket *Socket) {
	socket.sendMu.Lock()
}

func SendUnlocker(socket *Socket) {
	socket.sendMu.Unlock()
}

func PublicInit(socket *Socket, PubWSUrl string) {
	*socket = New(PubWSUrl)

	socket.privateConn = false
	socket.publicConn = true
	socket.bookConn = false
}

func BookInit(socket *Socket, BookWSUrl string) {
	*socket = New(BookWSUrl)

	socket.privateConn = false
	socket.publicConn = false
	socket.bookConn = true
}

func PrivateInit(socket *Socket, PrivWSUrl string) {
	*socket = New(PrivWSUrl)

	socket.privateConn = true
	socket.publicConn = false
	socket.bookConn = false
}

func IsPublicSocket(socket Socket) bool {
	if socket.publicConn == true && socket.privateConn == false && socket.bookConn == false {
		return true
	}
	return false
}

func IsPrivateSocket(socket Socket) bool {
	if socket.publicConn == false && socket.privateConn == true && socket.bookConn == false {
		return true
	}
	return false
}

func IsBookSocket(socket Socket) bool {
	if socket.publicConn == false && socket.privateConn == false && socket.bookConn == true {
		return true
	}
	return false
}

func IsInit(socket Socket) bool {
	if socket.publicConn == false && socket.privateConn == false && socket.bookConn == false {
		return false
	}
	return true
}

func (socket *Socket) GetIsConnected() bool {
	socket.isConnectedMu.RLock()
	defer socket.isConnectedMu.RUnlock()
	return socket.IsConnected
}

// Safe setter for IsConnected
func (socket *Socket) SetIsConnected(value bool) {
	socket.isConnectedMu.Lock()
	defer socket.isConnectedMu.Unlock()
	socket.IsConnected = value
}
