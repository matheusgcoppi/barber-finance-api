package mail

import (
	"bytes"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/require"
	"html/template"
	"os"
	"testing"
)

func TestSendEmailTemplate(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(".env file could not be loaded.")
	}
	sender := NewGmailSender(os.Getenv("EMAIL_SENDER_NAME"), os.Getenv("EMAIL_SENDER_ADDRESS"), os.Getenv("EMAIL_SENDER_PASSWORD"))

	subject := "Reset Password Token"
	q, _ := template.ParseFiles("forgot-password-template.html")
	var body bytes.Buffer
	err = q.Execute(&body, struct {
		Name  string
		Token string
	}{
		Name:  "Matheus",
		Token: "123456",
	})

	to := []string{"matheusscoppi22@gmail.com"}

	err = sender.SendEmail(subject, string(body.Bytes()), to, nil, nil, nil)
	require.NoError(t, err)
}
