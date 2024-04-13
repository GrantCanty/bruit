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
	var sec int64
	var s, m float64

	var stringTS string
	err := json.Unmarshal(d, &stringTS)
	if err != nil {
		err = json.Unmarshal(d, &sec)
		if err != nil {
			return err
		}
		u.Time = time.Unix(sec, 0).UTC()
	} else {
		tmp, err := strconv.ParseFloat(stringTS, 64)
		if err != nil {
			return err
		}
		s = tmp
		m = s - math.Floor(s)
		u.Time = time.Unix(int64(s), int64(m)).UTC()
	}
	u.Time.Unix()
	if err != nil {
		return err
	}

	return nil
}

/*func (u *UnixTime) UnmarshalJSON(d []byte) error {
	start := time.Now()

	var ts interface{}
	var sec, min float64

	ts = d

	switch tmp := ts.(type) {
	case string:
	json.Unmarshal(d, &tmp)
	log.Println(tmp)
	floatTime, err := strconv.ParseFloat(tmp, 64)
	if err != nil {
		return err
	}
	sec, min = math.Modf(floatTime)
	case float64:
	json.Unmarshal(d, &tmp)
	sec, min = math.Modf(tmp)
	case []uint8:
		log.Println(tmp)
		var num float64
		err := json.Unmarshal(d, &num)
		if err != nil {
			return err
		}
		sec, min = math.Modf(num)
	default:
		log.Println(reflect.TypeOf(ts))
		return errors.New("Unkown type of time in response")
	}

	u.Time = time.Unix(int64(sec), int64(min)).UTC()
	u.Time.Unix()

	log.Println("process time: ", time.Since(start))

	return nil
}*/
