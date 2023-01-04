package pgerr

import (
	"testing"

	"github.com/lib/pq"
)

func Test_UniqueColumn(t *testing.T) {
	pqErr := pq.Error{
		Code: "23505",
		Constraint: "users_email_key",
	}

	column, err := UniqueColumn(&pqErr)
	if err == nil {
		t.Errorf("UniqueColumn() = %v, want %v", nil, err)
	}

	if column != "email" {
		t.Errorf("UniqueColumn() = %v, want %v", "email", column)
	}
}

func Test_getUniqueColumn(t *testing.T) {
	type args struct {
		str string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single_args",
			args: args{
				str: "users_username_key",
			},
			want: "username",
		},

		{
			name: "double_args",
			args: args{
				str: "users_created_at_key",
			},
			want: "created at",
		},

		{
			name: "triple_args",
			args: args{
				str: "users_email_verified_at_key",
			},
			want: "email verified at",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUniqueColumn(tt.args.str); got != tt.want {
				t.Errorf("getUniqueColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}
