package repository

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/matheusgcoppi/barber-finance-api/database"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	"github.com/matheusgcoppi/barber-finance-api/mail"
	"github.com/matheusgcoppi/barber-finance-api/utils"
	"golang.org/x/crypto/bcrypt"
	_ "gorm.io/gorm"
	"html/template"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type DbRepository struct {
	Store *database.CustomDB
}

func (s *DbRepository) CreateUser(user *model.User) (*model.User, *model.Account, error) {
	newUser := &model.User{
		Active:   user.Active,
		Type:     user.Type,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	result := s.Store.Db.Create(&newUser)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "users_email_key") {
			return nil, nil, fmt.Errorf("email address '%s' is already in use", user.Email)
		} else {
			return nil, nil, result.Error
		}
	}
	if newUser.Type == model.Admin {
		newAccount := &model.Account{
			UserId:  newUser.ID,
			Balance: 0,
			User:    newUser,
		}

		account := s.Store.Db.Create(&newAccount)
		if account.Error != nil {
			return nil, nil, account.Error
		}
		return newUser, newAccount, nil
	}

	return newUser, nil, nil
}

func (s *DbRepository) LoginUser(email string, password string) (*model.User, error) {
	var user model.User
	result := s.Store.Db.First(&user, "email = ?", email)
	if gorm.IsRecordNotFoundError(result.Error) {
		return nil, fmt.Errorf("user not found")
	} else if result.Error != nil {
		return nil, result.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

func (s *DbRepository) GetUser() (error, []*model.User) {
	var users []*model.User
	result := s.Store.Db.Order("id").Find(&users)
	if result.Error != nil {
		return result.Error, nil
	}
	return nil, users
}

func (s *DbRepository) GetUserByID(id string) (*model.User, error) {
	var user *model.User
	result := s.Store.Db.First(&user, id)
	if id == "" {
		return nil, fmt.Errorf("User with id = " + id + " not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (s *DbRepository) DeleteUser(id string) error {
	result := s.Store.Db.Delete(&model.User{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("User with id = " + id + " not found")
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *DbRepository) UpdateUser(user *model.UserDTO, id string) (*model.User, error) {
	userById, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if (user.Active == true) || (user.Active == false) {
		userById.Active = user.Active
	}

	if user.Type != 0 {
		userById.Type = user.Type
	}

	if user.Username != "" {
		userById.Username = user.Username
	}

	if user.Email != "" {
		userById.Email = user.Email
	}

	if user.Password != "" {
		userById.Password = user.Password
	}

	s.Store.Db.Save(userById)
	return userById, nil
}
func (s *DbRepository) ForgotPassword(email string) error {
	err := godotenv.Load("./.env")
	if err != nil {
		return fmt.Errorf("error loading .env file")
	}
	var userId uint
	s.Store.Db.Raw("SELECT id FROM users WHERE email = ?", email).Scan(&userId)
	if userId == 0 {
		return fmt.Errorf("e-mail was not found")
	}

	fmt.Println(email)

	emailSenderName := os.Getenv("EMAIL_SENDER_NAME")
	emailSenderEmail := os.Getenv("EMAIL_SENDER_ADDRESS")
	emailSenderPassword := os.Getenv("EMAIL_SENDER_PASSWORD")

	if emailSenderName == "" || emailSenderEmail == "" || emailSenderPassword == "" {
		return fmt.Errorf("could not get email sender information from .env file")
	}

	sender := mail.NewGmailSender(emailSenderName, emailSenderEmail, emailSenderPassword)

	subject := "Reset Password Token"
	t, _ := template.ParseFiles("./mail/forgot-password-template.html")
	var body bytes.Buffer
	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("Subject: Forgot Password\n%s\n\n", headers)))

	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	var token string
	for i := 0; i < 7; i++ {
		randomInt := rng.Intn(10)
		token += string(table[randomInt])
	}

	var name string
	s.Store.Db.Raw("SELECT username FROM users WHERE id = ?", userId).Scan(&name)

	err = t.Execute(&body, struct {
		Name  string
		Token string
	}{
		Name:  name,
		Token: token,
	})
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	tokenEncrypted, err := utils.GetAESEncrypted(token, os.Getenv("KEY"), os.Getenv("IV"))
	if err != nil {
		panic(err)
	}

	var user model.User
	err = s.Store.Db.Raw("SELECT * FROM users WHERE id = ?", userId).Scan(&user).Error
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	var row model.UserToken
	err = s.Store.Db.Raw("SELECT * FROM users_token WHERE user_id = ?", userId).Scan(&row).Error
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	if row.UserID != 0 {
		row.Token = tokenEncrypted
		row.CreatedAt = time.Now()
		row.UpdatedAt = time.Now()
		s.Store.Db.Save(row)
	} else {
		newUserToken := &model.UserToken{
			UserID: userId,
			Token:  tokenEncrypted,
			User:   user,
		}

		resultToken := s.Store.Db.Create(&newUserToken)
		if resultToken.Error != nil {
			return fmt.Errorf("error creating UserToken %s", resultToken.Error)
		}
	}

	err = sender.SendEmail(subject, string(body.Bytes()), []string{email}, nil, nil, nil)

	return nil
}

func (s *DbRepository) ChangePassword(email, token string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty in cookie")
	} else if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	err := godotenv.Load("./.env")
	if err != nil {
		return err
	}

	var userInfo struct {
		token     string    `db:"token"`
		createdAt time.Time `db:"created_at"`
	}
	s.Store.Db.Raw("SELECT token, created_at "+
		"FROM users_token "+
		"WHERE user_id = (SELECT id FROM users WHERE email = ?)", email).Scan(&userInfo)

	tokenDecrypted, err := utils.GetAESDecrypted(userInfo.token, os.Getenv("KEY"), os.Getenv("IV"))
	if err != nil {
		return err
	}

	if strings.Compare(string(tokenDecrypted), token) != 0 {
		return fmt.Errorf("token is not valid")
	}

	currentTime := time.Now()
	duration := currentTime.Sub(userInfo.createdAt)

	if duration.Minutes() > 20 {
		fmt.Println("More than 20 minutes have passed since the token creation.")
	}

	return nil
}

func (s *DbRepository) TestResetPassword(token, password string) error {
	var id uint

	query := `SELECT id FROM users WHERE id = 
          (SELECT user_id FROM users_token 
           WHERE users_token.token = ?
           AND users_token.updated_at > now()::timestamp - INTERVAL '20 min')`
	err := s.Store.Db.Raw(query, token).Scan(&id).Error

	if err != nil {
		return err
	}

	if id == 0 {
		return fmt.Errorf("token to change the password was expired")
	}

	user, erre := s.GetUserByID(strconv.Itoa(int(id)))
	if err != nil {
		return fmt.Errorf(erre.Error())
	}

	if password != "" {
		user.Password = password
	}

	return nil
}

func NewUser(active bool, userType int, username, email, password string) *model.User {
	return &model.User{
		CustomModel: model.CustomModel{},
		Active:      active,
		Type:        model.Type(userType),
		Username:    username,
		Email:       email,
		Password:    password,
	}
}

func (s *DbRepository) CreateNewToken(userId uint, token string) *model.UserToken {
	return &model.UserToken{
		CustomModel: model.CustomModel{},
		UserID:      userId,
		Token:       token,
	}
}
