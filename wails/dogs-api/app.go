package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"log"
	"sort"
)


type RandomImageResponse struct {
	Message string
	Status string
}

type AllBreeds struct {
	Message map[string]map[string][]string
	Status string
}

type ImagesByBreed struct {
	Message []string
	Status string
}

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


func (a *App) GetRandomImageUrl() string {
	resp, err := http.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	
	randomImageData := RandomImageResponse{}
	json.Unmarshal(body, &randomImageData)

	return randomImageData.Message
}

func (a *App) GetBreedList() []string {
	var breeds []string 

	resp, err := http.Get("https://dog.ceo/api/breeds/list/all")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	allBreedsData := AllBreeds{}
	json.Unmarshal(body, &allBreedsData)
	for k := range allBreedsData.Message {
		breeds = append(breeds, k) 
	}

	sort.Strings(breeds)

	return breeds
}

func(a *App) GetImageUrlsByBreed(breed string) []string {
	url := fmt.Sprintf("%s%s%s%s", "https://dog.ceo/api/", "breed/", breed, "/images")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	ImagesByBreedData := ImagesByBreed{}
	json.Unmarshal(body, &ImagesByBreedData)

	return ImagesByBreedData.Message
}