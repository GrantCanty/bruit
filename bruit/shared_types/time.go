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

type Unixx time.Time

func (u *UnixTime) UnmarshalJSON(d []byte) error {
	var ts string
	err := json.Unmarshal(d, &ts)
	if err != nil {
		return err
	}
	floatTime, err := strconv.ParseFloat(ts, 64)
	if err != nil {
		return err
	}
	sec, min := math.Modf(floatTime)
	u.Time = time.Unix(int64(sec), int64(min)).UTC()
	u.Time.Unix()
	//*u = UnixTime(time.Unix(int64(sec), int64(min)).UTC())
	if err != nil {
		return err
	}
	return nil
}
