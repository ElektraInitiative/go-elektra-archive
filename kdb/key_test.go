package kdb_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	elektra "github.com/ElektraInitiative/go-elektra/kdb"
	. "github.com/ElektraInitiative/go-elektra/test"
)

func TestName(t *testing.T) {
	firstName := "user/tests/go-elektra/name/1"
	k, err := elektra.CreateKey(firstName)

	Check(t, err, "could not create key")
	Assert(t, k.Name() == firstName, "wrong key name")

	secondName := "user/tests/go-elektra/name/2"
	err = k.SetName(secondName)

	Check(t, err, "could not set key name")
	Assert(t, k.Name() == secondName, "could not set name")
}

func TestString(t *testing.T) {
	testValue := "Hello World"
	k, err := elektra.CreateKey("user/tests/go-elektra/string", testValue)

	Check(t, err, "could not create key")

	val := k.Value()

	Assertf(t, val == testValue, "Key.GetString() did not match %q", testValue)
}

func TestBoolean(t *testing.T) {
	k, err := elektra.CreateKey("user/tests/go-elektra/boolean")

	Check(t, err, "could not create key")

	testValue := true

	err = k.SetBoolean(testValue)

	Check(t, err, "SetBoolean failed")

	val := k.Boolean()

	Assertf(t, val == testValue, "Key.Boolean() %t did not match %t", val, testValue)

	testValue = !testValue

	err = k.SetBoolean(testValue)

	Check(t, err, "SetBoolean failed")

	val = k.Boolean()

	Assertf(t, val == testValue, "Key.Boolean() %t did not match %t", val, testValue)
}

func TestBytes(t *testing.T) {
	k, err := elektra.CreateKey("user/tests/go-elektra/bytes")

	Check(t, err, "could not create key")

	testValue := make([]byte, 10)
	rand.Read(testValue)

	err = k.SetBytes(testValue)

	Check(t, err, "")

	val := k.Bytes()

	Assertf(t, bytes.Compare(val, testValue) == 0, "Key.Bytes() %X did not match %X", val, testValue)
}

func TestMeta(t *testing.T) {
	k, err := elektra.CreateKey("user/tests/go-elektra/meta", "Hello World")

	Check(t, err, "could not create key")

	err = k.SetMeta("meta", "value")

	Check(t, err, "could not set meta")

	val := k.Meta("meta")

	Assert(t, val == "value", "Key.Meta() did not return the correct value")
}

func TestNamespace(t *testing.T) {
	key, _ := elektra.CreateKey("user/tests/go-elektra/namespace")

	namespace := key.Namespace()
	expected := "user"

	Assertf(t, namespace == expected, "Namespace be %q but is %q", expected, namespace)

	key, _ = elektra.CreateKey("/go-elektra/namespace")

	namespace = key.Namespace()
	expected = ""

	Assertf(t, namespace == expected, "Namespace be %q but is %q", expected, namespace)
}

var commonKeyNameTests = []struct {
	key1     string
	key2     string
	expected string
}{
	{"user/foo/bar", "user/foo/bar2", "user/foo"},
	{"proc/foo/bar", "user/foo/bar", "/foo/bar"},
	{"user/foo/bar", "user/bar/foo", "user"},
	{"proc/bar/foo", "user/foo/bar", ""},
}

func TestCommonKeyName(t *testing.T) {
	for _, test := range commonKeyNameTests {
		t.Run(fmt.Sprintf("(%q, %q)", test.key1, test.key2), func(t *testing.T) {
			key1, _ := elektra.CreateKey(test.key1)
			key2, _ := elektra.CreateKey(test.key2)

			commonName := elektra.CommonKeyName(key1, key2)

			Assertf(t, commonName == test.expected, "commonName should be %q but is %q", test.expected, commonName)
		})
	}
}
