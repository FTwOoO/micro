package cfg

import "testing"


type myConfig struct {
	ConfigurationImp
	X int
}
func TestNewConfiguration(t *testing.T) {
	cf := new(myConfig)
	err := NewConfiguration("parse_test.json", cf)
	if err != nil {
		t.Fatal(err)
	}

	v1 := cf.GetMongoDb().IsValid()
	if v1 == false {
		t.FailNow()
	}
	t.Log(cf)
}