package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type contextKey int

const (
	VarsKey contextKey = iota
	PathKey
	BaseUrlKey
)

func TestServer_handleGetUser(t *testing.T) {
	s := NewServer()

	ts := httptest.NewServer(makeHTTPHandleFunc(s.handleGetUser))
	nreq := 1000
	wg := &sync.WaitGroup{}

	for i := 0; i < nreq; i++ {
		wg.Add(1)
		go func(i int) {
			id := i%100 + 1
			url := fmt.Sprintf("%s/?id=%d", ts.URL, id)
			resp, err := http.Get(url)
			if err != nil {
				t.Error(err)
			}
			user := &User{}
			if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
				t.Error(err)
			}

			fmt.Printf("%+v\n", user)
			wg.Done()
		}(i)
		time.Sleep(time.Millisecond * 1)
	}
	wg.Wait()
	fmt.Printf(" times %v we hit the database:", s.dbHit)
}
