package main

import "testing"

// Тест всего и сразу
// func TestStringPrimitiveDecoder(t *testing.T) {
// 	tests := []struct {
// 		input       string
// 		expectedStr string
// 		expectedErr bool
// 	}{
// 		{"a4bc2d5e", "aaaabccddddde", false},
// 		{"abcd", "abcd", false},
// 		{"45", "", true},
// 		{"", "", false},
// 		{"qwe\\4\\5", "qwe45", false},
// 		{"qwe\\45", "qwe44444", false},
// 		{"qwe\\\\5", "qwe\\\\\\\\\\", false},
// 		{"abc\\", "", true},
// 	}

// 	for _, v := range tests {
// 		t.Run(v.input, func(t *testing.T) {
// 			str, err := StringPrimitiveDecoder(v.input)
// 			if (err != nil) != v.expectedErr {
// 				t.Errorf("StringPrimitiveDecoder() error %t, expected error: %t", err, v.expectedErr)
// 				return
// 			}
// 			if str != v.expectedStr {
// 				t.Errorf("StringPrimitiveDecoder() string = %q, expected string \"%q\"", str, v.expectedStr)
// 			}
// 		})
// 	}
// }

// Тест разбытый на одинаковые по функционалу тесты (худший вариант, представленного выше)
func Test1(t *testing.T) {
	input := "a4bc2d5e"
	expectedStr := "aaaabccddddde"
	expectedErr := false

	str, err := StringPrimitiveDecoder(input)

	if (err != nil) != expectedErr {
		t.Errorf("StringPrimitiveDecoder() error : %t, expected error: %t", err, expectedErr)
	}

	if str != expectedStr {
		t.Errorf("StringPrimitiveDecoder() string = %q, expected string = %q", str, expectedStr)
	}
}

func Test2(t *testing.T) {
	input := "abcd"
	expectedStr := "abcd"
	expectedErr := false

	str, err := StringPrimitiveDecoder(input)

	if (err != nil) != expectedErr {
		t.Errorf("StringPrimitiveDecoder() error : %t, expected error: %t", err, expectedErr)
	}

	if str != expectedStr {
		t.Errorf("StringPrimitiveDecoder() string = %q, expected string = %q", str, expectedStr)
	}
}

func Test3(t *testing.T) {
	input := "45"
	expectedStr := ""
	expectedErr := true

	str, err := StringPrimitiveDecoder(input)

	if (err != nil) != expectedErr {
		t.Errorf("StringPrimitiveDecoder() error : %t, expected error: %t", err, expectedErr)
	}

	if str != expectedStr {
		t.Errorf("StringPrimitiveDecoder() string = %q, expected string = %q", str, expectedStr)
	}
}

func Test4(t *testing.T) {
	input := ""
	expectedStr := ""
	expectedErr := false

	str, err := StringPrimitiveDecoder(input)

	if (err != nil) != expectedErr {
		t.Errorf("StringPrimitiveDecoder() error : %t, expected error: %t", err, expectedErr)
	}

	if str != expectedStr {
		t.Errorf("StringPrimitiveDecoder() string = %q, expected string = %q", str, expectedStr)
	}
}

func Test5(t *testing.T) {
	input := "qwe\\4\\5"
	expectedStr := "qwe45"
	expectedErr := false

	str, err := StringPrimitiveDecoder(input)

	if (err != nil) != expectedErr {
		t.Errorf("StringPrimitiveDecoder() error : %t, expected error: %t", err, expectedErr)
	}

	if str != expectedStr {
		t.Errorf("StringPrimitiveDecoder() string = %q, expected string = %q", str, expectedStr)
	}
}

func Test6(t *testing.T) {
	input := "qwe\\45"
	expectedStr := "qwe44444"
	expectedErr := false

	str, err := StringPrimitiveDecoder(input)

	if (err != nil) != expectedErr {
		t.Errorf("StringPrimitiveDecoder() error : %t, expected error: %t", err, expectedErr)
	}

	if str != expectedStr {
		t.Errorf("StringPrimitiveDecoder() string = %q, expected string = %q", str, expectedStr)
	}
}

func Test7(t *testing.T) {
	input := "qwe\\\\5"
	expectedStr := "qwe\\\\\\\\\\"
	expectedErr := false

	str, err := StringPrimitiveDecoder(input)

	if (err != nil) != expectedErr {
		t.Errorf("StringPrimitiveDecoder() error : %t, expected error: %t", err, expectedErr)
	}

	if str != expectedStr {
		t.Errorf("StringPrimitiveDecoder() string = %q, expected string = %q", str, expectedStr)
	}
}

func Test8(t *testing.T) {
	input := "abc\\"
	expectedStr := ""
	expectedErr := true

	str, err := StringPrimitiveDecoder(input)

	if (err != nil) != expectedErr {
		t.Errorf("StringPrimitiveDecoder() error : %t, expected error: %t", err, expectedErr)
	}

	if str != expectedStr {
		t.Errorf("StringPrimitiveDecoder() string = %q, expected string = %q", str, expectedStr)
	}
}
