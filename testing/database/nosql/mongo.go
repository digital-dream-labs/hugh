package nosql

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

func (s *Server) runMongo(c *Config) error {
	res, err := s.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "4.4.1-bionic",
		Env: []string{
			fmt.Sprintf("MONGO_INITDB_ROOT_USERNAME=%s", c.Username),
			fmt.Sprintf("MONGO_INITDB_ROOT_PASSWORD=%s", c.Password),
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

	hostname, port := getHostInfo(res, docker.Port("27017/tcp"))
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%v", c.Username, c.Password, hostname, port)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	s.DB = client.Database(c.Database)
	s.resources = append(s.resources, res)

	return nil
}
