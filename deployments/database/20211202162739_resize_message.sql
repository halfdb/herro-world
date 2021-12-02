-- +goose Up
ALTER TABLE `message`
    MODIFY COLUMN `content` VARBINARY(4096) NOT NULL;

-- +goose Down
ALTER TABLE `message`
    MODIFY COLUMN `content` VARBINARY(200) NOT NULL;
