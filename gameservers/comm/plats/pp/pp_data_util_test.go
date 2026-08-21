package pp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKey(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetKey()
			if !tt.wantErr(t, err, fmt.Sprintf("GetKey()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetKey()")
		})
	}
}

//func TestGetMgckey(t *testing.T) {
//	type args struct {
//		game   string
//		user   string
//		client *http.Client
//	}
//	tests := []struct {
//		name         string
//		args         args
//		wantMgckey   string
//		wantLocation *url.URL
//		wantErr      assert.ErrorAssertionFunc
//	}{
//		{},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotMgckey, gotLocation, err := GetMgckey(tt.args.game, tt.args.user, tt.args.client)
//			if !tt.wantErr(t, err, fmt.Sprintf("GetMgckey(%v, %v, %v)", tt.args.game, tt.args.user, tt.args.client)) {
//				return
//			}
//			assert.Equalf(t, tt.wantMgckey, gotMgckey, "GetMgckey(%v, %v, %v)", tt.args.game, tt.args.user, tt.args.client)
//			assert.Equalf(t, tt.wantLocation, gotLocation, "GetMgckey(%v, %v, %v)", tt.args.game, tt.args.user, tt.args.client)
//		})
//	}
//}
