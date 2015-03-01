package logger

/**
 * Created by Zf_D on 2015-02-28
 */
import (
	"testing"
	"time"
)

func testOut() {
	Debug("Debug", "1", 2, 3.3)
	Info("Info", []string{"A1","B1"})
	Warn("Warn", [][]string{{"A1","B1"},{"A2","B2"}})
	Error("Error", "~")
}

func Test_output(t *testing.T) {
	defer logFile.Close()
	SetLogLevel(Lv_Error)
	testOut()
	SetLogLevel(Lv_Debug)
	testOut()
}

func Test_sizeBackup(t *testing.T) {
	defer logFile.Close()
	SetFilePath("log\\testSize", "log.log")
	SetSizeBackup(2, 1 * KB)
	timer := time.NewTicker(1 * time.Second)
	i := 0
	for {
		select {
		case <-timer.C:
			for n:=0; n<=100; n++ {
				Debug(i, n)
				Info(i, n)
				Warn(i, n)
				Error(i, n)
			}
			i++
			if i > 5 {
				return
			}
		}
	}
}

func Test_dailyBackup(t *testing.T) {
	defer logFile.Close()
	SetFilePath("log\\testDaily", "log.log")
	SetDailyBackup(2)
	timer := time.NewTicker(1 * time.Second)
	i := 0
	for {
		select {
		case <-timer.C:
			Debug(i)
			Info(i)
			Warn(i)
			Error(i)
			i++
			if i > 5 {
				return
			}
		}
	}
}
