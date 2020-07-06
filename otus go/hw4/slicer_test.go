package hw4

import (
	"reflect"
	"testing"
)

func TestSlicer(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Test words 1",
			args: args{
				text: `один два, три четырЕ - пять. шЕсть, семь вОсемь девять дЕсять.
				Один два, три чЕтыре - пять. шесть, сЕмь восЕмь девять,
				один два, три четырЕ - пять. шесть, семь вОсЕмь -
				Один два, три чЕтыре - пять. шесть, сЕмь
				один два, три четыре - пять. шесть,
				Один два, три четырЕ - пять.
				один два, три чЕтыре
				Один два, три  -
				один два,
				Один`,
			},
			want:    []string{"один", "два", "три", "четыре", "пять", "шесть", "семь", "восемь", "девять", "десять"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Slicer(tt.args.text, 10)
			if (err != nil) != tt.wantErr {
				t.Errorf("Slicer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slicer() = %v, want %v", got, tt.want)
			}
		})
	}
}
