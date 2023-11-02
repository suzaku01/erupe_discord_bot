package main

import (
//add this
	"log"
	"net/http"
	"encoding/json"
	"database/sql"
)

func main() {
	logger.Info("Finished starting Erupe")
	//add this after "logger.Info("Finished starting Erupe")"
	RunMessageBot(channels)

}

//add this
func RunMessageBot(channels []*channelserver.Server) {	
    http.HandleFunc("/send", makeReceiveMessageHandler(channels))
	http.HandleFunc("/isalive", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server is alive!")
	})
	http.HandleFunc("/getplayers", GetPlayersHandler)
	
	log.Println("Server started on: http://localhost:9999")	//change here
	err1 := http.ListenAndServe(":9999", nil)
	if err1 != nil {
		log.Fatal("Error starting the server:", err1)
	}
}

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

func GetPlayersHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("postgres", "host=localhost user=postgres password=password dbname=erupe sslmode=disable")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	var totalPlayers int
	query := "SELECT SUM(current_players) FROM servers"
	err = db.QueryRow(query).Scan(&totalPlayers)
	if err != nil {
		log.Fatal("Failed to execute query:", err)
	}

	// レスポンスとして JSON を返します。
	responseData := map[string]int{"totalPlayers": totalPlayers}
	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to generate JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
