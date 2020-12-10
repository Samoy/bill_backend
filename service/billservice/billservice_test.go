/*
 * MIT License
 *
 * Copyright (c) 2020 Samoy
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package billservice

import (
	"github.com/Samoy/bill_backend/config"
	"github.com/Samoy/bill_backend/dao"
	"github.com/Samoy/bill_backend/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var userID uint = 10
var billID uint
var billTypeID uint = 1

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

func TestAddBill(t *testing.T) {
	type args struct {
		bill *models.Bill
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Add Bill",
			args{
				&models.Bill{
					Name:       "Test Bill",
					Amount:     decimal.NewFromFloat(1.00),
					Date:       time.Now(),
					BillTypeID: billTypeID,
					Remark:     "This is a remark that test adding bill ",
					UserID:     userID,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AddBill(tt.args.bill)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddBill() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				billID = tt.args.bill.ID
			}
		})
	}
}

func TestGetBill(t *testing.T) {
	type args struct {
		billID uint
		owner  uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Get Bill",
			args{
				billID,
				userID,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBill(tt.args.billID, tt.args.owner)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, got, "GetBill() test not pass")
		})
	}
}

func TestGetBillList(t *testing.T) {
	type args struct {
		userID    uint
		startTime string
		endTime   string
		page      int
		pageSize  int
		category  uint
		income    string
		sortKey   string
		asc       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Get Bill List",
			args{
				userID,
				"2019-01-10",
				"2021-01-10",
				0,
				10,
				billTypeID,
				"0",
				"updated_at",
				"0",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBillList(tt.args.userID, tt.args.startTime, tt.args.endTime, tt.args.page, tt.args.pageSize, tt.args.category, tt.args.income, tt.args.sortKey, tt.args.asc)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBillList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, got, "GetBillList() test not pass")
		})
	}
}

func TestUpdateBill(t *testing.T) {
	type args struct {
		billID uint
		userID uint
		data   map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Update Bill",
			args{
				billID,
				userID,
				map[string]interface{}{
					"amount": decimal.NewFromFloat(2.00),
					"name":   "After update bill name",
					"remark": "After update bill remark",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if bill, err := UpdateBill(tt.args.billID, tt.args.userID, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UpdateBill() error = %v, wantErr %v", err, tt.wantErr)
				assert.NotEmpty(t, bill, "Update bill service not pass")
			}
		})
	}
}

func TestDeleteBill(t *testing.T) {
	type args struct {
		billID uint
		userID uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Delete Bill",
			args{
				billID,
				userID,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteBill(tt.args.billID, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteBill() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func teardown() {
	_ = dao.CloseDB()
}
