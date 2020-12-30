package sql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	_ "github.com/go-sql-driver/mysql" // mysql
)

func (s *Server) runMySQL(c *Config) error {
	res, err := s.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			fmt.Sprintf("MYSQL_DATABASE=%s", c.Database),
			fmt.Sprintf("MYSQL_USER=%s", c.Username),
			fmt.Sprintf("MYSQL_PASSWORD=%s", c.Password),
			fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", c.Password),
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

	hostname, port := getHostInfo(res, docker.Port("3306/tcp"))

	var db *sql.DB
	// exponential backoff-retry, because the application in the container
	// might not be ready to accept connections yet
	if err = s.pool.Retry(func() error {
		var err2 error
		db, err2 = sql.Open(
			"mysql",
			fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.Username, c.Password, hostname, port, c.Database),
		)
		if err2 != nil {
			return err2
		}

		return db.Ping()
	}); err != nil {
		return fmt.Errorf("could not connect to docker mysql: %s", err)
	}

	s.DB = db
	s.resources = append(s.resources, res)

	return nil
}
