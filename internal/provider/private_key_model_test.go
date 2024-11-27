package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyModel_AttributeTypes(t *testing.T) {
	model := privateKeyModel{}

	expected := generateAttrTypesFromStruct(t, model)
	actual := model.AttributeTypes()

	assert.Equal(t, expected, actual, "AttributeTypes should return the correct attribute types")
}

func TestPrivateKeyModel_Attributes(t *testing.T) {
	model := privateKeyModel{}

	expected := generateAttrTypesFromStruct(t, model)
	actual := model.Attributes()

	for key := range expected {
		_, exists := actual[key]
		assert.True(t, exists, "Key %q should exist in Attributes", key)
	}
}
