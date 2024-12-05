package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"terraform-provider-coolify/internal/testutils"
)

func TestTeamModel_Attributes(t *testing.T) {
	model := teamModel{}

	expected := testutils.GenerateAttrTypesFromStruct(t, model)
	actual := model.FilterAttributes()

	for _, key := range teamsFilterNames {
		_, exists := actual[key]
		assert.True(t, exists, "Key %q should exist in actual attributes", key)
	}

	for key := range actual {
		_, exists := expected[key]
		assert.True(t, exists, "Key %q should exist in expected attributes", key)
	}
}
