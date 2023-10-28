package main

import (
//add this
	"log"
	"net/http"
	"encoding/json"
)

func main() {
	logger.Info("Finished starting Erupe")
	//add this after "logger.Info("Finished starting Erupe")"
	RunMessageBot(channels)

}

//add this
func RunMessageBot(channels []*channelserver.Server) {	
    http.HandleFunc("/send", makeReceiveMessageHandler(channels))
	
	log.Println("Server started on: http://localhost:9999")	//change here
	err1 := http.ListenAndServe(":9999", nil)
	if err1 != nil {
		log.Fatal("Error starting the server:", err1)
	}
}

//add this
func makeReceiveMessageHandler(channels []*channelserver.Server) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	
	// Decode request body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Failed to decode message", http.StatusBadRequest)
		return
	}

	message, ok := data["message"]
	if !ok {
		http.Error(w, "Message not found", http.StatusBadRequest)
		return
	}
	
	fmt.Println("Received Message:", message)

			for _, c := range channels {
				c.BroadcastChatMessage(message)
			}
}
}
