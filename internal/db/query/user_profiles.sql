-- name: GetUserProfile :one
SELECT *
FROM user_profiles
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: CreateUserProfile :one
INSERT INTO user_profiles (
  user_id,
  job_title,
  education,
  interests,
  relationship_goal,
  religion,
  smoking,
  drinking,
  height_cm,
  weight_kg,
  created_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: UpdateUserProfile :one
UPDATE user_profiles
SET
  job_title = COALESCE($2, job_title),
  education = COALESCE($3, education),
  interests = COALESCE($4, interests),
  relationship_goal = COALESCE($5, relationship_goal),
  religion = COALESCE($6, religion),
  smoking = COALESCE($7, smoking),
  drinking = COALESCE($8, drinking),
  height_cm = COALESCE($9, height_cm),
  weight_kg = COALESCE($10, weight_kg),
  updated_at = now()
WHERE user_id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteUserProfile :exec
UPDATE user_profiles
SET deleted_at = now()
WHERE user_id = $1 AND deleted_at IS NULL;
