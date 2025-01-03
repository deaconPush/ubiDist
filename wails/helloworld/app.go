package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
//	"log"
)


// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GenerateWord() (string, error) {
	// Make the GET request
	resp, err := http.Get("https://random-word-api.herokuapp.com/word")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err	
	}

	var words []string
	err = json.Unmarshal(body, &words)

	if err != nil {
		return "", err
	}	
	// log.Println("word: ", words[0])
	return words[0], nil
}