package convert

import (
	"math/big"
	"strconv"
)

// StrConvert 字符串的convert类型
type StrConvert string

// NewConvert 新建一个convert类型
func NewConvert(v string) *StrConvert {
	if v == "" {
		return nil
	}
	val := StrConvert(v)
	return &val
}

// Close 继承closer接口
func (s *StrConvert) Close() error {
	s = nil
	return nil
}

// String 返回原始对象
func (s *StrConvert) String() string {
	return string(*s)
}

// Bool 返回bool对象
func (s *StrConvert) Bool() (bool, error) {
	return strconv.ParseBool(s.String())
}

// Int 返回int对象
func (s *StrConvert) Int() (int, error) {
	v, e := strconv.ParseInt(s.String(), 10, 32)
	return int(v), e
}

// Int8 返回对应的int对象
func (s *StrConvert) Int8() (int8, error) {
	v, e := strconv.ParseInt(s.String(), 10, 8)
	return int8(v), e
}

// Int16 返回对应的int对象
func (s *StrConvert) Int16() (int16, error) {
	v, e := strconv.ParseInt(s.String(), 10, 16)
	return int16(v), e
}

// Int32 返回对应的int对象
func (s *StrConvert) Int32() (int32, error) {
	v, e := strconv.ParseInt(s.String(), 10, 32)
	return int32(v), e
}

// Int64 处理int64的对象,防止bigint的溢出
func (s *StrConvert) Int64() (int64, error) {
	v, e := strconv.ParseInt(s.String(), 10, 64)
	if e != nil {
		bigInt := &big.Int{}
		val, ok := bigInt.SetString(s.String(), 10)
		if !ok {
			return v, e
		}
		return val.Int64(), nil
	}
	return int64(v), e
}

// Uint 返回对应的int对象
func (s *StrConvert) Uint() (uint, error) {
	v, e := strconv.ParseUint(s.String(), 10, 64)
	return uint(v), e
}

// Uint8 返回对应的int对象
func (s *StrConvert) Uint8() (uint8, error) {
	v, e := strconv.ParseUint(s.String(), 10, 8)
	return uint8(v), e
}

// Uint16 返回对应的int对象
func (s *StrConvert) Uint16() (uint16, error) {
	v, e := strconv.ParseUint(s.String(), 10, 16)
	return uint16(v), e
}

// Uint32 返回对应的int对象
func (s *StrConvert) Uint32() (uint32, error) {
	v, e := strconv.ParseUint(s.String(), 10, 32)
	return uint32(v), e
}

// Uint64 返回对应的uint64对象
func (s *StrConvert) Uint64() (uint64, error) {
	v, e := strconv.ParseUint(s.String(), 10, 64)
	if e != nil {
		bigInt := &big.Int{}
		val, ok := bigInt.SetString(s.String(), 10)
		if !ok {
			return v, e
		}
		return val.Uint64(), nil
	}
	return uint64(v), e
}

// Float32 返回对应的float对象
func (s *StrConvert) Float32() (float32, error) {
	v, e := strconv.ParseFloat(s.String(), 32)
	return float32(v), e
}

// Float64 返回对应的float对象
func (s *StrConvert) Float64() (float64, error) {
	v, e := strconv.ParseFloat(s.String(), 3642)
	return float64(v), e
}
