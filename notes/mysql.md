install MySQL
```sh
sudo apt install mysql-server
# keep a mental note of the password you set for the `root` user
```

connect to MySQL as root user
```sh
sudo mysql

# if it doesn't work, try 👇
# mysql -u root -p


# mysql>
```

```sql
SHOW DATABASES;

USE my_database;

SHOW TABLES;

DESCRIBE my_table;
```

create a database in MySQL for our project

```sql
-- create a nwe UTF-8 `snippetbox` db
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE snippetbox;

CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

-- add an index on the `created` column
CREATE INDEX idx_snippets_created ON snippets(created);

DESC snippets;

-- add some placeholder entries
INSERT INTO snippets (title, content, created, expires) VALUES (
    'Take a short walk',
    'Step outside for 10 minutes to refresh your mind.',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Stretch your body',
    'Do a quick 2-minute stretch to ease muscle tension.',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);

SELECT * FROM snippets;

```

⚠️ Not a good idea to connect to MySQL as the `root` user from a web app.
We'll create a database user with restricted permissions on the database.

```sql
CREATE USER 'web'@'localhost';

GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';

-- make sure to have a secure password
ALTER USER 'web'@'localhost' IDENTIFIED BY 'ChangeME';

EXIT
```

Connect to the `snippetbox` database as the `web` user:
```sh
mysql -D snippetbox -u web -p

# mysql >
```

Test permissions:
```sql
INSERT INTO snippets (title, content, created, expires) VALUES (
    'Plan your day',
    'Write down your top 3 priorities for today.',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);

SELECT id, title, expires FROM snippets;

-- should be denied, as `web` is not authorized to drop tables.
DROP TABLE snippets;
```
