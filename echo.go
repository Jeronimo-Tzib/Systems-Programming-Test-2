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

//log a message to both console and the client's log file

func logMessage(clientAddr string, logFile *os.File, format string, args ...interface{}){
	message := fmt.Sprintf(format,args...)
	log.Printf("[%s]%s", clientAddr, message)
	if logFile != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(logFile, "[%s] %s\n", timestamp, message)
	}
}

//handle a client connection
func handleClient(conn net.Conn, wg *sync.WaitGroup, config Config) {
	defer wg.Done()
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()

	//create client log file
	logFile, err := createClientLogFile(clientAddr, config.LogDir)
	if err != nil {
		log.Printf("Error creating log file for %s: %v", clientAddr,err)
	} else{
		defer logFile.Close()
	}

	logMessage(clientAddr, logFile, "Connected")

	//set up a buffered reader
	reader := bufio.NewReader(conn)

	idleTimer := time.AfterFunc(inactivTimeout, func(){

		logMessage(clientAddr, logFile, "Disconnected due to inactivity")
		conn.Close()

	})

	defer idleTimer.Stop()

	for {
		//reset the timer on each iteration
		idleTimer.Reset(inactivTimeout)

		//set read deadline to handle potential blocking
		conn.SetReadDeadline(time.Now().Add(inactivTimeout))

		//read message from client (up to newline or max length)
		message, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection"){
				logMessage(clientAddr, logFile, "Disconnected")
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout(){
				logMessage(clientAddr, logFIle, "Disconnected due to timeout")
			}else {
				logMessage(clientAddr, logFile, "Error reading: %v", err)
			}
			break
			}

			//time whitespace and newlines
			message = strings.TrimSpace(message)

			// check message length
			if len(message) > maxMessageLength{
				respone := fmt.Sprintf("Error: Message too long (max %d bytes)\n", maxMessageLength)
				conn.Write([]byte(response))
				logMessage(clientAddr, logFile, "Mesage too long (%d bytes)", len(message))
				continue
			}

			//log the received message
			logMessage(clientAddr, logFile, "Received: %s", message)

			//process message
			response := processMessage(message, conn, clientAddr, logFile)

			//check if we should close the connection
			if response == "" {
				break // empty response idicates connection should be closed
			}

			//send response to client
			_, err = conn.Write([]byte(response + "\n"))
			if err != nil {
				logMessage(clientAddr, logFile, "Error writing: %v, err")
				break
				}
			}
		}

//process a message and generate response based on the message content
func processMessage(message string, conn net.Conn, clientAddr string, logFile *os.File) string{

//handle empty message
if message == ""{
	return "Tell me something..."
}

//check for command protocol messages
if strings.HasPrefix(message, "/"){
	cmd := strings.SplitN(message, " ", 2)
	switch cmd[0]{
	case "/time":
		return time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
	case "/quit":
		logMessage(clientAddr, logFile, "Client requested to quit")
		return "" // will close the connection
	case "/echo":
		if len(cmd) > 1 {
			return cmd[1] //echo back just the message part
		}
		return "" //no message to echo

	default:
		return "Unknown command: " + cmd[0]
	}
}

//handle personality mode
switch strings.ToLower(message){
case "hello":
	return "Greetings!"
case "bye":
	logMessage(clientAddr, logFile, "Client said bye")
	conn.Write([]byte("Goodbye!\n"))
	return "" //will close the connection
default:
	//default echo behavior
	return message
}

}

