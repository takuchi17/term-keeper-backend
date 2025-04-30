-- データベースの作成
CREATE DATABASE IF NOT EXISTS term_keeper_db_test;
USE term_keeper_db_test;

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


-- テストデータの挿入

-- ユーザーデータ挿入
INSERT INTO users (id, name, email, password) VALUES
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', '山田太郎', 'yamada@example.com', '$2a$10$abcdefghijklmnopqrstuv'),
('01HGDJ5HXZD3K6WFYS9JU0A1XG', '佐藤花子', 'sato@example.com', '$2a$10$wxyzabcdefghijklmnopqr'),
('01HGDJ5J8KF4L7XGZT0KV1B2YH', '鈴木一郎', 'suzuki@example.com', '$2a$10$rstuvwxyzabcdefghijklm');

-- カテゴリーデータ挿入
INSERT INTO categories (fk_user_id, name, hex_color_code) VALUES
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', 'プログラミング', '#FF5733'),
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', 'データベース', '#33A8FF'),
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', 'ネットワーク', '#33FF57'),
('01HGDJ5HXZD3K6WFYS9JU0A1XG', '機械学習', '#D433FF'),
('01HGDJ5HXZD3K6WFYS9JU0A1XG', 'クラウド', '#FFD633'),
('01HGDJ5J8KF4L7XGZT0KV1B2YH', 'セキュリティ', '#FF3333');

-- 用語データ挿入
INSERT INTO terms (fk_user_id, name, description) VALUES
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', 'SQL', 'Structured Query Language。リレーショナルデータベースの操作に使用される言語。'),
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', 'TCP/IP', 'インターネット通信の基盤となるプロトコル群。'),
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', 'Docker', 'コンテナ型の仮想化技術。'),
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', 'AWS', 'Amazonが提供するクラウドコンピューティングサービス。'),
('01HGDJ5GZRJ2J5VEXR8HT8V9WF', 'TLS', 'Transport Layer Security。通信の暗号化プロトコル。'),
('01HGDJ5HXZD3K6WFYS9JU0A1XG', 'Python', '汎用プログラミング言語の一つ。機械学習やデータ分析によく使われる。'),
('01HGDJ5HXZD3K6WFYS9JU0A1XG', 'Git', '分散型バージョン管理システム。'),
('01HGDJ5J8KF4L7XGZT0KV1B2YH', 'REST API', 'REpresentational State Transferに基づくAPI設計アーキテクチャ。'),
('01HGDJ5J8KF4L7XGZT0KV1B2YH', 'NoSQL', '非リレーショナルデータベース。'),
('01HGDJ5J8KF4L7XGZT0KV1B2YH', 'CI/CD', '継続的インテグレーション/継続的デリバリー。');

-- カテゴリーと用語の関連付け
INSERT INTO term_category_relations (fk_term_id, fk_category_id) VALUES
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
