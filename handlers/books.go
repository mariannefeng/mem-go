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

type Books struct{}

func NewBooksHandler() *Books {
	return &Books{}
}

func (*Books) Path() string {
	return "/books"
}

func (*Books) Handler() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		fmt.Println("either returns list of all books or creates a new book")
		return nil
	}
}

func (*Books) Methods() []string {
	return []string{
		http.MethodGet,
		http.MethodPost,
	}
}
