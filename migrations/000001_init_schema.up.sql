CREATE TABLE IF NOT EXISTS users (
    id int PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS segments (
    segment_name varchar(255) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS user_segments (
                                             user_id int NOT NULL,
                                             segment_name varchar(255) NOT NULL,
                                             PRIMARY KEY (user_id, segment_name),
                                             FOREIGN KEY (user_id) REFERENCES users(id),
                                             FOREIGN KEY (segment_name) REFERENCES segments(segment_name) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS history (
                                       user_id int,
                                       segment_name varchar(255),
                                       action_type varchar,
                                       action_timestamp timestamp DEFAULT now()
);

