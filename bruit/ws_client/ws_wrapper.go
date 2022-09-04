package ws_client

// GENERAL FUNCTIONS

/*func (socket *Socket) ConnectToServer(wg *sync.WaitGroup, server string, testing bool) {
	defer wg.Done()

	if server != "public" && server != "private" {
		return
		//return nil, errors.New("Server must be 'public' or 'private'")
	}
	//var socket Socket
	if server == "public" {
		*socket = New(stored_data.Pub_ws_url)
	} else {
		*socket = New(stored_data.Priv_ws_url)
	}

	socket.OnConnected = func(socket Socket) {
		log.Println("Connected to server")
	}
	socket.OnTextMessage = func(message string, socket Socket) {
		PubJsonDecoder(message, testing)
		log.Println(message)
	}

	socket.Connect()

	//return &socket, nil
	return
}

func ConnectToServer(wg *sync.WaitGroup, server string, testing bool) (*Socket, error) {
	defer wg.Done()

	if server != "public" && server != "private" {
		return nil, errors.New("Server must be 'public' or 'private'")
	}
	var socket Socket
	if server == "public" {
		socket = New(stored_data.Pub_ws_url)
	} else {
		socket = New(stored_data.Priv_ws_url)
	}

	socket.OnConnected = func(socket Socket) {
		log.Println("Connected to server")
	}
	socket.OnTextMessage = func(message string, socket Socket) {
		PubJsonDecoder(message, testing)
		log.Println(message)
	}

	socket.Connect()

	return &socket, nil
}

func ConnectToServer2(wg *sync.WaitGroup, server string, testing bool) *Socket {
	defer wg.Done()

	if server != "public" && server != "private" {
		panic(errors.New("Server must be 'public' or 'private'"))
		//return nil
	}
	var socket Socket
	if server == "public" {
		socket = New(stored_data.Pub_ws_url)
	} else {
		socket = New(stored_data.Priv_ws_url)
	}

	socket.OnConnected = func(socket Socket) {
		log.Println("Connected to server")
	}
	socket.OnTextMessage = func(message string, socket Socket) {
		PubJsonDecoder(message, testing)
		log.Println(message)
	}

	socket.Connect()

	return &socket
}

func (socket *Socket) Ping(num int) {
	sub, _ := json.Marshal(&types.Ping{
		Event: "ping",
		Id:    num,
	})
	socket.SendBinary(sub)
}

// PUBLIC FUNCTIONS

func (socket *Socket) SubscribeToOHLC(wg *sync.WaitGroup, pairs []string, interval int) {
	defer wg.Done()

	sub, _ := json.Marshal(&types.Subscribe{
		Event: "subscribe",
		Subscription: &types.OHLCSubscription{
			Interval: interval,
			Name:     "ohlc",
		},
		Pair: pairs,
	})
	socket.SendBinary(sub)
}

func (socket *Socket) PubListen(ctx context.Context, wg *sync.WaitGroup, ch chan interface{}, testing bool) {
	defer wg.Done()
	defer socket.Close()
	var res interface{}

	socket.receiveMu.Lock()
	socket.OnTextMessage = func(message string, socket Socket) {
		res = PubJsonDecoder(message, testing)
		ch <- res
	}
	socket.receiveMu.Unlock()

	<-ctx.Done()
	log.Println("closing public socket")
	return
}

func (socket *Socket) ListenToOHLC(ctx context.Context, wg *sync.WaitGroup, testing bool) {
	defer wg.Done()
	defer socket.Close()

	socket.receiveMu.Lock()
	socket.OnTextMessage = func(message string, socket Socket) {
		log.Println(message)
		bMessage := []byte(message)
		err := ohlcRespDecoder(bMessage, testing)
		if err != nil {
			_, err = serverConnectionStatusResponseDecoder(bMessage, testing)
		}
	}
	socket.receiveMu.Unlock()

	<-ctx.Done()
	log.Println("Closing OHLC Socket")
}

func BookJsonDecoder(response string, testing bool) []interface{} {
	var resp []interface{}
	byteResponse := []byte(response)
	//log.Println(byteResponse)
	_, err := hbResponseDecoder(byteResponse, testing)
	if err != nil {
		resp, err = rr()
	}

	return resp
}

func rr() ([]interface{}, error) {
	var resp []interface{}
	return resp, nil
}

func BookReader(response string, testing bool) {
	log.Println(response)
}

// PRIVATE FUNCTIONS

func (socket *Socket) SubscribeToOpenOrders(token string) {
	sub, _ := json.Marshal(&types.Subscribe{
		Event: "subscribe",
		Subscription: &types.NameAndToken{
			Name:  "openOrders",
			Token: token,
		},
	})
	socket.SendBinary(sub)
}

func (socket *Socket) CancelAll(token string) {
	sub, _ := json.Marshal(&types.Subscribe{
		Event: "cancelAll",
		Token: token,
	})
	socket.SendBinary(sub)
}

// find a way to ad tradeID
func (socket *Socket) CancelOrder(token string, tradeIDs []string) {
	sub, _ := json.Marshal(&types.CancelOrder{
		Event: "cancelOrder",
		Token: token,
		Txid:  tradeIDs,
	})
	socket.SendBinary(sub)
}

func (socket *Socket) AddOrder(token string, otype string, ttype string, pair string, vol string, price string, testing bool) {
	test := strconv.FormatBool(testing)
	data := &types.Order{
		WsToken:   token,
		Event:     "addOrder",
		OrderType: otype,
		TradeType: ttype,
		Pair:      pair,
		Volume:    vol,
		Price:     price,
		Validate:  test,
	}

	order, err := json.Marshal(data)
	if err != nil {
		return
	}
	socket.SendBinary(order)
}

func (socket *Socket) PrivListen(ctx context.Context, wg *sync.WaitGroup, ch chan interface{}, testing bool) {
	defer wg.Done()
	defer socket.Close()
	var res interface{}

	socket.OnTextMessage = func(message string, socket Socket) {
		res = PrivJsonDecoder(message, testing)
		ch <- res
	}
	<-ctx.Done()
	log.Println("closing private listen")
	return
}

// METHODS

func PubJsonDecoder(response string, testing bool) interface{} {
	var resp interface{}
	byteResponse := []byte(response)

	resp, err := ohlcResponseDecoder(byteResponse, testing)
	if err != nil {
		resp, err = hbResponseDecoder(byteResponse, testing)
		if err != nil {
			resp, err = serverConnectionStatusResponseDecoder(byteResponse, testing)
			if err != nil {
				resp, err = ohlcSubscriptionResponseDecoder(byteResponse, testing)
				if err != nil {
					log.Println(string("\033[31m"), "received response of unknown data type: ", response)
				}
			}
		}
	}

	//if testing == true {
	//	log.Println("asfasfdafdadf", reflect.TypeOf(resp), resp)
	//}
	return resp
}

func PrivJsonDecoder(response string, testing bool) interface{} {
	var resp interface{}
	byteResponse := []byte(response)

	resp, err := openOrdersResponseDecoder(byteResponse, testing)
	if err != nil {
		resp, err = hbResponseDecoder(byteResponse, testing)
		if err != nil {
			resp, err = cancelOrderResponseDecoder(byteResponse, testing)
			if err != nil {
				resp, err = serverConnectionStatusResponseDecoder(byteResponse, testing)
			}
		}
	}

	//if testing == true {
	//	log.Println(reflect.TypeOf(resp), resp)
	//}
	return resp
}

func serverConnectionStatusResponseDecoder(byteResponse []byte, testing bool) (*types.ServerConnectionStatusResponse, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing == true {
		log.Println("in connection status func")
	}
	var conn types.ServerConnectionStatusResponse
	err := decoder.Decode(&conn)
	if err != nil {
		if testing == true {
			log.Println("connection status error: ", err)
		}
		return nil, err
	}
	return &conn, err
}

func ohlcSubscriptionResponseDecoder(byteResponse []byte, testing bool) (*types.OHLCSuccessResponse, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing == true {
		log.Println("in ohlc subscription response func")
	}
	var ohlc types.OHLCSuccessResponse
	err := decoder.Decode(&ohlc)
	if err != nil {
		if testing == true {
			log.Println("ohlc subscription response error", err)
		}
		return nil, err
	}
	return &ohlc, err
}

func hbResponseDecoder(byteResponse []byte, testing bool) (*types.HeartBeat, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing == true {
		log.Println("in hb response func")
	}
	var heart types.HeartBeat
	err := decoder.Decode(&heart)
	if err != nil {
		if testing == true {
			log.Println(err)
		}
		return nil, err
	}
	return &heart, err
}

func ohlcResponseDecoder(byteResponse []byte, testing bool) (*types.OHLCResponse, error) {
	if testing == true {
		log.Println("in ohlcResponse func")
	}

	var resp types.OHLCResponse
	err := json.Unmarshal(byteResponse, &resp)
	if err != nil {
		if testing == true {
			log.Println("ohlcResponse error: ", err)
		}
		return nil, err
	}
	return &resp, err
}

func ohlcRespDecoder(byteResponse []byte, testing bool) error {
	if testing == true {
		log.Println("in ohlcResponse func")
	}

	var resp types.OHLCResponse
	err := json.Unmarshal(byteResponse, &resp)
	if err != nil {
		if testing == true {
			log.Println("ohlcResponse error: ", err)
		}
		return err
	}
	return err
}

func openOrdersResponseDecoder(byteResponse []byte, testing bool) (*types.OpenOrdersResponse, error) {
	if testing == true {
		log.Println("in openOrdersResponse func")
	}

	var resp types.OpenOrdersResponse
	err := json.Unmarshal(byteResponse, &resp)
	if err != nil {
		if testing == true {
			log.Println("openOrdersResponse error: ", err)
		}
		return nil, err
	}
	return &resp, err
}

func cancelOrderResponseDecoder(byteResponse []byte, testing bool) (*types.CancelOrderResponse, error) {
	if testing == true {
		log.Println("in open orders response func")
	}

	var resp types.CancelOrderResponse
	err := json.Unmarshal(byteResponse, &resp)
	if err != nil {
		if testing == true {
			log.Println("open orders response error: ", err)
		}
		return nil, err
	}
	return &resp, err
}*/
