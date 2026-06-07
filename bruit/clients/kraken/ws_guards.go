package kraken

import (
	web_socket "bruit/bruit/clients/kraken/web_socket_client"
	"bruit/bruit/ws_client"
	"errors"
	"log"
)

var ErrNotPubSocket error = errors.New("public socket function called on wrong socket")
var ErrNotBookSocket error = errors.New("book socket function called on wrong socket")
var ErrNotPrivSocket error = errors.New("private socket function called on wrong socket")

var ErrPubSocketNotInit error = errors.New("publicSocket is not initialized")
var ErrBookSocketNotInit error = errors.New("bookSocket is not initialized")
var ErrPrivSocketNotInit error = errors.New("privateSocket is not initialized")

var ErrFailedToConnectToPubServer error = errors.New("failed to connect to pub server")
var ErrFailedToConnectToBookServer error = errors.New("failed to connect to book server")
var ErrFailedToConnectToPrivServer error = errors.New("failed to connect to private server")

func IsPubSocketInit(client *web_socket.WebSocketClient) error {
	if !client.GetPubSocket().IsInit() {
		return ErrPubSocketNotInit
	}

	if !client.GetPubSocket().IsPublicSocket() {
		return ErrNotPubSocket
	}
	return nil // nil means that socket is init
}

func IsBookSocketInit(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !client.GetBookSocket().IsInit() {
		return ErrBookSocketNotInit
	}

	if !client.GetBookSocket().IsBookSocket() {
		return ErrNotBookSocket
	}
	return nil
}

func IsPrivSocketInit(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !client.GetPrivSocket().IsInit() {
		return ErrPrivSocketNotInit
	}

	if !client.GetPrivSocket().IsPrivateSocket() {
		return ErrNotPrivSocket
	}
	return nil
}

func PubSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	socket := client.GetPubSocket()
	ws_client.ReceiveLocker(socket)
	defer ws_client.ReceiveUnlocker(socket)

	if !socket.IsInit() {
		return ErrPubSocketNotInit
	}

	if !socket.IsPublicSocket() {
		return ErrNotPubSocket
	}

	if !socket.GetIsConnected() {
		socket.OnConnected = func(socket *ws_client.Socket) {
			log.Println("PubSocketGuard: Connected to pub server")
		}

		log.Println("Connecting to pub server...")
		socket.Connect()
		if !socket.GetIsConnected() {
			return ErrFailedToConnectToPubServer
		}
	}

	return nil
}

func BookSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	socket := client.GetBookSocket()
	ws_client.ReceiveLocker(socket)
	defer ws_client.ReceiveUnlocker(socket)

	if !socket.IsInit() {
		return ErrBookSocketNotInit
	}

	if !socket.IsBookSocket() {
		return ErrNotBookSocket
	}

	if !socket.GetIsConnected() {
		socket.OnConnected = func(socket *ws_client.Socket) {
			log.Println("BookSocketGuard: Connected to book server")
		}

		log.Println("Connecting to book server...")
		socket.Connect()
		if !socket.GetIsConnected() {
			return ErrFailedToConnectToBookServer
		}
	}
	return nil
}

func PrivSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	socket := client.GetPrivSocket()
	ws_client.ReceiveLocker(socket)
	defer ws_client.ReceiveUnlocker(socket)

	if !socket.IsInit() {
		return ErrPrivSocketNotInit
	}

	if !socket.IsPrivateSocket() {
		return ErrNotPrivSocket
	}

	if !socket.GetIsConnected() {
		socket.OnConnected = func(socket *ws_client.Socket) {
			log.Println("Connected to private server")
		}

		socket.Connect()
		if !socket.GetIsConnected() {
			return ErrFailedToConnectToPrivServer
		}
	}
	return nil
}

func AreChannelsInit(ws *web_socket.WebSocketClient) bool {
	return ws.GetPrivChan() != nil
}
