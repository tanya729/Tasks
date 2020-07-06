package hw6

import (
	"reflect"
	"testing"
)

func prepareList() *List {
	list := new(List)
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)
	list.PushFront(0)
	return list
}

func TestList_Len(t *testing.T) {
	list := prepareList()
	count := list.Len()
	if count != 4 {
		t.Error("Expected 4, got ", count)
	}
}

func TestList_First(t *testing.T) {
	list := prepareList()
	first := list.First()
	if first.Value() != 0 {
		t.Error("Expected 0, got ", first.data)
	}
}

func TestList_Last(t *testing.T) {
	list := prepareList()
	last := list.Last()
	if last.Value() != 3 {
		t.Error("Expected 3, got ", last.data)
	}
}

func TestList_PushBack(t *testing.T) {
	list := prepareList()
	testData := "some data"
	list.PushBack(testData)
	last := list.Last()
	if last.Value() != testData {
		t.Errorf("Expected %s, got %s", testData, last.data)
	}
}

func TestList_PushFront(t *testing.T) {
	list := prepareList()
	testData := "some data"
	list.PushFront(testData)
	first := list.First()
	if first.Value() != testData {
		t.Errorf("Expected %s, got %s", testData, first.data)
	}
}

func TestList_PushFront2(t *testing.T) {
	list := new(List)
	testData := "some data"
	list.PushFront(testData)
	first := list.First()
	if first.Value() != testData {
		t.Errorf("Expected %s, got %s", testData, first.data)
	}
}

func TestList_Remove(t *testing.T) {
	list := prepareList()
	testData := "some data"
	removed := list.PushBack(testData)
	testDataNext := "some data next"
	list.PushFront(testDataNext)
	list.Remove(removed)
	for i := list.First(); i != nil; i = i.Next() {
		if i.data == testData {
			t.Errorf("Got %s from removed item", testData)
		}
	}
}

func TestList_Remove2(t *testing.T) {
	list := prepareList()
	checker := list.PushBack("checker data")
	testData := "some data"
	removed := list.PushBack(testData)
	testDataNext := "some data next"
	list.PushBack(testDataNext)
	list.Remove(removed)
	if checker.next.Value() != testDataNext {
		t.Errorf("Expected %s, got %s", testDataNext, checker.next.Value())
	}
}

func TestList_Double_Remove(t *testing.T) {
	list := prepareList()
	testData := "some data"
	removed := list.PushBack(testData)
	testDataNext := "some data next"
	list.PushBack(testDataNext)
	list.Remove(removed)
	list.Remove(removed)
	if list.Len() != 5 {
		t.Errorf("Expected %d, got %d", 5, list.Len())
	}
}

func TestList_Remove_First(t *testing.T) {
	list := prepareList()
	first := list.First()
	firstValue := first.Value()
	list.Remove(first)
	if list.First().Value() == firstValue {
		t.Errorf("Not expected %s", firstValue)
	}
}

func TestList_Remove_All(t *testing.T) {
	list := prepareList()
	length := list.Len()
	for i := 0; i < length; i++ {
		item := list.First()
		list.Remove(item)
	}

	if list.Len() != 0 {
		t.Errorf("Not expected %d length", list.Len())
	}
	if list.First() != nil {
		t.Errorf("Not expected %v as first item", list.First())
	}
	if list.Last() != nil {
		t.Errorf("Not expected %v as last item", list.Last())
	}
}

func TestItem_Value(t *testing.T) {
	testData := "some data"
	item := Item{data: testData}
	if item.Value() != testData {
		t.Errorf("Expected %s, got %s", testData, item.Value())
	}
}

func TestItem_Prev(t *testing.T) {
	list := new(List)
	first := list.PushBack(1)
	second := list.PushBack(2)
	if !reflect.DeepEqual(first, second.Prev()) {
		t.Errorf("Expected %v, got %v", first, second.Prev())
	}
}

func TestItem_Next(t *testing.T) {
	list := new(List)
	first := list.PushBack(1)
	second := list.PushBack(2)
	if !reflect.DeepEqual(first.Next(), second) {
		t.Errorf("Expected %v, got %v", second, first.Next())
	}
}
