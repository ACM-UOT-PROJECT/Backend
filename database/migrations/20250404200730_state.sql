-- +goose Up
-- +goose StatementBegin
create table if not exists state (
	id integer primary key,
	winner text,
	tries integer,
	people integer
);

insert into state (winner, tries, people) values ("NO_ONE_WON_AT_ALL", 0, 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists state;
-- +goose StatementEnd
