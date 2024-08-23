package model

import (
	"database/sql"
	"fmt"
)

type Movies struct {
	Results []Movie
}

type Movie struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	Adult       bool    `json:"adult"`
	UserComment *string `json:"user_comment,omitempty"`
	CommentId   *int    `json:"comment_id,omitempty"`
}

type Movie_Comment struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id"`
	MovieId int    `json:"movie_id"`
	Comment string `json:"comment"`
}

func CreateComment(db *sql.DB, movie_id, user_id int, comment string) error {
	query := "INSERT INTO movies_comments (user_id, movie_id, comment) VALUES (?, ?, ?)"

	_, err := db.Exec(query, user_id, movie_id, comment)

	if err != nil {
		return err
	}
	return nil
}

func GetComments(db *sql.DB, user_id int) ([]Movie_Comment, error) {
	query := "SELECT id, movie_id, comment FROM movies_comments WHERE user_id = ?"

	var movies_comments []Movie_Comment
	rows, err := db.Query(query, user_id)

	if err != nil {
		return movies_comments, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Movie_Comment
		err := rows.Scan(&comment.Id, &comment.MovieId, &comment.Comment)
		if err != nil {
			return nil, err
		}
		movies_comments = append(movies_comments, comment)
	}
	return movies_comments, nil
}

func GetCommentsByMovieId(db *sql.DB, movie_id int) (Movie_Comment, error) {
	query := "SELECT comment FROM movies_comments WHERE movie_id = ?"

	var movies_comments Movie_Comment
	err := db.QueryRow(query, movie_id).Scan(&movies_comments.Comment)

	if err != nil {
		return movies_comments, err
	}

	return movies_comments, nil
}

func UpdateComments(db *sql.DB, user_id, commentId int, comment string) error {
	query := "UPDATE movies_comments SET comment = ? WHERE id = ? AND user_id = ?"

	_, err := db.Exec(query, comment, commentId, user_id)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}

func DeleteComment(db *sql.DB, user_id, commentId int) error {
	query := "DELETE FROM movies_comments WHERE id = ? AND user_id = ?"

	_, err := db.Exec(query, commentId, user_id)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}

func GetInternalPopularMovies(db *sql.DB) ([]Movie, error) {
	var movies []Movie
	query := "SELECT id FROM movies_viewed LIMIT 10"
	rows, err := db.Query(query)

	if err != nil {
		return movies, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.Id)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func IncrementRating(db *sql.DB, movie_id int) {
	query := "INSERT INTO movies_viewed (id, quantity_views) VALUES (?, 1) ON DUPLICATE KEY UPDATE quantity_views = quantity_views + 1"

	_, err := db.Exec(query, movie_id)
	if err != nil {
		fmt.Println("Error al Agregar visualizacion de pelicula")
		fmt.Print(err)
	}
}
