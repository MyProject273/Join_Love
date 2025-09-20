-- Insert roles
INSERT INTO roles (name) VALUES 
    ('user_premium'),
    ('admin');

-- Insert permissions
INSERT INTO permissions (name) VALUES
-- Users
('create_users'),
('read_users'),
('update_users'),
('delete_users'),

-- User Profiles
('create_user_profiles'),
('read_user_profiles'),
('update_user_profiles'),
('delete_user_profiles'),

-- Photos
('create_photos'),
('read_photos'),
('update_photos'),
('delete_photos'),

-- Likes
('create_likes'),
('read_likes'),
('delete_likes'),

-- Matches
('read_matches'),
('delete_matches'),

-- Messages
('create_messages'),
('read_messages'),
('delete_messages'),

-- Payments
('create_payments'),
('read_payments'),

-- Subscriptions
('create_subscriptions'),
('read_subscriptions'),
('cancel_subscriptions'),

-- Reports
('create_reports'),
('read_reports'),
('delete_reports'),

-- Privacy Settings
('update_privacy_settings'),
('read_privacy_settings'),

-- Notifications
('read_notifications'),
('delete_notifications');

-- Assign all permissions to admin
INSERT INTO role_permission (role_id, perm_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin';

-- Assign limited permissions to user_premium
INSERT INTO role_permission (role_id, perm_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name IN (
    -- User Profiles
    'create_user_profiles', 'read_user_profiles', 'update_user_profiles',
    -- Photos
    'create_photos', 'read_photos',
    -- Likes & Matches
    'create_likes', 'read_likes', 'read_matches',
    -- Messages
    'create_messages', 'read_messages',
    -- Payments & Subscriptions
    'create_payments', 'create_subscriptions', 'read_subscriptions', 'cancel_subscriptions',
    -- Reports
    'create_reports',
    -- Privacy
    'update_privacy_settings', 'read_privacy_settings',
    -- Notifications
    'read_notifications'
)
WHERE r.name = 'user';
