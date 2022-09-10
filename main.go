package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/joesantosio/example-order-book/config"
	"github.com/joesantosio/example-order-book/entity"
	"github.com/joesantosio/example-order-book/infrastructure/sqlite"
	"github.com/joesantosio/example-order-book/interfaces/cli"
)

func initRepos(dbPath string) (entity.Repositories, error) {
	if dbPath == "" {
		return nil, errors.New("DB_PATH not provided")
	}

	db, err := sqlite.Connect(dbPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error initializing db: %v", err))
	}

	repos, err := sqlite.InitRepos(db)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error initializing repos: %v", err))
	}

	return repos, nil
}

func main() {
	config, err := config.Get(os.Getenv)
	if err != nil {
		panic(err)
	}

	repos, err := initRepos(config.DBPath)
	if err != nil {
		panic(err)
	}
	defer repos.Close()

	err = cli.Init(os.Args, repos)
	if err != nil {
		panic(err)
	}
}
