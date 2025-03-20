package types

type Subscribe struct {
	Event        string      `json:"event"`
	Pair         []string    `json:"pair,omitempty"`
	Subscription interface{} `json:"subscription,omitempty"`
	Token        string      `json:"token,omitempty"`
}

type SubscribeV2 struct {
	Method string   `json:"method"`
	Params ParamsV2 `json:"params"`
}

type ParamsV2 struct {
	Channel string   `json:"channel"`
	Symbol  []string `json:"symbol"`
	Depth   int      `json:"depth"`
}

type NameAndToken struct {
	Name  string `json:"name"`
	Token string `json:"token,omitempty"`
}

type CancelOrder struct {
	Event string   `json:"event"`
	Token string   `json:"token"`
	Txid  []string `json:"txid"`
}

type OHLCSubscription struct {
	Interval int    `json:"interval"`
	Name     string `json:"name"`
}

type BookSubscription struct {
	Depth int    `json:"depth"`
	Name  string `json:"name"`
}

type Ping struct {
	Event string `json:"event"`
	Id    int    `json:"reqid"`
}

type Order struct {
	WsToken   string `json:"token"`
	Event     string `json:"event"`     //addOrder, editOrder...
	OrderType string `json:"ordertype"` //market, limit...
	TradeType string `json:"type"`      //buy, sell
	Pair      string `json:"pair"`
	Volume    string `json:"volume"`
	Price     string `json:"price,omitempty"`  //omit if market order
	Oflags    string `json:"oflags,omitempty"` //fcib: fee in base currency, post for limit orders...
	Validate  string `json:"validate,omitempty"`
}
