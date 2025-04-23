CREATE TABLE `attacks` (
    `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
    `method` TEXT NOT NULL,
    `target` TEXT NOT NULL,
    `port` INTEGER DEFAULT 0,
    `threads` INTEGER DEFAULT 3,
    `pps` INTEGER DEFAULT 250000,
    `parent` INTEGER NOT NULL,
    `duration` INTEGER,
    `type` INTEGER DEFAULT 1,
    `stopped` INTEGER DEFAULT 0,
    `date` INTEGER NOT NULL
);

CREATE TABLE `news` (
    `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
    `title` TEXT NOT NULL,
    `from` TEXT NOT NULL,
    `content` TEXT NOT NULL,
    `date` INTEGER NOT NULL
);

CREATE TABLE `sales` (
    `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
    `uniqid` TEXT NOT NULL,
    `amount` INTEGER NOT NULL,
    `crypto_amount` DOUBLE NOT NULL,
    `crypto_address` TEXT NOT NULL,
    `recieved` DOUBLE NOT NULL,
    `coin` TEXT NOT NULL,
    `status` TEXT NOT NULL,
    `product` TEXT NOT NULL,
    `parent` INTEGER NOT NULL,
    `date` INTEGER NOT NULL
);

CREATE TABLE `blacklists` (
    `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
    `host` TEXT NOT NULL
);

CREATE TABLE `tickets` (
    `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
    `user_id` INTEGER NOT NULL,
    `title` TEXT NOT NULL,
    `status` BLOB NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `username` TEXT NOT NULL
);

CREATE TABLE `messages` (
    `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
    `ticket_id` INTEGER NOT NULL,
    `user_id` INTEGER NOT NULL,
    `message` TEXT NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `tokens` (
    `id` INTEGER PRIMARY KEY AUTO_INCREMENT,  
    `token` TEXT NOT NULL,                    
    `expiry` DATETIME NOT NULL,
    `usernames` TEXT NOT NULL
);

CREATE TABLE `users` (
    `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
    `username` VARCHAR(100) NOT NULL,
    `key` BLOB NOT NULL,
    `salt` BLOB NOT NULL,
    `api` BLOB NOT NULL,
    `roles` TEXT NOT NULL,
    `expiry` INTEGER NOT NULL,
    `membership` VARCHAR(50) NOT NULL DEFAULT 'free',
    `concurrents` INTEGER DEFAULT 1,
    `servers` INTEGER DEFAULT 1,
    `duration` INTEGER DEFAULT 60,
    `balance` INTEGER DEFAULT 0,
    `apiReqs` INTEGER DEFAULT 1,
    `apiFails` INTEGER DEFAULT 1,
    `tele` VARCHAR(100) NOT NULL
);

CREATE TABLE `settings_conf` (
    `domain` VARCHAR(255) NOT NULL,
    `autobuy_flat` VARCHAR(255) NOT NULL,
    `freeuser` TINYINT(1) NOT NULL,
    `freeuser_roles` TEXT NOT NULL,
    `freeuser_cons` INTEGER DEFAULT 0,
    `freeuser_time` INTEGER DEFAULT 0,
    `fake_users` INTEGER DEFAULT 0,
    `fake_ongoing` INTEGER DEFAULT 0,
    `cnc` TINYINT(1) NOT NULL,
    `ssh` VARCHAR(255) NOT NULL,
    `telnet` VARCHAR(255) NOT NULL,
    `api` TINYINT(1) NOT NULL,
    `subnet1` VARCHAR(255),
    `subnet2` VARCHAR(255),
    `subnet3` VARCHAR(255),
    `subnet4` VARCHAR(255),
    `subnet5` VARCHAR(255),
    `subnet6` VARCHAR(255),
    `subnet7` VARCHAR(255),
    `subnet8` VARCHAR(255),
    `subnet9` VARCHAR(255)
);

INSERT INTO settings_conf (
    domain, autobuy_flat, freeuser, freeuser_roles,
    freeuser_cons, freeuser_time, fake_users, fake_ongoing,
    cnc, ssh, telnet, api,
    subnet1,  subnet2,    subnet3,
    subnet4,  subnet5,    subnet6,
    subnet7,  subnet8,    subnet9
) VALUES (
    'https://example.com/',   -- domain
    'USD',                    -- autobuy_flat
    1,                        -- freeuser
    'user,vip',               -- freeuser_roles
    5,                        -- freeuser_cons
    60,                       -- freeuser_time
    10,                       -- fake_users
    3,                        -- fake_ongoing
    1,                        -- cnc
    '1337',                   -- ssh
    '1227',                   -- telnet
    1,                        -- api

    'UDP',       -- subnet1
    'Layer 7',   -- subnet2
    'TCP',       -- subnet3
    'Game',      -- subnet4
    'Bypass',    -- subnet5
    '',          -- subnet6
    '',          -- subnet7
    '',          -- subnet8
    ''           -- subnet9
);


CREATE TABLE flood_methods (
  id          INTEGER PRIMARY KEY AUTO_INCREMENT,
  name        VARCHAR(255) NOT NULL,
  sname        VARCHAR(255) NOT NULL,
  type        INTEGER NOT NULL DEFAULT 0,
  description TEXT NOT NULL,
  subnet      INTEGER NOT NULL DEFAULT 0,
  mtype       INTEGER NOT NULL DEFAULT 0,
  vip         TINYINT(1) NOT NULL
);

INSERT INTO flood_methods
  (id, name, type, description, subnet, mtype, vip)
VALUES
  (42, 'UDP', 'Site name for UDP', 1, 'UDP METHOD', 1, 1, 0);

  CREATE TABLE `flood_apis` (
    id          INTEGER PRIMARY KEY AUTO_INCREMENT,
    name
    type
    slots
    url
    
  )