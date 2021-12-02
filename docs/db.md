# herro-worldとは

簡単なチャットサービスです。課題自身の詳細はこちら。https://moneyforward.kibe.la/notes/194818

要求仕様書はこちら。https://moneyforward.kibe.la/notes/194950

要するに
* ユーザーは他のユーザーをコンタクトに登録できます
* 登録した相手にメッセージを送れます
* 自分のチャット履歴を見ることもできます

DBについて、ユーザー、コンタクト、ユーザーの間のチャット、メッセージは保存する対象です。それに、グループやbinaryメッセージの対応を意識しつつDBを設計します。

# 初めに

ユーザー、コンタクト、チャット、メッセージ。この四つのデータはそれぞれの表で保存されてます。

他の三つがわかりやすいと思いますが、**チャット**について説明します。チャットというのは「会話が発生している場所」と思ってください。ですからメッセージは「相手に送る」ではなく、「自分と相手のチャットに投稿する」のです。グループ会話のためにこの設計にしました。

その上、`user_chat`の表を作成し、ユーザーとチャットの関連性を記録します。


# 共通

制約に記載がない場合は「NOT NULL」です。NULL可の場合は「NULL可」と記載します。

複数の表に共通するフィールド

| フィールド | タイプ | 制約等 | 説明 | 
| --- | --- | --- | --- |
| `created_at` | TIMESTAMP | | 作成した時のtimestamp |
| `updated_at` | TIMESTAMP | | 更新した時のtimestamp |
| `deleted_at` | TIMESTAMP | NULL可 | 削除した時のtimestamp |

# `user`表

ユーザー情報の表。`show_login_name`はログインネーム表示設定です。詳しくは要求仕様書を参照ください。

| フィールド | タイプ | 制約等 | 説明 | 
| --- | --- | --- | --- |
| `uid` | INT | AUTO_INCREMENT | 主キー |
| `login_name` | VARCHAR(40) | ユニーク | ログインネーム |
| `password` | VARCHAR(80) | | パスワード|
| `nickname` | VARCHAR(40) | | ニックネーム |
| `show_login_name` | BOOL | デフォルト値false | ログインネーム表示設定 |
| `public_key` | VARBINARY(300) | NULL可 | 公開キー |


```mysql
CREATE TABLE `user` (
-- 主キー
`uid` INT NOT NULL AUTO_INCREMENT,
-- ログインネーム
`login_name` VARCHAR(40) NOT NULL,
-- パスワード
`password` VARCHAR(80) NOT NULL,
-- ニックネーム
`nickname` VARCHAR(40) NOT NULL,
-- ログインネーム表示設定
`show_login_name` BOOL NOT NULL DEFAULT false,
-- 公開キー
`public_key` VARBINARY(300) NULL,

`created_at` TIMESTAMP NOT NULL,
`updated_at` TIMESTAMP NOT NULL,
`deleted_at` TIMESTAMP,
PRIMARY KEY (`uid`),
-- ログインネームにユニーク制約
UNIQUE (`login_name`)
);
```

# `chat`表

チャット情報の表。DMもグループもチャットと見なされ、表に`direct`というフィールドで区別する。つまりDM形のチャットは`direct == true`。

| フィールド | タイプ | 制約等 | 説明 | 
| --- | --- | --- | --- |
| `cid` | INT | AUTO_INCREMENT | 主キー|
| `direct` | BOOL | デフォルト値true | DMかどうか |
| `name` | VARCHAR(40) | NULL可 | グループ名、DMは設定不可 |

```mysql
CREATE TABLE `chat` (
-- 主キー
`cid` INT NOT NULL AUTO_INCREMENT,
-- DMかどうか
`direct` BOOL NOT NULL DEFAULT true,
-- グループ名、DMは設定不可
`name` VARCHAR(40),

`created_at` TIMESTAMP NOT NULL,
`updated_at` TIMESTAMP NOT NULL,
`deleted_at` TIMESTAMP,
PRIMARY KEY (`cid`)
);
```

# `contact`表

ユーザーのコンタクトを記録する表。

> もし
> * `(uid1, uid2)`が存在する
> * `uid1`から`uid2`へメッセージを送ったことがある
>
> なら
>
> * `(uid2, uid1)`は必ずある

※要求仕様書の「自動的コンタクト登録」を参照ください。

| フィールド | タイプ | 制約等 | 説明 | 
| --- | --- | --- | --- |
| `uid_self` | INT | | 自分のUID、主キーの一部 |
| `uid_other` | INT | | 相手のUID、主キーの一部 |
| `display_name` | VARCHAR(40) | NULL可 | 相手につけた表示名 |
| `cid` | INT | | DMのチャットID | 
| `blocked_at` | TIMESTAMP | NULL可 | ブロックした時のtimestamp |

```mysql
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
`deleted_at` TIMESTAMP,
`blocked_at` TIMESTAMP,
-- 両方のUIDの組み合わせが主キーになる
PRIMARY KEY (`uid_self`, `uid_other`)
);
```

# `user_chat`表

ユーザーが参加しているチャットの表。

`deleted_at == NULL`の場合、記録されたユーザー`uid`は`cid`というチャットに入ってる。そのチャットからの受信は可能。

`deleted_at != NULL`の場合、`uid`はもう（グループの場合）退室したまたは（DMの場合）相手を削除・ブロックした。

| フィールド | タイプ | 制約等 | 説明 | 
| --- | --- | --- | --- |
| `uid` | INT | | ユーザーのuid、主キーの一部 |
| `cid` | INT |  | チャットのcid、主キーの一部 |
| `created_at` | TIMESTAMP | | 入室した時のtimestamp |
| `deleted_at` | TIMESTAMP | NULL可 | 退室した時のtimestamp |


```mysql
CREATE TABLE `user_chat` (
-- ユーザーのuid
`uid` INT NOT NULL,
-- チャットのcid
`cid` INT NOT NULL,

`created_at` TIMESTAMP NOT NULL,
`deleted_at` TIMESTAMP,
-- 組み合わせが主キー
PRIMARY KEY (`uid`, `cid`)
);
```

# `message`表

メッセージの表。将来 binary メッセージ対応のため、`content`のタイプは`VARBINARY`、その`mime`も記録される。文字のインコーディングは UTF8。

| フィールド | タイプ | 制約等 | 説明 | 
| --- | --- | --- | --- |
| `mid` | INT | AUTO_INCREMENT | 主キー |
| `cid` | INT | | どのcidに投稿したのか |
| `uid` | INT | | 発信者のuid |
| `mime` | VARCHAR(40) | デフォルト値`'text/plain'` | 内容のmime |
| `content` | VARBINARY(4096) | | 内容 |
| `created_at` | TIMESTAMP | | 発信した時のtimestamp |

```mysql
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
`content` VARBINARY(4096) NOT NULL,
-- 発信したtimestamp
`created_at` TIMESTAMP NOT NULL,
`updated_at` TIMESTAMP NOT NULL,
`deleted_at` TIMESTAMP,
PRIMARY KEY (`mid`)
);
```