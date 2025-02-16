package users

import (
	"context"
	"crypto/rand"
	"math/big"
	"reflect"
	"testing"

	"github.com/Kapeland/task-Avito/internal/models/structs"
	"github.com/Kapeland/task-Avito/internal/storage"
	"github.com/Kapeland/task-Avito/internal/storage/db"
	"github.com/Kapeland/task-Avito/internal/utils/config"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
)

func TestNew(t *testing.T) {
	type args struct {
		db db.DBops
	}
	ctx := context.Background()
	if err := config.ReadLocalConfigYAML(); err != nil {
		t.Error(err)
	}
	cfg := config.GetConfig()
	logger.CreateLogger(&cfg)
	dbStor, err := storage.NewPostgresStorage(ctx)
	if err != nil {
		t.Error("NewPostgresStorage: " + err.Error())
	}
	tests := []struct {
		name string
		args args
		want *Repo
	}{
		{
			name: "Init DB",
			args: args{db: dbStor.DB},
			want: &Repo{db: dbStor.DB},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_BuyItem(t *testing.T) {
	type fields struct {
		db db.DBops
	}
	type args struct {
		ctx   context.Context
		item  string
		login string
	}
	ctx := context.Background()
	if err := config.ReadLocalConfigYAML(); err != nil {
		t.Error(err)
	}
	cfg := config.GetConfig()
	logger.CreateLogger(&cfg)
	dbStor, err := storage.NewPostgresStorage(ctx)
	if err != nil {
		t.Error("NewPostgresStorage: " + err.Error())
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Buy existing item, money enough",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx:   ctx,
				item:  "pen",
				login: "user1user1",
			},
			wantErr: false,
		},
		{
			name:   "Buy existing item, money not enough",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx:   ctx,
				item:  "pink-hoody",
				login: "user1user2",
			},
			wantErr: true,
		},
		{
			name:   "Buy not existing item, money enough",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx:   ctx,
				item:  "red-hoody",
				login: "user1user1",
			},
			wantErr: true,
		},
		{
			name:   "User not exists",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx:   ctx,
				item:  "pen",
				login: "user1user1337",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				db: tt.fields.db,
			}
			if err := r.BuyItemDB(tt.args.ctx, tt.args.item, tt.args.login); (err != nil) != tt.wantErr {
				t.Errorf("BuyItemDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_CreateUser(t *testing.T) {
	type fields struct {
		db db.DBops
	}
	type args struct {
		ctx  context.Context
		info structs.RegisterUserInfo
	}

	ctx := context.Background()
	if err := config.ReadLocalConfigYAML(); err != nil {
		t.Error(err)
	}
	cfg := config.GetConfig()
	logger.CreateLogger(&cfg)
	dbStor, err := storage.NewPostgresStorage(ctx)
	if err != nil {
		t.Error("NewPostgresStorage: " + err.Error())
	}

	tmpNumb, err := genInt(5)
	if err != nil {
		t.Error("genInt: " + err.Error())
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "New user",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx: ctx,
				info: structs.RegisterUserInfo{
					Login: "user1user" + tmpNumb,
					Pswd:  "Lhjxb[eq" + tmpNumb,
				},
			},
			wantErr: false,
		},
		{
			name:   "Existing login and pass",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx: ctx,
				info: structs.RegisterUserInfo{
					Login: "user1user1",
					Pswd:  "Lhjxb[eq1",
				},
			},
			wantErr: true,
		},
		{
			name:   "Existing login",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx: ctx,
				info: structs.RegisterUserInfo{
					Login: "user1user1",
					Pswd:  "Lhjxb[eq123",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				db: tt.fields.db,
			}

			if err := r.CreateUserDB(tt.args.ctx, tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("CreateUserDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_GetInfo(t *testing.T) {
	//TODO: Добавить ещё тестов
	type fields struct {
		db db.DBops
	}
	type args struct {
		ctx   context.Context
		login string
	}

	ctx := context.Background()
	if err := config.ReadLocalConfigYAML(); err != nil {
		t.Error(err)
	}
	cfg := config.GetConfig()
	logger.CreateLogger(&cfg)
	dbStor, err := storage.NewPostgresStorage(ctx)
	if err != nil {
		t.Error("NewPostgresStorage: " + err.Error())
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *structs.AccInfo
		wantErr bool
	}{
		{
			name:   "Get info of unexisting user",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx:   ctx,
				login: "user1user21313",
			},
			want:    &structs.AccInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				db: tt.fields.db,
			}
			got, err := r.GetInfoDB(tt.args.ctx, tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfoDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInfoDB() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_SendCoinTo(t *testing.T) {
	type fields struct {
		db db.DBops
	}
	type args struct {
		ctx       context.Context
		operation structs.SendCoinInfo
	}

	ctx := context.Background()
	if err := config.ReadLocalConfigYAML(); err != nil {
		t.Error(err)
	}
	cfg := config.GetConfig()
	logger.CreateLogger(&cfg)
	dbStor, err := storage.NewPostgresStorage(ctx)
	if err != nil {
		t.Error("NewPostgresStorage: " + err.Error())
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Send to existing, money enough (1)",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx: ctx,
				operation: structs.SendCoinInfo{
					From:   "user1user1",
					To:     "user1user2",
					Amount: 5,
				},
			},
			wantErr: false,
		},
		{
			name:   "Send to existing, money enough (2)",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx: ctx,
				operation: structs.SendCoinInfo{
					From:   "user1user2",
					To:     "user1user1",
					Amount: 5,
				},
			},
			wantErr: false,
		},
		{
			name:   "Send to not existing, money enough",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx: ctx,
				operation: structs.SendCoinInfo{
					From:   "user1user1",
					To:     "user1user2321",
					Amount: 5,
				},
			},
			wantErr: true,
		},
		{
			name:   "Send to not existing, money not enough",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx: ctx,
				operation: structs.SendCoinInfo{
					From:   "user1user1",
					To:     "user1user2321",
					Amount: 50000,
				},
			},
			wantErr: true,
		},
		{
			name:   "Send from not existing, money enough",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx: ctx,
				operation: structs.SendCoinInfo{
					From:   "user1user1123112",
					To:     "user1user2",
					Amount: 5,
				},
			},
			wantErr: true,
		},
		{
			name:   "Send from not existing, money not enough",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx: ctx,
				operation: structs.SendCoinInfo{
					From:   "user1user1123112",
					To:     "user1user2",
					Amount: 500000,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				db: tt.fields.db,
			}
			if err := r.SendCoinDB(tt.args.ctx, tt.args.operation); (err != nil) != tt.wantErr {
				t.Errorf("SendCoinDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_VerifyPassword(t *testing.T) {
	type fields struct {
		db db.DBops
	}
	type args struct {
		ctx  context.Context
		info structs.AuthUserInfo
	}

	ctx := context.Background()
	if err := config.ReadLocalConfigYAML(); err != nil {
		t.Error(err)
	}
	cfg := config.GetConfig()
	logger.CreateLogger(&cfg)
	dbStor, err := storage.NewPostgresStorage(ctx)
	if err != nil {
		t.Error("NewPostgresStorage: " + err.Error())
	}

	tmpNumb, err := genInt(5)
	if err != nil {
		t.Error("genInt: " + err.Error())
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:   "Matching pswd",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx: ctx,
				info: structs.AuthUserInfo{
					Login: "user1user1",
					Pswd:  "Lhjxb[eq1",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "Not matching pswd",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx: ctx,
				info: structs.AuthUserInfo{
					Login: "user1user1",
					Pswd:  "Lhjxb[eq12",
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name:   "Unexisting user",
			fields: fields{db: dbStor.DB},
			args: args{
				ctx: ctx,
				info: structs.AuthUserInfo{
					Login: "user1user" + tmpNumb,
					Pswd:  "Lhjxb[eq" + tmpNumb,
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				db: tt.fields.db,
			}
			got, err := r.VerifyPasswordDB(tt.args.ctx, tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPasswordDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyPasswordDB() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func genInt(length int) (string, error) {
	result := ""
	for {
		if len(result) >= length {
			return result, nil
		}
		num, err := rand.Int(rand.Reader, big.NewInt(int64(127)))
		if err != nil {
			return "", err
		}
		n := num.Int64()
		if n >= 48 && n <= 57 {
			result += string(rune(n))
		}
	}
}
