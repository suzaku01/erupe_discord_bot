package channelserver

import (
	//add this	
	"net/http"
	"strings"
)

// Start starts the session packet send and recv loop(s).
func (s *Session) Start() {
	go func() {
		s.logger.Debug("New connection", zap.String("RemoteAddr", s.rawConn.RemoteAddr().String()))
		// Unlike the sign and entrance server,
		// the client DOES NOT initalize the channel connection with 8 NULL bytes.
		go s.sendLoop()
		s.recvLoop()
	}()
	
	url := "http://localhost:9998/add_user"

	// POSTリクエストの実行
	resp, err := http.Post(url, "application/json", strings.NewReader("{}"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}
