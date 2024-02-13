package mail

import (
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(".env file could not be loaded.")
	}
	sender := NewGmailSender(os.Getenv("EMAIL_SENDER_NAME"), os.Getenv("EMAIL_SENDER_ADDRESS"), os.Getenv("EMAIL_SENDER_PASSWORD"))

	subject := "A test E-mail"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="https://google.com">google</a></p>
	`
	to := []string{"matheuscoppi22@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
