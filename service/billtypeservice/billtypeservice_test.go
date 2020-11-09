package billtypeservice

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

var billTypeID uint
var ownerID uint = 1

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestAddBillType(t *testing.T) {
	type args struct {
		billType *models.BillType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Add Bill Type",
			args{
				&models.BillType{
					Name:  "Test Bill Type",
					Image: "/testpath",
					Owner: ownerID,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddBillType(tt.args.billType); (err != nil) != tt.wantErr {
				t.Errorf("AddBillType() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				billTypeID = tt.args.billType.ID
			}
		})
	}
}

func TestGetBillType(t *testing.T) {
	type args struct {
		billTypeID uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Get Bill Type",
			args{
				billTypeID,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBillType(tt.args.billTypeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBillType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, got, "GetBillType() test not pass")
		})
	}
}

func TestUpdateBillType(t *testing.T) {
	type args struct {
		billTypeID uint
		billType   *models.BillType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Update Bill Type",
			args{
				billTypeID,
				&models.BillType{
					Name:  "After update Bill Type",
					Image: "/after_update_bill_type_path",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateBillType(tt.args.billTypeID, tt.args.billType); (err != nil) != tt.wantErr {
				t.Errorf("UpdateBillType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetBillTypeList(t *testing.T) {
	type args struct {
		ownerID uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Get Bill Type List",
			args{
				ownerID,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBillTypeList(tt.args.ownerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBillTypeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, got, "GetBillTypeList() test not pass")
		})
	}
}

func TestDeleteBillType(t *testing.T) {
	type args struct {
		billTypeID uint
		ownerID    uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Delete Bill Type",
			args{
				billTypeID,
				ownerID,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteBillType(tt.args.billTypeID, tt.args.ownerID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteBillType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func teardown() {
	_ = dao.CloseDB()
}
