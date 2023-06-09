-- +migrate Up
ALTER TABLE users ADD COLUMN privacy TEXT NOT NULL DEFAULT 'privacy' ;
