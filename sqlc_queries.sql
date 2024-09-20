-- name: CreateUser :one
INSERT INTO users (
	username, password
) VALUES (
	$1, $2
) RETURNING id;

-- name: GetPasswordFromUser :one
SELECT password FROM users WHERE username = $1;

-- name: GetRecordsFromUser :many
SELECT r.id, r.title, r.content, r.created_at
FROM users u
JOIN diary d ON u.id = d.diary_owner
JOIN records r ON d.id = r.diary_id
WHERE u.username = $1;