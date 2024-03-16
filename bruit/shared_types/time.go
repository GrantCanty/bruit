package shared_types

import (
	"encoding/json"
	"math"
	"time"
)

type UnixTime struct {
	time.Time
}

func (u *UnixTime) UnmarshalJSON(d []byte) error {
	var ts float64
	err := json.Unmarshal(d, &ts)
	if err != nil {
		return err
	}
	sec, min := math.Modf(ts)
	u.Time = time.Unix(int64(sec), int64(min)).UTC()
	u.Time.Unix()
	if err != nil {
		return err
	}
	return nil
}
