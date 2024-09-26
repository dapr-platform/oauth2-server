-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE o_client_info(
                              id VARCHAR(32) NOT NULL,
                              password VARCHAR(32) NOT NULL,
                              PRIMARY KEY (id)
);
insert into o_client_info (id,password) values ('dapr-client','123456');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS o_client_info cascade ;

-- +goose StatementEnd
