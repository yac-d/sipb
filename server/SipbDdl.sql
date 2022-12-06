CREATE DATABASE IF NOT EXISTS sipb;

DROP USER IF EXISTS 'sipb_server';
CREATE USER IF NOT EXISTS 'sipb_server' IDENTIFIED BY 'srvr';
GRANT ALL PRIVILEGES ON sipb.* TO 'sipb_server'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;

USE sipb;

DROP TABLE IF EXISTS fileitem;
DROP TABLE IF EXISTS item;

CREATE TABLE IF NOT EXISTS item (
    id CHAR(36) NOT NULL PRIMARY KEY,
    ts DATETIME DEFAULT CURRENT_TIMESTAMP,
    note VARCHAR(200) NOT NULL DEFAULT ""
);

CREATE TABLE IF NOT EXISTS fileitem (
    id CHAR(36) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(100) NOT NULL UNIQUE,
    bytes INT NOT NULL,
    mimetype VARCHAR(100) NOT NULL,
    CONSTRAINT FOREIGN KEY (id) REFERENCES item(id) ON DELETE CASCADE ON UPDATE CASCADE
);

DROP PROCEDURE IF EXISTS INSERT_FILE;
DROP PROCEDURE IF EXISTS NTH_MOST_RECENT_FILE;

DELIMITER //

CREATE PROCEDURE IF NOT EXISTS INSERT_FILE(uuid CHAR(36), name VARCHAR(100), location VARCHAR(100), bytes INT, mimetype VARCHAR(100), note VARCHAR(200))
    BEGIN
        INSERT INTO item(id, note) VALUES (uuid, note);
        INSERT INTO fileitem VALUES (uuid, name, location, bytes, mimetype);
	END
//

CREATE PROCEDURE IF NOT EXISTS NTH_MOST_RECENT_FILE(N INT)
    BEGIN
        SELECT fileitem.*, item.ts, item.note FROM fileitem INNER JOIN item ON item.id = fileitem.id
        ORDER BY item.ts DESC LIMIT 1 OFFSET N;
    END
//

DELIMITER ;
