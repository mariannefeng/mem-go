package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"mem.com/mem-go/db"
)

type GetBook struct {
	db *sql.DB
	q  *db.Queries
}

func NewGetBookHandler(db *sql.DB, q *db.Queries) *GetBook {
	return &GetBook{db: db, q: q}
}

type GetBookResp struct {
	ID      string       `json:"id"`
	Entries []*EntryResp `json:"entries"`
}

type EntryResp struct {
	ID      int32  `json:"id"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (g *GetBook) Handler() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		params := mux.Vars(r)
		bookID := params["book_id"]

		idInt, err := strconv.Atoi(bookID)
		if err != nil {
			return err
		}

		entries, err := g.q.GetEntriesByBook(r.Context(), int32(idInt))
		if err != nil {
			return err
		}

		var resp []*EntryResp
		for _, entry := range entries {
			resp = append(resp, &EntryResp{
				ID:      entry.ID,
				Type:    entry.Type,
				Content: entry.Content,
			})
		}

		payload, err := json.Marshal(&GetBookResp{
			ID:      bookID,
			Entries: resp,
		})

		if err != nil {
			return NewHTTPError(http.StatusInternalServerError, "failed parsing response", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)

		return nil
	}
}

func (*GetBook) Path() string {
	return "/books/{book_id}"
}

func (*GetBook) Methods() []string {
	return []string{http.MethodGet}
}

type Books struct {
	db *sql.DB
	q  *db.Queries
}

func NewBooksHandler(db *sql.DB, q *db.Queries) *Books {
	return &Books{db: db, q: q}
}

func (*Books) Path() string {
	return "/books"
}

type BookResp struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type CreateBookBody struct {
	Name string `json:"name"`
}

func (b *Books) Handler() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		switch r.Method {
		case http.MethodGet:
			books, err := b.q.GetBooks(r.Context())
			if err != nil {
				return err
			}

			var resp []*BookResp
			for _, book := range books {
				resp = append(resp, &BookResp{
					ID:   book.ID,
					Name: book.Name,
				})
			}

			payload, err := json.Marshal(resp)
			if err != nil {
				return NewHTTPError(http.StatusInternalServerError, "failed parsing response", err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(payload)
		case http.MethodPost:
			fmt.Println("creating a new book")

			var reqBody CreateBookBody

			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil {
				return NewHTTPError(http.StatusBadRequest, "invalid request body", err)
			}

			if _, err := b.q.CreateBook(r.Context(), reqBody.Name); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported request method")
		}

		return nil
	}
}

func (*Books) Methods() []string {
	return []string{
		http.MethodGet,
		http.MethodPost,
	}
}
