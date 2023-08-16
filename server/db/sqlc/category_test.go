package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCategories(t *testing.T) {
	categoryList, err := testQueries.GetCategoryList(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, categoryList)

	for _, category := range categoryList {
		require.NotZero(t, category.CategoryID)
		require.NotEmpty(t, category.Title)
		require.NotZero(t, category.CreatedAt)
		require.NotZero(t, category.UpdatedAt)
	}
}
