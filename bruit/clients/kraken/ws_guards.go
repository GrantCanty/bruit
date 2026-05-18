package kraken

import (
	web_socket "bruit/bruit/clients/kraken/web_socket_client"
	"bruit/bruit/ws_client"
	"errors"
	"log"
	"sync"
)

var pubSocketMutex sync.Mutex

func IsPubSocketInit(client *web_socket.WebSocketClient) error {
	if !client.GetPubSocket().IsInit() {
		return errors.New("krakenClient is not initialized")
	}

	if !client.GetPubSocket().IsPublicSocket() {
		return errors.New("public socket function called on a private socket")
	}
	return nil // nil means that socket is init
}

func IsBookSocketInit(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !client.GetBookSocket().IsInit() {
		return errors.New("bookSocket is not initialized")
	}

	if !client.GetBookSocket().IsBookSocket() {
		return errors.New("public/private socket function called on a book socket")
	}
	return nil
}

func IsPrivSocketInit(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !client.GetPrivSocket().IsInit() {
		return errors.New("krakenClient is not initialized")
	}

	if !client.GetPrivSocket().IsPrivateSocket() {
		return errors.New("public socket function called on a private socket")
	}
	return nil
}

func PubSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	pubSocketMutex.Lock()
	defer pubSocketMutex.Unlock()

	if !client.GetPubSocket().IsInit() {
		return errors.New("krakenClient is not initialized")
	}

	if !client.GetPubSocket().IsPublicSocket() {
		return errors.New("public socket function called on a private socket")
	}

	if !client.GetPubSocket().GetIsConnected() {
		client.GetPubSocket().OnConnected = func(socket *ws_client.Socket) {
			log.Println("PubSocketGuard: Connected to pub server")
		}

		log.Println("Connecting to pub server...")
		client.GetPubSocket().Connect()

	}

	return nil
}

func BookSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	ws_client.ReceiveLocker(client.GetBookSocket())

	if !client.GetBookSocket().IsInit() {
		return errors.New("bookSocket is not initialized")
	}

	if !client.GetBookSocket().IsBookSocket() {
		return errors.New("public/private socket function called on a book socket")
	}

	if !client.GetBookSocket().GetIsConnected() {
		client.GetBookSocket().OnConnected = func(socket *ws_client.Socket) {
			log.Println("BookSocketGuard: Connected to book server")
		}

		// add something here for locking access to the Book Socket
		log.Println("Connecting to book server...")
		client.GetBookSocket().Connect()

		ws_client.ReceiveUnlocker(client.GetBookSocket())
	} else {
		ws_client.ReceiveUnlocker(client.GetBookSocket())
	}
	return nil
}

func PrivSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	ws_client.ReceiveLocker(client.GetPrivSocket())

	if !client.GetPrivSocket().IsInit() {
		return errors.New("krakenClient is not initialized")
	}

	if !client.GetPrivSocket().IsPrivateSocket() {
		return errors.New("public socket function called on a private socket")
	}

	if !client.GetPrivSocket().GetIsConnected() {
		client.GetPrivSocket().OnConnected = func(socket *ws_client.Socket) {
			log.Println("Connected to private server")
		}

		client.GetPrivSocket().Connect()
		ws_client.ReceiveUnlocker(client.GetPrivSocket())
	} else {
		ws_client.ReceiveUnlocker(client.GetPrivSocket())
	}
	return nil
}

func AreChannelsInit(ws *web_socket.WebSocketClient) bool {
	return ws.GetPrivChan() != nil
}
