package solution

import (
	"reflect"
	"testing"
)

func Test_removeNoOrder(t *testing.T) {
	type args struct {
		s []int
		i int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "remove first element",
			args: args{s: []int{1, 2, 3}, i: 0},
			want: []int{3, 2},
		},
		{
			name: "remove last element",
			args: args{s: []int{1, 2, 3}, i: 2},
			want: []int{1, 2},
		},
		{
			name: "remove middle element",
			args: args{s: []int{1, 2, 3, 4}, i: 1},
			want: []int{1, 4, 3},
		},
		{
			name: "remove element from slice with duplicates",
			args: args{s: []int{1, 2, 2, 3}, i: 1},
			want: []int{1, 3, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeNoOrder(tt.args.s, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeNoOrder() = %v, want %v", got, tt.want)
			}
		})
	}

}
func Test_removeNoOrderMultiple(t *testing.T) {
	type args struct {
		s []int
		i []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "remove first two elements",
			args: args{s: []int{1, 2, 3}, i: []int{0, 0}},
			want: []int{2},
		},
		{
			name: "remove two last elements",
			args: args{s: []int{1, 2, 3}, i: []int{2, 1}},
			want: []int{1},
		},
		{
			name: "remove two middle elements",
			args: args{s: []int{1, 2, 3, 4}, i: []int{1, 1}},
			want: []int{1, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.s
			for _, v := range tt.args.i {
				got = removeNoOrder(got, v)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeNoOrder() = %v, want %v", got, tt.want)
			}
		})
	}

}
