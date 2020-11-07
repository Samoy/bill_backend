package userservice

import (
	"github.com/Samoy/bill_backend/config"
	"github.com/Samoy/bill_backend/dao"
	"github.com/Samoy/bill_backend/models"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setup() {
	config.Setup("../../app.ini")
	dao.Setup(config.DatabaseConf.Type,
		config.DatabaseConf.User,
		config.DatabaseConf.Password,
		config.DatabaseConf.Host,
		config.DatabaseConf.Name,
		config.DatabaseConf.TablePrefix,
	)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestRegister(t *testing.T) {
	type args struct {
		u *models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Register Service",
			args: args{
				u: &models.User{
					Username:  "testservice",
					Password:  "testservice",
					Telephone: "18212345678",
					Nickname:  "TestUser",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Register(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name     string
		args     args
		wantUser models.User
		wantErr  bool
	}{
		{
			"Login Service",
			args{
				"testservice",
				"testservice",
			},
			models.User{
				Username:  "testservice",
				Password:  "testservice",
				Telephone: "18212345678",
				Nickname:  "TestUser",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := Login(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, gotUser, "Login test not pass")
		})
	}
}

func Test_existUser(t *testing.T) {
	type args struct {
		username  string
		telephone string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Exist User Function",
			args{
				"exitusertest",
				"13112345678",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := existUser(tt.args.username, tt.args.telephone); (err != nil) != tt.wantErr {
				t.Errorf("existUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func teardown() {
	_ = dao.CloseDB()
}
