package basicfuncs

import "testing"

func TestGenerateRandomString(t *testing.T) {
	s := GenerateRandomString(10)
	if len(s) != 10 {
		t.Errorf("Expected string length 10, got %d", len(s))
	}
}

func TestConvertQueryParameterToInt(t *testing.T) {
	paramValue, err := ConvertQueryParameterToInt("10", 0)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if paramValue != 10 {
		t.Errorf("Expected paramValue 10, got %d", paramValue)
	}
	paramValue, err = ConvertQueryParameterToInt("", 10)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if paramValue != 10 {
		t.Errorf("Expected paramValue 10, got %d", paramValue)
	}
	_, err = ConvertQueryParameterToInt("abc", 0)
	if err == nil {
		t.Errorf("Expected error, got no error")
	}
}
func TestDeleteNonAlphanumeric(t *testing.T) {
	input := "Hello, World! 123"
	expected := "HelloWorld123"
	result := DeleteNonAlphanumeric(input)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	input = "Testing @#$%^&*()_+"
	expected = "Testing"
	result = DeleteNonAlphanumeric(input)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	input = "NoSpecialCharacters"
	expected = "NoSpecialCharacters"
	result = DeleteNonAlphanumeric(input)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}
