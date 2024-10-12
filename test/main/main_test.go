package main

import "testing"

func TestSumOfIntOrFloat(t *testing.T) {
	type args[K comparable, V interface{ int | float64 }] struct {
		m map[K]V
	}
	// 定义测试用例类型
	type testCase[K comparable, V interface{ int | float64 }] struct {
		name string
		args args[K, V]
		want V
	}
	// 定义一个存储测试用例的切片
	tests := []testCase[int, float64]{
		{"case1", args[int, float64]{map[int]float64{1: 1.0, 2: 2.0}}, 3.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumOfIntOrFloat(tt.args.m); got != tt.want {
				t.Errorf("SumOfIntOrFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sum(t *testing.T) {
	if v := sum(1, 2); v != 3 {
		t.Errorf("sum(1,2) = %v, want 3", v)
	}
}

func Test_sum1(t *testing.T) {
	// 定义一个测试用例类型
	type args struct {
		numbers []int
	}

	// 定义一个存储测试用例的切片
	tests := []struct {
		name string
		args args
		want int
	}{
		{"ceshi", args{numbers: []int{1, 2, 3}}, 6},
	}

	// 遍历执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sum(tt.args.numbers...); got != tt.want {
				t.Errorf("sum() = %v, want %v", got, tt.want)
			}
		})
	}
}
