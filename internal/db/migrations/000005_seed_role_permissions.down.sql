DELETE FROM role_permissions WHERE role_id IN (1, 2, 3);

DELETE FROM roles WHERE id IN (2, 3);

DELETE FROM permissions
WHERE id BETWEEN 1 AND 32;