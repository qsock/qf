package qlog

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func BenchmarkWriteSyncer(b *testing.B) {
	b.Run("write file with no buffer", func(b *testing.B) {
		file, err := ioutil.TempFile("", "log")
		assert.NoError(b, err)
		defer file.Close()
		defer os.Remove(file.Name())

		w := zapcore.AddSync(file)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w.Write([]byte("foobarbazbabble"))
			}
		})
	})
	b.Run("write file with buffer", func(b *testing.B) {
		file, err := ioutil.TempFile("", "log")
		assert.NoError(b, err)
		defer file.Close()
		defer os.Remove(file.Name())

		w, close := Buffer(zapcore.AddSync(file), 0, 0)
		defer close()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w.Write([]byte("foobarbazbabble"))
			}
		})
	})
}
