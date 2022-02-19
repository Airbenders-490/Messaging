CREATE KEYSPACE IF NOT EXISTS chat WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};

DROP TABLE IF EXISTS chat.messages;

CREATE TABLE IF NOT EXISTS chat.messages (
    room_id         text,
    from_student_id text,
    message_body    text,
    sent_timestamp  timestamp,
    PRIMARY KEY ( (room_id), sent_timestamp )
) WITH CLUSTERING ORDER BY (sent_timestamp DESC);

CREATE TABLE IF NOT EXISTS  chat.room (
    roomid text PRIMARY KEY,
    Name text,
    Admin text,
    students set<text>,
    deleted timestamp
);

CREATE TABLE IF NOT EXISTS chat.student_rooms (
    student text primary key,
    rooms set<text>
);

CREATE TABLE IF NOT EXISTS chat.student (
    student_id text PRIMARY KEY,
    first_name text,
    last_name text
);

INSERT INTO chat.room (roomid, admin, name, students) VALUES ('office', 'michael' , 'office', {'michael', 'toby', 'dwight'});
INSERT INTO chat.student_rooms (student, rooms) VALUES ('michael', {'office'});
INSERT INTO chat.student_rooms (student, rooms) VALUES ('toby', {'office'});
INSERT INTO chat.room (roomid, admin, name, students) VALUES ('allstars', 'dwight' , 'allstars', {'dwight', 'jim'});
INSERT INTO chat.student_rooms (student, rooms) VALUES ('dwight', {'allstars', 'office'});
INSERT INTO chat.student_rooms (student, rooms) VALUES ('jim', {'allstars'});
INSERT INTO chat.student (student_id) VALUES ('jim');
INSERT INTO chat.student (student_id) VALUES ('f882141d-07dd-460d-91da-877db8a3a670');
INSERT INTO chat.student (student_id) VALUES ('michael');
INSERT INTO chat.student (student_id) VALUES ('dwight');
INSERT INTO chat.student (student_id) VALUES ('toby');
INSERT INTO chat.student (student_id) VALUES ('pam');
INSERT INTO chat.student (student_id) VALUES ('kevin');
INSERT INTO chat.student (student_id) VALUES ('angela');