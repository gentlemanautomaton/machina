package commandoption_test

import "testing"

func TestParametersAdd(t *testing.T) {
	var params Parameters
	params.AddValue("test")
	params.Add("a", "1")
	want, got := "test,a=1", params.String()
	if want != got {
		t.Errorf("want \"%s\" (got \"%s\")", want, got)
	}
}

func TestParametersAddEmpty(t *testing.T) {
	var params Parameters
	params.Add("", "")
	params.AddValue("")
	if len(params) != 0 {
		t.Fail()
	}
}
