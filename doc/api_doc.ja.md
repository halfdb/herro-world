herro-world API設計書
=====

herro-worldに使うAPIです。

# 共通

* login以外のAPIにリクエストする時有効なトークンを付けないと`403`が出る
* 自分の見られないリソース（他人のコンタクトとか）も`403`

# アカウント

アカウント関連のAPI

## `POST /login`
ログイン


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`login_name` | `true` | |
|`password` | `true` | |

### 戻り値
- `200`: 成功。トークンを返す
- `403`: Authentication失敗。ログインしてリトライしてください。

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
## `POST /users/:uid`
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
- `403`: Authentication失敗。ログインしてリトライしてください。


# コンタクト管理

コンタクトを管理するAPI

## `GET /users/:uid/contacts`
自分のコンタクトを見る。削除・ブロックされたコンタクトは見られない。


### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/contacts.json
- `403`: Authentication失敗。ログインしてリトライしてください。

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
- `403`: Authentication失敗。ログインしてリトライしてください。

---
## `POST /users/:uid/contacts/:uid_other`
コンタクトなかの一つを更新。削除・ブロックされたコンタクトは更新できない。


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`display_name` | `false` |  |
|`deleted` | `false` | trueしか受け入れない |
|`blocked` | `false` | trueしか受け入れない |

### 戻り値
- `200`: 削除・ブロックされてないなら更新したコンタクトを返す。それ以外は空のレスポンスを返す。 https://raw.githubusercontent.com/halfdb/herro-world/main/schema/contact.json
- `400`: `display_name`, `deleted`, `blocked`中の一つ以上を指定ください
- `403`: Authentication失敗。ログインしてリトライしてください。

---
## `DELETE /users/:uid/contacts/:uid_other`
一つのコンタクトを削除。`deleted=true`をPOSTするのと同様。


### 戻り値
- `200`: 成功
- `403`: Authentication失敗。ログインしてリトライしてください。


# チャット

チャット関連のAPI

## `GET /users/:uid/chats`
チャットリストを見る


### 戻り値
- `200`: 自分のチャットの`cid`のリスト
- `403`: Authentication失敗。ログインしてリトライしてください。

---
## `GET /chats/:cid/messages`
チャット内のメッセージを見る。最新情報が先の順で。


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`page_size` | `false` | デフォルト値100。 |
|`page_token` | `false` | 前のレスポンスの`next_page_token` |

### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/messages.json
- `403`: Authentication失敗。ログインしてリトライしてください。

---
## `POST /chats/:cid/messages`
チャットにメッセージを投稿


### パラメーター

| フィールド | 必須 | コメント |
|---|---|---|
|`mime_id` | `false` |  |
|`content` | `true` | 内容のBASE64 |

### 戻り値
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/message.json
- `403`: Authentication失敗。ログインしてリトライしてください。


# 他

他のAPI

## `GET /mimes`
可能のmimeの一覧


### 戻り値
- `200`: 成功

### 戻り値の例
```
[{100: "text/plain"}]
```
