CREATE KEYSPACE IF NOT EXISTS chat WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};

DROP TABLE IF EXISTS chat.messages;

CREATE TABLE IF NOT EXISTS chat.messages (
    message_id      UUID,
    from_student_id text,
    to_student_id   text,
    message_body    text,
    team_id         text,
    sent_timestamp  timestamp,
    PRIMARY KEY ( (from_student_id, to_student_id), sent_timestamp )
) WITH CLUSTERING ORDER BY (sent_timestamp ASC);
