CREATE TABLE posts(
    post_id VARCHAR(255) PRIMARY KEY UNIQUE,
    title VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    allow_comments BOOLEAN NOT NULL default true,
    created_at TIMESTAMP with time zone NOT NULL,
    updated_at TIMESTAMP with time zone
);

CREATE TABLE comments(
    comment_id VARCHAR(255) NOT NULL UNIQUE,
    post_id VARCHAR(255) REFERENCES posts(post_id) ON DELETE CASCADE NOT NULL,
    parent_id VARCHAR(255),
    content TEXT NOT NULL,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP with time zone NOT NULL,
    updated_at TIMESTAMP with time zone
);