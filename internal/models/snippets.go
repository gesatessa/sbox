package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB // this will hold the database connection pool.
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	q := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(q, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// id is int64, but we want to return an int, so we need to convert it.
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	q := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// QueryRow() returns a pointer to a sql.Row struct, which represents a single row returned from the database.
	row := m.DB.QueryRow(q, id)

	var s Snippet

	// row.Scan() method is used to read the values from the `sql.Row`
	// and assign them to the fields of the Snippet struct.
	// The driver will automatically convert the database types to the appropriate Go types based on the destination variables.
	// The arguments to Scan() are pointers to the variables where the data should be stored.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// 📢 to encapsulate the odel completely, we define ErrNoRecord
		// This makes our handlers more flexible and decoupled from the database layer,
		// as they can check for ErrNoRecord instead of sql.ErrNoRows (database-specific error).
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	// return the filled-in Snippet struct.
	return s, nil
}

// why not return []*Snippet instead of []Snippet?
// Because we are not mutating the Snippet structs after scanning them from the database,
// there is no need to use pointers. Returning a slice of Snippet values is more straightforward and efficient in this case.
// The database/sql package will handle the conversion of database types to Go types when scanning the rows,
// so we can simply return a slice of Snippet values without needing to use pointers.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	q := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	// make sure to close the rows result set before the Latest() method returns.
	// This will help to prevent resource leaks and ensure that all database connections are properly released.
	// The `defer` statement should be placed immediately after checking for errors from the Query() method
	// to ensure that it is executed regardless of how the function exits (e.g., due to an error or a successful return).
	// ⚠️ If we place it befor checking for errors, and the Query() method returns an error,
	// the rows variable will be nil, and calling rows.Close() will cause a panic.
	defer rows.Close()

	// var snippets []*Snippet
	var snippets []Snippet

	for rows.Next() {
		var s Snippet
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// we still need to check for errors after iterating over the rows,
	// as there may have been an error during the iteration process
	// (e.g., a network issue or a problem with the database connection).
	// the snippets slice may have been partially filled with valid data,
	// but we want to ensure that any errors are properly handled and reported.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
