package kraken

import (
	web_socket "bruit/bruit/clients/kraken/web_socket_client"
	"bruit/bruit/ws_client"
	"errors"
	"log"
)

func IsPubSocketInit(client web_socket.WebSocketClient) error {
	if !ws_client.IsInit(client.GetPubSocket()) {
		return errors.New("KrakenClient is not initialized")
	}

	if !ws_client.IsPublicSocket(client.GetPubSocket()) {
		return errors.New("Public socket function called on a private socket")
	}
	return nil // nil means that socket is init
}

func IsBookSocketInit(client web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !ws_client.IsInit(client.GetBookSocket()) {
		return errors.New("bookSocket is not initialized")
	}

	if !ws_client.IsBookSocket(client.GetBookSocket()) {
		return errors.New("Public/private socket function called on a book socket")
	}
	return nil
}

func IsPrivSocketInit(client web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !ws_client.IsInit(client.GetPrivSocket()) {
		return errors.New("KrakenClient is not initialized")
	}

	if !ws_client.IsPrivateSocket(client.GetPrivSocket()) {
		return errors.New("Public socket function called on a private socket")
	}
	return nil
}

func PubSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	ws_client.ReceiveLocker(client.GetPubSocketPointer())

	if !ws_client.IsInit(client.GetPubSocket()) {
		return errors.New("KrakenClient is not initialized")
	}

	if !ws_client.IsPublicSocket(client.GetPubSocket()) {
		return errors.New("Public socket function called on a private socket")
	}

	if !client.GetPubSocket().IsConnected {
		client.GetPubSocketPointer().OnConnected = func(socket ws_client.Socket) {
			log.Println("Connected to public server")
		}

		log.Println("Connecting...")
		client.GetPubSocketPointer().Connect()

		ws_client.ReceiveUnlocker(client.GetPubSocketPointer())
	} else {
		ws_client.ReceiveUnlocker(client.GetPubSocketPointer())
	}

	return nil
}

func BookSocketGuard(client *web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	ws_client.ReceiveLocker(client.GetBookSocketPointer())

	if !ws_client.IsInit(client.GetBookSocket()) {
		return errors.New("bookSocket is not initialized")
	}

	if !ws_client.IsBookSocket(client.GetBookSocket()) {
		return errors.New("Public/private socket function called on a book socket")
	}

	if !client.GetBookSocket().IsConnected {
		client.GetBookSocketPointer().OnConnected = func(socket ws_client.Socket) {
			log.Println("BookSocketGuard: Connected to book server")
		}

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
		return errors.New("KrakenClient is not initialized")
	}

	if !ws_client.IsPrivateSocket(client.GetPrivSocket()) {
		return errors.New("Public socket function called on a private socket")
	}

	if !client.GetPrivSocket().IsConnected {
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
	if ws.GetPrivChan() != nil && ws.GetBookChan() != nil {
		return true
	} else {
		return false
	}
}
