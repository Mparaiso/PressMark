-- +migrate Up
-- Initial migration, articles
CREATE TABLE ARTICLES (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    TITLE VARCHAR(255) NOT NULL,
    CONTENT TEXT NOT NULL,
    CREATED TIMESTAMP DEFAULT( DATETIME('now') ), 
    UPDATED TIMESTAMP DEFAULT( DATETIME('now') )
);
INSERT INTO ARTICLES(TITLE,CONTENT) VALUES(
    "This is a first article",
    "this is the content of the first article"
);

-- +migrate Down
-- remove articles
DROP TABLE ARTICLES;