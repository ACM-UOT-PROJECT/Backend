-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists state (
	id integer primary key,
	winner text,
	tries integer,
	people integer
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists state;
-- +goose StatementEnd
