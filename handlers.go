package channelserver

import (
	"net/http"
)

func logoutPlayer(s *Session) {
	s.server.Lock()
	if _, exists := s.server.sessions[s.rawConn]; exists {
		delete(s.server.sessions, s.rawConn)
	}
	s.rawConn.Close()
	delete(s.server.objectIDs, s)
	s.server.Unlock()

	for _, stage := range s.server.stages {
		// Tell sessions registered to disconnecting players quest to unregister
		if stage.host != nil && stage.host.charID == s.charID {
			for _, sess := range s.server.sessions {
				for rSlot := range stage.reservedClientSlots {
					if sess.charID == rSlot && sess.stage != nil && sess.stage.id[3:5] != "Qs" {
						sess.QueueSendMHF(&mhfpacket.MsgSysStageDestruct{})
					}
				}
			}
		}
		for session := range stage.clients {
			if session.charID == s.charID {
				delete(stage.clients, session)
			}
		}
	}

	_, err := s.server.db.Exec("UPDATE sign_sessions SET server_id=NULL, char_id=NULL WHERE token=$1", s.token)
	if err != nil {
		panic(err)
	}

	_, err = s.server.db.Exec("UPDATE servers SET current_players=$1 WHERE server_id=$2", len(s.server.sessions), s.server.ID)
	if err != nil {
		panic(err)
	}

	var timePlayed int
	var sessionTime int
	_ = s.server.db.QueryRow("SELECT time_played FROM characters WHERE id = $1", s.charID).Scan(&timePlayed)
	sessionTime = int(TimeAdjusted().Unix()) - int(s.sessionStart)
	timePlayed += sessionTime

	var rpGained int
	if mhfcourse.CourseExists(30, s.courses) {
		rpGained = timePlayed / 900
		timePlayed = timePlayed % 900
		s.server.db.Exec("UPDATE characters SET cafe_time=cafe_time+$1 WHERE id=$2", sessionTime, s.charID)
	} else {
		rpGained = timePlayed / 1800
		timePlayed = timePlayed % 1800
	}

	s.server.db.Exec("UPDATE characters SET time_played = $1 WHERE id = $2", timePlayed, s.charID)

	treasureHuntUnregister(s)

	if s.stage == nil {
		return
	}

	s.server.BroadcastMHF(&mhfpacket.MsgSysDeleteUser{
		CharID: s.charID,
	}, s)

	s.server.Lock()
	for _, stage := range s.server.stages {
		if _, exists := stage.reservedClientSlots[s.charID]; exists {
			delete(stage.reservedClientSlots, s.charID)
		}
	}
	s.server.Unlock()

	removeSessionFromSemaphore(s)
	removeSessionFromStage(s)

	saveData, err := GetCharacterSaveData(s, s.charID)
	if err != nil || saveData == nil {
		s.logger.Error("Failed to get savedata")
		return
	}
	saveData.RP += uint16(rpGained)
	if saveData.RP >= s.server.erupeConfig.GameplayOptions.MaximumRP {
		saveData.RP = s.server.erupeConfig.GameplayOptions.MaximumRP
	}
	saveData.Save(s)
	
	s.logger.Info(fmt.Sprintf("[%s] Disconnected", s.Name))
	url := "http://localhost:9998/delete_user"

	// POSTリクエストの実行
	resp, err := http.Post(url, "application/json", strings.NewReader("{}"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}