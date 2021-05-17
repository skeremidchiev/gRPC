package mapstorage

import "testing"

func TestStore(t *testing.T) {
	ns := NewStorage()
	value := "AyRTjlxZJj"

	ns.Store(value)

	b := ns.exists(value)
	if b != true {
		t.Error("Exists failed!")
	}

	v, err := ns.GetRandom()
	if err != nil {
		t.Error("Store failed with: ", err)
	}
	if v != value {
		t.Error("Invalid value")
	}
}

func TestDelete(t *testing.T) {
	ns := NewStorage()
	value := "AyRTjlxZJj"

	ns.Store(value)

	b := ns.exists(value)
	if b != true {
		t.Error("Value should exist!")
	}

	err := ns.Delete(value)
	if err != nil {
		t.Error("Store failed with: ", err)
	}

	b = ns.exists(value)
	if b != false {
		t.Error("Value should not exist!")
	}
}

func TestGetAll(t *testing.T) {
	ns := NewStorage()

	values := []string{"AyRTjlxZJj", "SpGMoDuAFK"}

	for i := 0; i < len(values); i++ {
		ns.Store(values[i])
	}

	vs, err := ns.GetAll()
	if err != nil {
		t.Error("GetAll failed with: ", err)
	}

	for i := 0; i < len(values); i++ {
		if vs[i] != values[i] {
			t.Error("Mismatching values")
			return
		}
	}
}

func TestStoreErrorCase(t *testing.T) {
	ns := NewStorage()
	value := "AyRTjlxZJj"

	ns.Store(value)
	err := ns.Store(value)
	if err == nil {
		t.Error("Store of the same element should fail!")
	}
}

func TestDeleteErrorCase(t *testing.T) {
	ns := NewStorage()
	value := "AyRTjlxZJj"

	err := ns.Delete(value)
	if err == nil {
		t.Error("Delete of non existing element should fail!")
	}
}

func TestGetAllErrorCase(t *testing.T) {
	ns := NewStorage()

	_, err := ns.GetAll()
	if err == nil {
		t.Error("GetAll should fail if storage is empty!")
	}
}

func TestGetRandomErrorCase(t *testing.T) {
	ns := NewStorage()

	_, err := ns.GetRandom()
	if err == nil {
		t.Error("GetRandom should fail if storage is empty!")
	}
}
