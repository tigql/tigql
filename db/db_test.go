package db_test

import (
	"testing"

	"github.com/tigql/tigql/db"

	"github.com/stretchr/testify/require"
)

func TestOpenAndClose(t *testing.T) {
	t.Parallel()
	got, err := db.Open(currentTestDatasource)
	require.NoError(t, err)
	require.Equal(t, currentTestDatasource.DBName, got.Schema.Name)
	err = got.Close()
	require.NoError(t, err)
}

func TestDB_CheckVersion(t *testing.T) {
	err := currentTestDB.CheckVersion()
	require.NoError(t, err)
}
