package database

import (
	"github.com/CollCaz/UniSite/database/gen/unicontentdb/public/model"
	t "github.com/CollCaz/UniSite/database/gen/unicontentdb/public/table"
	s "github.com/go-jet/jet/v2/postgres"
)

type User struct {
	Id       int32
	Username string
	Email    string
	Roles    []string
}

type joinedApiUserModel struct {
	model.APIUser
	Roles []model.Role
}

func (d *DataService) scanUser(stmt s.Statement) (User, error) {
	dest := joinedApiUserModel{}

	err := stmt.Query(d.db, &dest)
	if err != nil {
		d.logger.Error(err.Error())
		d.logger.Info(stmt.DebugSql())
		return User{}, err
	}

	roles := []string{}
	for _, role := range dest.Roles {
		roles = append(roles, role.Name)
	}

	res := User{
		Id:       dest.ID,
		Username: dest.UserName,
		Email:    dest.Email,
		Roles:    roles,
	}

	return res, nil
}

type RegisterUserArgs struct {
	Username        string
	Email           string
	Password        string
	PasswordConfirm string
}

func (d *DataService) RegisterUser(args RegisterUserArgs) (User, error) {
	stmt := t.APIUser.
		INSERT(
			t.APIUser.UserName,
			t.APIUser.Email,
			t.APIUser.PassHash,
		).
		VALUES(
			args.Username,
			args.Email,
			args.Password,
		).
		RETURNING(
			t.APIUser.ID,
		)

	return d.scanUser(stmt)
}

type LoginUserArgs struct {
	Email    string
	Password string
}

func (d *DataService) LoginUser(args LoginUserArgs) (User, error) {
	stmt := s.
		SELECT(
			t.APIUser.ID,
			t.APIUser.UserName,
			t.APIUser.Email,
			t.Role.Name,
		).
		FROM(
			t.APIUser.INNER_JOIN(t.UserRoles, t.UserRoles.UserID.EQ(t.APIUser.ID)).
				INNER_JOIN(t.Role, t.UserRoles.RoleID.EQ(t.Role.ID)),
		).
		WHERE(
			s.AND(
				t.APIUser.Email.EQ(s.String(args.Email)),
				t.APIUser.PassHash.EQ(s.String(args.Password)),
			),
		)

	return d.scanUser(stmt)
}
