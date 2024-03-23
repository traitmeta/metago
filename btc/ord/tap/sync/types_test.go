package sync

import (
	"reflect"
	"testing"
)

func TestParseElementFromString(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name    string
		args    args
		want    *Element
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				content: "cardinalalphabet.1=a,2=b,3=c,4=d,5=e,6=f,7=g,8=h,9=i,10=j,11=k,12=l,13=m,14=n,15=o,16=p,17=q,18=r,19=s,20=t,21=u,22=v,23=w,24=x,25=y,26=z,16.element",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseElementFromString(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseElementFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseElementFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
