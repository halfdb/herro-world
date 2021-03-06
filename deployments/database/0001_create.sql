-- +goose Up

CREATE TABLE `user` (
-- 主キー
`uid` INT NOT NULL AUTO_INCREMENT,
-- ログインネーム
`login_name` VARCHAR(40) NOT NULL,
-- パスワード
`password` VARCHAR(80) NOT NULL,
-- ニックネーム
`nickname` VARCHAR(40),
-- ログインネーム表示設定
`show_login_name` BOOL NOT NULL DEFAULT false,

`created_at` TIMESTAMP NOT NULL,
`updated_at` TIMESTAMP NOT NULL,
`deleted_at` TIMESTAMP NULL,
PRIMARY KEY (`uid`),
-- ログインネームにユニーク制約
UNIQUE (`login_name`)
);

CREATE TABLE `chat` (
-- 主キー
`cid` INT NOT NULL AUTO_INCREMENT,
-- DMかどうか
`direct` BOOL NOT NULL DEFAULT true,
-- グループ名、DMは設定不可
`name` VARCHAR(40),

`created_at` TIMESTAMP NOT NULL,
`updated_at` TIMESTAMP NOT NULL,
`deleted_at` TIMESTAMP NULL,
PRIMARY KEY (`cid`)
);

CREATE TABLE `contact` (
-- 自分のUID
`uid_self` INT NOT NULL,
-- 相手のUID
`uid_other` INT NOT NULL,
-- 相手につけた表示名
`display_name` VARCHAR(40),
-- DMのチャットID
`cid` INT NOT NULL,

`created_at` TIMESTAMP NOT NULL,
`updated_at` TIMESTAMP NOT NULL,
`deleted_at` TIMESTAMP NULL,
`blocked_at` TIMESTAMP,
-- 両方のUIDの組み合わせが主キーになる
PRIMARY KEY (`uid_self`, `uid_other`),
CONSTRAINT contact_ibfk_1 FOREIGN KEY (`uid_self`) REFERENCES `user`(`uid`),
CONSTRAINT contact_ibfk_2 FOREIGN KEY (`uid_other`) REFERENCES `user`(`uid`),
CONSTRAINT contact_ibfk_3 FOREIGN KEY (`cid`) REFERENCES `chat`(`cid`)
);

CREATE TABLE `user_chat` (
-- ユーザーのuid
`uid` INT NOT NULL,
-- チャットのcid
`cid` INT NOT NULL,

`created_at` TIMESTAMP NOT NULL,
`deleted_at` TIMESTAMP NULL,
-- 組み合わせが主キー
PRIMARY KEY (`uid`, `cid`),
CONSTRAINT user_chat_ibfk_1 FOREIGN KEY (`uid`) REFERENCES `user`(`uid`),
CONSTRAINT user_chat_ibfk_2 FOREIGN KEY (`cid`) REFERENCES `chat`(`cid`)
);

CREATE TABLE `message` (
-- 主キー
`mid` INT NOT NULL AUTO_INCREMENT,
-- どのcidに投稿したのか
`cid` INT NOT NULL,
-- 発信者
`uid` INT NOT NULL,
-- 内容のmime
`mime` VARCHAR(40) NOT NULL DEFAULT 'text/plain',
-- 内容
`content` VARBINARY(200) NOT NULL,
-- 発信したtimestamp
`created_at` TIMESTAMP NOT NULL,
`updated_at` TIMESTAMP NOT NULL,
`deleted_at` TIMESTAMP NULL,
PRIMARY KEY (`mid`),
CONSTRAINT message_ibfk_1 FOREIGN KEY (`uid`) REFERENCES `user`(`uid`),
CONSTRAINT message_ibfk_2 FOREIGN KEY (`cid`) REFERENCES `chat`(`cid`)
);


-- +goose Down

drop table if exists `contact`;
drop table if exists `user_chat`;
drop table if exists `message`;
drop table if exists `chat`;
drop table if exists `user`;

