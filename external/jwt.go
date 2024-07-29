package external

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(body interface{}) (string, error) {

	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_MINS_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["iss"] = "https://test.com"
	claims["iat"] = time.Now().Unix()
	claims["nbf"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	claims["body"] = body

	fmt.Println(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}
