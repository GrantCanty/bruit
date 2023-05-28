package kraken

import (
	web_socket "bruit/bruit/clients/kraken/web_socket_client"
	"bruit/bruit/ws_client"
	"errors"
)

/*func socketGuard(socket ws_client.Socket) error { // checks if socket is init. returns error if not init
	if !ws_client.IsInit(socket) {
		return errors.New("KrakenClient is not initialized")
	}
	return nil
}*/

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

func PubSocketGuard(client web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !ws_client.IsInit(client.GetPubSocket()) {
		return errors.New("KrakenClient is not initialized")
	}

	if !ws_client.IsPublicSocket(client.GetPubSocket()) {
		return errors.New("Public socket function called on a private socket")
	}

	if !client.GetPubSocket().IsConnected {
		client.GetPubSocketPointer().Connect()
	}
	return nil
}

func BookSocketGuard(client web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !ws_client.IsInit(client.GetBookSocket()) {
		return errors.New("bookSocket is not initialized")
	}

	if !ws_client.IsBookSocket(client.GetBookSocket()) {
		return errors.New("Public/private socket function called on a book socket")
	}

	if !client.GetBookSocket().IsConnected {
		client.GetBookSocketPointer().Connect()
	}
	return nil
}

func PrivSocketGuard(client web_socket.WebSocketClient) error { // checks if socket is init and public. returns error if either is not true
	if !ws_client.IsInit(client.GetPrivSocket()) {
		return errors.New("KrakenClient is not initialized")
	}

	if !ws_client.IsPrivateSocket(client.GetPrivSocket()) {
		return errors.New("Public socket function called on a private socket")
	}

	/*if !client.GetPrivSocket().IsConnected {
		log.Println("priv socket not connected")
		client.GetPrivSocketPointer().Connect()
	}
	log.Println("priv socket connected")*/
	return nil
}

func AreChannelsInit(ws *web_socket.WebSocketClient) bool {
	if ws.GetPubChan() != nil && ws.GetPrivChan() != nil && ws.GetBookChan() != nil {
		return true
	} else {
		return false
	}
}
