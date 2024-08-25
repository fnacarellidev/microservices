-- name: CreateUser :one
INSERT INTO users (
	username, password
) VALUES (
	$1, $2
) RETURNING id;

-- name: GetPasswordFromUser :one
SELECT password FROM users WHERE username = $1;

-- name: CreatePost :one
INSERT INTO posts (
	post_owner, content
) VALUES (
	$1, $2
) RETURNING id;

-- name: GetPostsFromUser :many
SELECT content, created_at FROM posts WHERE post_owner = $1;
