package repo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_EntityTypesGetAll_RepairsMissingItemType(t *testing.T) {
	ctx := context.Background()
	g, err := tRepos.Groups.GroupCreate(ctx, "location-only", uuid.Nil)
	require.NoError(t, err)
	location, err := tRepos.EntityTypes.GetDefault(ctx, g.ID, true)
	require.NoError(t, err)
	first, err := tRepos.EntityTypes.GetAll(ctx, g.ID)
	require.NoError(t, err)
	require.Len(t, first, 2)
	var itemID uuid.UUID
	for _, typ := range first {
		if typ.IsLocation {
			require.Equal(t, location.ID, typ.ID)
		} else {
			itemID = typ.ID
		}
	}
	require.NotEqual(t, uuid.Nil, itemID)
	second, err := tRepos.EntityTypes.GetAll(ctx, g.ID)
	require.NoError(t, err)
	require.Equal(t, first, second, "Repeated reads must not duplicate default types")
	entity, err := tRepos.Entities.Create(ctx, g.ID, EntityCreate{Name: "Document case", EntityTypeID: itemID, Quantity: 1})
	require.NoError(t, err)
	require.Equal(t, "Document case", entity.Name)
}
