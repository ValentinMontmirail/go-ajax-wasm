
CREATE TABLE authors (
    -- 'id' column: stores unique identifiers for each author
    id integer PRIMARY KEY AUTOINCREMENT,  -- 'integer' indicates the data type is an integer
                                          -- 'PRIMARY KEY' means 'id' is the unique identifier for each record
                                          -- 'AUTOINCREMENT' automatically increments this value for new rows

    -- 'name' column: stores the name of the authors
    name text NOT NULL,                   -- 'text' indicates the data type is text
                                          -- 'NOT NULL' means this field cannot be left empty

    -- 'bio' column: stores the biography of the authors
    bio text NOT NULL                     -- Similarly, 'text' for string data and
                                          -- 'NOT NULL' ensures this field must have a value
);
