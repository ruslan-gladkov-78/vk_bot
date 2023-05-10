package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SessionData struct {
	Response struct {
		Key    string `json:"key"`
		Server string `json:"server"`
		TS     string `json:"ts"`
	} `json:"response"`
}

type LongPollUpdate struct {
	TS      string `json:"ts"`
	Updates []struct {
		GroupID int    `json:"group_id"`
		Type    string `json:"type"`
		EventID string `json:"event_id"`
		V       string `json:"v"`
		Object  struct {
			Message struct {
				Date                  int    `json:"date"`
				FromID                int    `json:"from_id"`
				ID                    int    `json:"id"`
				ConversationMessageId int    `json:"conversation_message_id"`
				Text                  string `json:"text"`
			} `json:"message"`
		} `json:"object"`
	} `json:"updates"`
}

func GetLongPollSessionData(groupID string) SessionData {
	method := "groups.getLongPollServer"

	v := CreateURLValues()
	v.Add("group_id", groupID)

	u := fmt.Sprintf("%s/method/%s?%s", vkURL, method, v.Encode())
	resp, _ := http.Get(u)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	sd := SessionData{}
	_ = json.Unmarshal(body, &sd)
	return sd
}

func AccessLongPoll(sd SessionData) {
	ts := sd.Response.TS
	v := url.Values{}
	v.Add("act", "a_check")
	v.Add("key", sd.Response.Key)
	v.Add("ts", ts)
	v.Add("wait", "25")
	for {
		v.Set("ts", ts)
		u := fmt.Sprintf("%s?%s", sd.Response.Server, v.Encode())

		resp, _ := http.Get(u)

		body, _ := io.ReadAll(resp.Body)
		var lpUpdate LongPollUpdate
		_ = json.Unmarshal(body, &lpUpdate)
		_, err := ProcessUpdate(lpUpdate)
		if err != nil {
			fmt.Println(err.Error())
		}
		ts = lpUpdate.TS
		_ = resp.Body.Close()
	}
}
