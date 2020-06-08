package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"golang.org/x/net/context"
)

// APIKeys : format for json file for api credentials
type APIKeys struct {
	Domain string `json:"domain"`
	Key    string `json:"api_key"`
}

func getAPIKeys(path string) (APIKeys, error) {
	var api APIKeys

	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return api, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &api)

	return api, nil
}

// SendCode : sending email with MailGun API
func SendCode(to string, code int) error {

	apiKeys, err := getAPIKeys("mailgun.json")

	mg := mailgun.NewMailgun(apiKeys.Domain, apiKeys.Key)
	fmt.Println("to:", to)
	fmt.Println("passcode:", code)

	codeString := strconv.Itoa(code)

	sender := "mail@kahvipatel.com"
	subject := "Your passcode"
	body := "Hi, enter " + codeString
	recipient := to

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)
	message.SetTemplate("password-template")
	message.AddTemplateVariable("passcode", codeString)
	message.AddTemplateVariable("link", "https://kahvipatel.com/book-list")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		return err
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return nil
}
