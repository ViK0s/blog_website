package models

import (
	"database/sql"
	"errors"
	"log"
)

// ts should have a tag string or table or smth
type Project struct {
	ID      int
	Title   string
	Content string
	Url     string
	ImgUrl  string
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type ProjectModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *ProjectModel) Insert(title string, content string, imgurl string) (int, error) {
	stmt := `INSERT INTO projects (title, content, url, imgurl)`

	result, err := m.DB.Exec(stmt, title, content, imgurl)
	if err != nil {
		return 0, err
	}
	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil

}

// var ErrNoRecord = errors.New("models: no matching record found")

// This will return a specific snippet based on its id.
func (m *ProjectModel) Get(id int) (*Project, error) {
	stmt := `SELECT id, title, content, url, imgurl FROM projects`
	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	row := m.DB.QueryRow(stmt, id)
	// Initialize a pointer to a new zeroed Snippet struct.
	s := &Project{}
	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Url, &s.ImgUrl)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own ErrNoRecord error
		// instead (we'll create this in a moment).
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	// If everything went OK then return the Snippet object.
	return s, nil
}

// This will return the the most recent blogpost
func (m *ProjectModel) Latest() []*Project {
	stmt := `SELECT id, title, content, url, imgurl FROM projects
	ORDER BY id DESC LIMIT 10`
	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result of
	// our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		print("smth went wrong")
		return nil
	}

	defer rows.Close()

	snippets := []*Project{}

	for rows.Next() {
		s := &Project{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Url, &s.ImgUrl)
		if err != nil {
			log.Fatal(err)
		}
		snippets = append(snippets, s)
	}

	return snippets
}
