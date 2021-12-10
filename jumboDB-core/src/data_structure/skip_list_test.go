package data_structure

import (
	"JumboDB/jumboDB-core/src/protocol"
	"github.com/google/go-cmp/cmp"
	"testing"
)

/*
Unit test for skipList constructor
*/
func TestNewSkipList(t *testing.T) {
	skipList := NewSkipList(10)
	if skipList.MaxLevel != 10 {
		t.Errorf("Error in maxlevel of skip list, not matched the argument")
	}
	if skipList.Size != 0 {
		t.Errorf("Error in size of skip list, the initial size should be 0")
	}
	if skipList.Data == nil {
		t.Errorf("Error in data of skip list, the initial data should contain head")
	}
	if skipList.Data.Type != HEAD {
		t.Errorf("Error in data of skip list, the initial data should have key head")
	}
}

/*
Unit test for skipList put
*/
func TestSkipListPut(t *testing.T) {
	skipList := NewSkipList(10)
	skipList.Put(NewOperation("9", "9", PUT, 1))
	skipList.Put(NewOperation("2", "2", PUT, 2))
	skipList.Put(NewOperation("5", "5", PUT, 3))
	sortedKeySlice := skipList.GetDataKeySlice()
	expectedKeySlice := []string{"2", "5", "9"}
	if !cmp.Equal(sortedKeySlice, expectedKeySlice) {
		t.Errorf("Error in put of skip list, the data structure is not sorted expected [%v] but got [%v]",
			expectedKeySlice, sortedKeySlice)
	}
	skipList.Put(NewOperation("d", "d", PUT,4 ))
	skipList.Put(NewOperation("e", "e", PUT, 5))
	skipList.Put(NewOperation("a", "a", PUT, 6))
	sortedKeySlice = skipList.GetDataKeySlice()
	expectedKeySlice = []string{"2", "5", "9","a","d","e"}
	if !cmp.Equal(sortedKeySlice, expectedKeySlice) {
		t.Errorf("Error in put of skip list, the data structure is not sorted expected [%v] but got [%v]",
			expectedKeySlice, sortedKeySlice)
	}
}

/*
Unit test for skipList get
*/
func TestSkipListGet(t *testing.T) {
	skipList := NewSkipList(10)
	skipList.Put(NewOperation("9", "9", PUT, 1))
	skipList.Put(NewOperation("2", "2", PUT, 2))
	skipList.Put(NewOperation("5", "5", PUT, 3))

	getHelper(t, skipList, "2", NewOperation("2", "2", PUT, 4))
	getHelper(t, skipList, "9", NewOperation("9", "9", PUT, 5))
	getHelper(t, skipList, "5", NewOperation("5", "5", PUT, 6))

	skipList.Put(NewOperation("d", "d", PUT, 7))
	skipList.Put(NewOperation("e", "e", PUT, 8))
	skipList.Put(NewOperation("a", "a", PUT, 9))

	getHelper(t, skipList, "d", NewOperation("d", "d", PUT, 10))
	getHelper(t, skipList, "e", NewOperation("e", "e", PUT, 11))
	getHelper(t, skipList, "a", NewOperation("a", "a", PUT, 12))
}

func getHelper(t *testing.T, skipList *SkipList, key string, expectedValue *Operation) {
	val, err := skipList.Get(key)
	if err != nil {
		t.Errorf("Error in get of skip list, expected [%s] get [error]", expectedValue)
		return
	}
	if !cmp.Equal(val, expectedValue) {
		t.Errorf("Error in get of skip list, expected [%s] get [%s]",expectedValue, val)
		return
	}
}

