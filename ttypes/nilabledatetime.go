// A quickly mysql access component.
// Copyright 2023 The daog Authors. All rights reserved.

package ttypes

import (
	"bytes"
	"database/sql"
	"github.com/rolandhe/daog"
	"strings"
	"time"
)

// NilableDatetime 可以为空的时间类型，扩展 sql.NullTime, 支持格式及json输出，格式由 DatetimeFormat 变量指定
type NilableDatetime struct {
	sql.NullTime
}

// FromDatetime 转换 time.Time 类型为 NilableDatetime 类型，返回值为指针，如果接收变量为 NilableDatetime 类型，需要使用加*号来解引用: *nilableDatetime
func FromDatetime(d time.Time) *NilableDatetime {
	return &NilableDatetime{
		sql.NullTime{d, true},
	}
}

// String 实现 fmt.Stringer 接口
func (ndt NilableDatetime) String() string {
	if !ndt.Valid {
		return "<nil>"
	}
	return ndt.Time.Format(DatetimeFormat)
}

// UnmarshalJSON 实现 json.Unmarshaler
func (s *NilableDatetime) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		s.Valid = false
		return nil
	}
	if bytes.Compare(b, nullJsonValue) == 0 {
		s.Valid = false
		return nil
	}
	value := strings.Trim(string(b), `"`)       //get rid of "
	t, err := time.Parse(DatetimeFormat, value) //parse time
	if err != nil {
		daog.GLogger.SimpleLogError(err)
		return err
	}
	s.Time = t
	s.Valid = true
	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (s NilableDatetime) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return []byte(`"` + s.Time.Format(DatetimeFormat) + `"`), nil
}

func (s *NilableDatetime) ToTimePointer() *time.Time {
	if !s.Valid {
		return nil
	}
	return &s.Time
}
