package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

type ActionStruct struct {
	Type  string `json:"type"`
	Link  string `json:"link,omitempty"`
	Label string `json:"label"`
}

type Button struct {
	Action ActionStruct `json:"action"`
	Color  string       `json:"color,omitempty"`
}

type Keyboard struct {
	Inline  bool       `json:"inline"`
	OneTime bool       `json:"one_time"`
	Buttons [][]Button `json:"buttons"`
}

var (
	mainKeyboard = Keyboard{
		Inline:  false,
		OneTime: false,
		Buttons: [][]Button{
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Расскажи анекдот",
				}, Color: "primary"},
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Покажи котика",
				}, Color: "primary"},
			},
			{

				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Который час?",
				}, Color: "primary"},
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Мне нужна помощь",
				}, Color: "primary"},
			},
		},
	}
	anecdoteKeyboard = Keyboard{
		Inline:  false,
		OneTime: true,
		Buttons: [][]Button{
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Про Петьку и Василия Ивановича",
				}, Color: "primary"},
			},
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Про Вовочку",
				}, Color: "primary"},
			},
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Про армию",
				}, Color: "primary"},
			},
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Вернуться",
				}, Color: "primary"},
			},
		},
	}
	catKeyboard = Keyboard{
		Inline:  false,
		OneTime: true,
		Buttons: [][]Button{
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Рыженькие",
				}, Color: "primary"},
			},
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Серенькие",
				}, Color: "primary"},
			},
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Лысенькие",
				}, Color: "primary"},
			},
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Вернуться",
				}, Color: "primary"},
			},
		},
	}
	timeKeyboard = Keyboard{
		Inline:  false,
		OneTime: true,
		Buttons: [][]Button{
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Санкт-Петербург",
				}, Color: "primary"},
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Екатеринбург",
				}, Color: "primary"},
			},
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Новосибирск",
				}, Color: "primary"},
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Владивосток",
				}, Color: "primary"},
			},
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Другой",
				}, Color: "primary"},
			},
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Вернуться",
				}, Color: "primary"},
			},
		},
	}
	helpKeyboard = Keyboard{
		Inline:  false,
		OneTime: false,
		Buttons: [][]Button{
			{
				Button{Action: ActionStruct{
					Type:  "open_link",
					Link:  "https://nasiliu.net/",
					Label: "Абьюз в семье",
				}},
			},
			{
				Button{Action: ActionStruct{
					Type:  "open_link",
					Link:  "https://psysovet.ru/",
					Label: "Депрессивное состояние",
				}},
			},
			{
				Button{Action: ActionStruct{
					Type:  "open_link",
					Link:  "https://ya-roditel.ru",
					Label: "Проблемы с детьми",
				}},
			},
			{
				Button{Action: ActionStruct{
					Type:  "text",
					Link:  "",
					Label: "Вернуться",
				}, Color: "primary"},
			},
		},
	}
)

func SendKeyboard(id int, text string, keyboard Keyboard) ([]byte, error) {
	method := "messages.send"
	keyboardBytes, _ := json.Marshal(keyboard)

	v := CreateURLValues()
	v.Add("message", text)
	v.Add("keyboard", string(keyboardBytes))
	v.Add("user_id", strconv.Itoa(id))
	v.Add("random_id", strconv.Itoa(rand.Intn(MaxInt)))

	u := fmt.Sprintf("%s/method/%s?%s", vkURL, method, v.Encode())
	resp, err := http.Get(u)

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, err
}
