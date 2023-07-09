-- +migrate Up
CREATE TABLE event_member (
  eventId INTEGER NOT NULL,
  memberId INTEGER NOT NULL,
  option  INTEGER NOT NULL,
  FOREIGN KEY (eventId) REFERENCES events (eventId),
  FOREIGN KEY (memberId) REFERENCES users (userId)
);
