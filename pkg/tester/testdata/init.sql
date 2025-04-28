-- データベースの作成
CREATE DATABASE IF NOT EXISTS term_keeper_db_test;

USE term_keeper_db_test;

-- テーブル作成
CREATE TABLE IF NOT EXISTS users (
      id CHAR(26),
      name VARCHAR(32),
      email VARCHAR(255) NOT NULL UNIQUE,
      password VARCHAR(255) NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS categories (
      id INT AUTO_INCREMENT,
      user_id CHAR(26),
      name VARCHAR(100),
      hex_color_code CHAR(6),
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
      PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS term_details (
      id INT AUTO_INCREMENT,
      name VARCHAR(255) NOT NULL,
      description VARCHAR(500),
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS term_category_relations (
      term_id INT NOT NULL,
      category_id INT NOT NULL,
      PRIMARY KEY(term_id, category_id),
      FOREIGN KEY (term_id) REFERENCES term_details(id) ON DELETE CASCADE,
      FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_term_relations (
      id INT NOT NULL,
      user_id CHAR(26) NOT NULL,
      FOREIGN KEY (id) REFERENCES term_details(id) ON DELETE CASCADE,
      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
      PRIMARY KEY(id, user_id)
);

-- テストデータの挿入

-- ユーザーデータ挿入
INSERT INTO users (id, name, email, password) VALUES
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', '山田太郎', 'yamada@example.com', '$2a$10$abcdefghijklmnopqrstuv'),
('01HGDJ5HXZD3K6WFYS9JU0A1XG', '佐藤花子', 'sato@example.com', '$2a$10$wxyzabcdefghijklmnopqr'),
('01HGDJ5J8KF4L7XGZT0KV1B2YH', '鈴木一郎', 'suzuki@example.com', '$2a$10$rstuvwxyzabcdefghijklm');

-- カテゴリーデータ挿入
INSERT INTO categories (user_id, name, hex_color_code) VALUES
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', 'プログラミング', 'FF5733'),
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', 'データベース', '33A8FF'),
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', 'ネットワーク', '33FF57'),
('01HGDJ5HXZD3K6WFYS9JU0A1XG', '機械学習', 'D433FF'),
('01HGDJ5HXZD3K6WFYS9JU0A1XG', 'クラウド', 'FFD633'),
('01HGDJ5J8KF4L7XGZT0KV1B2YH', 'セキュリティ', 'FF3333');

-- 用語データ挿入
INSERT INTO term_details (name, description) VALUES
('SQL', 'Structured Query Language。リレーショナルデータベースの操作に使用される言語。'),
('TCP/IP', 'インターネット通信の基盤となるプロトコル群。'),
('Docker', 'コンテナ型の仮想化技術。'),
('AWS', 'Amazonが提供するクラウドコンピューティングサービス。'),
('TLS', 'Transport Layer Security。通信の暗号化プロトコル。'),
('Python', '汎用プログラミング言語の一つ。機械学習やデータ分析によく使われる。'),
('Git', '分散型バージョン管理システム。'),
('REST API', 'REpresentational State Transferに基づくAPI設計アーキテクチャ。'),
('NoSQL', '非リレーショナルデータベース。'),
('CI/CD', '継続的インテグレーション/継続的デリバリー。');

-- ユーザーと用語の関連付け
INSERT INTO user_term_relations (id, user_id) VALUES
(1, '01HGDJ5GZRJ2J5VEXR8HT8V9WF'),
(2, '01HGDJ5GZRJ2J5VEXR8HT8V9WF'),
(3, '01HGDJ5GZRJ2J5VEXR8HT8V9WF'),
(4, '01HGDJ5GZRJ2J5VEXR8HT8V9WF'),
(5, '01HGDJ5GZRJ2J5VEXR8HT8V9WF'),
(1, '01HGDJ5HXZD3K6WFYS9JU0A1XG'),
(4, '01HGDJ5HXZD3K6WFYS9JU0A1XG'),
(6, '01HGDJ5HXZD3K6WFYS9JU0A1XG'),
(7, '01HGDJ5HXZD3K6WFYS9JU0A1XG'),
(5, '01HGDJ5J8KF4L7XGZT0KV1B2YH'),
(8, '01HGDJ5J8KF4L7XGZT0KV1B2YH'),
(9, '01HGDJ5J8KF4L7XGZT0KV1B2YH'),
(10, '01HGDJ5J8KF4L7XGZT0KV1B2YH');

-- カテゴリーと用語の関連付け
INSERT INTO term_category_relations (term_id, category_id) VALUES
(1, 2), -- SQL → データベース
(2, 3), -- TCP/IP → ネットワーク
(3, 1), -- Docker → プログラミング
(4, 5), -- AWS → クラウド
(5, 6), -- TLS → セキュリティ
(6, 1), -- Python → プログラミング
(6, 4), -- Python → 機械学習
(7, 1), -- Git → プログラミング
(8, 1), -- REST API → プログラミング
(8, 3), -- REST API → ネットワーク
(9, 2), -- NoSQL → データベース
(10, 1); -- CI/CD → プログラミング