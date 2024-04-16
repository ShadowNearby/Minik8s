package storage

import (
	"context"
	"errors"
	"testing"
)

type MyStruct struct {
	Field1 string
	Field2 int
}

func TestStorage(t *testing.T) {
	// test create
	etcdStorage := CreateEtcdStorage([]string{"localhost:2380"})
	var key1 = "key1"
	var key2 = "key2"

	var expectedErr error = nil
	myStruct := &MyStruct{Field1: "Hello", Field2: 42}
	err := etcdStorage.Put(context.Background(), key1, &myStruct)

	if !errors.Is(expectedErr, err) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	myStruct5 := &MyStruct{Field1: "Hello2", Field2: 43}
	err = etcdStorage.Put(context.Background(), key2, &myStruct5)
	if !errors.Is(expectedErr, err) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	// test get
	var myStruct2 MyStruct
	err = etcdStorage.Get(context.Background(), key1, &myStruct2)
	if !errors.Is(expectedErr, err) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if myStruct2.Field1 != myStruct.Field1 {
		t.Errorf("Expected %v, got %v", myStruct.Field1, myStruct2.Field1)
	}
	if myStruct2.Field2 != myStruct.Field2 {
		t.Errorf("Expected %v, got %v", myStruct.Field2, myStruct2.Field2)
	}

	// test getList
	//var myStructList []MyStruct
	//err = etcdStorage.GetList(context.Background(), "myStruct", &myStructList)
	//if err != expectedErr {
	//	t.Errorf("Expected error %v, got %v", expectedErr, err)
	//}
	//if len(myStructList) != 2 {
	//	for key, myStruct := range myStructList {
	//		fmt.Println(key, myStruct)
	//	}
	//	t.Errorf("Expected %v, got %v", 1, len(myStructList))
	//}
	//
	//// test update
	//myStruct.Field1 = "World"
	//err = etcdStorage.GuaranteedUpdate(context.Background(), "myStruct", &myStruct)
	//if err != expectedErr {
	//	t.Errorf("Expected error %v, got %v", expectedErr, err)
	//}

	// test get
	var myStruct3 MyStruct
	err = etcdStorage.Get(context.Background(), key1, &myStruct3)
	if !errors.Is(expectedErr, err) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if myStruct3.Field1 != myStruct.Field1 {
		t.Errorf("Expected %v, got %v", myStruct.Field1, myStruct3.Field1)
	}

	// test delete
	err = etcdStorage.Delete(context.Background(), key1)
	if !errors.Is(expectedErr, err) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	err = etcdStorage.Delete(context.Background(), key2)
	if !errors.Is(expectedErr, err) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	// test get
	var myStruct4 MyStruct
	err = etcdStorage.Get(context.Background(), key1, &myStruct4)
	if errors.Is(expectedErr, err) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}
