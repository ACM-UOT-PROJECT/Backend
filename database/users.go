package database

import (
	"backend/database/gen/model"
	t "backend/database/gen/table"

	s "github.com/go-jet/jet/v2/sqlite"
)

type User struct {
	Id       int32
	Username string
	Token    string
}

func (d *DataService) scanUser(stmt s.Statement) (User, error) {
	dest := model.User{}

	err := stmt.Query(d.db, &dest)
	if err != nil {
		d.logger.Error(err.Error())
		d.logger.Info(stmt.DebugSql())
		return User{}, err
	}

	res := User{
		Id:       *dest.ID,
		Username: dest.UserName,
		Token:    dest.Token,
	}

	return res, nil
}

type RegisterUserArgs struct {
	UserName string
	Token    string `validate:"jwt"`
}

func (d *DataService) RegisterUser(args RegisterUserArgs) (User, error) {
	stmt := t.User.
		INSERT(
			t.User.UserName,
			t.User.Token,
		).
		VALUES(
			args.UserName,
			args.Token,
		).
		RETURNING(
			t.User.ID,
		).
		ON_CONFLICT(
			t.User.Token,
		).
		DO_UPDATE(
			s.SET(
				t.User.ID.SET(t.User.ID),
				t.User.UserName.SET(s.String(args.UserName)),
			),
		).
		RETURNING(
			t.User.AllColumns,
		)

	return d.scanUser(stmt)
}

type LoginUserArgs struct {
	Token string
}

func (d *DataService) LoginUser(args LoginUserArgs) (User, error) {
	stmt := s.
		SELECT(
			t.User.ID,
			t.User.UserName,
		).
		FROM(
			t.User,
		).
		WHERE(
			t.User.Token.EQ(s.String(args.Token)),
		)
	return d.scanUser(stmt)
}
