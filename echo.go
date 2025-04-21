//Owner: Jeronimo Tzib

// Filename: main.go
// Purpose: This program demonstrates how to create a TCP network connection using Go

package main

//imported necessary packages

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (

	maxMessageLength = 1024 //max allowed message length
	inactivTimeout = 30 * time.Second //set the inactivity timeout (30 seconds)

)

//configuration struct to hold server settings
type Config struct{

Port int 
logDir string

}

//create a logging directory if it doesn't exist
func ensLogDir(logDir string) error {
	return os.MkdirAll(logDir, 0755)
}

//create a log file for a specific client
func createClientLogFile(clientAddr string, logDir string)(*os.File, error){

		//cleans the client address to make it suitable for a filename
		safeAddr := string.Replace(strings.Replace(clientAddr, ":", "_", -1), ".", "-", -1)
	logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", safeAddr))
return os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func logMessage(clientAddr string, logFile *os.File, format string, args ...interface{}){
	message := fmt.Sprintf(format,args...)
	log.Printf("[%s]%s", clientAddr, message)
	if logFile != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(logFile, "[%s] %s\n", timestamp, message)
	}
}

func worker(wg *sync.WaitGroup, tasks chan string, dialer net.Dialer) {
	defer wg.Done()
	maxRetries := 3
    for addr := range tasks {
		var success bool
		for i := range maxRetries {      
		conn, err := dialer.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			fmt.Printf("Connection to %s was successful\n", addr)
			success = true
			break
		}
		backoff := time.Duration(1<<i) * time.Second
		fmt.Printf("Attempt %d to %s failed. Waiting %v...\n", i+1,  addr, backoff)
		time.Sleep(backoff)
	    }
		if !success {
			fmt.Printf("Failed to connect to %s after %d attempts\n", addr, maxRetries)
		}
	}
}

func main() {

	var wg sync.WaitGroup
	tasks := make(chan string, 100)

    target := "scanme.nmap.org"

	dialer := net.Dialer {
		Timeout: 5 * time.Second,
	}
  
	workers := 100

    for i := 1; i <= workers; i++ {
		wg.Add(1)
		go worker(&wg, tasks, dialer)
	}

	ports := 512

	for p := 1; p <= ports; p++ {
		port := strconv.Itoa(p)
        address := net.JoinHostPort(target, port)
		tasks <- address
	}
	close(tasks)
	wg.Wait()
}