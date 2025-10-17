package part

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	repoStub "github.com/alexander-kartavtsev/starship/inventory/internal/repository_stub/part"
	servStub "github.com/alexander-kartavtsev/starship/inventory/internal/service_stub/part"
)

func TestGetPart(t *testing.T) {
	r := repoStub.NewRepository()
	s := servStub.NewService(r)

	t.Run("GetPart_correct", func(t *testing.T) {
		partUuid := "part_uuid_1"
		expected := model.Part{
			Uuid:         "part_uuid_1",
			Dimensions:   &model.Dimensions{},
			Manufacturer: &model.Manufacturer{},
		}

		partRes, err := s.Get(t.Context(), partUuid)

		assert.NoError(t, err)
		assert.Equal(t, expected, partRes)
	})

	t.Run("GetPart_notFoundErr", func(t *testing.T) {
		partUuid := "part_uuid_not_exists"
		expected := model.Part{}

		partRes, err := s.Get(t.Context(), partUuid)

		assert.Error(t, err)
		assert.Equal(t, expected, partRes)
	})
}
