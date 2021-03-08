package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var services map[string]string
var mux sync.Mutex

func loadService() {
	var serviceFilePath string
	if runtime.GOOS == "windows" {
		serviceFilePath = os.Getenv("WINDIR") + "\\system32\\drivers\\etc\\services"
	} else {
		serviceFilePath = "/etc/services"
	}
	content, err := os.ReadFile(serviceFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot find %s file.", serviceFilePath)
		os.Exit(1)
	}
	services = make(map[string]string)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		words := strings.Fields(line)
		services[words[1]] = words[0]
	}
}

func GetServiceName(protocol string, port uint16) string {
	mux.Lock()
	defer mux.Unlock()
	if services == nil {
		loadService()
	}

	val, ok := services[strconv.FormatUint(uint64(port), 10)+"/"+protocol]
	if !ok {
		return "<unknown>"
	}
	return val
}
