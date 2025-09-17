-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    user_name
) VALUES (
    $1,$2,$3
)
RETURNING *;


-- name: GetListUser :many
SELECT * FROM  users 
LIMIT $1
OFFSET $2;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1 ;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = $1 ;

-- name: UpdateUser :one
UPDATE users
SET user_name = COALESCE(sqlc.narg(user_name), user_name),
    phone = COALESCE(sqlc.narg(phone), phone),
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    gender = COALESCE(sqlc.narg(gender), gender),
    birthdate = COALESCE(sqlc.narg(birthdate), birthdate),
    avatar_url = COALESCE(sqlc.narg(avatar_url), avatar_url),
    bio = COALESCE(sqlc.narg(bio), bio)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateUserActiveStatus :one
UPDATE users
SET is_active = $1
WHERE id = $2
RETURNING *;

-- name: UpdateUserVerifiedStatus :one
UPDATE users
SET is_verified = $1
WHERE id = $2
RETURNING *;

-- name: UpdateUserLastLogin :exec
Update users
SET last_login = $1
WHERE id = $2;

-- name: CheckUserVerified :one
SELECT is_verified
FROM users
WHERE id = $1
  AND deleted_at IS NULL;