package kraken

import (
	web_socket "bruit/bruit/clients/kraken/web_socket_client"
	"bruit/bruit/ws_client"
	"errors"
	"log"
)

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
	socket := client.GetPubSocket()
	ws_client.ReceiveLocker(socket)
	defer ws_client.ReceiveUnlocker(socket)

	if !socket.IsInit() {
		return errors.New("krakenClient is not initialized")
	}

	if !socket.IsPublicSocket() {
		return errors.New("public socket function called on a private socket")
	}

	if !socket.GetIsConnected() {
		socket.OnConnected = func(socket *ws_client.Socket) {
			log.Println("PubSocketGuard: Connected to pub server")
		}

		log.Println("Connecting to pub server...")
		socket.Connect()
		if !socket.GetIsConnected() {
			return errors.New("failed to connect to pub server")
		}
	}

	return nil
}

func BookSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	socket := client.GetBookSocket()
	ws_client.ReceiveLocker(socket)
	defer ws_client.ReceiveUnlocker(socket)

	if !socket.IsInit() {
		return errors.New("bookSocket is not initialized")
	}

	if !socket.IsBookSocket() {
		return errors.New("public/private socket function called on a book socket")
	}

	if !socket.GetIsConnected() {
		socket.OnConnected = func(socket *ws_client.Socket) {
			log.Println("BookSocketGuard: Connected to book server")
		}

		log.Println("Connecting to book server...")
		socket.Connect()
		if !socket.GetIsConnected() {
			return errors.New("failed to connect to book server")
		}
	}
	return nil
}

func PrivSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	socket := client.GetPrivSocket()
	ws_client.ReceiveLocker(socket)
	defer ws_client.ReceiveUnlocker(socket)

	if !socket.IsInit() {
		return errors.New("krakenClient is not initialized")
	}

	if !socket.IsPrivateSocket() {
		return errors.New("public socket function called on a private socket")
	}

	if !socket.GetIsConnected() {
		socket.OnConnected = func(socket *ws_client.Socket) {
			log.Println("Connected to private server")
		}

		socket.Connect()
		if !socket.GetIsConnected() {
			return errors.New("failed to connect to private server")
		}
	}
	return nil
}

func AreChannelsInit(ws *web_socket.WebSocketClient) bool {
	return ws.GetPrivChan() != nil
}
