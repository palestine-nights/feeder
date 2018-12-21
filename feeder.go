package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

import (
	"github.com/palestine-nights/backend/src/db"
)

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (user *User) MustGetToken() string {
	userBody, _ := json.Marshal(user)
	body := bytes.NewBuffer(userBody)

	response, err := http.Post(
		"https://auth.palestinenights.com/auth",
		"application/json",
		body,
	)

	handleError(err)

	var data interface{}
	if response.StatusCode == http.StatusOK {
		bodyAsByteArray, err := ioutil.ReadAll(response.Body)

		handleError(err)

		err = json.Unmarshal(bodyAsByteArray, &data)

		handleError(err)

		return data.(map[string]interface{})["token"].(string)
	}

	panic(fmt.Sprintf("Bad status %v", response.StatusCode))
}

func FeedMenu(token string) {
	ApiUrl := "https://api.palestinenights.com/menu"
	files, err := filepath.Glob("./data/menu/*.json")

	handleError(err)

	for _, file := range files {
		jsonFile, err := os.Open(file)

		handleError(err)

		byteValue, _ := ioutil.ReadAll(jsonFile)

		menuItems := make([]db.MenuItem, 0)
		json.Unmarshal(byteValue, &menuItems)

		for _, item := range menuItems {
			j, _ := json.Marshal(item)
			body := bytes.NewBuffer(j)

			req, _ := http.NewRequest(http.MethodPost, ApiUrl, body)
			req.Header.Set("Authorization", "Bearer " + token)
			response, err := http.DefaultClient.Do(req)

			handleError(err)

			if response.StatusCode == http.StatusCreated {
				_, err = ioutil.ReadAll(response.Body)

				handleError(err)

				fmt.Println(response.StatusCode)
			} else {
				fmt.Println(fmt.Sprintf("Bad status %v", response.StatusCode))
			}
		}
	}
}

func main() {
	user := User{}

	flag.StringVar(&user.UserName, "username", "", "Administrator Username")
	flag.StringVar(&user.Password, "password", "", "Administrator Password")

	flag.Parse()

	if user.UserName == "" {
		panic("Username is not specified")
	}

	if user.Password == "" {
		panic("Password is not specified")
	}

	token := user.MustGetToken()

	FeedMenu(token)
}
