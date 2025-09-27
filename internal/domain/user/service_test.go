package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/stretchr/testify/assert"
)

type userMatcher struct {
	expected User
}

func (m userMatcher) Matches(x interface{}) bool {
	u, ok := x.(User)
	if !ok {
		return false
	}
	// сравниваем по логину, например, а хеш пароля игнорируем
	return u.Login == m.expected.Login
}

func (m userMatcher) String() string {
	return fmt.Sprintf("matches user with login %s", m.expected.Login)
}

func TestService_RegisterUser(t *testing.T) {
	type args struct {
		user  User
		times int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"Success", args{User{"test", "test"}, 1}, 1, false},
		{"Empty password", args{User{"test", ""}, 0}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorer(ctrl)
			m.EXPECT().RegisterUser(gomock.Any(), userMatcher{tt.args.user}).
				Return(tt.want, nil).Times(tt.args.times)

			cfg := &config.Config{}
			s := NewService(cfg, m)

			id, err := s.RegisterUser(context.Background(), tt.args.user)
			if id != tt.want {
				assert.Equal(t, tt.want, id)
			}
			if err != nil {
				if tt.wantErr {
					assert.Error(t, ErrEmptyPass, err)
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
