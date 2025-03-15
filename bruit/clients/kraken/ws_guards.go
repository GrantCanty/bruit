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
	if !ws_client.IsInit(client.GetPubSocket()) {
		return errors.New("krakenClient is not initialized")
	}

	if !ws_client.IsPublicSocket(client.GetPubSocket()) {
		return errors.New("public socket function called on a private socket")
	}
	return nil // nil means that socket is init
}

func IsBookSocketInit(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !ws_client.IsInit(client.GetBookSocket()) {
		return errors.New("bookSocket is not initialized")
	}

	if !ws_client.IsBookSocket(client.GetBookSocket()) {
		return errors.New("public/private socket function called on a book socket")
	}
	return nil
}

func IsPrivSocketInit(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !ws_client.IsInit(client.GetPrivSocket()) {
		return errors.New("krakenClient is not initialized")
	}

	if !ws_client.IsPrivateSocket(client.GetPrivSocket()) {
		return errors.New("public socket function called on a private socket")
	}
	return nil
}

func PubSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	pubSocketMutex.Lock()
	defer pubSocketMutex.Unlock()

	if !ws_client.IsInit(client.GetPubSocket()) {
		return errors.New("krakenClient is not initialized")
	}

	if !ws_client.IsPublicSocket(client.GetPubSocket()) {
		return errors.New("public socket function called on a private socket")
	}

	if !client.GetPubSocketPointer().GetIsConnected() {
		client.GetPubSocketPointer().OnConnected = func(socket ws_client.Socket) {
			log.Println("PubSocketGuard: Connected to pub server")
		}

		log.Println("Connecting to pub server...")
		client.GetPubSocketPointer().Connect()

	}

	return nil
}

func BookSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	ws_client.ReceiveLocker(client.GetBookSocketPointer())

	if !ws_client.IsInit(client.GetBookSocket()) {
		return errors.New("bookSocket is not initialized")
	}

	if !ws_client.IsBookSocket(client.GetBookSocket()) {
		return errors.New("public/private socket function called on a book socket")
	}

	if !client.GetBookSocketPointer().GetIsConnected() {
		client.GetBookSocketPointer().OnConnected = func(socket ws_client.Socket) {
			log.Println("BookSocketGuard: Connected to book server")
		}

		// add something here for locking access to the Book Socket
		log.Println("Connecting to pub server...")
		client.GetBookSocketPointer().Connect()

		ws_client.ReceiveUnlocker(client.GetBookSocketPointer())
	} else {
		ws_client.ReceiveUnlocker(client.GetBookSocketPointer())
	}
	return nil
}

func PrivSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	ws_client.ReceiveLocker(client.GetPrivSocketPointer())

	if !ws_client.IsInit(client.GetPrivSocket()) {
		return errors.New("krakenClient is not initialized")
	}

	if !ws_client.IsPrivateSocket(client.GetPrivSocket()) {
		return errors.New("public socket function called on a private socket")
	}

	if !client.GetPrivSocketPointer().GetIsConnected() {
		client.GetPrivSocketPointer().OnConnected = func(socket ws_client.Socket) {
			log.Println("Connected to private server")
		}

		client.GetPrivSocketPointer().Connect()
		ws_client.ReceiveUnlocker(client.GetPrivSocketPointer())
	} else {
		ws_client.ReceiveUnlocker(client.GetPrivSocketPointer())
	}
	return nil
}

func AreChannelsInit(ws *web_socket.WebSocketClient) bool {
	return ws.GetPrivChan() != nil
}
