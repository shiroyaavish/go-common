package datatypes

import (
	"database/sql/driver"
	"time"
)

// UnixTimestampMillis Define custom type to represent Unix timestamp in milliseconds
type UnixTimestampMillis struct {
	time.Time
}

// Scan converts the database field to UnixTimestampMillis type
func (u *UnixTimestampMillis) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	unixTime := value.(int64) / 1000 // value is in ms
	u.Time = time.Unix(unixTime, 0)
	return nil
}

// Value converts UnixTimestampMillis type to a value that can be stored in the database
func (u UnixTimestampMillis) Value() (driver.Value, error) {
	if u.IsZero() {
		return nil, nil
	}
	return u.UnixNano() / int64(time.Millisecond), nil
}

// UnixTimestampNanos Define custom type to represent Unix timestamp in nanoseconds
type UnixTimestampNanos struct {
	time.Time
}

// Scan converts the database field to UnixTimestampNanos type
func (u *UnixTimestampNanos) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	unixTime := value.(int64) / 1000 // value is in ms
	u.Time = time.Unix(unixTime, 0)
	return nil
}

// Value converts UnixTimestampNanos type to a value that can be stored in the database
func (u UnixTimestampNanos) Value() (driver.Value, error) {
	if u.IsZero() {
		return nil, nil
	}
	return u.UnixNano(), nil
}
