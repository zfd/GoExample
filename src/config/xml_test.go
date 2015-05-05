package config

import (
	"testing"
	"fmt"
)

func Test_readXml(t *testing.T) {
	const result =  `map[server>address:127.0.0.1 server>port:8080 log>file>name:log.log log>file>size:1024]`
	readXml(configMap, `example.xml`)
	if fmt.Sprintf("%s", configMap) != result {
		t.Error(fmt.Sprintf("%s", configMap), " \nNot Equal\n", result)
	}
}
