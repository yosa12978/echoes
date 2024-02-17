CREATE TABLE posts (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL
);

CREATE TABLE links (
    id VARCHAR(36) PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    created TIMESTAMP NOT NULL
);

CREATE TABLE comments (
    id VARCHAR(36) PRIMARY KEY,
    email TEXT NOT NULL,
    name TEXT NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    postId VARCHAR(36)
);

ALTER TABLE comments ADD FOREIGN KEY (postId) REFERENCES posts (id);

CREATE TABLE accounts (
    id VARCHAR(36) PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    created TIMESTAMP NOT NULL
);