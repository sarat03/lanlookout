package check

import "testing"

func TestJoinNames(t *testing.T) {

	name, dns := joinNames([]string{"nas.lan.", "nas.local."})
	if name != "nas.lan" || dns != "nas.lan nas.local" {
		t.Fatalf("got %q, %q", name, dns)
	}

	name, dns = joinNames(nil)
	if name != "" || dns != "" {
		t.Fatalf("empty input: got %q, %q", name, dns)
	}
}
