package env

import (
	"testing"

	assert "github.com/blendlabs/go-assert"
	util "github.com/blendlabs/go-util"
)

func TestNewVarsFromEnvironment(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(NewVarsFromEnvironment())
}

func TestVarsSet(t *testing.T) {
	assert := assert.New(t)

	vars := Vars{
		"Foo": "baz",
	}

	vars.Set("Foo", "bar")
	assert.Equal("bar", vars.String("Foo"))

	vars.Set("NotFoo", "buzz")
	assert.Equal("buzz", vars.String("NotFoo"))
}

func TestEnvBool(t *testing.T) {
	assert := assert.New(t)

	vars := Vars{
		"true":  "true",
		"1":     "1",
		"yes":   "yes",
		"false": "false",
	}

	assert.True(vars.Bool("true"))
	assert.True(vars.Bool("1"))
	assert.True(vars.Bool("yes"))
	assert.False(vars.Bool("false"))
	assert.False(vars.Bool("no"))

	// Test Set False
	assert.False(vars.Bool("false"))

	// Test Unset Default
	assert.False(vars.Bool("0"))

	// Test Unset User Default
	assert.True(vars.Bool("0", true))
}

func TestEnvInt(t *testing.T) {
	assert := assert.New(t)

	vars := Vars{
		"One": "1",
		"Two": "2",
		"Foo": "Bar",
	}

	assert.Equal(1, vars.Int("One"))
	assert.Equal(2, vars.Int("Two"))
	assert.Zero(vars.Int("Foo"))
	assert.Zero(vars.Int("Baz"))
	assert.Equal(4, vars.Int("Baz", 4))
}

func TestEnvInt64(t *testing.T) {
	assert := assert.New(t)

	vars := Vars{
		"One": "1",
		"Two": "2",
		"Foo": "Bar",
	}

	assert.Equal(1, vars.Int64("One"))
	assert.Equal(2, vars.Int64("Two"))
	assert.Zero(vars.Int64("Foo"))
	assert.Zero(vars.Int64("Baz"))
	assert.Equal(4, vars.Int64("Baz", 4))
}

func TestEnvBytes(t *testing.T) {
	assert := assert.New(t)

	vars := Vars{
		"Foo": "abcdef",
	}

	assert.Equal("abcdef", string(vars.Bytes("Foo")))
	assert.Nil(vars.Bytes("NotFoo"))
	assert.Equal("Bar", string(vars.Bytes("NotFoo", []byte("Bar"))))
}

func TestEnvBase64(t *testing.T) {
	assert := assert.New(t)

	testValue := util.Base64.Encode([]byte("this is a test"))
	vars := Vars{
		"Foo": string(testValue),
	}

	assert.Equal("this is a test", string(vars.Base64("Foo")))
	assert.Equal("this is not a test", string(vars.Base64("Bar", []byte("this is not a test"))))
}

func TestEnvHasKey(t *testing.T) {
	assert := assert.New(t)

	vars := Vars{
		"test1": "foo",
		"test2": "bar",
		"test3": "baz",
		"test4": "buzz",
	}

	assert.True(vars.HasVar("test1"))
	assert.False(vars.HasVar("notTest1"))
}

func TestEnvHasAnyKeys(t *testing.T) {
	assert := assert.New(t)

	vars := Vars{
		"test1": "foo",
		"test2": "bar",
		"test3": "baz",
		"test4": "buzz",
	}

	assert.True(vars.HasAnyVars("test1"))
	assert.True(vars.HasAnyVars("test1", "test2", "test3", "test4"))
	assert.True(vars.HasAnyVars("test1", "test2", "test3", "notTest4"))
	assert.False(vars.HasAnyVars("notTest1", "notTest2"))
	assert.False(vars.HasAnyVars())
}

func TestEnvHasAllKeys(t *testing.T) {
	assert := assert.New(t)

	vars := Vars{
		"test1": "foo",
		"test2": "bar",
		"test3": "baz",
		"test4": "buzz",
	}

	assert.True(vars.HasAllVars("test1"))
	assert.True(vars.HasAllVars("test1", "test2", "test3", "test4"))
	assert.False(vars.HasAllVars("test1", "test2", "test3", "notTest4"))
	assert.False(vars.HasAllVars())
}

func TestVarsKeys(t *testing.T) {
	assert := assert.New(t)

	vars := Vars{
		"test1": "foo",
		"test2": "bar",
		"test3": "baz",
		"test4": "buzz",
	}

	keys := vars.Vars()
	assert.Len(keys, 4)
	assert.Any(keys, func(v interface{}) bool { return v.(string) == "test1" })
	assert.Any(keys, func(v interface{}) bool { return v.(string) == "test2" })
	assert.Any(keys, func(v interface{}) bool { return v.(string) == "test3" })
	assert.Any(keys, func(v interface{}) bool { return v.(string) == "test4" })
}

func TestEnvUnion(t *testing.T) {
	assert := assert.New(t)

	vars1 := Vars{
		"test3": "baz",
		"test4": "buzz",
	}

	vars2 := Vars{
		"test1": "foo",
		"test2": "bar",
	}

	union := vars1.Union(vars2)

	assert.Len(union, 4)
	assert.True(union.HasAllVars("test1", "test3"))
}

type readInto struct {
	Test1 string  `env:"test1"`
	Test2 int     `env:"test2"`
	Test3 float64 `env:"test3"`
}

func TestEnvReadInto(t *testing.T) {
	assert := assert.New(t)

	vars1 := Vars{
		"test1": "foo",
		"test2": "1",
		"test3": "2.0",
	}

	var obj readInto
	err := vars1.ReadInto(&obj)
	assert.Nil(err)
	assert.Equal("foo", obj.Test1)
	assert.Equal(1, obj.Test2)
	assert.Equal(2.0, obj.Test3)
}
