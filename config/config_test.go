package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSetup(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Config setup",
			args: args{
				path: "../app.ini",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Setup(tt.args.path)
			assert.Equal(t, "debug", AppConf.RunMode, "App RunMode not correct")
			assert.Equal(t, "5ac8853dca0ef68e1bcc86ffab2a0a27", AppConf.JwtSecret, "App JwtSecret not correct")
			assert.Equal(t, 5000, ServerConf.HTTPPort, "ServerConf HTTPPort not correct")
			assert.Equal(t, 60*time.Second, ServerConf.WriteTimeout, "ServerConf WriteTimeout not correct")
			assert.Equal(t, 60*time.Second, ServerConf.ReadTimeout, "ServerConf ReadTimeout not correct")
			assert.Equal(t, "mysql", DatabaseConf.Type, "DatabaseConf Type not correct")
			assert.Equal(t, "root", DatabaseConf.User, "DatabaseConf User not correct")
			assert.Equal(t, "asd1234.", DatabaseConf.Password, "DatabaseConf Password not correct")
			assert.Equal(t, "127.0.0.1:3306", DatabaseConf.Host, "DatabaseConf Host not correct")
			assert.Equal(t, "bill", DatabaseConf.Name, "DatabaseConf Name not correct")
			assert.Equal(t, "t_", DatabaseConf.TablePrefix, "DatabaseConf TablePrefix not correct")
		})
	}
}
