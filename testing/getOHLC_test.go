package testing

import (
	"bruit/bruit/clients/kraken"
	rest "bruit/bruit/clients/kraken/rest_client"
	"bruit/bruit/clients/kraken/types"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOHLC(t *testing.T) {
	resp := "{'error':[],'result':{'KAVAUSD':[[1710519360,'0.9535','0.9538','0.9535','0.9538','0.9536','14380.69231581',15]],,'last':1710562440}}"
	//var expectedResp types.RestResp

	/*expectedResp, err := json.Marshal(resp)
	if err != nil {
		t.Errorf("expected err to be nil but got %v", err)
	}*/

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, resp)
	}))
	defer svr.Close()

	v := rest.ReturnNonceValues()

	var c kraken.KrakenClient

	var s *types.RestResp

	res, err := c.Rest.PublicRequest(svr.URL, v, s)
	if err != nil {
		t.Errorf("expected err to be nil but got %v", err)
	}

	/*if res != expectedResp {
		t.Errorf("expected res to be %v but got %v", expectedResp, res)
	}*/

	fmt.Println(res, resp)

}
