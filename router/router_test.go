package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Samoy/bill_backend/config"
	"github.com/Samoy/bill_backend/dao"
	v1 "github.com/Samoy/bill_backend/router/api/v1"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var (
	w          *httptest.ResponseRecorder
	token      string
	billTypeID uint
	billID     uint
)

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
			code, data := testNormalApi(tt.args.mode, "/ping", http.MethodGet, nil)
			assert.Equal(t, http.StatusOK, code)
			// 测试注册
			code, data = testNormalApi(tt.args.mode, "/api/v1/user/register", http.MethodPost, &v1.RegisterBody{
				Username:  "testrouter",
				Password:  "testrouter",
				Telephone: "13712345678",
				Nickname:  "Testtest",
			})
			assert.Equal(t, http.StatusOK, code)
			assert.NotNil(t, *data, "Register test not pass")
			//测试登录
			code, data = testNormalApi(tt.args.mode, "/api/v1/user/login", http.MethodPost, v1.LoginBody{
				Username: "testrouter",
				Password: "testrouter",
			})
			assert.Equal(t, http.StatusOK, code, "login test not pass")
			assert.NotNil(t, *data, "Login test not pass")
			token = (*data)["data"].(map[string]interface{})["token"].(string)
			//测试添加账单类型
			err, code, data := testFileApi(
				tt.args.mode,
				"/api/v1/bill_type",
				http.MethodPost,
				"image",
				"../testdata/bill_type_test.png",
				map[string]interface{}{
					"name": "Test Add BillType",
				},
			)
			assert.Nil(t, err, "Add bill type test not pass")
			assert.Equal(t, http.StatusOK, code, "Add bill type test not pass")
			assert.NotEmpty(t, *data, "Add bill type test not pass")
			billTypeID = uint((*data)["data"].(map[string]interface{})["id"].(float64))
			//测试获取账单类型
			code, data = testNormalApi(tt.args.mode, fmt.Sprintf("/api/v1/bill_type?bill_type_id=%d", billTypeID), http.MethodGet, nil)
			assert.Equal(t, http.StatusOK, code, "Get bill type test not pass")
			assert.NotEmpty(t, *data, "Get bill type test not pass")
			//测试更新账单类型
			err, code, data = testFileApi(
				tt.args.mode,
				"/api/v1/bill_type",
				http.MethodPut,
				"image",
				"../testdata/update_bill_type_test.png",
				map[string]interface{}{
					"name":         "Test Update BillType",
					"bill_type_id": billTypeID,
				},
			)
			assert.Nil(t, err, "Update bill type test not pass")
			assert.Equal(t, http.StatusOK, code, "Update bill type test not pass")
			assert.NotEmpty(t, *data, "Update bill type test not pass")
			//测试获取账单类型列表
			code, data = testNormalApi(tt.args.mode, "/api/v1/bill_type_list", http.MethodGet, nil)
			assert.Equal(t, http.StatusOK, code, "Get bill type list test not pass")
			assert.NotEmpty(t, *data, "Get bill type list test not pass")
			//测试删除账单类型
			code, data = testNormalApi(tt.args.mode, "/api/v1/bill_type", http.MethodDelete, map[string]interface{}{
				"bill_type_id": billTypeID,
			})
			assert.Equal(t, http.StatusOK, code, "Delete bill type test not pass")
			//测试添加账单
			code, data = testNormalApi(tt.args.mode, "/api/v1/bill", http.MethodPost, v1.BillBody{
				Name:   "Test Add Bill Router",
				Amount: decimal.NewFromFloat(1.00),
				Income: false,
				TypeID: billTypeID,
				Remark: "This is a  add bill router test",
			})
			assert.Equal(t, http.StatusOK, code, "Add bill test not pass")
			assert.NotEmpty(t, *data, "Add bill test not pass")
			billID = uint((*data)["data"].(map[string]interface{})["id"].(float64))
			//测试获取账单
			testNormalApi(tt.args.mode, fmt.Sprintf("/api/v1/bill?bill_id=%d", billID), http.MethodGet, nil)
			//测试更新账单
			code, data = testNormalApi(tt.args.mode, "/api/v1/bill", http.MethodPut, v1.UpdateBillBody{
				BillID: billID,
				BillBody: v1.BillBody{
					Name:   "Test Update Bill Router",
					Amount: decimal.NewFromFloat(2.00),
					Income: false,
					TypeID: billTypeID,
					Remark: "This is a  update bill router test",
				},
			})
			assert.Equal(t, http.StatusOK, code, "Update bill test not pass")
			assert.NotEmpty(t, *data, "Update bill test not pass")
			//测试获取账单列表
			code, data = testNormalApi(tt.args.mode,
				fmt.Sprintf(
					"/api/v1/bill_list?start_time=%s&end_time=%s&page=%d&page_size=%d&type=%d&income=%d&sort_key=%s&asc=%d",
					"2006-01-01 00:00:00", "2030-12-31 23:59:59", 1, 10, billTypeID, 0, "amount", 0,
				),
				http.MethodGet, nil)
			assert.Equal(t, http.StatusOK, code, "Get bill list test not pass")
			assert.NotEmpty(t, *data, "Get bill list test not pass")
			//测试删除账单
			code, data = testNormalApi(tt.args.mode, "/api/v1/bill", http.MethodDelete, v1.DeleteBillBody{
				BillID: billID,
			})
			assert.Equal(t, http.StatusOK, code, "Delete bill test not pass")
		})
	}
}

func testNormalApi(mode, path, method string, body interface{}) (code int, data *map[string]interface{}) {
	r := InitRouter(mode)
	w = httptest.NewRecorder()
	reqBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Add("token", token)
	r.ServeHTTP(w, req)
	res := w.Body
	code = w.Code
	_ = json.Unmarshal(res.Bytes(), &data)
	_ = w.Result().Body.Close()
	return
}

func testFileApi(mode, path, method, fileKey, filePath string, params map[string]interface{}) (err error, code int, data *map[string]interface{}) {
	r := InitRouter(mode)
	w = httptest.NewRecorder()
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileKey, filepath.Base(filePath))
	if err != nil {
		return
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, fmt.Sprintf("%v", val))
	}
	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("token", token)
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
