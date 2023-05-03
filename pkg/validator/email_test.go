package validator

import "testing"

func TestValidateEmail(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid email address", args: args{s: "test123@gmail.com"}, wantErr: false},
		{name: "invalid email address", args: args{s: "test1906gmail.com"}, wantErr: true},
		{name: "invalid email address host", args: args{s: "test1906@gmailcom"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateEmail(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
