-- +goose Up
ALTER TABLE posts DROP CONSTRAINT posts_status_check;
ALTER TABLE posts ADD CONSTRAINT posts_status_check CHECK (status IN ('draft', 'public', 'private'));
ALTER TABLE posts ALTER COLUMN status SET DEFAULT 'draft';

-- +goose Down
ALTER TABLE posts DROP CONSTRAINT posts_status_check;
ALTER TABLE posts ADD CONSTRAINT posts_status_check CHECK (status IN ('public', 'private'));
ALTER TABLE posts ALTER COLUMN status SET DEFAULT 'public';