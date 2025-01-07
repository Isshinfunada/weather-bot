package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAreaRepoTest(t *testing.T) (repository.AreaRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := repository.NewAreaRepository(db)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

func TestFindHierarchyByClass20ID_Success(t *testing.T) {
	repo, mock, cleanup := setupAreaRepoTest(t)
	defer cleanup()

	ctx := context.Background()
	class20ID := "20"

	// 正規表現でクエリをマッチさせる
	query := regexp.QuoteMeta(`
         SELECT
             c20.id, c20.name, c20.en_name, c20.parent_id,
             c15.id, c15.name, c15.en_name, c15.parent_id,
             c10.id, c10.name, c10.en_name, c10.parent_id,
             o.id, o.name, o.en_name, o.parent_id,
             ct.id, ct.name, ct.en_name, ct.office_name
         FROM area_class20 c20
         JOIN area_class15 c15 ON c20.parent_id = c15.id
         JOIN area_class10 c10 ON c15.parent_id = c10.id
         JOIN area_offices o ON c10.parent_id = o.id
         JOIN area_centers ct ON o.parent_id = ct.id
         WHERE c20.id = $1
	`)

	// モックデータ行を作成
	rows := sqlmock.NewRows([]string{
		"id", "name", "en_name", "parent_id", // c20.*
		"id", "name", "en_name", "parent_id", // c15.*
		"id", "name", "en_name", "parent_id", // c10.*
		"id", "name", "en_name", "parent_id", // o.*
		"id", "name", "en_name", "office_name", // ct.*
	}).AddRow(
		"20", "Area20", "Area20_EN", "15", // c20
		"15", "Area15", "Area15_EN", "10", // c15
		"10", "Area10", "Area10_EN", "5", // c10
		"5", "Office", "Office_EN", "1", // o
		"1", "Center", "Center_EN", "Center_Office", // ct
	)

	mock.ExpectQuery(query).WithArgs(class20ID).WillReturnRows(rows)

	hierarchy, err := repo.FindHierarchyByClass20ID(ctx, class20ID)
	assert.NoError(t, err)
	require.NotNil(t, hierarchy)

	// 結果の検証
	assert.Equal(t, "20", hierarchy.Class20.ID)
	assert.Equal(t, "Area15", hierarchy.Class15.Name)
	assert.Equal(t, "Area10_EN", hierarchy.Class10.EnName)
	assert.Equal(t, "Office", hierarchy.Office.Name)
	assert.Equal(t, "Center_Office", hierarchy.Center.OfficeName)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindHierarchyByClass20ID_NoRows(t *testing.T) {
	repo, mock, cleanup := setupAreaRepoTest(t)
	defer cleanup()

	ctx := context.Background()
	class20ID := "999"

	query := regexp.QuoteMeta(`
         SELECT
             c20.id, c20.name, c20.en_name, c20.parent_id,
             c15.id, c15.name, c15.en_name, c15.parent_id,
             c10.id, c10.name, c10.en_name, c10.parent_id,
             o.id, o.name, o.en_name, o.parent_id,
             ct.id, ct.name, ct.en_name, ct.office_name
         FROM area_class20 c20
         JOIN area_class15 c15 ON c20.parent_id = c15.id
         JOIN area_class10 c10 ON c15.parent_id = c10.id
         JOIN area_offices o ON c10.parent_id = o.id
         JOIN area_centers ct ON o.parent_id = ct.id
         WHERE c20.id = $1
	`)

	// QueryRow で no rows を返す場合
	mock.ExpectQuery(query).WithArgs(class20ID).WillReturnError(sql.ErrNoRows)

	hierarchy, err := repo.FindHierarchyByClass20ID(ctx, class20ID)
	assert.NoError(t, err)
	assert.Nil(t, hierarchy)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindHierarchyByClass20ID_QueryError(t *testing.T) {
	repo, mock, cleanup := setupAreaRepoTest(t)
	defer cleanup()

	ctx := context.Background()
	class20ID := "20"

	query := regexp.QuoteMeta(`
         SELECT
             c20.id, c20.name, c20.en_name, c20.parent_id,
             c15.id, c15.name, c15.en_name, c15.parent_id,
             c10.id, c10.name, c10.en_name, c10.parent_id,
             o.id, o.name, o.en_name, o.parent_id,
             ct.id, ct.name, ct.en_name, ct.office_name
         FROM area_class20 c20
         JOIN area_class15 c15 ON c20.parent_id = c15.id
         JOIN area_class10 c10 ON c15.parent_id = c10.id
         JOIN area_offices o ON c10.parent_id = o.id
         JOIN area_centers ct ON o.parent_id = ct.id
         WHERE c20.id = $1
	`)

	mock.ExpectQuery(query).WithArgs(class20ID).WillReturnError(errors.New("query failed"))

	hierarchy, err := repo.FindHierarchyByClass20ID(ctx, class20ID)
	assert.Error(t, err)
	assert.Nil(t, hierarchy)
	assert.Contains(t, err.Error(), "query error")

	require.NoError(t, mock.ExpectationsWereMet())
}
