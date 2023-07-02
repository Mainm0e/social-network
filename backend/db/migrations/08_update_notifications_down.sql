-- +migrate Down
ALTER TABLE notifications DROP COLUMN groupId;