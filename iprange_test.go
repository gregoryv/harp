package warp

import "testing"

func Test_IPRange(t *testing.T) {
	cases := []struct {
		name string
		in   string
		exp  int
	}{
		{in: "192.168.1.1", exp: 1},
		{in: "192.168.1.*", exp: 255},
		{in: "192.168.1.3-4", exp: 2},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := IPRange(c.in)
			if err != nil {
				t.Fatal(err)
			}
			if got := len(res); got != c.exp {
				t.Errorf("got %v, exp %v", got, c.exp)
			}
		})
	}
}
