package config

import (
	"sync"
	"os"
	"strings"
)

const (
	configFileDir = `example.xml`
)

var (
	once *sync.Once
	configLock *sync.RWMutex
	configMap = make(map[string]string)
)

func Init() {
	once.Do(execute)
}

func execute() {
	readXml(configMap, getPjPath()+`\`+configFileDir)
}

func getPjPath() string {
	pjPath, _ := os.Getwd()
	return strings.SplitN(pjPath, `\src\`, -1)[0]
}

func GetConfig(key string) string {
	configLock.Lock()
	defer configLock.Unlock()
	return configMap[key]
}

func SetConfig(key string, value string) {
	configLock.RLock()
	defer configLock.RUnlock()
	configMap[key] = value
}
