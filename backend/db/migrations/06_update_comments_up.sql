-- +migrate Up
ALTER TABLE comments ADD COLUMN image TEXT NULL ;