/*
Unit test for skipList del
*/
func TestSkipListDel(t *testing.T) {
	skipList := NewSkipList(10)
	skipList.Put(NewOperation("9", "9", PUT, 13))
	skipList.Put(NewOperation("2", "2", PUT, 14))
	skipList.Put(NewOperation("5", "5", PUT, 15))

	getHelper(t, skipList, "2", NewOperation("2", "2", PUT, 16))
	skipList.Del("2", 17)

	getErrorHelper(t, skipList, "2")

	getHelper(t, skipList, "9", NewOperation("9", "9", PUT, 18))
	getHelper(t, skipList, "5", NewOperation("5", "5", PUT, 19))

	skipList.Put(NewOperation("d", "d", PUT, 20))
	skipList.Put(NewOperation("e", "e", PUT, 21))
	skipList.Put(NewOperation("a", "a", PUT, 22))

	getHelper(t, skipList, "d", NewOperation("d", "d", PUT, 23))
	getHelper(t, skipList, "e", NewOperation("e", "e", PUT, 24))
	getHelper(t, skipList, "a", NewOperation("a", "a", PUT, 25))

	skipList.Del("d", 26)
	skipList.Del("9", 27)
	getErrorHelper(t, skipList, "d")
	getErrorHelper(t, skipList, "9")
}

func getErrorHelper(t *testing.T, skipList *SkipList, key string) {
	_, err := skipList.Get(key)
	if err == nil {
		t.Errorf("Error in get not exist data from skipList, it suppose get an err but did not")
		return
	}
	return
}

/*
Unit test for skipList getall
*/
func TestSkipListGetAll(t *testing.T) {
	skipList := NewSkipList(10)
	skipList.Put(NewOperation("9", "9", PUT, 1))
	skipList.Put(NewOperation("2", "2", PUT, 2))
	skipList.Put(NewOperation("5", "5", PUT,3 ))
	skipListDataSlice := skipList.GetAll()

	dataSlice := make([]protocol.Payload, 3)
	dataSlice[0] = *protocol.NewPayload("2", "2",4 )
	dataSlice[1] = *protocol.NewPayload("5", "5", 5)
	dataSlice[2] = *protocol.NewPayload("9", "9", 6)

	if !cmp.Equal(skipListDataSlice, dataSlice) {
		t.Errorf("Error in getall of skip list, the data structure is not sorted expected [%v] but got [%v]",
			dataSlice, skipListDataSlice)
	}
}

/*
Unit test for skipList mix
*/
func TestSkipListMix(t *testing.T) {
	skipList := NewSkipList(10)
	skipList.Put(NewOperation("9", "9", PUT, 7))
	skipList.Put(NewOperation("2", "2", PUT, 8))
	skipList.Put(NewOperation("5", "5", PUT, 9))

	skipList.Put(NewOperation("d", "d", PUT, 10))
	skipList.Put(NewOperation("e", "e", PUT, 11))
	skipList.Put(NewOperation("a", "a", PUT, 12))

	getHelper(t, skipList, "d", NewOperation("d", "d", PUT, 13))
	getHelper(t, skipList, "e", NewOperation("e", "e", PUT, 14))
	getHelper(t, skipList, "a", NewOperation("a", "a", PUT, 15))

	getHelper(t, skipList, "d", NewOperation("d", "d", PUT, 16))
	getErrorHelper(t, skipList, "200")
	getHelper(t, skipList, "e", NewOperation("e", "e", PUT,17))
	getErrorHelper(t, skipList, "304")
	getHelper(t, skipList, "a", NewOperation("a", "a", PUT, 18))
	getErrorHelper(t, skipList, "231")

	skipList.Del("a", 19)
	skipList.Del("e", 20)
	skipList.Del("d", 21)

	skipListDataSlice := skipList.GetAll()

	dataSlice := make([]protocol.Payload, 3)
	dataSlice[0] = *protocol.NewPayload("2", "2", 22)
	dataSlice[1] = *protocol.NewPayload("5", "5", 23)
	dataSlice[2] = *protocol.NewPayload("9", "9", 24)

	if !cmp.Equal(skipListDataSlice, dataSlice) {
		t.Errorf("Error in getall of skip list, the data structure is not sorted expected [%v] but got [%v]",
			dataSlice, skipListDataSlice)
	}
}
