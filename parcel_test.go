package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())

	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	get, err := store.Get(id)
	require.NoError(t, err)
	get.Number = parcel.Number
	assert.Equal(t, get, parcel)

	err = store.Delete(id)
	require.NoError(t, err)
	deleteTest, err := store.Get(id)
	require.Error(t, err)
	require.Empty(t, deleteTest)
	defer db.Close()
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)
	check, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, check.Address)
	defer db.Close()
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	checkStatus, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, checkStatus.Status, ParcelStatusSent)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		assert.NotEmpty(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)

	assert.Equal(t, len(parcels), len(storedParcels))

	for _, parcel := range storedParcels {

		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
