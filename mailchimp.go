package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const BaseUrl string = "https://mandrillapp.com/api/1.0/"
const SendEndpoint string = "messages/send.json"

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
	To          []To         `json:"to"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type To struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

func getConfiguration() Configuration {
	file, _ := os.Open("configuration.json")
	decoder := json.NewDecoder(file)
	configuration := &Configuration{}
	decoder.Decode(&configuration)
	return *configuration
}

func createAttachment(mime, name, filepath string) Attachment {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	content := base64.URLEncoding.EncodeToString(file)
	attachment := Attachment{
		Type:    mime,
		Name:    name,
		Content: content,
	}
	return attachment
}

func main() {
	configuration := getConfiguration()

	attachments := make([]Attachment, 1)
	attachment := createAttachment("text/plain", "mailchimp.go", "mailchimp.go")

	if len(os.Args) != 2 {
		os.Stderr.Write([]byte("Not enough arguments\n"))
		os.Stderr.Write([]byte(fmt.Sprintf("usage: %s email\n", os.Args[0])))
		return
	}

	recipients := make([]To, 1)
	recipient := To{
		Email: os.Args[1],
		Name:  "Mailchimp",
		Type:  "to",
	}

	message := Message{
		Html:        "Hey there, here is my Go code that uses the Mandrill API to send itself as an attachment in this email. I'm sending this as an application for your Software Engineer Intern position. The code for this can also be found on my Github at github.com/dkua/mandrill_application. Cheers!",
		Text:        "Hey there, here is my Go code that uses the Mandrill API to send itself as an attachment in this email. I'm sending this as an application for your Software Engineer Intern position. The code for this can also be found on my Github at github.com/dkua/mandrill_application. Cheers!",
		Subject:     "Code Sample - David Kua",
		FromEmail:   configuration.Username,
		To:          append(recipients, recipient),
		Attachments: append(attachments, attachment),
	}

	email := Email{
		Key:     configuration.Password,
		Message: message,
	}

	content, err := json.Marshal(email)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(string(content))
	resp, err := http.Post(BaseUrl+SendEndpoint, "application/json", bytes.NewBuffer(content))
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(contents))
}
