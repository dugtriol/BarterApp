package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/dugtriol/BarterApp/internal/pkg/db"
	"github.com/dugtriol/BarterApp/internal/pkg/repository"
	"github.com/dugtriol/BarterApp/internal/pkg/repository/postgresql"
	"github.com/gorilla/mux"
)

const (
	port          = ":9000"
	queryParamKey = "key"
)

type server1 struct {
	repo *postgresql.UserRepo
}

type addArticleRequest struct {
	Name   string `json:"name"`
	Rating int64  `json:"rating"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.NewDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	articleRepo := postgresql.NewArticles(database)
	implementation := server1{repo: articleRepo}
	router := createRouter(implementation)
	http.Handle("/", router)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func createRouter(implementation server1) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc(
		"/article", func(w http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodPost:
				implementation.Create(w, req)
			case http.MethodPut:
				implementation.Update(w, req)
			default:
				fmt.Println("error")
			}

		},
	)

	router.HandleFunc(
		fmt.Sprintf("/article/{%s:[0-9]+}", queryParamKey), func(w http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodGet:
				implementation.Get(w, req)
			case http.MethodDelete:
				implementation.Delete(w, req)
			default:
				fmt.Println("error")
			}
		},
	)
	return router
}

func (s *server1) Create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var unm addArticleRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	articleRepo := &repository.User{
		Name:   unm.Name,
		Rating: unm.Rating,
	}
	id, err := s.repo.Add(req.Context(), articleRepo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	articleRepo.ID = id
	articleJson, _ := json.Marshal(articleRepo)
	w.Write(articleJson)
}

func (s *server1) Update(_ http.ResponseWriter, req *http.Request) {
	println("update")

}

func (s *server1) Delete(w http.ResponseWriter, req *http.Request) {
	fmt.Println("delete")

}

func (s *server1) Get(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	article, err := s.repo.GetByID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNoFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	articleJson, _ := json.Marshal(article)
	w.Write(articleJson)
}
