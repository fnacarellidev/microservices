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

-- name: GetIdFromUser :one
SELECT id FROM users WHERE username = $1;

-- name: GetDiaryFromUser :one
SELECT d.*
FROM diary d
JOIN users u ON d.diary_owner = u.id
WHERE u.id = $1;

-- name: CreateDiaryForUser :exec
INSERT INTO diary (diary_owner)
VALUES ($1);

-- name: CreateRecordOnUserDiary :exec
INSERT INTO records (diary_id, title, content)
VALUES ($1, $2, $3);
