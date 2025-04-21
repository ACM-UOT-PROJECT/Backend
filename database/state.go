package database

import (
	"backend/database/gen/model"
	t "backend/database/gen/table"

	s "github.com/go-jet/jet/v2/sqlite"
)

type State struct {
	Winner string
	Tries  int
	People int
}

func (d *DataService) scanState(stmt s.Statement) (State, error) {
	dest := model.State{}

	err := stmt.Query(d.db, &dest)
	if err != nil {
		d.logger.Debug(stmt.DebugSql())
		d.logger.Error(err.Error())
		return State{}, err
	}

	r := State{
		Winner: *dest.Winner,
		Tries:  int(*dest.Tries),
		People: int(*dest.People),
	}

	return r, nil
}

func (d *DataService) GetState() (State, error) {
	stmt := s.SELECT(
		t.State.AllColumns,
	).FROM(
		t.State,
	).LIMIT(1)

	return d.scanState(stmt)
}

type SetStateArgs struct {
	Winner string `validate:"min=1"`
	Tries  int32  `validate:"min=0"`
	People int32  `validate:"min=0"`
}

func (d *DataService) SetState(args SetStateArgs) (State, error) {
	stmt := t.State.UPDATE(
		t.State.AllColumns,
	).MODEL(model.State{
		Winner: &args.Winner,
		Tries:  &args.Tries,
		People: &args.People,
	}).WHERE(
		t.State.ID.EQ(s.Int(1)),
	).RETURNING(
		t.State.AllColumns.Except(t.State.ID),
	)

	return d.scanState(stmt)
}

type SetWinnerArgs struct {
	Winner string `validate:"min=1"`
}

func (d *DataService) SetWinner(args SetWinnerArgs) (State, error) {
	stmt := t.State.UPDATE(
		t.State.Winner,
	).SET(
		args.Winner,
	).WHERE(
		t.State.ID.EQ(s.Int(1)),
	).RETURNING(t.State.AllColumns)

	return d.scanState(stmt)
}

type SetTriesArgs struct {
	Tries int32 `validate:"min=0"`
}

func (d *DataService) SetTries(args SetTriesArgs) (State, error) {
	stmt := t.State.UPDATE(
		t.State.Tries,
	).SET(
		args.Tries,
	).WHERE(
		t.State.ID.EQ(s.Int(1)),
	).RETURNING(t.State.AllColumns)

	return d.scanState(stmt)
}

type SetPeopleArgs struct {
	People int32 `validate:"min=0"`
}

func (d *DataService) SetPeople(args SetPeopleArgs) (State, error) {
	stmt := t.State.UPDATE(
		t.State.People,
	).SET(
		args.People,
	).WHERE(
		t.State.ID.EQ(s.Int(1)),
	).RETURNING(t.State.AllColumns)

	return d.scanState(stmt)
}

func (d *DataService) IncrementTries() (State, error) {
	stmt := t.State.UPDATE(
		t.State.Tries,
	).SET(
		t.State.Tries.ADD(s.Int(1)),
	).WHERE(
		t.State.ID.EQ(s.Int(1)),
	).RETURNING(
		t.State.AllColumns,
	)

	return d.scanState(stmt)
}

func (d *DataService) IncrementPeople() (State, error) {
	stmt := t.State.UPDATE(
		t.State.People,
	).SET(
		t.State.People.ADD(s.Int(1)),
	).WHERE(
		t.State.ID.EQ(s.Int(1)),
	).RETURNING(t.State.AllColumns)

	return d.scanState(stmt)
}
