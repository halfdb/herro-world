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
- `403`: Forbidden.

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
- `200`: OK.
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
- `403`: Forbidden.


# Contact management

APIs to manage contacts.

## `GET /users/:uid/contacts`
View one's contacts. Deleted or blocked contacts are invisible.


### Return
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/contacts.json
- `403`: Forbidden.

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
- `403`: Forbidden.

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
- `403`: Forbidden.

---
## `DELETE /users/:uid/contacts/:uid_other`
Delete a contact. A shortcut to POST `deleted=true`


### Return
- `200`: OK.
- `403`: Forbidden.


# Chatting

Chatting related APIs

## `GET /users/:uid/chats`
View chat list


### Return
- `200`: Returns a list of `cid`s of the chats that the user is in.
- `403`: Forbidden.

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
- `403`: Forbidden.

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
- `403`: Forbidden.


# Group

Group related APIs

## `POST /chats`
Create group


### Params

| Field name | Mandatory | Comment |
|---|---|---|
|`uids` | `true` | UIDs of the members. At least 3. In the form of `uids=1&uids=2&uids=3`. |
|`name` | `false` |  |

### Return
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/chat.json
- `403`: Adding users who are not in contacts is not allowed.

---
## `GET /chats/:cid/members`
Check group members


### Return
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/users.json

---
## `POST /chats/:cid/members`
Add members


### Params

| Field name | Mandatory | Comment |
|---|---|---|
|`uids` | `true` | UIDs of the added members. At least 1. In the form of `uids=1&uids=2&uids=3`. |

### Return
- `200`: https://raw.githubusercontent.com/halfdb/herro-world/main/schema/users.json
- `403`: Adding users who are not in contacts is not allowed.

