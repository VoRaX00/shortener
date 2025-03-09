package main

import (
	"database/sql"
	"fmt"
	"github.com/VoRaX00/shortener/internal/config"
	_ "github.com/VoRaX00/shortener/migrations"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type PgConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	Username       string `yaml:"user"`
	Password       string `yaml:"password"`
	Database       string `yaml:"dbname"`
	SSLMode        string `yaml:"ssl_mode"`
	IsDrop         bool   `yaml:"is_drop"`
	MigrationsPath string `yaml:"migrations_path"`
}

const postgresConfigPath = "./cmd/migrator/postgres.yml"

func main() {
	//if err := godotenv.Load(); err != nil {
	//	panic(err)
	//}

	cfg := config.MustConfig[PgConfig](postgresConfigPath)
	//cfg.Password = os.Getenv("POSTGRES_PASSWORD")

	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	if cfg.IsDrop {
		if err = goose.DownTo(db, cfg.MigrationsPath, 0); err != nil {
			panic(err)
		}
	}

	if err = goose.Up(db, cfg.MigrationsPath); err != nil {
		panic(err)
	}
}
