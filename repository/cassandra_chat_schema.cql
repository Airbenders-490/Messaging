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

-- INSERT INTO chat.messages (room_id, sent_timestamp, from_student_id, message_body) VALUES(9548c1de-8f2d-4fb0-8549-22265ebc1e89,23456792,'zubair', 'helloo there friend');
-- INSERT INTO chat.participants (room_id, students) VALUES(9548c1de-8f2d-4fb0-8549-22265ebc1e89, { 0f81db5c-8e99-4c62-8a34-d84a81de414e, db483bf6-db69-4df1-98c5-e40fcd46e5dc, db483bf6-db69-4df1-98c5-e40fcd46e5dc } );
-- UPDATE chat.participants SET students = students - {0f81db5c-8e99-4c62-8a34-d84a81de414f} WHERE room_id =  9548c1de-8f2d-4fb0-8549-22265ebc1e90;