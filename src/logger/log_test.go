package logger

import (
	"testing"
)

func Test_logger(t *testing.T) {
	Debug("D", "I", "W", "E")
	Debug(Level_D,Level_I,Level_W,Level_E)
	for i := 1; i <= 4; i++ {
		SetLogLevel(i)
		Debug(i)
		Info(i)
		Warn(i)
		Error(i)
	}
}
