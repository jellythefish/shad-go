package varfmt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	var interfaceSlice []interface{} = make([]interface{}, 1002)
	interfaceSlice[10] = 1
	interfaceSlice[100] = 1
	interfaceSlice[1000] = 1

	for _, tc := range []struct {
		format string
		args   []interface{}
		result string
	}{
		{
			format: "{}",
			args:   []interface{}{0},
			result: "0",
		},
		{
			format: "{0} {0}",
			args:   []interface{}{1},
			result: "1 1",
		},
		{
			format: "{1} {5}",
			args:   []interface{}{0, 1, 2, 3, 4, 5, 6},
			result: "1 5",
		},
		{
			format: "{} {} {} {} {}",
			args:   []interface{}{0, 1, 2, 3, 4},
			result: "0 1 2 3 4",
		},
		{
			format: "{} {0} {0} {0} {}",
			args:   []interface{}{0, 1, 2, 3, 4},
			result: "0 0 0 0 4",
		},
		{
			format: "Hello, {2}",
			args:   []interface{}{0, 1, "World"},
			result: "Hello, World",
		},
		{
			format: "check {10}",
			args:   interfaceSlice[:11],
			result: "check 1",
		},
		{
			format: "check {100}",
			args:   interfaceSlice[:101],
			result: "check 1",
		},
		{
			format: "check {1000}",
			args:   interfaceSlice[:1001],
			result: "check 1",
		},
	} {
		t.Run(tc.result, func(t *testing.T) {
			require.Equal(t, tc.result, Sprintf(tc.format, tc.args...))
		})
	}
}

func BenchmarkFormat(b *testing.B) {
	for _, tc := range []struct {
		name   string
		format string
		args   []interface{}
	}{
		{
			name:   "small int",
			format: "{}",
			args:   []interface{}{42},
		},
		{
			name:   "small string",
			format: "{} {}",
			args:   []interface{}{"Hello", "World"},
		},
		{
			name:   "big",
			format: strings.Repeat("{0}{1}", 1000),
			args:   []interface{}{42, 43},
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = Sprintf(tc.format, tc.args...)
			}
		})
	}
}

func BenchmarkSprintf(b *testing.B) {
	for _, tc := range []struct {
		name   string
		format string
		args   []interface{}
	}{
		{
			name:   "small",
			format: "%d",
			args:   []interface{}{42},
		},
		{
			name:   "small string",
			format: "%v %v",
			args:   []interface{}{"Hello", "World"},
		}, {
			name:   "big",
			format: strings.Repeat("%[0]v%[1]v", 1000),
			args:   []interface{}{42, 43},
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = fmt.Sprintf(tc.format, tc.args...)
			}
		})
	}
}
