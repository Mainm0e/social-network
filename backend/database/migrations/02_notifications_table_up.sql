-- +migrate Up
CREATE TABLE notifications
(
  notificationId INTEGER NULL,
  receiverId         INTEGER NOT NULL,
  senderId           INTEGER NOT NULL,
  type          TEXT    NOT NULL, --following request, group invitation,  requests to join the group, event

  content        TEXT    NOT NULL,
  creationTime   TEXT    NOT NULL,
  PRIMARY KEY (notificationId AUTOINCREMENT),
  FOREIGN KEY (receiverId) REFERENCES users (userId),
  FOREIGN KEY (senderId) REFERENCES users (userId)
);
