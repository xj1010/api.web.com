package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Time is alias type for time.Time
type MyTime struct {
	time.Time
}

const (
	timeFormat = "2006-01-02 15:04:05"
	zone        = "Asia/Shanghai"
)

/*// UnmarshalJSON implements json unmarshal interface.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = Time(now)
	return
}*/

// MarshalJSON implements json marshal interface.
func (t MyTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(timeFormat))
	return []byte(formatted), nil

}

func (t MyTime) String() string {
	return t.Time.Format(timeFormat)
}

func (t MyTime) local() time.Time {
	loc, _ := time.LoadLocation(zone)
	return t.Time.In(loc)
}

// Value ...
func (t MyTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time 注意是指针类型 method
func (t *MyTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = MyTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}