CREATE DATABASE IF NOT EXISTS term_keeper_db;

USE term_keeper_db;

CREATE TABLE IF NOT EXISTS users (
      id CHAR(26),
      name VARCHAR(32),
      email VARCHAR(255) NOT NULL UNIQUE,
      password VARCHAR(255) NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS categories (
      id INT AUTO_INCREMENT,
      user_id CHAR(26),
      name VARCHAR(100),
      hex_color_code CHAR(6),
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
      PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS word_details (
      id INT AUTO_INCREMENT,
      name VARCHAR(255) NOT NULL,
      description VARCHAR(500),
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS word_category_relations (
      word_id INT NOT NULL,
      category_id INT NOT NULL,
      PRIMARY KEY(word_id, category_id),
      FOREIGN KEY (word_id) REFERENCES word_details(id) ON DELETE CASCADE,
      FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_word_relations (
      id INT NOT NULL,
      user_id CHAR(26) NOT NULL,
      FOREIGN KEY (id) REFERENCES word_details(id) ON DELETE CASCADE,
      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
      PRIMARY KEY(id, user_id)
);
