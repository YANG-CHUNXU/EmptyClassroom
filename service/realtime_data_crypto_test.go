package service

import "testing"

func TestEncodeLoginPassword(t *testing.T) {
	got, err := encodeLoginPassword("123456")
	if err != nil {
		t.Fatalf("encodeLoginPassword() error = %v", err)
	}
	const want = "NGZSNVNOdG5kdDNwMDZYSm0rSi81dz09"
	if got != want {
		t.Fatalf("encodeLoginPassword() = %q, want %q", got, want)
	}
}

func TestResolveRealtimeCampusID(t *testing.T) {
	testCases := []struct {
		name string
		id   int
		want string
	}{
		{name: "west", id: 1, want: "01"},
		{name: "shahe", id: 2, want: "04"},
		{name: "fallback", id: 4, want: "04"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := resolveRealtimeCampusID(tc.id); got != tc.want {
				t.Fatalf("resolveRealtimeCampusID() = %q, want %q", got, tc.want)
			}
		})
	}
}
