package urlvaluer

import (
	"fmt"
	"github.com/go-test/deep"
	"net/url"
	"testing"
)

type RandomStruct struct {
	Value int
}

type CustomType struct {
	value int
}

func (ct *CustomType) String() string {
	return fmt.Sprintf("%d", ct.value)
}

type anotherCustomType struct {
	value bool
}

func (ct anotherCustomType) String() string {
	return fmt.Sprintf("%v", ct.value)
}

func toStringPtr(s string) *string {
	return &s
}

func testMarshal(t *testing.T, testStruct any, expected url.Values, isErrExpected bool) {
	urlValues, err := Marshal(testStruct)
	if err == nil && isErrExpected {
		t.Fatal("error expected but there is no error")
	}
	if err != nil && !isErrExpected {
		t.Fatalf("unexpected error: %s", err)
	}

	if notEquals := deep.Equal(expected, urlValues); notEquals != nil {
		t.Fatalf("unexpected result, %v", notEquals)
	}
}

func TestMarshal(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		testStruct := struct {
			FirstName string `url:"first_name"`
		}{"John"}
		expected := url.Values{"first_name": []string{"John"}}

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("String Pointer", func(t *testing.T) {
		testStruct := struct {
			FirstName *string `url:"FirstName"`
		}{toStringPtr("John")}
		expected := url.Values{"FirstName": []string{"John"}}

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Uint", func(t *testing.T) {
		testStruct := struct {
			Age uint `url:"Age"`
		}{55}
		expected := url.Values{"Age": []string{"55"}}

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Bool", func(t *testing.T) {
		testStruct := struct {
			IsClient bool `url:"IsClient"`
			IsLegit  bool `url:"IsLegit"`
		}{true, false}
		expected := url.Values{
			"IsClient": []string{"true"},
			"IsLegit":  []string{"false"},
		}

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Float32", func(t *testing.T) {
		testStruct := struct {
			Value float32 `url:"ehh"`
		}{0.1232131}
		expected := url.Values{"ehh": []string{"0.1232131"}}

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Any Float32", func(t *testing.T) {
		testStruct := struct {
			Value any `url:"Value"`
		}{1.2}
		expected := url.Values{"Value": []string{"1.2"}}

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("[]bool", func(t *testing.T) {
		testStruct := struct {
			Bools []bool `url:"Bools"`
		}{[]bool{true, false, true, false}}
		expected := url.Values{
			"Bools": []string{"true", "false", "true", "false"},
		}

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Int Zero", func(t *testing.T) {
		testStruct := struct {
			Value int `url:"Value"`
		}{0}
		expected := url.Values{"Value": []string{"0"}}

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("[2]int", func(t *testing.T) {
		testStruct := struct {
			Ints [2]int `url:"Ints"`
		}{[2]int{5, 20}}
		expected := url.Values{
			"Ints": []string{"5", "20"},
		}

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Custom Type with .String() interface", func(t *testing.T) {
		testStruct := struct {
			CustomType CustomType `url:"CustomType"`
		}{CustomType{value: 1231}}
		expected := url.Values{"CustomType": []string{"1231"}}

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Custom Type Ptr .String() interface", func(t *testing.T) {
		testStruct := struct {
			CustomType anotherCustomType `url:"CustomType"`
		}{anotherCustomType{value: true}}
		expected := url.Values{"CustomType": []string{"true"}}

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Ptr null", func(t *testing.T) {
		testStruct := struct {
			Value *string `url:"Value"`
		}{nil}
		expected := make(url.Values)

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Not Tagged", func(t *testing.T) {
		testStruct := struct {
			Value string
		}{"SomeValue"}
		expected := make(url.Values)

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Not Public", func(t *testing.T) {
		testStruct := struct {
			value string `url:"Value"`
		}{"SomeValue"}
		expected := make(url.Values)

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Ignore tag", func(t *testing.T) {
		testStruct := struct {
			Value string `url:"-"`
		}{"SomeValue"}
		expected := make(url.Values)

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Map Ignore", func(t *testing.T) {
		testStruct := struct {
			Map map[string]string `url:"Value"`
		}{map[string]string{"Value": "SomeValue"}}
		expected := make(url.Values)

		testMarshal(t, &testStruct, expected, false)
	})

	t.Run("Struct Ignore", func(t *testing.T) {
		testStruct := struct {
			Struct RandomStruct `url:"Struct"`
		}{Struct: RandomStruct{-112}}
		expected := make(url.Values)

		testMarshal(t, &testStruct, expected, false)
	})
}
