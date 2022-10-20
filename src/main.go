package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/shirou/gopsutil/process"
)

const (
	filePath    string = "FILE_PATH"
	processName string = "PROCESS_NAME"
)

func main() {

	if isEnvExist(filePath) == false || isEnvExist(processName) == false {
		panic("You should set FILE_PATH and PROCESS_NAME ENVs!!")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if !isValidEvent(event) {
					continue
				}
				log.Println("event:", event)
				KillProcess(os.Getenv(processName))
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(os.Getenv(filePath))
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}

func KillProcess(name string) {
	processes, err := process.Processes()
	if err != nil {
		e := fmt.Errorf("%v", err)
		fmt.Printf("Error: %v\n", e)
	}
	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			e := fmt.Errorf("%v", err)
			fmt.Printf("Error: %v\n", e)
		}
		if n == name {
			p.Kill()
			fmt.Printf("Process %v has been killed. \n", p.Pid)
		}
	}
}

func isEnvExist(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return false
}

func isValidEvent(event fsnotify.Event) bool {
	if event.Op&fsnotify.Create != fsnotify.Create {
		return false
	}
	if filepath.Base(event.Name) != "..data" {
		return false
	}
	return true
}
