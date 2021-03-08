package main

import (
	"reflect"
	"testing"
)

func TestFPortParse(t *testing.T) {
	testsValid := [][]interface{}{
		{FPort("10"), []uint16{10}},
		{FPort("10,20"), []uint16{10, 20}},
		{FPort("10~12"), []uint16{10, 11, 12}},
		{FPort("10~12,5"), []uint16{10, 11, 12, 5}},
	}

	for _, test := range testsValid {
		ports, _ := test[0].(FPort).parse()
		if !reflect.DeepEqual(ports, test[1]) {
			t.Errorf("Get %v, want %v", ports, test[1])
		}
	}

	testsInvalid := []FPort{
		FPort("-1"),
		FPort("65536"),
		FPort("10,-10"),
		FPort("-10~10"),
		FPort("10~12~58"),
	}

	for _, test := range testsInvalid {
		_, err := test.parse()
		if err == nil {
			t.Errorf("Get <nil>, want error")
		}
	}
}
