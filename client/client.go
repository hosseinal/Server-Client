package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
)

const chanSize = 5
const minRand = 1000
const maxRand = 9999

func main() {

	reqChn := make(chan int, chanSize)

	fmt.Println("for exit input 0")
	fmt.Println("for send request input 1")
	fmt.Println("to find out how many request can you send input 2")

	for {
		var input string
		_, _ = fmt.Scanln(&input)

		if input == "0" {
			break
		} else if input == "1" {
			select {
			case reqChn <- 1:
				go sendRequest(reqChn)
				fmt.Println("--------------------")
				fmt.Println("sending request ...")
			default:
				fmt.Println("--------------------")
				fmt.Println("you requested too much")
			}
		}else if input=="2" {
			fmt.Println("you can send ",chanSize-len(reqChn))
		}
	}
}

func sendRequest(rc chan int) {
	fmt.Println("--------------------")
	fmt.Println("start sending request...")
	conn, err := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	if err != nil {
		log.Panic("you cannot connect to server")
	}

	randNumber := rand.Intn(maxRand - minRand)
	randNumber = randNumber + minRand

	fmt.Println("--------------------")
	fmt.Println("start sending request " , randNumber)

	_, err = fmt.Fprintln(conn, randNumber)
	if err != nil {
		fmt.Println("error in sending data")
	}

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	line := scanner.Text()
	fmt.Println("--------------------")
	fmt.Println("recive :" ,line)
	<-rc
}
