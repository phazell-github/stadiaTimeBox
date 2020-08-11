package main

import (
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	//fire up chrome to stadia
	cmd := exec.Command("google-chrome", "https://stadia.google.com/home")
	cmd.Run()

	// fire up a server
	http.Handle("/", http.FileServer(http.Dir(".")))
	timeLeft := 5400
	http.HandleFunc("/api", getTime(&timeLeft))

	// start a countdown timer
	ticker := time.NewTicker(1000 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				timeLeft--
			}
		}
	}()

	serve()
	time.Sleep(5400000 * time.Millisecond)
	ticker.Stop()
	done <- true

	// shutdown stadia

}

func getTime(time *int) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(strconv.Itoa(*time)))
	}

}

func serve() {
	err := http.ListenAndServe(":6357", nil)
	if err != nil {
		log.Fatal("ListenAndServe failed ", err)
	}
}
