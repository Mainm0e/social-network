-- +migrate Up
ALTER TABLE groups ADD COLUMN creationTime TEXT NOT NULL  DEFAULT "0000-00-00 00:00:00";  