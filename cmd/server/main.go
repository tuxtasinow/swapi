package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
	"github.com/tuxtasinow/swapi/internal/client/http/swapi"
)

const httpPort = ":3000"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	httpClient := &http.Client{Timeout: time.Second * 10}
	swapiClient := swapi.NewClient(httpClient)

	r.Get("/{planets}", func(w http.ResponseWriter, r *http.Request) {
		planets := chi.URLParam(r, "planets")
		id, err := strconv.ParseInt(planets, 10, 64)
		if err != nil {
			http.Error(w, "invalid planet id", http.StatusBadRequest)
			return
		}

		fmt.Printf("Request planets: %s\n", planets)

		res, err := swapiClient.GetPlanet(id)
		if err != nil {
			log.Println(err)
		}

		raw, err := json.Marshal(res)
		if err != nil {
			log.Println(err)
		}

		_, err = w.Write([]byte(raw))
		if err != nil {
			log.Println(err)
		}
	})

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	jobs, err := initJob(s)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		fmt.Println("starting server on port " + httpPort)
		err = http.ListenAndServe(httpPort, r)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()

		fmt.Printf("starting job: %v\n", jobs[0].ID())
		s.Start()
	}()

	wg.Wait()
}

func initJob(scheduler gocron.Scheduler) ([]gocron.Job, error) {
	j, err := scheduler.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("hello")
			},
		),
	)

	if err != nil {
		return nil, err
	}

	return []gocron.Job{j}, nil
}
