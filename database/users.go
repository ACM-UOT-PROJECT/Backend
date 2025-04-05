package database

import (
	"backend/database/gen/model"
	t "backend/database/gen/table"

	s "github.com/go-jet/jet/v2/sqlite"
)

type User struct {
	Id       int32
	Username string
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
	}

	return res, nil
}

type RegisterUserArgs struct {
	Username string
	Token    string `validate:"jwt"`
}

func (d *DataService) RegisterUser(args RegisterUserArgs) (User, error) {
	stmt := t.User.
		INSERT(
			t.User.UserName,
			t.User.IPAddress,
		).
		VALUES(
			args.Username,
			args.Token,
		).
		RETURNING(
			t.User.ID,
		)

	return d.scanUser(stmt)
}

type LoginUserArgs struct {
	IpAddress string `validate:"ip"`
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
			t.User.IPAddress.EQ(s.String(args.IpAddress)),
		)
	return d.scanUser(stmt)
}
