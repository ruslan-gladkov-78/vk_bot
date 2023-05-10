package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func CreateURLValues() url.Values {
	values := url.Values{}
	values.Add("access_token", token)
	values.Add("v", version)
	return values
}

func ProcessUpdate(update LongPollUpdate) ([]byte, error) {
	if len(update.Updates) == 0 {
		return nil, nil
	}
	msgText := update.Updates[0].Object.Message.Text
	switch {
	case msgText == "Начать":
		return ProcessStart(update)
	case msgText == "Расскажи анекдот":
		return ProcessAnecdote(update)
	case msgText == "Покажи котика":
		return ProcessCat(update)
	case msgText == "Который час?":
		return ProcessTime(update)
	case msgText == "Мне нужна помощь":
		return ProcessHelp(update)
	case msgText == "Вернуться":
		return ProcessReturn(update)
	case msgText == "Про Петьку и Василия Ивановича" || msgText == "Про Вовочку" ||
		msgText == "Про армию":
		return SendAnecdote(update)
	case msgText == "Санкт-Петербург" || msgText == "Екатеринбург" ||
		msgText == "Новосибирск" || msgText == "Владивосток":
		return SendTime(update)
	case msgText == "Другой":
		return SendMessage(update.Updates[0].Object.Message.FromID, timeZoneMessage)
	case strings.Contains(msgText, "GMT"):
		return ParseAndSendTime(update)
	case msgText == "Рыженькие" || msgText == "Серенькие" || msgText == "Лысенькие":
		_, err := SendCat(update)
		if err != nil {
			fmt.Println(err.Error())
		}
		return SendKeyboard(update.Updates[0].Object.Message.FromID, "Держите!", mainKeyboard)
	}
	return nil, nil
}

func ProcessStart(update LongPollUpdate) ([]byte, error) {
	userID := update.Updates[0].Object.Message.FromID
	return SendKeyboard(userID, helloMessage, mainKeyboard)
}

func ProcessAnecdote(update LongPollUpdate) ([]byte, error) {
	userID := update.Updates[0].Object.Message.FromID
	return SendKeyboard(userID, anecdoteMessage, anecdoteKeyboard)
}

func ProcessCat(update LongPollUpdate) ([]byte, error) {
	userID := update.Updates[0].Object.Message.FromID
	return SendKeyboard(userID, catMessage, catKeyboard)
}

func ProcessTime(update LongPollUpdate) ([]byte, error) {
	userID := update.Updates[0].Object.Message.FromID
	return SendKeyboard(userID, timeMessage, timeKeyboard)
}

func ProcessHelp(update LongPollUpdate) ([]byte, error) {
	userID := update.Updates[0].Object.Message.FromID
	return SendKeyboard(userID, helpMessage, helpKeyboard)
}

func ProcessReturn(update LongPollUpdate) ([]byte, error) {
	userID := update.Updates[0].Object.Message.FromID
	return SendKeyboard(userID, returnMessage, mainKeyboard)
}

func SendAnecdote(update LongPollUpdate) ([]byte, error) {
	userID := update.Updates[0].Object.Message.FromID
	topic := update.Updates[0].Object.Message.Text
	anecdote := anecdotes[topic][rand.Intn(len(anecdotes[topic]))]
	return SendKeyboard(userID, anecdote, mainKeyboard)
}

func SendTime(update LongPollUpdate) ([]byte, error) {
	userID := update.Updates[0].Object.Message.FromID
	city := update.Updates[0].Object.Message.Text
	diff, ok := timeDifference[city]
	if ok {
		currentTime := time.Now().Local().Add(time.Duration(diff) * time.Hour)
		t := currentTime.Format("15:04:05")
		return SendKeyboard(userID, t, mainKeyboard)
	}
	return nil, nil
}

func ParseAndSendTime(update LongPollUpdate) ([]byte, error) {
	userID := update.Updates[0].Object.Message.FromID
	msgText := update.Updates[0].Object.Message.Text
	diff, err := strconv.Atoi(msgText[3:])
	if err != nil {
		return nil, err
	}
	currentTime := time.Now().Local().Add(time.Duration(diff-3) * time.Hour)
	t := currentTime.Format("15:04:05")
	return SendKeyboard(userID, t, mainKeyboard)
}

func SendCat(update LongPollUpdate) ([]byte, error) {
	userID := update.Updates[0].Object.Message.FromID
	catType := update.Updates[0].Object.Message.Text
	catID := catsImages[catType][rand.Intn(len(catsImages[catType]))]
	return SendImage(userID, catID)
}

func SendMessage(id int, text string) ([]byte, error) {
	method := "messages.send"

	v := CreateURLValues()
	v.Add("user_id", strconv.Itoa(id))
	v.Add("random_id", strconv.Itoa(rand.Intn(MaxInt)))
	v.Add("message", text)

	u := fmt.Sprintf("%s/method/%s?%s", vkURL, method, v.Encode())
	resp, err := http.Get(u)

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, err

}

func SendImage(userID int, imgID string) ([]byte, error) {
	method := "messages.send"

	v := CreateURLValues()
	v.Add("user_id", strconv.Itoa(userID))
	v.Add("random_id", strconv.Itoa(rand.Intn(MaxInt)))
	v.Add("attachment", "photo"+imgID)

	u := fmt.Sprintf("%s/method/%s?%s", vkURL, method, v.Encode())
	resp, err := http.Get(u)

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, err

}
