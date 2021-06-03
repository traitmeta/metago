package time

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Timestamp time.Time
type Timestamp time.Time

// Unix get the unix time
func (t *Timestamp) Unix() int64 {
	tt := time.Time(*t)
	if tt.IsZero() {
		return 0
	}

	return tt.Unix()
}

// MarshalJSON transfer timestamp to byte
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	stamp := fmt.Sprint(ts)
	return []byte(stamp), nil
}

// UnmarshalJSON transfer byte to timestamp
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(int64(ts), 0))
	return nil
}

// GetBSON implements bson.Getter
func (t Timestamp) GetBSON() (interface{}, error) {
	return struct {
		Value int64
	}{
		Value: time.Time(t).Unix(),
	}, nil
}

// SetBSON implements bson.Setter
func (t *Timestamp) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		Value int64
	})
	if err := raw.Unmarshal(decoded); err != nil {
		return err
	}

	*t = Timestamp(time.Unix(decoded.Value, 0))
	return nil
}

func (t *Timestamp) Value() (driver.Value, error) {
	return time.Time(*t).Format("2006-01-02 15:04:05"), nil
}

func (t *Timestamp) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		*t = Timestamp(vt)
	default:
		return errors.New("type error")
	}

	return nil
}
