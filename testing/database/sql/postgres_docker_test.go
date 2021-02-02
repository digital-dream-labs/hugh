package sql

import (
	"testing"

	"github.com/ory/dockertest"
)

func TestServer_runPostgres(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	type args struct {
		c *Config
	}
	tests := []struct {
		name    string
		server  Server
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			server: Server{
				pool: pool,
			},
			args: args{
				c: &Config{
					Username: "test",
					Password: "test",
					Database: "test",
					Type:     Postgres,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.server.runPostgres(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Server.runMySQL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.server.DB == nil {
				t.Errorf("database is nil")
			}
		})
	}
}
