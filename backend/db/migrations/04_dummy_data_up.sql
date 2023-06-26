-- +migrate Up
INSERT INTO users (NickName, firstName, lastName, birthDate, email, password, aboutMe, avatar, creationTime)
VALUES ('johnny', 'John', 'Doe', '1990-05-15', 'john.doe@example.com', 'password123', 'About John', NULL, '2023-05-31 10:00:00'),
       ('jane23', 'Jane', 'Smith', '1988-08-22', 'jane.smith@example.com', 'pass321word', 'About Jane', NULL, '2023-05-31 11:00:00'),
       ('sam87', 'Sam', 'Johnson', '1995-02-10', 'sam.johnson@example.com', 'testpass', 'About Sam', NULL, '2023-05-31 12:00:00');

INSERT INTO groups (creatorId, title, description)
VALUES (1, 'Group 1', 'Description for Group 1'),
       (2, 'Group 2', 'Description for Group 2'),
       (1, 'Group 3', 'Description for Group 3');

INSERT INTO posts (userId, title, content, creationTime, status, groupId,image)
VALUES (1, 'Post 1', 'Content of Post 1', '2023-05-31 13:00:00', 'public', 0,""),
       (2, 'Post 2', 'Content of Post 2', '2023-05-31 14:00:00', 'private', 0,""),
       (1, 'Post 3', 'Content of Post 3', '2023-05-31 15:00:00', 'group', 1,"");

INSERT INTO comments (userId, postId, content, creationTime)
VALUES (2, 1, 'Comment on Post 1', '2023-05-31 16:00:00'),
       (1, 3, 'Comment on Post 3', '2023-05-31 17:00:00');

INSERT INTO follow (followerId, followeeId, status)
VALUES (1, 2, 'following'),
       (2, 1, 'following'),
       (1, 3, 'following');

INSERT INTO group_member (userId, groupId, status)
VALUES (1, 1, 'member'),
       (2, 1, 'member'),
       (1, 2, 'member');

INSERT INTO semiPrivate (postId, userId)
VALUES (1, 2),
       (3, 1);

-- +migrate Up
INSERT INTO notifications (receiverId, senderId, type, content, creationTime)
VALUES (1, 2, 'following request', 'You have a new follower request', '2023-05-31 18:00:00'),
       (2, 1, 'group invitation', 'You have been invited to join Group 2', '2023-05-31 19:00:00');

INSERT INTO events (creatorId, receiverId, groupId, title, content, creationTime)
VALUES (1, 2, 1, 'Event 1', 'Content of Event 1', '2023-05-31 20:00:00'),
       (2, 1, 2, 'Event 2', 'Content of Event 2', '2023-05-31 21:00:00');

INSERT INTO messages (senderId, receiverId, messageContent, sendTime, seen)
VALUES (1, 2, 'Message 1', '2023-05-31 22:00:00', 0),
       (2, 1, 'Message 2', '2023-05-31 23:00:00', 0);

-- +migrate Down
DELETE FROM messages;

-- Remove the previously inserted rows from other tables (follow, group_member, semiPrivate, notifications, events)
DELETE FROM follow WHERE followerId IN (1, 2) AND followeeId IN (2, 1,3) AND status = 'following';
DELETE FROM group_member WHERE userId IN (1, 2) AND groupId IN (1, 2) AND status = 'member';
DELETE FROM semiPrivate WHERE postId IN (1, 3) AND userId IN (1, 2);
DELETE FROM notifications WHERE receiverId IN (1, 2) AND senderId IN (2, 1);
DELETE FROM events WHERE creatorId IN (1, 2) AND receiverId IN (1, 2) AND groupId IN (1, 2);

-- Remove the previously inserted rows from the main tables (users, groups, posts, comments)
DELETE FROM users WHERE NickName IN ('johnny', 'jane23', 'sam87') AND firstName IN ('John', 'Jane', 'Sam');
DELETE FROM groups WHERE title IN ('Group 1', 'Group 2', 'Group 3');
DELETE FROM posts WHERE title IN ('Post 1', 'Post 2', 'Post 3');
DELETE FROM comments WHERE content IN ('Comment on Post 1', 'Comment on Post 3');
