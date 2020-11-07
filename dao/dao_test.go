package dao

import (
	"github.com/Samoy/bill_backend/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setup() {
	config.Setup("../app.ini")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestSetup(t *testing.T) {
	type args struct {
		dbType        string
		dbUser        string
		dbPassword    string
		dbHost        string
		dbName        string
		dbTablePrefix string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Setup DB",
			args{
				config.DatabaseConf.Type,
				config.DatabaseConf.User,
				config.DatabaseConf.Password,
				config.DatabaseConf.Host,
				config.DatabaseConf.Name,
				config.DatabaseConf.TablePrefix,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Setup(tt.args.dbType, tt.args.dbUser, tt.args.dbPassword, tt.args.dbHost, tt.args.dbName, tt.args.dbTablePrefix)
			assert.NotNil(t, DB, "Setup database test not pass")
		})
	}
}

func TestCloseDB(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Close DB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestSetup(t)
			err := CloseDB()
			assert.Nil(t, err, "Close database test not pass")
		})
	}
}
