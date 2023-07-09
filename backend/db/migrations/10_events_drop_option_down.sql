-- +migrate Down
ALTER TABLE events ADD COLUMN option TEXT,
ALTER TABLE events ADD COLUMN receiverId INTEGER,
-- Add foreign key constraint
ALTER TABLE events ADD FOREIGN KEY (receiverId) REFERENCES users (userId);
