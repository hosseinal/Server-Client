package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

const maxsize = 500
const maxwaiting = 5


type item struct {
	waitedTime  int
	value       int
	currentPoll int
	conn        net.Conn
}

type threadPollHelper struct {
	workerSize int
	mutex      sync.Mutex
	chn        chan bool
	items      []item
}

func main() {

	var firstPollItems []item
	var secondPollItems []item
	var doneItems []item

	firstChn := make(chan bool, maxsize)
	secondChn := make(chan bool, maxsize)


	poll1 := threadPollHelper{
		workerSize: 10,
		mutex:      sync.Mutex{},
		chn:        firstChn,
		items:      firstPollItems,
	}

	poll2 := threadPollHelper{
		workerSize: 0,
		mutex:      sync.Mutex{},
		chn:        secondChn,
		items:      secondPollItems,
	}

	threadPoll1(&poll1, &poll2)
	threadPoll2(&poll2 , &doneItems)

	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic("we got an error here")
	}
	defer li.Close()

	//we do not have while in golang
	go func() {
		for {
			conn, err := li.Accept()
			if err != nil {
				log.Println("here we miss a connection request ...")
			}
			scanner := bufio.NewScanner(conn)
			scanner.Scan()
			line := scanner.Text()
			value, err := strconv.Atoi(line)
			if err!= nil{
				fmt.Println("send bad request ...")
				continue
			}
			item := item{
				waitedTime: 0,
				value:      value,
				conn:       conn,
			}
			poll1.items = append(poll1.items, item)
			firstChn <- true
		}
	}()


	fmt.Println("to find out how many is all accepted request input 1 ")
	fmt.Println("to find out how many is all done request input 2 ")
	fmt.Println("to find out how many is all waited request in poll 1 input 3 ")
	fmt.Println("to find out how many is all waited request in poll 2 input 4 ")

	for {
		var input string
		_, _ = fmt.Scanln(&input)

		if input=="1"{
			var sum int
			poll1.mutex.Lock()
			sum = sum + len(poll1.items)
			poll1.mutex.Unlock()

			poll2.mutex.Lock()
			sum = sum +len(doneItems)
			sum = sum +len(poll1.items)
			poll2.mutex.Unlock()

			fmt.Println(sum)

		}else if input=="2"{
			poll2.mutex.Lock()
			fmt.Println(len(doneItems))
			poll2.mutex.Unlock()
		}else if input=="3"{
			poll1.mutex.Lock()
			fmt.Println(len(poll1.items))
			poll1.mutex.Unlock()
		}else if input=="4"{
			poll2.mutex.Lock()
			fmt.Println(len(poll2.items))
			poll2.mutex.Unlock()

		}
	}

}

func threadPoll1(helper *threadPollHelper, dist *threadPollHelper) {
	for i := 0; i < 5; i++ {
		go func() {
			for _ = range helper.chn {
				helper.mutex.Lock()
				i := scheduler(helper.items)
				citem := helper.items[i]
				helper.items = append(helper.items[:i], helper.items[i+1:]...)
				helper.mutex.Unlock()

				//do your job
				time.Sleep(400 * time.Millisecond)
				//end your job

				helper.mutex.Lock()
				dist.items = append(dist.items, citem)
				helper.mutex.Unlock()
				dist.chn <- true
			}
		}()
	}
}

func threadPoll2(helper *threadPollHelper , doneItems *[]item) {
	for i := 0; i < 6; i++ {
		go func() {
			for _ = range helper.chn {
				helper.mutex.Lock()
				i := scheduler(helper.items)
				citem := helper.items[i]
				helper.items = append(helper.items[:i], helper.items[i+1:]...)
				*doneItems = append(*doneItems, citem)
				helper.mutex.Unlock()

				//do your job
				time.Sleep(1000 * time.Millisecond)
				//end your job


				_ , err := fmt.Fprintln(citem.conn, citem.value)
				if err != nil {
					continue
				}
				fmt.Println("done")
			}
		}()
	}
}

func scheduler(items []item) int {
	if len(items) < 1 {
		return -1
	}
	//check for too much waited item
	for i, v := range items {
		if v.waitedTime > maxwaiting {
			return i
		}
	}

	//
	for i, v := range items {
		if v.value%2 == 0 {
			return i
		}
	}

	return 0
}
