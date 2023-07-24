-- +migrate Up
CREATE TABLE users
(
  userId       INTEGER NULL,
  NickName     TEXT    NULL,
  firstName    TEXT    NOT NULL,
  lastName     TEXT    NOT NULL,
  birthDate    TEXT    NOT NULL,
  email        TEXT    NOT NULL,
  password     TEXT    NOT NULL,
  aboutMe      TEXT    NULL,
  -- NOT YET
  avatar               NULL,
  creationTime TEXT    NOT NULL,
  PRIMARY KEY (userId AUTOINCREMENT)
);

CREATE TABLE groups
(
  groupId     INTEGER NOT NULL,
  creatorId   INTEGER NOT NULL,
  title       TEXT    NOT NULL,
  description TEXT    NOT NULL,
  PRIMARY KEY (groupId),
  FOREIGN KEY (creatorId) REFERENCES users (userId)
);

CREATE TABLE posts
(
  postId       INTEGER NULL,
  userId       INTEGER NOT NULL,
  title        TEXT    NOT NULL,
  content      TEXT    NOT NULL,
  creationTime TEXT    NOT NULL,
  -- group, public, private, semi-private
  status       TEXT    NOT NULL,
  -- NOT YET
  image        NULL,
  -- if posting in a group
  groupId      INTEGER NULL DEFAULT -1,
  PRIMARY KEY (postId AUTOINCREMENT),
  FOREIGN KEY (userId) REFERENCES users (userId),
  FOREIGN KEY (groupId) REFERENCES groups (groupId)
);

CREATE TABLE comments
(
  commentId    INTEGER NULL,
  userId       INTEGER NOT NULL,
  postId       INTEGER NOT NULL,
  content      TEXT    NOT NULL,
  creationTime TEXT    NOT NULL,
  PRIMARY KEY (commentId AUTOINCREMENT),
  FOREIGN KEY (userId) REFERENCES users (userId),
  FOREIGN KEY (postId) REFERENCES posts (postId)
);

-- follower and followee relationship
CREATE TABLE follow
(
  followId  INTEGER NULL,
  followerId INTEGER NOT NULL,
  followeeId INTEGER NOT NULL,
  -- follower, pending
  status     TEXT    NOT NULL,
  PRIMARY KEY (followId AUTOINCREMENT),
  FOREIGN KEY (followerId) REFERENCES users (userId)
  FOREIGN KEY (followeeId) REFERENCES users (userId)

);

-- group-members is an N-to-N relationship, so we use a third table for it
CREATE TABLE group_member
(
  id     INTEGER NULL,
  userId  INTEGER NOT NULL,
  groupId INTEGER NOT NULL,
  -- invited, requester, member
  status  TEXT    NOT NULL,
  PRIMARY KEY (id AUTOINCREMENT),
  FOREIGN KEY (userId) REFERENCES users (userId),
  FOREIGN KEY (groupId) REFERENCES groups (groupId)
);

-- OLD VERSION
CREATE TABLE messages
(
  messageId      INTEGER NULL,
  senderId       INTEGER NOT NULL,
  receiverId     INTEGER NOT NULL,
  messageContent TEXT    NOT NULL,
  sendTime       TEXT    NOT NULL,
  seen           INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (messageId AUTOINCREMENT),
  FOREIGN KEY (senderId) REFERENCES users (userId),
  FOREIGN KEY (receiverId) REFERENCES users (userId)
);

-- selected followers can see a post
CREATE TABLE semiPrivate
(
  postId INTEGER NOT NULL,
  userId INTEGER NOT NULL,
  FOREIGN KEY (postId) REFERENCES posts (postId),
  FOREIGN KEY (userId) REFERENCES users (userId)
);

