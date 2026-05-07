package config

import "fmt"

type Postgres struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
}

func (p Postgres) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		p.Host, p.Port, p.DB, p.User, p.Password,
	)
}

func loadPostgres() (Postgres, error) {
	host, err := getMustEnv("POSTGRES_HOST")
	if err != nil {
		return Postgres{}, err
	}
	port, err := getMustEnv("POSTGRES_PORT")
	if err != nil {
		return Postgres{}, err
	}
	db, err := getMustEnv("POSTGRES_DB")
	if err != nil {
		return Postgres{}, err
	}
	user, err := getMustEnv("POSTGRES_USER")
	if err != nil {
		return Postgres{}, err
	}
	password, err := getMustEnv("POSTGRES_PASSWORD")
	if err != nil {
		return Postgres{}, err
	}

	return Postgres{
		Host:     host,
		Port:     port,
		DB:       db,
		User:     user,
		Password: password,
	}, nil
}
