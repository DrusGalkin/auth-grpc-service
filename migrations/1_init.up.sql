CREATE TABLE oldmine.users
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL,
    email         VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role          ENUM('admin', 'moder', 'user') NOT NULL DEFAULT 'user',
    UNIQUE KEY unique_username (username),
    UNIQUE KEY unique_email (email)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_0900_ai_ci
COMMENT='Пользователи';

CREATE INDEX idx_users_username ON oldmine.users (username);
CREATE INDEX idx_users_email ON oldmine.users (email);