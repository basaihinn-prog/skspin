package jiliDemo

import (
	"testing"
)

func TestGetDemoUrl(t *testing.T) {
	type args struct {
		name   string
		gameId int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				name:   "",
				gameId: 51,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDemoUrl(tt.args.name, tt.args.gameId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDemoUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDemoUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccount(t *testing.T) {
	type args struct {
		accountName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				accountName: "Account=xiaoxiang2",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAccount(tt.args.accountName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExchangeJili(t *testing.T) {
	type args struct {
		account string
		amount  float64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				account: "xiaoxiang",
				amount:  99999.99,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExchangeJili(tt.args.account, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExchangeJili() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExchangeJili() got = %v, want %v", got, tt.want)
			}
		})
	}
}
