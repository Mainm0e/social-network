-- +migrate Up
ALTER TABLE notifications ADD COLUMN groupId INTEGER  DEFAULT 0;  