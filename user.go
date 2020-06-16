package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

const (
	signingKey = "rtd4409kl1"
)

// User struct for user object
type User struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Email       string `json:"email,omitempty"`
	MobileNo    string `json:"mobile_no,omitempty"`
	UserType    string `json:"user_type,omitempty"`
	UserStatus  string `json:"user_status,omitempty"`
	DateCreated string `json:"date_created,omitempty"`
}

func register(c echo.Context) error {
	u := &User{}

	if err := c.Bind(u); err != nil {
		return err
	}

	var e string
	row := db.QueryRow("SELECT username FROM users WHERE username = ?", u.Username)

	err := row.Scan(&e)
	if err == sql.ErrNoRows {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 5)
		if err != nil {
			return err
		}

		u.Password = string(hash)

		insertUserSQL := `INSERT INTO users(username, password) VALUES(?,?)`

		stmt, err := db.Prepare(insertUserSQL)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(u.Username, u.Password)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"success": "user created"})
	}

	return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"status": "user already exist"})
}

func login(c echo.Context) error {
	u := &User{}

	if err := c.Bind(u); err != nil {
		return err
	}

	var id, passwd string
	row := db.QueryRow("SELECT id, password from users WHERE username = ?", u.Username)
	err := row.Scan(&id, &passwd)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "user not found"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwd), []byte(u.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "usename or password not correct"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        id,
		"username":  u.Username,
		"user_type": "learner",
		"expiresAt": time.Now().Add(15 * time.Minute),
	})

	tokenString, err := token.SignedString([]byte(signingKey))
	return c.JSON(http.StatusAccepted, map[string]interface{}{"token": tokenString})
}

func getProfile(c echo.Context) error {
	tokenString := strings.Split(c.Request().Header.Get("Authorization"), " ")
	token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unxpected signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// user is valid
		// check the token in db and return the requested result
		return c.JSON(http.StatusOK, claims)
	}

	return nil
}
