package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHandleGetUser(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(http.HandlerFunc(s.handlerGetUser))
	noreq := 1000
	wg := &sync.WaitGroup{}

	for i := 0; i < noreq; i++ {
		wg.Add(1)
		go func(i int) {
			id := i%100 + 1
			url := fmt.Sprintf("%s?id=%d", ts.URL, id)
			resp, err := http.Get(url)
			if err != nil {
				t.Error(err)
			}

			user := &User{}
			if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
				log.Println(user)
				t.Error(err)
			}
			fmt.Printf("USER: %+v\n", user)
			wg.Done()
		}(i)
		time.Sleep(1 * time.Millisecond)
	}
	wg.Wait()
	fmt.Println("times hit database: ", s.dbhits)
}
