-- +migrate Up
PRAGMA foreign_keys = 0;

-- Create a temporary table with the desired structure
CREATE TABLE events_backup (
    eventId INTEGER PRIMARY KEY AUTOINCREMENT,
    creatorId INTEGER NOT NULL,
    groupId INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    creationTime TEXT NOT NULL
);

-- Copy the data from the original table to the backup table
INSERT INTO events_backup SELECT eventId, creatorId, groupId, title, content, creationTime FROM events;

-- Drop the original table
DROP TABLE events;

-- Rename the backup table to the original table name
ALTER TABLE events_backup RENAME TO events;

PRAGMA foreign_keys = 1;