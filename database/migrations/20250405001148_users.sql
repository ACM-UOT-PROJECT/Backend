-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists user (
	id integer primary key,
	user_name text not null,
	token text not null,
	created_at datetime default CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists user;
-- +goose StatementEnd
