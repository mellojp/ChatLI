package api

import (
	"chatli/data"
	"fmt"
	"time"

	"bytes"
	"encoding/json"
	"net/http"
)

const url = "https://disposable-chat.onrender.com"

func CreateSession(name string) (*data.Session, error) {
	msg := map[string]string{
		"username": name,
	}
	body, _ := json.Marshal(msg)
	resp, err := http.Post(url+"/sessions", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := data.Session{
		Id:           "",
		Username:     name,
		JoinedRooms:  []string{},
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return &res, nil
}

func CreateRoom(s data.Session) (*data.Room, error) {
	reqUrl := url + "/rooms"
	req, _ := http.NewRequest("POST", reqUrl, nil)
	req.Header.Add("Authorization", "Bearer "+s.Id)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := data.Room{
		Id:           "",
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
		ActiveUsers:  []string{},
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return &res, nil
}

func JoinRoom(s data.Session, roomId string) error {
	reqUrl := url + "/rooms/" + roomId + "/join"
	req, _ := http.NewRequest("POST", reqUrl, nil)
	req.Header.Add("Authorization", "Bearer "+s.Id)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("erro ao entrar na sala: status %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	return nil
}

func LoadChatMessages(s data.Session, roomId string, limit int) ([]data.Message, error) {
	reqUrl := fmt.Sprintf("%s/rooms/%s/messages?limit=%d", url, roomId, limit)
	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.Header.Add("Authorization", "Bearer "+s.Id)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("erro na requisição: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	type response struct {
		Messages []data.Message `json:"messages"`
	}
	var res response
	json.NewDecoder(resp.Body).Decode(&res)
	return res.Messages, nil
}
