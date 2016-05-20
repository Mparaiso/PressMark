-- +migrate Up
-- Initial migration, articles
CREATE TABLE ARTICLES (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP DEFAULT( DATETIME('now') ), 
    updated TIMESTAMP DEFAULT( DATETIME('now') )
);
INSERT INTO ARTICLES(TITLE,CONTENT) VALUES(
    "This is a first article",
    "this is the content of the first article"
);

-- +migrate Down
-- remove articles
DROP TABLE ARTICLES;