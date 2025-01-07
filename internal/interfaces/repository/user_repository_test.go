package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockDB(t *testing.T) (repository.UserRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := repository.NewUserRepository(db)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

func TestCreateUser_Succeess(t *testing.T) {
	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()
	user := &entity.User{
		LINEUserID:     "U123",
		SelectedAreaID: 0110000,
		NotifyTime:     time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Local),
		IsActive:       true,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
	    INSERT INTO users (line_user_id, selected_area_id, notify_time, is_active, created_at, updated_at)
	    VALUES ($1, $2, $3, $4, $5, $6)
	    RETURNING id
	`)).
		WithArgs(user.LINEUserID, user.SelectedAreaID, user.NotifyTime, user.IsActive, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	created, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)
	assert.Equal(t, 1, created.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_QueryError(t *testing.T) {
	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	ctx := context.Background()
	user := &entity.User{
		LINEUserID:     "U123",
		SelectedAreaID: 0110000,
		NotifyTime:     time.Now(),
		IsActive:       true,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO users (line_user_id, selected_area_id, notify_time, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`)).
		WithArgs(user.LINEUserID, user.SelectedAreaID, user.NotifyTime, user.IsActive, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("insert failed"))

	_, err := repo.CreateUser(ctx, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "faild to insert user")
}

func TestFindUserByID_Success(t *testing.T) {
	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	ctx := context.Background()
	query := `
		SELECT
            id, line_user_id, selected_area_id, notify_time,
            is_active, created_at, updated_at
        FROM users
        WHERE id = $1
        LIMIT 1
	`
	rows := sqlmock.NewRows([]string{
		"id", "line_user_id", "selected_area_id", "notify_time", "is_active", "created_at", "updated_at",
	}).AddRow(1, "U123", 0110000, time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.FindUserByID(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
}

func TestFindUserByID_NotFound(t *testing.T) {
	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	ctx := context.Background()
	query := `
		SELECT
            id, line_user_id, selected_area_id, notify_time,
            is_active, created_at, updated_at
        FROM users
        WHERE id = $1
        LIMIT 1
	`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.FindUserByID(ctx, 999)
	require.NoError(t, err)
	assert.Nil(t, user)
}

func TestFindUserByLINEUserID_Success(t *testing.T) {
	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	ctx := context.Background()
	query := `
		SELECT
            id, line_user_id, selected_area_id, notify_time,
            is_active, created_at, updated_at
        FROM users
        WHERE line_user_id = $1
        LIMIT 1
	`
	rows := sqlmock.NewRows([]string{
		"id", "line_user_id", "selected_area_id", "notify_time", "is_active", "created_at", "updated_at",
	}).AddRow(1, "U123", 0110000, time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("U123").
		WillReturnRows(rows)

	user, err := repo.FindUserByLINEUserID(ctx, "U123")
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
}

func TestFindUserByLINEUserID_NotFound(t *testing.T) {
	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	ctx := context.Background()
	query := `
		SELECT
            id, line_user_id, selected_area_id, notify_time,
            is_active, created_at, updated_at
        FROM users
        WHERE line_user_id = $1
        LIMIT 1
	`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("U456").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.FindUserByLINEUserID(ctx, "U456")
	require.NoError(t, err)
	assert.Nil(t, user)
}

func TestUpdateUser_Success(t *testing.T) {
	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	ctx := context.Background()
	user := &entity.User{
		ID:             1,
		SelectedAreaID: 0120200,
		NotifyTime:     time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
		IsActive:       false,
	}

	query := `
		UPDATE users
		SET
			selected_area_id = $1,
			notify_time = $2,
			is_active = $3,
			updated_at = $4
		WHERE id = $5
	`
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(user.SelectedAreaID, user.NotifyTime, user.IsActive, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateUser(ctx, user)
	require.NoError(t, err)
}

func TestUpdateUser_NoRows(t *testing.T) {
	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	ctx := context.Background()
	user := &entity.User{
		ID:             999,
		SelectedAreaID: 0120200,
		NotifyTime:     time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
		IsActive:       false,
	}

	query := `
		UPDATE users
		SET
			selected_area_id = $1,
			notify_time = $2,
			is_active = $3,
			updated_at = $4
		WHERE id = $5
	`
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(user.SelectedAreaID, user.NotifyTime, user.IsActive, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // no rows affected

	err := repo.UpdateUser(ctx, user)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no rows updated")
}

func TestDeleteUser_Success(t *testing.T) {
	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	ctx := context.Background()
	query := `DELETE FROM users WHERE id = $1`
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteUser(ctx, 1)
	require.NoError(t, err)
}

func TestDeleteUser_NoRows(t *testing.T) {
	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	ctx := context.Background()
	query := `DELETE FROM users WHERE id = $1`
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0)) // no rows deleted

	err := repo.DeleteUser(ctx, 999)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no rows deleted")
}
