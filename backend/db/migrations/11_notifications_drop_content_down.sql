-- +migrate Down
ALTER TABLE notifications ADD COLUMN content TEXT  NULL;