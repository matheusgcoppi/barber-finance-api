package mail

import (
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/matheusgcoppi/barber-finance-api/utils"
	"github.com/stretchr/testify/require"
	"html/template"
	"math/rand"
	"os"
	"testing"
	"time"
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
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	var token string
	for i := 0; i < 6; i++ {
		randomInt := rng.Intn(10)
		token += string(table[randomInt])
	}

	//key := os.Getenv("KEY")
	//iv := os.Getenv("IV")

	fmt.Println(token)
	encryptedToken, err := utils.GetAESEncrypted(token, os.Getenv("KEY"), os.Getenv("IV"))
	if err != nil {
		log.Fatal("error encrypting token")
	}
	err = q.Execute(&body, struct {
		Name  string
		Token string
	}{
		Name:  "Matheus",
		Token: encryptedToken,
	})

	to := []string{"matheusscoppi22@gmail.com"}

	err = sender.SendEmail(subject, string(body.Bytes()), to, nil, nil, nil)
	require.NoError(t, err)
}
