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
    `username` TEXT NOT NULL,
    `key` BLOB NOT NULL,
    `salt` BLOB NOT NULL,
    `api` BLOB NOT NULL,
    `roles` TEXT NOT NULL,
    `expiry` INTEGER NOT NULL,
    `membership` TEXT NOT NULL DEFAULT 'free',
    `concurrents` INTEGER DEFAULT 1,
    `servers` INTEGER DEFAULT 1,
    `duration` INTEGER DEFAULT 60,
    `balance` INTEGER DEFAULT 0,
    `apiReqs` INTEGER DEFAULT 1,
    `apiFails` INTEGER DEFAULT 1,
    `tele` TEXT NOT NULL
);
