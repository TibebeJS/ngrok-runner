package notifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Auth struct {
	Endpoint  string
	Email     string
	Password  string
	FieldName string
}

type WebHook struct {
	Endpoint string
	Auth     Auth
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (b *WebHook) Login() (string, error) {

	url := b.Auth.Endpoint

	var data map[string]string = map[string]string{
		"email":    b.Auth.Email,
		"password": b.Auth.Password,
	}

	reqBody, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	// 	var jsonStr = []byte(fmt.Sprintf(
	// 		`{
	// 	"email": "%s"
	// 	"password": "%s"
	// }`, b.Auth.Email, b.Auth.Password))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"content-type": []string{"application/json"},
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return "", errors.New(string(body))
	}

	response := AuthResponse{}

	_ = json.Unmarshal([]byte(body), &response)

	return response.Token, nil
}

func (b *WebHook) Notify(message string) error {
	url := b.Endpoint

	token, err := b.Login()

	if err != nil {
		return nil
	}

	// var jsonStr = []byte(message)

	var data map[string]string = map[string]string{
		b.Auth.FieldName: message,
	}

	reqBody, err := json.Marshal(data)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return errors.New(string(body))
	}

	return nil
}
