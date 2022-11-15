package decoders

import "log"

func BookJsonDecoder(response string, testing bool) []interface{} {
	var resp []interface{}
	//byteResponse := []byte(response)
	log.Println(response)

	/*_, err := hbResponseDecoder(byteResponse, testing)
	if err != nil {
		resp, err = rr()
	}*/

	return resp
}
