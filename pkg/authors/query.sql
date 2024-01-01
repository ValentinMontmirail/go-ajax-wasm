-- name: GetAuthor :one
-- This query is named 'GetAuthor' and is expected to return one result (':one').
-- It selects all columns (*) from the 'authors' table.
-- 'WHERE id = ?' filters the records to find the author with a specific 'id'.
-- 'LIMIT 1' ensures that only one record is returned, even if there are multiple matches (which shouldn't happen with a primary key).
SELECT * FROM authors
WHERE id = ? LIMIT 1;

-- name: ListAuthors :many
-- This query is named 'ListAuthors' and can return many results (':many').
-- It selects all columns (*) from the 'authors' table.
-- 'ORDER BY ID' sorts the results by 'ID'.
SELECT * FROM authors
ORDER BY ID DESC;

-- name: CreateAuthor :one
-- This query is named 'CreateAuthor' and returns the newly created record (':one').
-- It inserts a new row into the 'authors' table with values for 'name' and 'bio'.
-- The '?' placeholders are for parameter substitution with actual values.
-- 'RETURNING *' returns all columns of the newly inserted row.
INSERT INTO authors (
  name, bio
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateAuthor :exec
-- This query is named 'UpdateAuthor' and is an action query (':exec').
-- It updates the 'name' and 'bio' columns in the 'authors' table.
-- 'WHERE id = ?' specifies which author to update, based on their 'id'.
-- The '?' are placeholders for the new values of 'name', 'bio', and the 'id' of the author to update.
UPDATE authors
set name = ?,
bio = ?
WHERE id = ?;

-- name: DeleteAuthor :exec
-- This query is named 'DeleteAuthor' and is an action query (':exec').
-- It deletes a row from the 'authors' table.
-- 'WHERE id = ?' specifies the 'id' of the author to delete.
-- The '?' is a placeholder for the 'id' of the author to be deleted.
DELETE FROM authors
WHERE id = ?;
