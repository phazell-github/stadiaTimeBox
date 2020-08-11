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
	go cmd.Run()

	timeLeft := 5400

	// fire up a server
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/api", getTime(&timeLeft))
	go serve()

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

	time.Sleep(5400 * time.Second)
	ticker.Stop()
	done <- true

	// shutdown stadia
	funsOver := exec.Command("killall", "chrome")
	funsOver.Run()
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
