package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var wg sync.WaitGroup

var goodReq int32
var badReq int32

var c = http.Client{
	Transport: &http.Transport{
		DisableKeepAlives:true,
		Dial: func(netw, addr string) (net.Conn, error) {
			deadline := time.Now().Add(15 * time.Second)
			c, err := net.DialTimeout(netw, addr, time.Second*10)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)

			return c, nil
		},
	},
}

func main() {

	//for i := 0; i < 100; i++ {
	//	wg.Add(1)
	//	if i%3 == 0 {
	//		go req("web", i)
	//	} else {
	//		go req("phone", i)
	//	}
	//
	//}

	for i := 0; i < 1000; i++ {
		//time.Sleep(time.Millisecond)
		wg.Add(1)
		go getReq(i)
	}

	wg.Wait()
	fmt.Printf("good: %d, bad: %d\n", goodReq, badReq)
}

func req(cc string, index int) {
	defer wg.Done()
	body := []byte(fmt.Sprintf(`
{
	"orderId":1,
	"userId": 1334,
	"amount": 100.01,
	"channel": "%s"
}
		`, cc) )

	reader := bytes.NewReader(body)

	response, err := c.Post("http://localhost:9091/order", "application/json", reader)
	if err != nil {
		fmt.Printf("index: %d, %s\n", index, err)
		atomic.AddInt32(&badReq, 1)
		return
	}
	if response != nil {
		defer response.Body.Close()

		res, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("index: %d, %s\n", index, res)
		atomic.AddInt32(&goodReq, 1)
	}

}

func getReq(index int)  {
	defer wg.Done()
	response, err := c.Get("http://localhost:9091/test")
	if err != nil {
		fmt.Printf("index: %d, %s\n", index, err)
		atomic.AddInt32(&badReq, 1)
		return
	}

	defer response.Body.Close()
	fmt.Printf("index: %d, %s\n", index, "good")
	atomic.AddInt32(&goodReq, 1)
}
