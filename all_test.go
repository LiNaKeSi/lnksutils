package utils

import "testing"

func TestGetFreePort(t *testing.T) {
	port, err := GetFreePort("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if port == 0 {
		t.Fatal("Failed Get TCP port")
	}

	port, err = GetFreePort("")
	if err != nil {
		t.Fatal(err)
	}
	if port == 0 {
		t.Fatal("Failed Get TCP port")
	}
}
