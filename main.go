package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func sendTextFromFile(conn net.Conn, filePath string, interval time.Duration, preserveFormat bool) {
	defer conn.Close()

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if preserveFormat {
		for scanner.Scan() {
			line := scanner.Text()

			for _, word := range strings.Fields(line) {
				conn.Write([]byte(word + " "))

				time.Sleep(interval)
			}

			conn.Write([]byte("\n"))
		}
	}

	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		for _, word := range words {
			conn.Write([]byte(word + " "))

			time.Sleep(interval)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	host := "127.0.0.1"
	port := 8080
	filepath := "shakespeare.txt"
	wordsPerMinute := 250
	interval := time.Second / time.Duration(wordsPerMinute/60)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	log.Printf("Server running on %s:%d", host, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error: %s", err)
			continue
		}

		log.Printf("Connected to %s", conn.RemoteAddr())

		go sendTextFromFile(conn, filepath, interval, true)
	}

}
