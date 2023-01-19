package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/mkaiho/go-auth-api/entity"
	portmocks "github.com/mkaiho/go-auth-api/mocks/usecase/port"
	"github.com/mkaiho/go-auth-api/usecase/port"
	"github.com/stretchr/testify/assert"
)

func Test_userInteractor_GetUser(t *testing.T) {
	type mockUserGetReturn struct {
		user *entity.User
		err  error
	}
	type mockReturn struct {
		userGet *mockUserGetReturn
	}
	type args struct {
		ctx   context.Context
		input GetUserInput
	}
	tests := []struct {
		name       string
		args       args
		mockReturn mockReturn
		want       *entity.User
		wantErr    bool
	}{
		{
			name: "return user",
			args: args{
				ctx: context.Background(),
				input: GetUserInput{
					ID: "test_user_id_001",
				},
			},
			mockReturn: mockReturn{
				userGet: &mockUserGetReturn{
					user: &entity.User{
						ID:    "test_user_id_001",
						Name:  "test_user_001",
						Email: "test_001@example.com",
					},
					err: nil,
				},
			},
			want: &entity.User{
				ID:    "test_user_id_001",
				Name:  "test_user_001",
				Email: "test_001@example.com",
			},
			wantErr: false,
		},
		{
			name: "return error when user creation failed",
			args: args{
				ctx: context.Background(),
				input: GetUserInput{
					ID: "test_user_id_001",
				},
			},
			mockReturn: mockReturn{
				userGet: &mockUserGetReturn{
					user: &entity.User{
						ID:    "test_user_id_001",
						Name:  "test_user_001",
						Email: "test_001@example.com",
					},
					err: errors.New("failed to create user"),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := portmocks.NewUserGateway(t)
			if tt.mockReturn.userGet != nil {
				users.
					On("Get", tt.args.ctx, tt.args.input.ID).
					Return(
						tt.mockReturn.userGet.user,
						tt.mockReturn.userGet.err,
					).
					Times(1)
			}
			it := &userInteractor{
				users: users,
			}
			got, err := it.GetUser(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got, "userInteractor.GetUser() = %v, want %v", got, tt.want)
		})
	}
}

func Test_userInteractor_CreateUser(t *testing.T) {
	type mockUserCreateReturn struct {
		user *entity.User
		err  error
	}
	type mockUserCredsCreateReturn struct {
		creds *entity.UserCredential
		err   error
	}
	type mockReturn struct {
		userCreate      *mockUserCreateReturn
		userCredsCreate *mockUserCredsCreateReturn
	}

	type args struct {
		ctx   context.Context
		input CreateUserInput
	}
	tests := []struct {
		name       string
		args       args
		mockReturn mockReturn
		want       *entity.User
		wantErr    bool
	}{
		{
			name: "return created user",
			args: args{
				ctx: context.Background(),
				input: CreateUserInput{
					Name:     "test_user_001",
					Email:    "test_001@example.com",
					Password: "test_pass",
				},
			},
			mockReturn: mockReturn{
				userCreate: &mockUserCreateReturn{
					user: &entity.User{
						ID:    "test_user_id_001",
						Name:  "test_user_001",
						Email: "test_001@example.com",
					},
					err: nil,
				},
				userCredsCreate: &mockUserCredsCreateReturn{
					creds: &entity.UserCredential{
						ID:       "test_user_creds_001",
						UserID:   "test_user_id_001",
						Email:    "test_001@example.com",
						Password: "test_pass",
					},
					err: nil,
				},
			},
			want: &entity.User{
				ID:    "test_user_id_001",
				Name:  "test_user_001",
				Email: "test_001@example.com",
			},
			wantErr: false,
		},
		{
			name: "return error when user creation failed",
			args: args{
				ctx: context.Background(),
				input: CreateUserInput{
					Name:     "test_user_001",
					Email:    "test_001@example.com",
					Password: "test_pass",
				},
			},
			mockReturn: mockReturn{
				userCreate: &mockUserCreateReturn{
					user: nil,
					err:  errors.New("failed to create user"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "return error when user credentials creation failed",
			args: args{
				ctx: context.Background(),
				input: CreateUserInput{
					Name:     "test_user_001",
					Email:    "test_001@example.com",
					Password: "test_pass",
				},
			},
			mockReturn: mockReturn{
				userCreate: &mockUserCreateReturn{
					user: &entity.User{
						ID:    "test_user_id_001",
						Name:  "test_user_001",
						Email: "test_001@example.com",
					},
					err: nil,
				},
				userCredsCreate: &mockUserCredsCreateReturn{
					creds: nil,
					err:   errors.New("failed to create user credentials"),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := portmocks.NewUserGateway(t)
			if tt.mockReturn.userCreate != nil {
				users.
					On("Create", tt.args.ctx, port.UserCreateInput{
						Name:  tt.args.input.Name,
						Email: tt.args.input.Email,
					}).
					Return(
						tt.mockReturn.userCreate.user,
						tt.mockReturn.userCreate.err,
					).
					Times(1)
			}
			userCreds := portmocks.NewUserCredentialGateway(t)
			if tt.mockReturn.userCredsCreate != nil {
				userCreds.
					On("Create", tt.args.ctx, port.UserCredentialCreateInput{
						UserID:   tt.mockReturn.userCreate.user.ID,
						Email:    tt.mockReturn.userCreate.user.Email,
						Password: tt.args.input.Password,
					}).
					Return(
						tt.mockReturn.userCredsCreate.creds,
						tt.mockReturn.userCredsCreate.err,
					).
					Times(1)
			}

			it := &userInteractor{
				users:     users,
				userCreds: userCreds,
			}
			got, err := it.CreateUser(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got, "userInteractor.CreateUser() = %v, want %v", got, tt.want)
		})
	}
}
