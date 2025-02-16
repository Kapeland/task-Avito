package auth

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
	"github.com/gofrs/uuid"
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

func TestRepo_CreateUserSecret(t *testing.T) {
	type fields struct {
		db db.DBops
	}
	type args struct {
		ctx        context.Context
		userSecret *structs.UserSecret
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

	tmpKey, _ := genInt(64)
	tmpUUID, _ := uuid.NewV4()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Unexisting user",
			fields: fields{dbStor.DB},
			args: args{
				ctx: ctx,
				userSecret: &structs.UserSecret{
					Login:     "user22user22",
					Secret:    "741852963",
					SessionID: "123123",
				},
			},
			wantErr: true,
		},
		{
			name:   "Existing user, not duplicate",
			fields: fields{dbStor.DB},
			args: args{
				ctx: ctx,
				userSecret: &structs.UserSecret{
					Login:     "user1user2",
					Secret:    tmpKey,
					SessionID: tmpUUID.String(),
				},
			},
			wantErr: false,
		},
		{
			name:   "Existing user, but duplicate",
			fields: fields{dbStor.DB},
			args: args{
				ctx: ctx,
				userSecret: &structs.UserSecret{
					Login:     "user1user2",
					Secret:    tmpKey,
					SessionID: tmpUUID.String(),
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
			if err := r.CreateUserSecret(tt.args.ctx, tt.args.userSecret); (err != nil) != tt.wantErr {
				t.Errorf("CreateUserSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_GetSecretByLoginAndSession(t *testing.T) {
	type fields struct {
		db db.DBops
	}
	type args struct {
		ctx    context.Context
		lgnSsn structs.UserSecret
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
		want    *structs.UserSecret
		wantErr bool
	}{
		{
			name:   "Existing login and session",
			fields: fields{dbStor.DB},
			args: args{
				ctx: ctx,
				lgnSsn: structs.UserSecret{
					Login:     "user1user1",
					SessionID: "d7da5bab-6992-4b5b-8de5-3d4d54359747",
				},
			},
			want: &structs.UserSecret{
				Login:     "user1user1",
				Secret:    "0ryXaQIkiEliUegl30f76L3w3e7pvf0eqaqEY4QECQkN5bXU1NdLEvbWe1uDMM6s",
				SessionID: "d7da5bab-6992-4b5b-8de5-3d4d54359747",
			},
			wantErr: false,
		},
		{
			name:   "Existing login, but not session",
			fields: fields{dbStor.DB},
			args: args{
				ctx: ctx,
				lgnSsn: structs.UserSecret{
					Login:     "user1user1",
					SessionID: "d7da5bab-6992-4b5b-8de5",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "Existing session, but not login",
			fields: fields{dbStor.DB},
			args: args{
				ctx: ctx,
				lgnSsn: structs.UserSecret{
					Login:     "user1user11111",
					SessionID: "d7da5bab-6992-4b5b-8de5-3d4d54359747",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "Unexisting both",
			fields: fields{dbStor.DB},
			args: args{
				ctx: ctx,
				lgnSsn: structs.UserSecret{
					Login:     "user1user111111",
					SessionID: "d7da5bab-6992-4b5b-8de5",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				db: tt.fields.db,
			}
			got, err := r.GetSecretByLoginAndSession(tt.args.ctx, tt.args.lgnSsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecretByLoginAndSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSecretByLoginAndSession() got = %v, want %v", got, tt.want)
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
