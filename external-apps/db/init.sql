CREATE DATABASE IF NOT EXISTS term_keeper_db;

USE term_keeper_db;

-- テーブル作成
CREATE TABLE IF NOT EXISTS users (
      id CHAR(26) NOT NULL,
      name VARCHAR(32) NOT NULL,
      email VARCHAR(255) NOT NULL UNIQUE,
      password VARCHAR(255) NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS categories (
      id INT AUTO_INCREMENT NOT NULL,
      fk_user_id CHAR(26) NOT NULL,
      name VARCHAR(100) NOT NULL,
      hex_color_code CHAR(7),
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      FOREIGN KEY (fk_user_id) REFERENCES users(id) ON DELETE CASCADE,
      PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS terms (
      id INT AUTO_INCREMENT NOT NULL,
      fk_user_id CHAR(26) NOT NULL,
      name VARCHAR(255) NOT NULL,
      description VARCHAR(500),
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      FOREIGN KEY (fk_user_id) REFERENCES users(id) ON DELETE CASCADE,
      PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS term_category_relations (
      fk_term_id INT NOT NULL,
      fk_category_id INT NOT NULL,
      FOREIGN KEY (fk_term_id) REFERENCES terms(id) ON DELETE CASCADE,
      FOREIGN KEY (fk_category_id) REFERENCES categories(id) ON DELETE CASCADE,
      PRIMARY KEY(fk_term_id, fk_category_id)
);
