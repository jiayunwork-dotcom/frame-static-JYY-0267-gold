package format

import "testing"

func TestFixed(t *testing.T) {
	if Fixed(1.5, 2) != "1.50" {
		t.Fatalf("got %q", Fixed(1.5, 2))
	}
}

func TestSigned(t *testing.T) {
	if Signed(-3.0, 1) != "-3.0" {
		t.Fatalf("got %q", Signed(-3.0, 1))
	}
	if Signed(3.0, 1) != "+3.0" {
		t.Fatalf("got %q", Signed(3.0, 1))
	}
}

func TestEng(t *testing.T) {
	if Eng(0) != "0" {
		t.Fatalf("zero -> %q", Eng(0))
	}
	if Eng(1500) != "1.500k" {
		t.Fatalf("1500 -> %q", Eng(1500))
	}
	if Eng(0.002) != "2.000m" {
		t.Fatalf("0.002 -> %q", Eng(0.002))
	}
}
