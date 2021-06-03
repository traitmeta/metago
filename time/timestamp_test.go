package time

import (
	"database/sql/driver"
	"reflect"
	"testing"
	"time"
)

func TestTimestamp_GetBSON(t *testing.T) {
	tests := []struct {
		name    string
		t       Timestamp
		want    interface{}
		wantErr bool
	}{
		{
			name: "test timestamp get bson success",
			t:    Timestamp(time.Unix(19999999, 0)),
			want: struct {
				Value int64
			}{
				Value: time.Unix(19999999, 0).Unix(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.GetBSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		t       Timestamp
		want    []byte
		wantErr bool
	}{
		{
			name:    "test timestamp marshal json",
			t:       Timestamp(time.Unix(19999999, 0)),
			want:    []byte{49, 57, 57, 57, 57, 57, 57, 57},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_Scan(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		t       Timestamp
		args    args
		want    Timestamp
		wantErr bool
	}{
		{
			name: "test scan true",
			t:    Timestamp{},
			args: args{
				v: time.Unix(19999999, 0),
			},
			want:    Timestamp(time.Unix(19999999, 0)),
			wantErr: false,
		},
		{
			name: "test scan false",
			t:    Timestamp{},
			args: args{
				v: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.Scan(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.t, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", tt.t, tt.want)
			}
		})
	}
}

func TestTimestamp_Unix(t *testing.T) {
	tests := []struct {
		name string
		t    Timestamp
		want int64
	}{
		{
			name: "test get time unix",
			t:    Timestamp(time.Unix(19999999, 0)),
			want: 19999999,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Unix(); got != tt.want {
				t.Errorf("Unix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		t       Timestamp
		args    args
		want    Timestamp
		wantErr bool
	}{
		{
			name: "test unmarshal json success",
			t:    Timestamp{},
			args: args{
				b: []byte{49, 57, 57, 57, 57, 57, 57, 57},
			},
			want:    Timestamp(time.Unix(19999999, 0)),
			wantErr: false,
		},
		{
			name: "test unmarshal json fail",
			t:    Timestamp{},
			args: args{
				b: []byte{49, 57, 57, 57, 57, 57, 57, 0x5F},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.t, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", tt.t, tt.want)
			}
		})
	}
}

func TestTimestamp_Value(t *testing.T) {
	tests := []struct {
		name    string
		t       Timestamp
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "test timestamp value",
			t:       Timestamp(time.Unix(19999999, 0)),
			want:    time.Unix(19999999, 0).Format("2006-01-02 15:04:05"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
		})
	}
}
