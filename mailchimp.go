package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Email struct {
	Key     string  `json:"key"`
	Message Message `json:"message"`
}
type Message struct {
	Html        string       `json:"html"`
	Text        string       `json:"text"`
	Subject     string       `json:"subject"`
	FromEmail   string       `json:"from_email"`
	FromName    string       `json:"from_name"`
	To          string       `json:"to"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func getConfiguration() Configuration {
	file, _ := os.Open("configuration.json")
	decoder := json.NewDecoder(file)
	configuration := &Configuration{}
	decoder.Decode(&configuration)
	return *configuration
}

func main() {
	configuration := getConfiguration()

	attachments := make([]Attachment, 0)
	attachment := Attachment{
		Type:    "text/plain",
		Name:    "mailchimp.go",
		Content: "blahblah",
	}

  if len(os.Args) != 2 {
      os.Stderr.Write([]byte("Not enough arguments\n"))
      os.Stderr.Write([]byte(fmt.Sprintf("usage: %s email\n", os.Args[0])))
      return
  }

	message := Message{
		Html:        "Hey there, here is my Go code that uses the Mandrill API to send itself as an attachment in this email. I'm sending this as an application for your Software Engineer Intern position. For more of my code you can check out my Github: github.com/dkua. Cheers!",
		Text:        "Hey there, here is my Go code that uses the Mandrill API to send itself as an attachment in this email. I'm sending this as an application for your Software Engineer Intern position. For more of my code you can check out my Github: github.com/dkua. Cheers!",
		Subject:     "Code Sample - David Kua",
		FromEmail:   configuration.Username,
		To:          os.Args[1],
		Attachments: append(attachments, attachment),
	}

	email := Email{
		Key:     configuration.Password,
		Message: message,
	}

	req, err := json.Marshal(email)
	if err != nil {
		fmt.Println("Error:", err)
	}
  fmt.Println(string(req))
}
