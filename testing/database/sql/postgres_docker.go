package sql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // import postgres
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

func (s *Server) runPostgres(c *Config) error {
	res, err := s.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "9.6-alpine",
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", c.Database),
			fmt.Sprintf("POSTGRES_USER=%s", c.Username),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", c.Password),
		},
		NetworkID: "bridge",
	})
	if err != nil {
		log.Fatalf("Could not run database resource: %s", err)
	}

	// set the resource to expire in 15min in-case the tests die
	if err := res.Expire(900); err != nil {
		log.Fatal(err)
	}

	hostname, port := getHostInfo(res, docker.Port("5432/tcp"))

	var db *sql.DB
	// exponential backoff-retry, because the application in the container
	// might not be ready to accept connections yet
	if err = s.pool.Retry(func() error {
		var err2 error
		db, err2 = sql.Open(
			"postgres",
			fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Username, c.Password, hostname, port, c.Database),
		)
		if err2 != nil {
			return err2
		}

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker postgres: %s", err)
	}

	s.DB = db
	s.resources = append(s.resources, res)

	return nil
}
