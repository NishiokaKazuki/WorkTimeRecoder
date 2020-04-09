DELETE FROM users;

INSERT INTO users(id, name, hash, disabled) VALUES
    (1,'nishioka', 'U010TJV6G1E', false);

DELETE FROM work_times;

INSERT INTO work_times(id, user_id, content, supplement, is_finished, disabled, started_at, finished_at) VALUES
    (1, 1, 'てすとおお', '', true, false, '2020-04-01 10:00:00', '2020-04-01 18:30:00'),
    (2, 1, 'てすとお２', '買出時間含む', true, false, '2020-04-03 8:30:00', '2020-04-01 17:15:00');

INSERT INTO work_times(work_time_id, is_finished, disabled, started_at, finished_at) VALUES
    (1, true, false, '2020-04-01 12:00:00', '2020-04-01 12:30:00'),
    (2, true, false, '2020-04-03 13:00:00', '2020-04-03 14:00:00');