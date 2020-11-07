package router

import (
	"bytes"
	"encoding/json"
	"github.com/Samoy/bill_backend/config"
	"github.com/Samoy/bill_backend/dao"
	v1 "github.com/Samoy/bill_backend/router/api/v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var w *httptest.ResponseRecorder

func setup() {
	config.Setup("../app.ini")
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

func TestInitRouter(t *testing.T) {
	type args struct {
		mode string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Router Init",
			args: args{
				mode: "debug",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试ping
			code, data := testApi(tt.args.mode, "/ping", http.MethodGet, nil)
			assert.Equal(t, http.StatusOK, code)
			// 测试注册
			code, data = testApi(tt.args.mode, "/api/v1/user/register", http.MethodPost, &v1.RegisterBody{
				Username:  "testtest",
				Password:  "testtest",
				Telephone: "13712345678",
				Nickname:  "Testtest",
			})
			assert.Equal(t, http.StatusOK, code)
			assert.NotNil(t, *data, "Register test not pass")
			//测试登录
			code, _ = testApi(tt.args.mode, "/api/v1/user/login", http.MethodPost, v1.LoginBody{
				Username: "testtest",
				Password: "testtest",
			})
			assert.Equal(t, http.StatusOK, code)
			assert.NotNil(t, *data, "Login test not pass")
		})
	}
}

func testApi(mode, path, method string, body interface{}) (code int, data *map[string]interface{}) {
	r := InitRouter(mode)
	w = httptest.NewRecorder()
	reqBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewReader(reqBody))
	r.ServeHTTP(w, req)
	res := w.Body
	code = w.Code
	_ = json.Unmarshal(res.Bytes(), &data)
	_ = w.Result().Body.Close()
	return
}

func teardown() {
	_ = w.Result().Body.Close()
}
