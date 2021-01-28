package qlog

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

// 应用唯一标识符
func FieldAid(value string) Field {
	return String("aid", value)
}

// 模块
func FieldMod(value string) Field {
	value = strings.Replace(value, " ", ".", -1)
	return String("mod", value)
}

// 依赖的实例名称。以mysql为例，"dsn = "root:juno@tcp(127.0.0.1:3306)/juno?charset=utf8"，addr为 "127.0.0.1:3306"
func FieldAddr(value string) Field {
	return String("addr", value)
}

// FieldAddrAny ...
func FieldAddrAny(value interface{}) Field {
	return Any("addr", value)
}

// FieldName ...
func FieldName(value string) Field {
	return String("name", value)
}

// FieldType ...
func FieldType(value string) Field {
	return String("type", value)
}

// FieldCode ...
func FieldCode(value int32) Field {
	return Int32("code", value)
}

// 耗时时间
func FieldCost(value time.Duration) Field {
	return String("cost", fmt.Sprintf("%.3f", float64(value.Round(time.Microsecond))/float64(time.Millisecond)))
}

// FieldKey ...
func FieldKey(value string) Field {
	return String("key", value)
}

// 耗时时间
func FieldKeyAny(value interface{}) Field {
	return Any("key", value)
}

// FieldValue ...
func FieldValue(value string) Field {
	return String("value", value)
}

// FieldValueAny ...
func FieldValueAny(value interface{}) Field {
	return Any("value", value)
}

// FieldErrKind ...
func FieldErrKind(value string) Field {
	return String("errKind", value)
}

// FieldErr ...
func FieldErr(err error) Field {
	return zap.Error(err)
}

// FieldErr ...
func FieldStringErr(err string) Field {
	return String("err", err)
}

// FieldExtMessage ...
func FieldExtMessage(vals ...interface{}) Field {
	return zap.Any("ext", vals)
}

// FieldStack ...
func FieldStack(value []byte) Field {
	return ByteString("stack", value)
}

// FieldMethod ...
func FieldMethod(value string) Field {
	return String("method", value)
}

// FieldEvent ...
func FieldEvent(value string) Field {
	return String("event", value)
}
