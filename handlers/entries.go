package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"mem.com/mem-go/db"
)

type CreateEntry struct {
	db *sql.DB
	q  *db.Queries
}

func NewEntryHandler(db *sql.DB, q *db.Queries) *CreateEntry {
	return &CreateEntry{db: db, q: q}
}

type CreateEntryBody struct {
	Content string  `json:"content"`
	Type    string  `json:"type"`
	Key     *string `json:"key,omitempty"`
}

func (c *CreateEntry) Handler() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		params := mux.Vars(r)
		bookID := params["book_id"]

		var reqBody CreateEntryBody

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			return NewHTTPError(http.StatusBadRequest, "invalid request body", err)
		}

		idInt, err := strconv.Atoi(bookID)
		if err != nil {
			return err
		}

		var key sql.NullString

		if reqBody.Key != nil {
			key = sql.NullString{
				Valid:  true,
				String: *reqBody.Key,
			}
		}

		if _, err := c.q.CreateEntry(r.Context(), db.CreateEntryParams{
			BookID:  int32(idInt),
			Type:    reqBody.Type,
			Content: reqBody.Content,
			Key:     key,
		}); err != nil {
			return err
		}

		return nil
	}
}

func (*CreateEntry) Path() string {
	return "/books/{book_id}/entries"
}

func (*CreateEntry) Methods() []string {
	return []string{
		http.MethodPost,
	}
}

type DeleteEntry struct {
	db *sql.DB
	q  *db.Queries
}

func NewDeleteEntryHandler(db *sql.DB, q *db.Queries) *DeleteEntry {
	return &DeleteEntry{db: db, q: q}
}

func (c *DeleteEntry) Handler() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		params := mux.Vars(r)
		bookID := params["book_id"]
		entryID := params["entry_id"]

		idInt, err := strconv.Atoi(entryID)
		if err != nil {
			return err
		}

		bookIdInt, err := strconv.Atoi(bookID)
		if err != nil {
			return err
		}

		return c.q.DeleteEntry(r.Context(), db.DeleteEntryParams{
			ID:     int32(idInt),
			BookID: int32(bookIdInt),
		})

		// TODO: delete the corresponding entry in s3
	}
}

func (*DeleteEntry) Path() string {
	return "/books/{book_id}/entries/{entry_id}"
}

func (*DeleteEntry) Methods() []string {
	return []string{
		http.MethodDelete,
	}
}
