package logger

/**
 * Created by Zf_D on 2015-02-28
 */
import (
    "fmt"
    "log"
    "os"
    "runtime"
    "strings"
    "strconv"
    "time"
    "sync"
    "io"
)

const (
    _ = iota //日志等级
    Lv_Debug //1
    Lv_Info //2
    Lv_Warn //3
    Lv_Error //4
)

const (
    _ = int64(1) << (iota * 10) //大小
    KB //1024
    MB //1048576
    GB //1073741824
    TB //1099511627776
)

const DateFormat = "2006-01-02 15-04-05"

var (
    projectPath = GetProjectPath() //项目路径
    fileDir string = "log" //文件夹
    fileName string = "log.log" //文件名
    consoleFlag bool = true //是否输出控制台
    logLevel int = Lv_Debug //日志等级
    backupType int = 1 //1为按大小备份，2为按时间备份
    backupSize int64 = 1 * MB //备份大小
    logFile *os.File //日志文件
    logger *log.Logger //日志对象
    lock *sync.RWMutex = new(sync.RWMutex) //锁
)

//设置文件名及路径
func SetFilePath(arg_fileDir string, arg_fileName string) {
    if arg_fileDir != "" {
        fileDir = arg_fileDir
    }
    if arg_fileName != "" {
        fileName = arg_fileName
        logger = getLogger()
    }
}

//设置等级
func SetLogLevel(arg_logLevel int) {
    if arg_logLevel >= Lv_Debug && arg_logLevel <= Lv_Error {
        logLevel = arg_logLevel
    }
}

//是否输出控制台
func SetConsoleFlag(arg_consoleFlag bool) {
    consoleFlag = arg_consoleFlag
}

//获取项目路径
func GetProjectPath() string {
    path, _ := os.Getwd()
    return strings.SplitN(path, "src", -1)[0]
}

func getLogger() *log.Logger {
    logFileDir := projectPath + "\\" + fileDir
    logFilePath := logFileDir + "\\" + fileName

    //创建文件夹
    err := os.MkdirAll(logFileDir, os.ModePerm)
    if err != nil {
        log.Fatal(err.Error())
        return nil
    }

    //打开或创建文件
    logFile, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
    if err != nil {
        log.Fatal(err.Error())
        return nil
    }

    return log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func console(s string) {
    if consoleFlag {
        _, file, line, _ := runtime.Caller(2)
        file = file[strings.LastIndex(file, "/")+1:]
        fmt.Println(file+":"+strconv.Itoa(line), s)
    }
}

func output(level int, suffix string, v []interface{}) {
    if logLevel <= level {
        lock.Lock()
        defer lock.Unlock()
        if logger == nil {
            logger = getLogger()
        }
        s := suffix
        for i := 0; i < len(v); i++ {
            s += " " + fmt.Sprintf("%v", v[i])
        }
        logger.Output(2, s)
        console(s)
    }
}

func Debug(v ...interface{}) {
    output(Lv_Debug, "[D]", v)
}

func Info(v ...interface{}) {
    output(Lv_Info, "[I]", v)
}

func Warn(v ...interface{}) {
    output(Lv_Warn, "[W]", v)
}

func Error(v ...interface{}) {
    output(Lv_Error, "[E]", v)
}

//设置按大小备份
func SetSizeBackup(tempSec int, size int64) {
    backupType = 1
    backupSize = size
    if tempSec > 0 {
        go func() {
            timer := time.NewTicker((time.Duration)(tempSec) * time.Second)
            for {
                select {
                case <-timer.C:
                    if 1 != backupType {
                        return
                    }
                    if fileSize() >= backupSize {
                        backup()
                    }
                }
            }
        }()
    }
}

//设置按时间备份
func SetDailyBackup(tempSec int) {
    backupType = 2
    if tempSec > 0 {
        go func() {
            timer := time.NewTicker((time.Duration)(tempSec) * time.Second)
            for {
                select {
                case <-timer.C:
                    if 2 != backupType {
                        return
                    }
                    backup()
                }
            }
        }()
    }
}

func fileSize() int64 {
    logFilePath := projectPath + "\\" + fileDir + "\\" + fileName
    fileInfo, err := os.Stat(logFilePath)
    if err != nil {
        log.Fatal(err.Error())
        return 0
    }
    return fileInfo.Size()
}

func copyFile(fromFilePath string, toFilePath string) {
    fromFile, err := os.Open(fromFilePath)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer fromFile.Close()

    toFile, err := os.Create(toFilePath)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer toFile.Close()

    _, err = io.Copy(toFile, fromFile)
    if err != nil {
        log.Fatal(err.Error())
    }
}

func backup() {
    lock.Lock()
    defer lock.Unlock()

    logFileDir := projectPath + "\\" + fileDir
    backupDir := logFileDir + "\\backup"
    logFilePath := logFileDir + "\\" +fileName
    backupFilePath := backupDir + "\\" + time.Now().Format(DateFormat) + ".log"

    //创建文件夹
    err := os.MkdirAll(backupDir, os.ModePerm)
    if err != nil {
        log.Fatal(err.Error())
    }

    //先关闭日志文件
    if logFile != nil {
        logFile.Close()
    }

    //备份旧的日志
    copyFile(logFilePath, backupFilePath)

    //清空，重新写入
    logFile, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_TRUNC|os.O_APPEND|os.O_CREATE, os.ModePerm)
    if err != nil {
        log.Fatal(err.Error())
    }
    logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
}
