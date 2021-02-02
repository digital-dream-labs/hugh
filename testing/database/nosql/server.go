package nosql

import (
	"errors"
	"log"

	"github.com/ory/dockertest"
	"go.mongodb.org/mongo-driver/mongo"
)

// Server is the struct returned on creation
type Server struct {
	DB        *mongo.Database
	pool      *dockertest.Pool
	resources []*dockertest.Resource
}

// Type defines the type of database you'd like
type Type int

const (
	Mongo Type = iota
)

// Config contains info about the test DB you're trying to get
type Config struct {
	Type     Type
	Username string
	Password string
	Database string
}

// NewServer accepts a config
func NewServer(c *Config) (*Server, error) {
	var (
		s   Server
		err error
	)

	s.pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	switch c.Type {
	case Mongo:
		if err := s.runMongo(c); err != nil {
			return nil, err
		}
		// todo:  add more
	default:
		return nil, errors.New("invaild type")
	}

	return &s, nil
}

// Kill executes clean-up functions for the Server
func (s *Server) Kill() {
	for _, v := range s.resources {
		if err := s.pool.Purge(v); err != nil {
			log.Fatal(err)
		}
	}
}
