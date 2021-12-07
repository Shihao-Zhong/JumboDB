package data_structure

import (
	"os"
	"testing"
)

/*
Unit test for brand new SSTable constructor
*/
func TestGetLevelSSTable(t *testing.T) {
	path := "test_get_sstable_1.config"
	err := os.Remove(path)
	ssTable := GetLevelSSTable(path)
	if ssTable.TotalLevel != 0 {
		t.Errorf("Error in total level of sstable, not matched the argument")
	}
	if ssTable.Sequence != 0 {
		t.Errorf("Error in sequence of sstable, the initial sequence should be 0")
	}
	err = os.Remove(path)
	if err != nil {
		t.Errorf("Error in remove test file, may need manually fix")
	}

}

/*
Unit test for exist new SSTable constructor
*/
func TestExistGetLevelSSTable(t *testing.T) {
	path := "test_get_sstable_2.config"
	err := os.Remove(path)
	ssTable := GetLevelSSTable(path)
	if ssTable.TotalLevel != 0 {
		t.Errorf("Error in total level of sstable, not matched the argument")
	}
	if ssTable.Sequence != 0 {
		t.Errorf("Error in sequence of sstable, the initial sequence should be 0")
	}
	ssTable.TotalLevel = 10
	ssTable.ConfigToJsonFile()
	ssTable = GetLevelSSTable(path)
	if ssTable.TotalLevel != 10 {
		t.Errorf("Error in total level of sstable, the new config is not persist")
	}
	err = os.Remove(path)
	if err != nil {
		t.Errorf("Error in remove test file, may need manually fix")
	}
}
