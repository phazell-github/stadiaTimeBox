package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	//fire up chrome to stadia
	// start a countdown timer
	ticker := time.NewTicker(1000 * time.Millisecond)
	done := make(chan bool)
	timeLeft := 123

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println(timeLeft, t)
				timeLeft--
			}
		}
	}()

	time.Sleep(123000 * time.Millisecond)
	ticker.Stop()
	done <- true

	// fire up a server
	// - serve page
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("timerPage.html"))))
	serve()

	// - API to deliver remaining time
	http.HandleFunc("/api", getTime(timeLeft))

}

func getTime(time int) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(strconv.Itoa(time)))
	}

}

func serve() {
	err := http.ListenAndServe(":6357", nil)
	if err != nil {
		log.Fatal("ListenAndServe failed ", err)
	}
}
