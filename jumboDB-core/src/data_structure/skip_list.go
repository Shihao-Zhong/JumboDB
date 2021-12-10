package data_structure

import (
	"JumboDB/jumboDB-core/src/protocol"
	"errors"
	"log"
	"math/rand"
	"time"
)

type SkipList struct {
	Data       *SkipListNode
	Size       int
	MaxLevel   int
	CurLevel   int
	GlobalLock bool
}

type SkipListNode struct {
	Key   string
	Value *Operation
	Next  []*SkipListNode
	Type  string
}

const (
	HEAD string = "HEAD"
	TAIL string = "TAIL"
	DATA string = "DATA"
)

func newLevel() bool {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(100) < 50
}

func max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func (i *SkipList) Put(value *Operation) {
	// Create new node
	node := NewSkipListNode(value.Key, value)

	// First calculate how many level index
	level := 0
	for newLevel() && level < i.MaxLevel-1 {
		level += 1
	}

	i.CurLevel = max(i.CurLevel, level)
	//log.Printf("Incoming node with key [%s] level [%d]", key, level)
	node.Next = make([]*SkipListNode, i.MaxLevel)
	// add the counter and insert node
	i.Size += i.Data.Put(node, level, i.CurLevel)
	//log.Printf("start put with level [%d]", level)

	//log.Printf("Size of the skipList = [%d]", i.Size)
}

func (i *SkipListNode) Put(node *SkipListNode, level int, currentLevel int) int {
	//log.Printf("Put with key [%s] level [%d] cur [%s] currentLevel [%d]", node.Key, level, i.Key, currentLevel)

	// if exist then replace node
	if i.Next[currentLevel].Key == node.Key {
		node.Next[currentLevel] = i.Next[currentLevel].Next[currentLevel]
		i.Next[currentLevel] = node
		if currentLevel > 0 {
			i.Put(node, level, currentLevel-1)
		}
		return 0
	} else if i.Next[currentLevel].Type == TAIL || compareSmaller(node.Key, i.Next[currentLevel].Key) {
		// if input is smaller than next node or the last one, create new node
		if level >= currentLevel {
			prevNext := i.Next[currentLevel]
			node.Next[currentLevel] = prevNext
			i.Next[currentLevel] = node
		}

		if currentLevel > 0 {
			i.Put(node, level, currentLevel-1)
		}
	} else {
		i.Next[currentLevel].Put(node, level, currentLevel)
	}
	return 1
}

func (i *SkipList) Get(key string) (*Operation, error) {
	//log.Printf("Start get with key [%s] and level[%d]", key, i.CurLevel-1)
	return i.Data.Get(key, i.CurLevel)
}

func (i *SkipListNode) Get(key string, level int) (*Operation, error) {
	//log.Printf("Get with key [%s] level [%d] in node [%s]", key, level, i.Key)
	if i.Key == key {
		if i.Value.Operation != DEL {
			return i.Value, nil
		}
		return nil, errors.New("element not exist")
	}
	// if reach to the last one or next is greater than key go next level

	//log.Printf("current = %v next = %v", i, i.Next)
	if i.Next[level].Type == TAIL || compareSmaller(key, i.Next[level].Key) {
		if level == 0 {
			return nil, errors.New("element not exist")
		} else {
			return i.Get(key, level-1)
		}
	}
	return i.Next[level].Get(key, level)
}

func (i *SkipList) Print() {
	log.Printf("\nCurrent skiplist status: \n")
	for currentLevel := i.MaxLevel - 1; currentLevel >= 0; currentLevel-- {
		i.Data.PrintLevel(currentLevel)
	}
}

func (i *SkipListNode) PrintLevel(level int) {
	if i == nil {
		log.Println()
		return
	}
	log.Printf("{key = [%s], value = [%s]}", i.Key, i.Value)
	i.Next[level].PrintLevel(level)
}

func (i *SkipList) GetDataKeySlice() []string {
	var data []string
	cur := i.Data.Next[0]
	for cur.Type != TAIL {
		//log.Printf("key = %s cnt = %d size = %d\n",cur.Key, cnt, i.Size)
		data = append(data, cur.Key)
		cur = cur.Next[0]
	}
	return data
}

func (i *SkipList) GetAll() []protocol.Payload {
	var data []protocol.Payload
	//data := make([]protocol.Payload, i.Size)
	cur := i.Data.Next[0]
	for cur.Type != TAIL {
		if cur.Value.Operation != DEL {
			data = append(data, *protocol.NewPayload(cur.Key, cur.Value.Value, cur.Value.TransactionId))
		}
		cur = cur.Next[0]
	}
	return data
}

func (i *SkipList) Del(key string, transactionId int) {
	delOperation := NewOperation(key, "", DEL, transactionId)
	i.Put(delOperation)
}

func compareSmaller(s1 string, s2 string) bool {
	return s1 < s2
}

func NewSkipListNode(key string, value *Operation) *SkipListNode {
	node := new(SkipListNode)
	node.Key = key
	node.Value = value
	node.Next = []*SkipListNode{}
	node.Type = DATA
	return node
}

func NewSkipList(maxLevel int) *SkipList {
	// need remove
	skipList := new(SkipList)
	skipList.Size = 0
	head := NewSkipListNode("<HEAD>", NewOperation("", "", "", 0))
	head.Type = HEAD
	head.Next = make([]*SkipListNode, maxLevel)
	tail := NewSkipListNode("<TAIL>", NewOperation("", "", "", 0))
	tail.Type = TAIL
	tail.Next = make([]*SkipListNode, maxLevel)
	for i := range head.Next {
		head.Next[i] = tail
	}
	skipList.Data = head
	skipList.MaxLevel = maxLevel
	skipList.GlobalLock = false
	return skipList
}

func (i *SkipList) toFile(path string) int {
	line := 0
	writer := OpenFileWithWriter(path)
	cur := i.Data.Next[0]
	for cur.Type != TAIL {
		log.Printf("cur = [%s]", cur.Value.toString())
		writer.Write(cur.Value.OperationToJson())
		writer.WriteString("\n")
		cur = cur.Next[0]
		line += 1
	}
	writer.Flush()
	return line
}
