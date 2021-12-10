package data_structure

import (
	"JumboDB/jumboDB-core/src/config"
	"bufio"
	"github.com/bits-and-blooms/bloom/v3"
	"io"
	"log"
	"sort"
)

type MergeRow struct {
	Operation *Operation
	Index     int
}

func multiWayMerge(tableList []*Table, fileName string, fileNumber int, config *config.TomlConfig) *Table {
	// use naive solution for now
	writer := OpenFileWithWriter(fileName)
	var scanners []*bufio.Reader
	for _, table := range tableList {
		scanners = append(scanners, OpenFileWithReader(table.FilePath))
	}
	bloomFilter := bloom.NewWithEstimates(uint(config.Storage.MemoryTableSize*fileNumber), config.Storage.BloomFilterFalsePositiveRate)
	var unCompleteReadFileNumber = fileNumber
	if unCompleteReadFileNumber == 0 {
		log.Printf("Invalid file number in level, something wrong in multiway merge")
		return nil
	}
	// use map and ordered slice to simulate the min heap
	//mergeMap := make(map[string]*MergeRow, unCompleteReadFileNumber)
	mergeSlice := make([]*MergeRow, 0, unCompleteReadFileNumber)

	// put all data into map
	for idx, scanner := range scanners {
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				unCompleteReadFileNumber -= 1
			} else {
				log.Fatalln(err)
			}
		}
		row := NewMergeRow(line, idx)
		//mergeMap[row.Operation.Key] = row
		mergeSlice = append(mergeSlice, row)
	}

	for len(mergeSlice) > 0 {
		// first sort the slice by key
		sort.Sort(MergeRowComparer(mergeSlice))
		//sort.Strings(mergeSlice)

		// get smallest row
		//currentWriteRow := mergeMap[mergeSlice[0]]
		currentWriteRow := mergeSlice[0]
		bfKey := []byte(currentWriteRow.Operation.Key)
		if !bloomFilter.Test(bfKey) {
			bloomFilter.Add([]byte(currentWriteRow.Operation.Key))
			writer.Write(currentWriteRow.Operation.OperationToJson())
			writer.WriteString("\n")
		}

		// remove the row from map and slice
		//delete(mergeMap, mergeSlice[0])
		mergeSlice = mergeSlice[1:]

		if unCompleteReadFileNumber > 0 {
			line, err := scanners[currentWriteRow.Index].ReadString('\n')
			if err != nil {
				if err == io.EOF {
					unCompleteReadFileNumber -= 1
					continue
				} else {
					log.Fatalln(err)
				}
			}
			newRow := NewMergeRow(line, currentWriteRow.Index)
			//mergeMap[newRow.Operation.Key] = newRow
			mergeSlice = append(mergeSlice, newRow)
		}

	}
	writer.Flush()
	return NewTable(fileName, bloomFilter)
}

func NewMergeRow(line string, index int) *MergeRow {
	opt := NewOperationFromJson(line)
	row := new(MergeRow)
	row.Operation = opt
	row.Index = index
	return row
}

type MergeRowComparer []*MergeRow

func (r MergeRowComparer) Len() int { return len(r) }
func (r MergeRowComparer) Less(i, j int) bool {
	if r[i].Operation.Key != r[j].Operation.Key {
		return r[i].Operation.Key < r[j].Operation.Key
	}
	return r[i].Index < r[j].Index
}
func (r MergeRowComparer) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
