package auth

import (
	"saiStorageMongo/src/github.com/satori/go.uuid"
	"fmt"
	"crypto"
	"time"
	"encoding/base64"
)

const DEFUALT_TOKEN_EXPIRATION_DATE = 60 * 60 * 24 * 365

var (
	TokenCollection = "tokens"
	Tokens          = make(map[string]Token)
)

type Token struct {
	ID              string                `json:"_id" bson:"-"`
	Token           string                `json:"token" bson:"token"`
	EXPIRATION_TIME string                `json:"expiration_time" bson:"expiration_time"`
	UserID          string                `json:"user_id" bson:"user_id"`
	Permissions     map[string]Permission `json:"permissions" bson:"permissions"`
}

func CreateNewToken() *Token {
	uuidToken, _ := uuid.NewV4()

	token := &Token{
		Token:           base64.StdEncoding.EncodeToString([]byte(crypto.MD5.New().Sum([]byte(fmt.Sprint(uuidToken))))),
		EXPIRATION_TIME: time.Now().Add(time.Second * DEFUALT_TOKEN_EXPIRATION_DATE).Format("2006-01-02 15:04:05"),
	}

	return token
}

func CreateUserToken(user *User) *Token {
	token := CreateNewToken()
	token.UserID = user.ID
	return token
}

func (token *Token) VerifyToken() bool {
	expirationTime, err := time.Parse("2006-01-02 15:04:05", token.EXPIRATION_TIME)
	if err != nil {
		return false
	}

	return time.Since(expirationTime).Seconds() >= 0
}

func SyncTokens(tokens []Token) {
	for _, token := range tokens {
		Tokens[token.Token] = token
	}
}

func RemoveToken(token string) {
	delete(Tokens, token)
}

func AddToken(token *Token) {
	Tokens[token.Token] = *token
}
