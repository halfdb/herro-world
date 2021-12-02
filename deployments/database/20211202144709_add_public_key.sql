-- +goose Up
ALTER TABLE `user`
ADD COLUMN `public_key` VARBINARY(300) AFTER `show_login_name`;

-- +goose Down
ALTER TABLE `user`
DROP COLUMN `public_key`;
