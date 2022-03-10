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
    deleted timestamp,
    class text
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

-- michael is eaf54fae-1ab8-4b5a-8047-51904f6ae884
-- dwight is 172ff420-f0eb-4d75-a26b-8d058a8499ec
INSERT INTO chat.room (roomid, admin, name, students, class) VALUES ('office', 'eaf54fae-1ab8-4b5a-8047-51904f6ae884', 'office', {'eaf54fae-1ab8-4b5a-8047-51904f6ae884', 'toby', '172ff420-f0eb-4d75-a26b-8d058a8499ec'}, 'soen490');
INSERT INTO chat.student_rooms (student, rooms) VALUES ('eaf54fae-1ab8-4b5a-8047-51904f6ae884', {'office'});
INSERT INTO chat.student_rooms (student, rooms) VALUES ('toby', {'office'});
INSERT INTO chat.room (roomid, admin, name, students, class) VALUES ('allstars', '172ff420-f0eb-4d75-a26b-8d058a8499ec' , 'allstars', {'172ff420-f0eb-4d75-a26b-8d058a8499ec', 'jim'}, 'soen385');
INSERT INTO chat.student_rooms (student, rooms) VALUES ('172ff420-f0eb-4d75-a26b-8d058a8499ec', {'allstars', 'office'});
INSERT INTO chat.student_rooms (student, rooms) VALUES ('jim', {'allstars'});
INSERT INTO chat.student (student_id, first_name, last_name) VALUES ('jim', 'jim', 'halpert');
INSERT INTO chat.student (student_id, first_name, last_name) VALUES ('eaf54fae-1ab8-4b5a-8047-51904f6ae884', 'michael', 'scott');
INSERT INTO chat.student (student_id, first_name, last_name) VALUES ('172ff420-f0eb-4d75-a26b-8d058a8499ec', 'dwight', 'schrute');
INSERT INTO chat.student (student_id, first_name, last_name) VALUES ('toby', 'toby', 'flenderson');
INSERT INTO chat.student (student_id, first_name, last_name) VALUES ('pam', 'pam', 'beesly');
INSERT INTO chat.student (student_id, first_name, last_name) VALUES ('kevin', 'kevin', 'malone');
INSERT INTO chat.student (student_id, first_name, last_name) VALUES ('angela', 'angela', 'martin');