package shared_types

import (
	"encoding/json"
	"math"
	"strconv"
	"time"
)

type UnixTime struct {
	time.Time
}

//starts by trying to unmarsal using timestamp as a string. if it fails, trying to unmarshal using ts as float
func (u *UnixTime) UnmarshalJSON(d []byte) error {
	var stringTS string
	err := json.Unmarshal(d, &stringTS)
	if err != nil {
		var floatTS float64
		err = json.Unmarshal(d, &floatTS)
		if err != nil {
			return err
		}
		sec, min := math.Modf(floatTS)
		u.Time = time.Unix(int64(sec), int64(min)).UTC()
		u.Time.Unix()
		if err != nil {
			return err
		}
		return nil
	}

	floatTime, err := strconv.ParseFloat(stringTS, 64)
	if err != nil {
		return err
	}
	sec, min := math.Modf(floatTime)
	u.Time = time.Unix(int64(sec), int64(min)).UTC()
	u.Time.Unix()
	if err != nil {
		return err
	}
	return nil

	// below is old func

	//var ts string
	/*var ts float64
	err := json.Unmarshal(d, &ts)
	if err != nil {
		log.Println("error hereeee")
		return err
	}
	floatTime, err := strconv.ParseFloat(ts, 64)
	if err != nil {
		return err
	}
	sec, min := math.Modf(ts)
	u.Time = time.Unix(int64(sec), int64(min)).UTC()
	u.Time.Unix()
	if err != nil {
		return err
	}
	return nil*/
}
