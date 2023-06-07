-- +migrate Up
ALTER TABLE users ADD COLUMN private TEXT NOT NULL DEFAULT 'private' ;
