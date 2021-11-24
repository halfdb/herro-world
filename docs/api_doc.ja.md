herro-world API設計書
=====

herro-worldに使うAPIです。

# 共通

* `POST /login`, `POST /users/`以外のAPIにリクエストする時有効なトークンを付けないと`400`が出る
* 解析できないリクエストには`400`
* 自分の見ることのできないリソース（他人のコンタクトとか）は`403`

# アカウント

アカウント関連のAPI

## `POST /login`
ログイン。成功した場合はJWTトークンを返します。その後のリクエストのheaderに`Authorization: Bearer <token_string>`を追加する必要があります。


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`login_name` | `true` | |
|`password` | `true` | |

### 戻り値
- `200`: 成功。トークンを返す
- `401`: 認証失敗。
### 戻り値の例
```
{ "token": "token_string" }
```
---
## `POST /users`
登録


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`login_name` | `true` | |
|`nickname` | `true` | |
|`password` | `true` | |

### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/user.json
- `409`: コンフリクトのため、ユーザー作成できなかった。

---
## `GET /users/:uid`
ユーザー情報


### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/user.json

---
## `PATCH /users/:uid`
自分の情報を更新


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`nickname` | `false` |  |
|`show_login_name` | `false` |  |
|`password` | `false` |  |

### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/user.json
- `400`: `nickname`, `show_login_name`, `password`中の一つ以上を指定ください


# コンタクト管理

コンタクトを管理するAPI

## `GET /users/:uid/contacts`
自分のコンタクトを見る。削除・ブロックされたコンタクトは見れない。


### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/contacts.json

---
## `POST /users/:uid/contacts`
コンタクトにユーザー追加


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`uid` | `true` | |
|`display_name` | `false` |  |

### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/contact.json
- `409`: コンフリクトのため、コンタクト作成できなかった。

---
## `PATCH /users/:uid/contacts/:uid_other`
コンタクトなかの一つを更新。削除・ブロックされたコンタクトは更新できない。


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`display_name` | `true` |  |

### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/contact.json
- `404`: 存在しない。
---
## `DELETE /users/:uid/contacts/:uid_other`
一つのコンタクトを削除。


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`blocked` | `false` | ブロックもするかどうか |

### 戻り値
- `200`: 成功。
- `404`: 存在しない。

# チャット

チャット関連のAPI

## `GET /users/:uid/chats`
チャットリストを見る


### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/chats.json

---
## `GET /chats/:cid/messages`
チャット内のメッセージを見る。作成日時が新しいものから順に取得します。


### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/messages.json

---
## `POST /chats/:cid/messages`
チャットにメッセージを投稿


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`mime` | `false` |  |
|`content` | `true` | 内容のBASE64 |

### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/message.json
- `403`: 自分をブロックした相手にDM送るのは禁止
- `413`: contentの長さは上限を超えた


# グループ

グループ関連のAPI

## `POST /chats`
グループを作る


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`uids` | `true` | 参加者のUID。三人以上。形は`uids=1&uids=2&uids=3`。 |
|`name` | `false` |  |

### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/chat.json
- `403`: 自分のコンタクトにないユーザーを追加するのは禁止されます

---
## `GET /chats/:cid/members`
グループメンバーを確認


### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/users.json

---
## `POST /chats/:cid/members`
グループメンバーを追加


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`uids` | `true` | 追加する参加者のUID。一人以上。形は`uids=1&uids=2&uids=3`。 |

### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/users.json
- `403`: 自分のコンタクトにないユーザーを追加するのは禁止されます

