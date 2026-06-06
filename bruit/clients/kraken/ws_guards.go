package kraken

import (
	web_socket "bruit/bruit/clients/kraken/web_socket_client"
	"bruit/bruit/ws_client"
	"errors"
	"log"
)

var KrakenClientNotInit error = errors.New("krakenClient is not initialized")

var NotPublicSocket error = errors.New("public socket function called on wrong socket")
var NotBookSocket error = errors.New("book socket function called on wrong socket")
var NotPrivSocket error = errors.New("private socket function called on wrong socket")

var PublicSocketNotInit error = errors.New("publicSocket is not initialized")
var BookSocketNotInit error = errors.New("bookSocket is not initialized")
var PrivSocketNotInit error = errors.New("privateSocket is not initialized")

func IsPubSocketInit(client *web_socket.WebSocketClient) error {
	if !client.GetPubSocket().IsInit() {
		return KrakenClientNotInit
	}

	if !client.GetPubSocket().IsPublicSocket() {
		return NotPublicSocket
	}
	return nil // nil means that socket is init
}

func IsBookSocketInit(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !client.GetBookSocket().IsInit() {
		return BookSocketNotInit
	}

	if !client.GetBookSocket().IsBookSocket() {
		return NotBookSocket
	}
	return nil
}

func IsPrivSocketInit(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !client.GetPrivSocket().IsInit() {
		return KrakenClientNotInit
	}

	if !client.GetPrivSocket().IsPrivateSocket() {
		return NotPrivSocket
	}
	return nil
}

func PubSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	socket := client.GetPubSocket()
	ws_client.ReceiveLocker(socket)
	defer ws_client.ReceiveUnlocker(socket)

	if !socket.IsInit() {
		return KrakenClientNotInit
	}

	if !socket.IsPublicSocket() {
		return NotPublicSocket
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
		return BookSocketNotInit
	}

	if !socket.IsBookSocket() {
		return NotBookSocket
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
		return KrakenClientNotInit
	}

	if !socket.IsPrivateSocket() {
		return NotPrivSocket
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
