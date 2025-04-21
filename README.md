# Enhanced TCP Echo Server

This project enhances a basic TCP echo server by adding concurrency, logging, command handling, timeout management, and other useful features. The server allows multiple clients to connect simultaneously and includes personality features for a more interactive experience.

## Features

- Concurrent Client Handling: Supports multiple clients simultaneously using goroutines  
- Connection Logging: Logs all connections and disconnections with timestamps  
- Message Logging: Saves each client's messages to individual log files  
- Graceful Disconnect Handling: Properly handles client disconnections without crashing  
- Input Processing: Trims whitespace from messages before processing  
- Configurable Port: Server port can be set via command-line flag  
- Inactivity Timeout: Automatically disconnects inactive clients after 30 seconds  
- Message Length Protection: Rejects messages longer than 1024 bytes  
- Custom Responses: Special responses for keywords like "hello" and "bye"  
- Command Protocol: Supports commands like /time, /quit, and /echo  

## How to Run the Server

1. Make sure you have Go installed on your system.  
2. Clone this repository or download the source code.  
3. Navigate to the project directory.  
4. Run the server:

   bash
   go run echo.go --port 4000

Optional command-line flags:  
- --port: Set the listening port (default: 4000)  
- --logdir: Set the directory for client logs (default: "logs")

## Testing the Server

1. Connect to the server using netcat or telnet:

   bash
   nc localhost 4000

  If you do not have netcat on your system run this in your terminal to install it
 
  bash
  sudo apt update
  sudo apt install netcat

  Run this command to confirm the installation

  nc -h


2. Send messages to test echo functionality.  

3. Try special commands:  
   - /time - Get current server time  
   - /echo message - Echo back just the message  
   - /quit - Close the connection  

4. Try personality features:  
   - Send "hello" → Server responds with "Greetings!"  
   - Send empty message → Server responds with "Tell me something..."  
   - Send "bye" → Server responds with "Goodbye!" and closes the connection  

## Demo Video

[Watch the Echo Server Demo on YouTube](https://youtu.be/aOqkaPPcA4Y)

## Educational Reflections

Most Educationally Enriching Feature:
The part where I learned the most was when I was required to implement concurrency using goroutines. Learning how to safely handle multiple client connections simultaneously provided deep insights into Go's concurrency model. The combination of goroutines with Wait Groups demonstrated how Go makes concurrent programming more approachable while ensuring proper resource management.

Feature Requiring Most Research:
Implementing the inactivity timeout required the most research. Understanding how to properly use timers, connection deadlines, and how to handle them in conjunction with I/O operations was challenging. I had to research Go's time package extensively, particularly the AfterFunc method and how to reset timers. To ensure this feature worked correctly without causing resource leaks or premature disconnections, I had to carefully study of Go's documentation and community resources.

## Future Improvements

- Add authentication functionality  
- Implement secure connections (TLS)  
- Create admin commands for server management  
- Add support for file transfers  
- Develop a GUI client for easier interaction
