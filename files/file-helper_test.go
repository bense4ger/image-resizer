package files

import "testing"

func TestMustOK(t *testing.T) {
	helper := &FSHelper{
		WorkingDir: "foo",
		Extension:  "bar",
	}

	h, err := Must(helper)
	if h == nil || err != nil {
		t.Errorf("Expected helper to be created: %s", err.Error())
	}

}

func TestMustErr(t *testing.T) {
	helper := &FSHelper{}
	_, err := Must(helper)
	if err == nil {
		t.Errorf("Expected Must to return an error")
	}
}
