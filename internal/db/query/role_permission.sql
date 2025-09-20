-- name: CreateRole :one
INSERT INTO roles (name)
VALUES ($1)
RETURNING id, name;

-- name: GetRoleByID :one
SELECT id, name
FROM roles
WHERE id = $1;

-- name: GetRoleByName :one
SELECT id, name
FROM roles
WHERE name = $1;

-- name: ListRoles :many
SELECT id, name
FROM roles
ORDER BY id;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = $1;

-- ===============================

-- name: CreatePermission :one
INSERT INTO permissions (name)
VALUES ($1)
RETURNING id, name;

-- name: GetPermissionByID :one
SELECT id, name
FROM permissions
WHERE id = $1;

-- name: ListPermissions :many
SELECT id, name
FROM permissions
ORDER BY id;

-- name: DeletePermission :exec
DELETE FROM permissions
WHERE id = $1;

-- ===============================

-- name: AssignPermissionToRole :exec
INSERT INTO role_permission (role_id, perm_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemovePermissionFromRole :exec
DELETE FROM role_permission
WHERE role_id = $1 AND perm_id = $2;

-- name: ListPermissionsByRole :many
SELECT p.id, p.name
FROM permissions p
JOIN role_permission rp ON rp.perm_id = p.id
WHERE rp.role_id = $1;

-- ===============================

-- name: AssignRoleToUser :exec
INSERT INTO user_role (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveRoleFromUser :exec
DELETE FROM user_role
WHERE user_id = $1 AND role_id = $2;

-- name: ListRolesByUser :many
SELECT r.id, r.name
FROM roles r
JOIN user_role ur ON ur.role_id = r.id
WHERE ur.user_id = $1;
