package lnksutils

import (
	"testing"
)

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

func TestIsDirEmpty(t *testing.T) {
	if IsDirEmpty(".") == true {
		t.Fatal("current directory shouldn't be empty")
	}
	if IsDirEmpty("doesntexists") == false {
		t.Fatal("doesntexists directory shouldn't be empty")
	}
	if IsDirEmpty("aa") == false {
		t.Fatal("doesntexists directory shouldn't be empty")
	}

}
