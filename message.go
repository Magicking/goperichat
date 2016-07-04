package main

import (
	"encoding/json"
)

// layer 1 => CommandFrame
type CommandFrame struct {
/*
Kind 1 => Payload = SessionMessage, Body = empty
Kind 2 => Payload = empty, Body = CommandFrame
Kind 4 => Payload = empty, Body = RoomMsg
*/
	Kind int `json:"kind"`
	Payload string `json:"payload"`
	Signature string `json:"signature"`
	Body string `json:"body"`
}

// kind: 4
// layer: 2
type RoomMsg struct {
	Occupancy int `json:"occupancy"`
	Room string `json:"room"`
	TotalParticipants int `json:"total_participants"`
}

type SenderT struct {
	UserId string `json:"user_id"`
	ParticipantIndex int `json:"participant_index"`
}
// kind: 1
// layer: 2
type SessionMessage struct {
	Body string `json:"body"`
	Bt int `json:"bt"` //Body Type
	Oa bool `json:"oa"` //OAuth ???
	Room string `json:"room"`
	Sender SenderT `json:"sender"`
	Ssid string `json:"ssid"`
	St int `json:"st"`//Always 0 ?
	Timestamp int `json:"timestamp"`
}

/*
type 1: Chat msg
Username,DisplayName,Body
type 2:
type 3:
type 4: Geoloc frame
Heading, Lat, Lng, 

*/
type PeriMessage struct {
	Type int `json:"type"`
	Body string `json:"body"`
	DisplayName string `json:"displayName"`
	Heading float64 `json:"heading"`
	Initials string `json:"initials"`
	JuryDuration int `json:"jury_duration"`
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	Ntpforbroadcasterframe int `json:"ntpForBroadcasterFrame"`
	Ntpforliveframe int `json:"ntpForLiveFrame"`
	ParticipantIndex int `json:"participant_index"`
	Profileimageurl string `json:"profileImageURL"`
	Remoteid string `json:"remoteID"`
	ReportType int `json:"report_type"`
	SentenceDuration int `json:"sentence_duration"`
	SentenceType int `json:"sentence_type"`
	Timestamp int `json:"timestamp"`
	Username string `json:"username"`
	Uuid string `json:"uuid"`
	Verdict int `json:"verdict"`
	V int `json:"v"`
}

func FrameFilter(frame []byte) (PeriMessage, error) {
	var cf CommandFrame
	var sm SessionMessage
	var pm PeriMessage

	err := json.Unmarshal(frame, &cf)
	if err != nil {
		return pm, err
	}
	if cf.Kind == 1 {
		err = json.Unmarshal([]byte(cf.Payload), &sm)
		if err != nil {
			return pm, err
		}
		if sm.Bt == 1 {
			err = json.Unmarshal([]byte(sm.Body), &pm)
			return pm, nil
		}
	}
	return pm, nil
}
