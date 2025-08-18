package store

import (
	"errors"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
)

func TestStore_Put(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	originalIsUnlucky := isUnlucky
	isUnlucky = func() bool { return false }
	defer func() { isUnlucky = originalIsUnlucky }()

	mock.ExpectExec("INSERT INTO store").
		WithArgs("foo", "bar").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	store := NewStore(mock)

	err = store.Put("foo", "bar")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStore_Put_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	originalIsUnlucky := isUnlucky
	isUnlucky = func() bool { return false }
	defer func() { isUnlucky = originalIsUnlucky }()

	mock.ExpectExec("INSERT INTO store").
		WithArgs("foo", "bar").
		WillReturnError(errors.New("error"))

	store := NewStore(mock)

	err = store.Put("foo", "bar")
	require.EqualError(t, err, "error")
}

func TestStore_Put_Unlucky(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	isUnlucky = func() bool { return true }
	defer func() { isUnlucky = func() bool { return false } }()

	store := NewStore(mock)

	err = store.Put("foo", "bar")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStore_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	originalIsUnlucky := isUnlucky
	isUnlucky = func() bool { return false }
	defer func() { isUnlucky = originalIsUnlucky }()

	mock.ExpectExec("DELETE FROM store").
		WithArgs("foo").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	store := NewStore(mock)
	err = store.Delete("foo")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStore_Delete_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	originalIsUnlucky := isUnlucky
	isUnlucky = func() bool { return false }
	defer func() { isUnlucky = originalIsUnlucky }()

	mock.ExpectExec("DELETE FROM store").
		WithArgs("foo").
		WillReturnError(errors.New("error"))

	store := NewStore(mock)
	err = store.Delete("foo")
	require.EqualError(t, err, "error")
}

func TestStore_Delete_NoSuchKey(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	originalIsUnlucky := isUnlucky
	isUnlucky = func() bool { return false }
	defer func() { isUnlucky = originalIsUnlucky }()

	mock.ExpectExec("DELETE FROM store").
		WithArgs("foo").
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	store := NewStore(mock)
	err = store.Delete("foo")

	require.Error(t, err)
	require.EqualError(t, err, "No such key: foo")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStore_Get(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	originalIsUnlucky := isUnlucky
	isUnlucky = func() bool { return false }
	defer func() { isUnlucky = originalIsUnlucky }()

	rows := pgxmock.NewRows([]string{"value"}).AddRow("bar")

	mock.ExpectQuery("SELECT value FROM store WHERE key = \\$1").
		WithArgs("foo").
		WillReturnRows(rows)

	store := NewStore(mock)
	result, err := store.Get("foo")
	require.NoError(t, err)
	require.Equal(t, "bar", result)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStore_Get_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	originalIsUnlucky := isUnlucky
	isUnlucky = func() bool { return false }
	defer func() { isUnlucky = originalIsUnlucky }()

	mock.ExpectQuery("SELECT value FROM store WHERE key = \\$1").
		WithArgs("foo").
		WillReturnError(errors.New("error"))

	store := NewStore(mock)
	_, err = store.Get("foo")
	require.Error(t, err)
	require.EqualError(t, err, "error")
}

func TestStore_Dumb(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	originalIsUnlucky := isUnlucky
	isUnlucky = func() bool { return false }
	defer func() { isUnlucky = originalIsUnlucky }()

	rows := pgxmock.NewRows([]string{"key", "value"}).
		AddRow("foo", "bar").
		AddRow("baz", "qux")

	mock.ExpectQuery("SELECT key, value FROM store").
		WillReturnRows(rows)

	store := NewStore(mock)
	result, err := store.Dumb()
	require.NoError(t, err)

	expected := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}
	require.Equal(t, expected, result)
	require.NoError(t, mock.ExpectationsWereMet())

}

func TestStore_Dumb_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	originalIsUnlucky := isUnlucky
	isUnlucky = func() bool { return false }
	defer func() { isUnlucky = originalIsUnlucky }()

	mock.ExpectQuery("SELECT key, value FROM store").
		WillReturnError(errors.New("error"))

	store := NewStore(mock)
	_, err = store.Dumb()

	require.EqualError(t, err, "query failed: error")
	require.Error(t, err)
}

func TestStore_Dumb_ScanError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	originalIsUnlucky := isUnlucky
	isUnlucky = func() bool { return false }
	defer func() { isUnlucky = originalIsUnlucky }()

	rows := pgxmock.NewRows([]string{"wrong"}).
		AddRow("oops")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT key, value FROM store")).
		WillReturnRows(rows)

	store := NewStore(mock)
	_, err = store.Dumb()

	require.Error(t, err)
	require.Contains(t, err.Error(), "error scanning row")
	require.NoError(t, mock.ExpectationsWereMet())
}
