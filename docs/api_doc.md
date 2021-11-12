herro-world API設計書
=====

Document for Chat APIs

# Common

* Accessing APIs other than `/login` without a valid token results in `403`
* Accessing unauthorized resources (e.g. others' contacts) results in `403` as well

# Account

Account related APIs

## `POST /login`
Login with login_name and password


### Params

| Field name | Mandatory | Comment |
|---|---|---|
|`login_name` | `true` | |
|`password` | `true` | |

### Return
- `200`: Successful login. Returns token.
- `403`: Auth failed. Login and try again.

### Return example
```
{ "token": "token_string" }
```
---
## `POST /users`
Register


### Params

| Field name | Mandatory | Comment |
|---|---|---|
|`login_name` | `true` | |
|`nickname` | `true` | |
|`password` | `true` | |

### Return
- `200`: Success.
- `409`: Unable to create user due to conflict.

---
## `GET /users/:uid`
Read a user's info.


### Return
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/user.json

---
## `POST /users/:uid`
Update user info of oneself.


### Params

| Field name | Mandatory | Comment |
|---|---|---|
|`nickname` | `false` |  |
|`show_login_name` | `false` |  |
|`password` | `false` |  |

### Return
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/user.json
- `400`: At least one of `nickname`, `show_login_name`, `password` must be specified.
- `403`: Auth failed. Login and try again.


# Contact management

APIs to manage contacts.

## `GET /users/:uid/contacts`
View one's contacts. Deleted or blocked contacts are invisible.


### Return
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/contacts.json
- `403`: Auth failed. Login and try again.

---
## `POST /users/:uid/contacts`
Add a user into contact.


### Params

| Field name | Mandatory | Comment |
|---|---|---|
|`uid` | `true` | |
|`display_name` | `false` |  |

### Return
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/contact.json
- `403`: Auth failed. Login and try again.

---
## `POST /users/:uid/contacts/:uid_other`
Update one of one's contacts. Deleted or blocked contacts cannot be updated.


### Params

| Field name | Mandatory | Comment |
|---|---|---|
|`display_name` | `false` |  |
|`deleted` | `false` | The only allowed value to POST is true. |
|`blocked` | `false` | The only allowed value to POST is true. |

### Return
- `200`: The contact if not deleted or blocked. Empty else. https://raw.githubusercontent.com/halfdb/herro-world/main/schema/contact.json
- `400`: At least one of `display_name`, `deleted`, `blocked` must be specified.
- `403`: Auth failed. Login and try again.

---
## `DELETE /users/:uid/contacts/:uid_other`
Delete a contact. A shortcut to POST `deleted=true`


### Return
- `200`: Success.
- `403`: Auth failed. Login and try again.


# Chatting

Chatting related APIs

## `GET /users/:uid/chats`
View chat list


### Return
- `200`: Returns a list of `cid`s of the chats that the user is in.
- `403`: Auth failed. Login and try again.

---
## `GET /chats/:cid/messages`
View messages in a chat. Newest messages first by default.


### Params

| Field name | Mandatory | Comment |
|---|---|---|
|`page_size` | `false` | Default 100. |
|`page_token` | `false` | `next_page_token` from a previous response |

### Return
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/messages.json
- `403`: Auth failed. Login and try again.

---
## `POST /chats/:cid/messages`
Post new message into a chat.


### Params

| Field name | Mandatory | Comment |
|---|---|---|
|`mime_id` | `false` |  |
|`content` | `true` | BASE64 of content. |

### Return
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/message.json
- `403`: Auth failed. Login and try again.


# Misc

Other APIs.

## `GET /mimes`
List available MIMEs


### Return
- `200`: Success.

### Return example
```
[{100: "text/plain"}]
```
