SET CHARSET utf8;

DROP TABLE IF EXISTS users;

CREATE TABLE users
(
    id               bigint unsigned AUTO_INCREMENT,
    name             text NOT NULL,
    hash             VARCHAR(12) unique NOT NULL,
    disabled         boolean DEFAULT false,
    created_at       timestamp NOT NULL DEFAULT current_timestamp,
    updated_at       timestamp NOT NULL DEFAULT current_timestamp on update current_timestamp,
    PRIMARY KEY (id)
) DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

DROP TABLE IF EXISTS work_times;

CREATE TABLE work_times
(
    id               bigint unsigned AUTO_INCREMENT,
    user_id          bigint unsigned NOT NULL,
    content          VARCHAR(100) unique NOT NULL,
    supplement       text NOT NULL,
    is_finished      boolean DEFAULT false,
    disabled         boolean DEFAULT false,
    started_at       timestamp NOT NULL DEFAULT current_timestamp,
    finished_at      timestamp NOT NULL DEFAULT current_timestamp on update current_timestamp,
    FOREIGN KEY (user_id)
    REFERENCES users(id),
    PRIMARY KEY (id)
) DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

DROP TABLE IF EXISTS work_rests;

CREATE TABLE work_rests
(
    id               bigint unsigned AUTO_INCREMENT,
    work_time_id     bigint unsigned NOT NULL,
    is_finished      boolean DEFAULT false,
    disabled         boolean DEFAULT false,
    started_at       timestamp NOT NULL DEFAULT current_timestamp,
    finished_at      timestamp NOT NULL DEFAULT current_timestamp on update current_timestamp,
    FOREIGN KEY (work_time_id)
    REFERENCES work_times(id),
    PRIMARY KEY (id)
) DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

DROP TABLE IF EXISTS session_work_times;

CREATE TABLE session_work_times
(
    work_time_id     bigint unsigned NOT NULL,
    hash             VARCHAR(31) unique NOT NULL,
    disabled         boolean   NOT NULL DEFAULT false,
    created_at       timestamp NOT NULL DEFAULT current_timestamp,
    updated_at       timestamp NOT NULL DEFAULT current_timestamp,
)