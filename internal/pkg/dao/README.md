# Main task of this module

* Mainly works the wrapper of sqlboiler's `models` package
* Also handles constraints that are not reflected in the `models` package
    * e.g. A direct chat should be created before user A adds user B to the contact
    * e.g. When a chat is created as a row in the `chat` table, multiple `user_chat` rows also need to be created
    * e.g. When user A deletes user B from the contacts, A also quits the direct chat between them.

# Verbs

## Reading verbs

## `Fetch`
* A function to get a row or multiple rows using their identifiers (e.g. primary key)
* Should include a `withDeleted` param unless it's not necessary
* Should not assume the existence of the required value
* Must return `sql.ErrNoRows` as an error if no row is found

## `Exist`
* A function to check if an item exist
* Resembles the `Fetch` function
* Must return false and no error when the specified row is not found
* Must return `bool`

## `Lookup`
* A function to get row(s) using some conditions other than identifiers
* Usually involves complex sql queries
* Should include `withDeleted<TableName>` params to decide whether deleted rows in the `<TableName>` would be considered
* Should not assume the existence of the required value
* Must return `sql.ErrNoRows` as the error if the row is not found
* The type of return value depends on the lookup

## Writing verbs

* In this project, writing functions should write one row at a time for data security. Batch writing functions, for example, are exceptions.
* Writing functions must accept a `*sql.Tx` parameter and use it as the sole method to operate the database.

## `Create`
* A function to create the specified ONE row
* May skip checking if all the constraints are fulfilled before inserting into the database
* Should accept DAO objects instead of specifications

## `Update`
* A function to update the specified ONE row
* Should assume the specified row exists
* Should accept DAO objects instead of identifiers
* Should not be used to delete or restore

## `Delete`
* A function to delete the specified ONE row
* Should assume the specified row exists and is not deleted
* Should accept DAO objects instead of identifiers

## `Restore`
* A function to restore ONE deleted row
* Should assume the specified row exists and is deleted
* Should accept DAO objects instead of identifiers

# Recommended Practise

## Read before write

Get a DAO object before updating it. This practise confirms the existence of the target row and constrains the number of target rows to 1.
