-- +migrate Up
ALTER TABLE messages DROP COLUMN seen ;
ALTER TABLE messages ADD COLUMN msgType TEXT NOT NULL DEFAULT 'privateMsg' ;