package main

import (
	"net/http"
	"encoding/json"
)

// We put only data we need in the struct
type PeriscopeMeta struct {
	ChatToken string `json:"chat_token"`
	Broadcastinfo PeriscopeBroadcast `json:"broadcast"`
}
type PeriscopeBroadcast struct {
	Id string `json:"id"`
}

type ChanPerms struct {
	PB uint `json:"pb"`
	CM uint `json:"cm"`
}

type PeriscopeMetaChat struct {
	Subscriber string `json:"subscriber"`
	Publisher string `json:"publisher"`
	AuthToken string `json:"auth_token"`
	SignerKey string `json:"signer_key"`
	Channel string `json:"channel"`
	ShouldVerifySignature bool `json:"should_verify_signature"`
	AccessToken string `json:"access_token"`
	Endpoint string `json:"endpoint"`
	RoomId string `json:"room_id"`
	ParticipantIndex uint `json:"participant_index"`
	ReadOnly bool `json:"read_only"`
	ShouldLog bool `json:"should_log"`
	ChanPerm ChanPerms `json:"chan_perms"`
}

var urlVideoPublic = "https://api.periscope.tv/api/v2/accessVideoPublic?broadcast_id="
var urlChatPublic = "https://api.periscope.tv/api/v2/accessChatPublic?chat_token="

func toJsonOrPanic(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func getJson(out interface{}, baseUrl, id string) error {
	url := baseUrl + id
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

    err = json.NewDecoder(res.Body).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

func GetPeriscopeMeta(id string) (*PeriscopeMeta, error) {
	var pm PeriscopeMeta

	err := getJson(&pm, urlVideoPublic, id)
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

func GetPeriscopeMetaChat(token string) (*PeriscopeMetaChat, error) {
	var cm PeriscopeMetaChat

	err := getJson(&cm, urlChatPublic, token)
	if err != nil {
		return nil, err
	}
	return &cm, nil
}
