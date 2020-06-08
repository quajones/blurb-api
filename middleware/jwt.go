package middleware

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	//jwt.StandardClaims
}

func GenerateToken(userClaims Claims) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["acc"] = userClaims
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	t, err := token.SignedString([]byte(os.Getenv("BLURB_JWT")))
	if err != nil {
		return "", err
	}
	return t, nil
}
func JWTCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		account := claims["acc"].(interface{})
		c.Set("user", account)
		c.Set("userId", account.(map[string]interface{})["id"].(string))
		return next(c)
	}
}
