package database

import "wasm/pkg/authors"

func (db *MyDB) GetAllAuthors() ([]authors.Author, error) {
	// Create a new instance of the 'authors' queries with the current database connection
	queries := authors.New(db.db)

	// Begin a new transaction
	tx, err := db.db.Begin()
	if err != nil {
		// If there's an error starting a transaction, return the error
		return nil, err
	}

	// Execute the 'ListAuthors' query within the context of the transaction
	authors, err := queries.WithTx(tx).ListAuthors(db.ctx)
	if err != nil {
		// If there's an error executing 'ListAuthors', attempt to rollback the transaction
		err2 := tx.Rollback()
		if err != nil {
			// If there's an error during rollback, return this error
			return nil, err2
		}
		// If rollback is successful, return the original error from 'ListAuthors'
		return nil, err
	}

	// If 'ListAuthors' executes successfully, commit the transaction
	err = tx.Commit()
	if err != nil {
		// If there's an error during commit, return the error
		return nil, err
	}

	// Return the list of authors and a nil error on successful execution
	return authors, nil
}

func (db *MyDB) CreateAuthor(name string, bio string) (authors.Author, error) {
	// defaultAuthor is a default instance of authors.Author struct
	defaultAuthor := authors.Author{}

	// Create a new instance of the 'authors' queries with the current database connection
	queries := authors.New(db.db)

	// Begin a new transaction
	tx, err := db.db.Begin()
	if err != nil {
		// If there's an error starting a transaction, return the default author object and the error
		return defaultAuthor, err
	}

	// Execute the 'CreateAuthor' query within the context of the transaction
	// Pass 'name' and 'bio' as parameters to the query
	authors, err := queries.WithTx(tx).CreateAuthor(db.ctx, authors.CreateAuthorParams{
		Name: name,
		Bio:  bio,
	})
	if err != nil {
		// If there's an error executing 'CreateAuthor', attempt to rollback the transaction
		err2 := tx.Rollback()
		if err != nil {
			// If there's an error during rollback, return the default author object and this error
			return defaultAuthor, err2
		}
		// If rollback is successful, return the default author object and the original error
		return defaultAuthor, err
	}

	// If 'CreateAuthor' executes successfully, commit the transaction
	err = tx.Commit()
	if err != nil {
		// If there's an error during commit, return the default author object and the error
		return defaultAuthor, err
	}

	// Return the newly created author object and a nil error on successful execution
	return authors, nil
}
