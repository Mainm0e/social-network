-- +migrate Up
CREATE TABLE events (
  eventId INTEGER NULL,
  creatorId INTEGER NOT NULL,
  receiverId INTEGER NOT NULL,
  groupId INTEGER NOT NULL,
  title    TEXT    NOT NULL,
  content TEXT    NOT NULL,
  creationTime   TEXT    NOT NULL,
  option  INTEGER NULL,
  PRIMARY KEY (eventId AUTOINCREMENT),
  FOREIGN KEY (creatorId) REFERENCES users (userId),
  FOREIGN KEY (receiverId) REFERENCES users (userId),
  FOREIGN KEY (groupId) REFERENCES groups (groupId)
);
