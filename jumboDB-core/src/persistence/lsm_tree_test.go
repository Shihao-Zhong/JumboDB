package persistence

import (
	"JumboDB/jumboDB-core/src/config"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLSMTree(t *testing.T) {
	// first create three tables in disk

	config.Path = "test_config.toml"
	lsmTree := NewLSMTree()

	lsmTree.Put("k1", "v1")
	lsmTree.Put("k2", "v1")
	lsmTree.Put("k3", "v1")
	time.Sleep(2 * time.Second)

	lsmTree.Put("k1", "v2")
	lsmTree.Put("k4", "v1")
	lsmTree.Del("k2")
	time.Sleep(2 * time.Second)

	lsmTree.Put("k5", "v3")
	lsmTree.Put("k2", "v3")
	lsmTree.Put("k4", "v2")

	time.Sleep(2 * time.Second)

	if lsmTree.Get("k1") != "v2" {
		t.Errorf("Error in compaction, loss some data with key [%s]", "k1")
	}

	if lsmTree.Get("k2") != "v3" {
		t.Errorf("Error in compaction, loss some data with key [%s]", "k2")
	}

	if lsmTree.Get("k3") != "v1" {
		t.Errorf("Error in compaction, loss some data with key [%s]", "k3")
	}

	if lsmTree.Get("k4") != "v2" {
		t.Errorf("Error in compaction, loss some data with key [%s]", "k4")
	}

	if lsmTree.Get("k5") != "v3" {
		t.Errorf("Error in compaction, loss some data with key [%s]", "k5")
	}

	d, _ := os.Open(".")
	defer d.Close()
	files, _ := d.Readdir(-1)
	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".jumbo" {
				os.Remove(file.Name())
				fmt.Println("Deleted ", file.Name())
			}
		}
	}
	os.Remove("ssTableIndex.index")
	os.Remove("lsm.index")
}
