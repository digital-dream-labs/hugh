package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// New initializes a new database connection
func New(o ...Option) (*mongo.Database, error) {
	cfg := opts{}

	for _, i := range o {
		i(&cfg)
	}

	if cfg.direct && cfg.cluster {
		return nil, errors.New("invalid configuration -- either direct or cluster, not both")
	}

	if cfg.cluster && cfg.port != 0 {
		return nil, errors.New("invalid configuration -- either port or cluster, not both")
	}

	var prefix string
	if cfg.cluster {
		prefix = "mongodb+srv"
	} else {
		prefix = "mongodb"
	}

	uri := fmt.Sprintf(
		"%s://%s:%s@%s",
		prefix,
		cfg.username,
		cfg.password,
		cfg.host,
	)

	if cfg.port != 0 {
		uri += fmt.Sprintf(":%v", cfg.port)
	}

	if cfg.authdbname != "" {
		uri += fmt.Sprintf("/%s", cfg.authdbname)

	}

	switch {
	case cfg.direct:
		uri += "/?connect=direct"
	case cfg.retryWrites && cfg.writeConcern == "":
		uri += "/?retryWrites=true"
	case !cfg.retryWrites && cfg.writeConcern != "":
		uri += "/?w=" + cfg.writeConcern
	case cfg.retryWrites && cfg.writeConcern != "":
		uri += `/?retryWrites=true&w=` + cfg.writeConcern
	}

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("uri is %s: %v", uri, err)
	}

	db := client.Database(cfg.name)
	return db, nil
}
