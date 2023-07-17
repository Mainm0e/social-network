-- +migrate Down
ALTER TABLE messages ADD COLUMN seen INTEGER  NULL;
ALTER TABLE messages DROP COLUMN msgType;

