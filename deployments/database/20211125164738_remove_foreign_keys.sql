-- +goose Up

ALTER TABLE `contact`
    DROP FOREIGN KEY contact_ibfk_1;
ALTER TABLE `contact`
    DROP FOREIGN KEY contact_ibfk_2;
ALTER TABLE `contact`
    DROP FOREIGN KEY contact_ibfk_3;

ALTER TABLE `message`
    DROP FOREIGN KEY message_ibfk_1;
ALTER TABLE `message`
    DROP FOREIGN KEY message_ibfk_2;

ALTER TABLE `user_chat`
    DROP FOREIGN KEY user_chat_ibfk_1;
ALTER TABLE `user_chat`
    DROP FOREIGN KEY user_chat_ibfk_2;


-- +goose Down

ALTER TABLE `contact`
    ADD CONSTRAINT contact_ibfk_1 FOREIGN KEY (`uid_self`) REFERENCES `user`(`uid`);
ALTER TABLE `contact`
    ADD CONSTRAINT contact_ibfk_2 FOREIGN KEY (`uid_other`) REFERENCES `user`(`uid`);
ALTER TABLE `contact`
    ADD CONSTRAINT contact_ibfk_3 FOREIGN KEY (`cid`) REFERENCES `chat`(`cid`);

ALTER TABLE `message`
    ADD CONSTRAINT message_ibfk_1 FOREIGN KEY (`uid`) REFERENCES `user`(`uid`);
ALTER TABLE `message`
    ADD CONSTRAINT message_ibfk_2 FOREIGN KEY (`cid`) REFERENCES `chat`(`cid`);

ALTER TABLE `user_chat`
    ADD CONSTRAINT user_chat_ibfk_1 FOREIGN KEY (`uid`) REFERENCES `user`(`uid`);
ALTER TABLE `user_chat`
    ADD CONSTRAINT user_chat_ibfk_2 FOREIGN KEY (`cid`) REFERENCES `chat`(`cid`);
