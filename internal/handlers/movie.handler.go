package handlers

import (
	"c07_practica/pkg/model"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func MovieRouterHandlers(router *mux.Router, base_url string, api_key string, secret string, db *sql.DB) {
	router.Handle("/movies/popular", model.AuthMiddleware(getPopularMovies(base_url, api_key), secret)).Methods("GET")
	router.Handle("/movies/popular/internal", model.AuthMiddleware(getInternalPopularMovies(db, api_key, base_url), secret)).Methods("GET")
	router.Handle("/movies/{id}", model.AuthMiddleware(getMovie(db, base_url, api_key), secret)).Methods("GET")

	router.Handle("/movies/comment/{movie_id}", model.AuthMiddleware(createComment(db, api_key, base_url), secret)).Methods("POST")
	router.Handle("/movies/comment/{comment_id}", model.AuthMiddleware(updateComment(db), secret)).Methods("PUT")
	router.Handle("/movies/comment/{comment_id}", model.AuthMiddleware(deleteComment(db), secret)).Methods("DELETE")

	router.Handle("/movies/match/comments", model.AuthMiddleware(getUserComments(db, api_key, base_url), secret)).Methods("GET")
}

// API

func getMovie(db *sql.DB, base_url, api_key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer model.IncrementRating(db, id)

		movie, err := getMovieById(id, api_key, base_url)
		if err != nil {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		getCommentMovieById, err := model.GetCommentsByMovieId(db, id)
		if err == nil {
			movie.UserComment = &getCommentMovieById.Comment
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movie)

	}
}

func getInternalPopularMovies(db *sql.DB, api_key string, base_url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var movies_response []model.Movie
		movies, err := model.GetInternalPopularMovies(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, movie_id := range movies {
			movie, err := getMovieById(movie_id.Id, api_key, base_url)
			if err != nil {
				http.Error(w, "Movie not found", http.StatusNotFound)
				return
			}
			movies_response = append(movies_response, *movie)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movies_response)
	}
}

func getPopularMovies(base_url, api_key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := base_url + "/movie/popular"
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("accept", "application/json")
		req.Header.Add("Authorization", api_key)

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var movies model.Movies

		err = json.Unmarshal(body, &movies)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode((movies))
	}

}

// DATABASE

func createComment(db *sql.DB, api_key string, base_url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["movie_id"])

		var comment model.Movie_Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		movie, err := getMovieById(id, api_key, base_url)

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		user_id, ok := model.ObtenerIdJWT(r.Context())

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = model.CreateComment(db, movie.Id, int(user_id), comment.Comment)

		if err != nil {
			http.Error(w, "Cannot Create Comment", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Tu comentario fue creado con exito!")
	}
}

func updateComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		comment_id, err := strconv.Atoi(vars["comment_id"])

		var comment model.Movie_Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user_id, ok := model.ObtenerIdJWT(r.Context())

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = model.UpdateComments(db, int(user_id), comment_id, comment.Comment)

		if err != nil {
			http.Error(w, "Cannot Update Comment", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Tu comentario fue actualizado con exito!")
	}
}

func deleteComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		comment_id, err := strconv.Atoi(vars["comment_id"])

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user_id, ok := model.ObtenerIdJWT(r.Context())

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = model.DeleteComment(db, int(user_id), comment_id)

		if err != nil {
			http.Error(w, "Cannot Update Comment", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Tu comentario fue Eliminado con exito!")
	}
}

func getUserComments(db *sql.DB, api_key string, base_url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user_id, ok := model.ObtenerIdJWT(r.Context())

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		movies_comments, err := model.GetComments(db, int(user_id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var movies_with_comments []*model.Movie
		for _, comments := range movies_comments {
			movie, err := getMovieById(comments.MovieId, api_key, base_url)

			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Movie not found", http.StatusNotFound)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			userComment := comments.Comment
			commentId := comments.Id

			movie.UserComment = &userComment
			movie.CommentId = &commentId

			movies_with_comments = append(movies_with_comments, movie)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movies_with_comments)
	}

}

// UTILS

func getMovieById(movie_id int, api_key string, base_url string) (*model.Movie, error) {
	url := base_url + "/movie/" + strconv.Itoa(movie_id)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return &model.Movie{}, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", api_key)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return &model.Movie{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &model.Movie{}, err
	}
	var movie model.Movie

	err = json.Unmarshal(body, &movie)
	if err != nil {
		return &model.Movie{}, err
	}
	return &movie, nil
}
