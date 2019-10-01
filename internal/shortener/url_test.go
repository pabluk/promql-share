package url

import (
	"github.com/alicebob/miniredis"
	"testing"
)

func TestGenerateShareID(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	s.Set("counter", "0")

	cases := []struct {
		in, want string
	}{
		{"graph?g0.range_input=1h&g0.expr=go_threads&g0.tab=1", "6"},
		{"graph?g0.range_input=1h&g0.expr=prometheus_http_requests_total&g0.tab=0", "s"},
		{"graph?g0.range_input=1h&g0.expr=prometheus_http_requests_total&g0.tab=1", "C"},
	}
	for _, c := range cases {
		got := generateShareID(c.in, s.Addr())
		if got != c.want {
			t.Errorf("GenerateShareID(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
