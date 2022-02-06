-- CREATE Tables
CREATE TABLE CATEGORIES (
    id SERIAL,
    category_name VARCHAR(50),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY(id, category_name),
    UNIQUE(id)
);

CREATE TABLE AUTHORS (
    id SERIAL,
    firstname VARCHAR(50),
    lastname VARCHAR(50),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY(id, firstname, lastname),
    UNIQUE(id)
);

CREATE TABLE BOOKS (
    id SERIAL,
    category_id INT,
    author_id INT,
    title VARCHAR(50),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY(id, category_id, author_id, title),
    CONSTRAINT fk_category
      FOREIGN KEY(category_id) 
	  REFERENCES CATEGORIES(id)
      ON DELETE SET NULL,
    CONSTRAINT fk_author
      FOREIGN KEY(author_id) 
	  REFERENCES AUTHORS(id)
      ON DELETE SET NULL
);

-- Create trigger for updating
CREATE OR REPLACE FUNCTION update_changetimestamp_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now(); 
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_timestamp BEFORE UPDATE
ON CATEGORIES FOR EACH ROW EXECUTE PROCEDURE update_changetimestamp_column();

CREATE TRIGGER update_timestamp BEFORE UPDATE
ON AUTHORS FOR EACH ROW EXECUTE PROCEDURE update_changetimestamp_column();

CREATE TRIGGER update_timestamp BEFORE UPDATE
ON BOOKS FOR EACH ROW EXECUTE PROCEDURE update_changetimestamp_column();

-- Insert Initial recoreds
INSERT INTO 
  AUTHORS(firstname,lastname,created_at) 
VALUES
  ('Tim','Davis',CURRENT_TIMESTAMP),
  ('Bob','Smith',CURRENT_TIMESTAMP),
  ('Isabella','Johnson',CURRENT_TIMESTAMP),
  ('Jack','Williams',CURRENT_TIMESTAMP),
  ('Emma','Brown',CURRENT_TIMESTAMP);

INSERT INTO
  CATEGORIES(category_name, created_at)
VALUES
  ('Biographies',CURRENT_TIMESTAMP),
  ('Business',CURRENT_TIMESTAMP),
  ('Technology',CURRENT_TIMESTAMP),
  ('Education',CURRENT_TIMESTAMP),
  ('Mystery',CURRENT_TIMESTAMP);

INSERT INTO BOOKS
 	(category_id, author_id, title, created_at)
SELECT
	cat.id as category_id, auth.id as author_id, 'How to success',CURRENT_TIMESTAMP
FROM categories cat
CROSS JOIN
	authors auth
WHERE
	 cat.category_name = 'Business'  AND auth.lastname in ('Johnson', 'Williams','Davis') ;

UPDATE BOOKS
SET title = 'I nailed it!'
FROM authors
WHERE
  BOOKS.author_id = authors.id
  AND AUTHORS.lastname = 'Davis';

UPDATE BOOKS
SET title = 'If you want to be Johnson'
FROM authors
WHERE
  BOOKS.author_id = authors.id
  AND AUTHORS.lastname = 'Johnson';

UPDATE BOOKS
SET title = 'The House of Williams'
FROM authors
WHERE
  BOOKS.author_id = authors.id
  AND AUTHORS.lastname = 'Williams';

