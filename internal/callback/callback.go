package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Callback struct {
	EventType string `json:"event-type"`
	MsgId     string `json:"msg-id"`
	Link      string `json:"link"`
	ClickTime int64  `json:"click-time"`
	UserAgent string `json:"user-agent"`
	ClientIp  string `json:"client-ip"`
}

func Clicked(url string, data *Callback) error {

	json_data, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println("Callback failed: " + err.Error())
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		r, _ := io.ReadAll(resp.Body)
		fmt.Println("Callback failed: " + string(r))
	}

	return nil
}
