package main

import (
	"fmt"
	"log"
	"os"

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
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
					KillProcess(os.Getenv(processName))
				}
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
