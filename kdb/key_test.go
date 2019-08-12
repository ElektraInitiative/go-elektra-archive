package kdb_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/ElektraInitiative/go-elektra"
)

func TestName(t *testing.T) {
	kdb := elektra.New()

	firstName := "user/hello_world"
	secondName := "user/hello_world2"

	k, err := kdb.CreateKey(firstName)

	Check(t, err, "could not create key")
	Assert(t, k.Name() == firstName, "wrong key name")

	err = k.SetName(secondName)

	Check(t, err, "could not set key name")
	Assert(t, k.Name() == secondName, "could not set name")
}

func TestString(t *testing.T) {
	kdb := elektra.New()

	testValue := "Hello World"
	k, err := kdb.CreateKey("user/hello_world", testValue)

	Check(t, err, "could not create key")

	val := k.Value()

	Assertf(t, val == testValue, "Key.GetString() did not match %q", testValue)
}

func TestBoolean(t *testing.T) {
	kdb := elektra.New()

	k, err := kdb.CreateKey("user/hello_world")

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
	kdb := elektra.New()

	k, err := kdb.CreateKey("user/hello_world")

	Check(t, err, "could not create key")

	testValue := make([]byte, 10)
	rand.Read(testValue)

	err = k.SetBytes(testValue)

	Check(t, err, "")

	val := k.Bytes()

	Assertf(t, bytes.Compare(val, testValue) == 0, "Key.Bytes() %q did not match %q", val, testValue)
}

func TestMeta(t *testing.T) {
	kdb := elektra.New()

	k, err := kdb.CreateKey("user/hello_world", "Hello World")

	Check(t, err, "could not create key")

	err = k.SetMeta("meta", "value")

	Check(t, err, "could not set meta")

	val := k.Meta("meta")

	Assert(t, val == "value", "Key.Meta() did not return the correct value")
}
