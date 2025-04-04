package database

import (
	"database/sql"
	"log/slog"
	"os"
)

type DataService struct {
	db     *sql.DB
	logger *slog.Logger
}

type NewDataServiceArgs struct {
	Db     *sql.DB
	Logger *slog.Logger
}

func NewDataService(args NewDataServiceArgs) *DataService {
	if args.Logger == nil {
		args.Logger = slog.Default()
	}
	if args.Db == nil {
		args.Logger.Error("Did not provide valid db connection.")
		os.Exit(1)
	}
	return &DataService{
		db:     args.Db,
		logger: args.Logger,
	}
}
