package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	OP_FORWARD    = '>'
	OP_BACK       = '<'
	OP_INCREMENT  = '+'
	OP_DECREMENT  = '-'
	OP_PRINT      = '.'
	OP_INPUT      = ','
	OP_LOOP_BEGIN = '['
	OP_LOOP_END   = ']'
)

var (
	memCells    uint64 = 30000
	memData     []uint8
	memPointer  uint32
	progData    []uint8
	progOps     uint32 = 0
	progPointer uint32 = 0
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:")
		fmt.Println("   bfuck [file name]")
		return
	} else if os.Args[1] == "-h" {
		fmt.Println("Go BrainFuck interpreter v0.1.1")
		fmt.Println("   Input is buffered. Input must end with LF (10).")
		fmt.Println("   Input \"321\" gives four values: 51, 50, 49 and 10.")
		fmt.Println("   Data is stored as an array of uint8.")
		fmt.Println("  ", memCells, " memory cells available. Using memory beyound limit will cause error.")
		fmt.Println("   Characters other than standard BrainFuck language commands are ignored.")
		return
	}

	b, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		fmt.Println("Error. File does not exist!")
		return
	}

	progData = make([]uint8, 256)
	for _, char := range b {
		if progOps == uint32(len(progData)) {
			progData = append(progData, make([]uint8, 128)...)
		}
		if char == OP_FORWARD || char == OP_BACK ||
			char == OP_INCREMENT || char == OP_DECREMENT ||
			char == OP_PRINT || char == OP_INPUT ||
			char == OP_LOOP_BEGIN || char == OP_LOOP_END {
			progData[progOps] = char
			progOps += 1
		}
	}

	memData = make([]uint8, memCells)
	stdin := bufio.NewReader(os.Stdin)

	for progPointer < progOps {
		switch progData[progPointer] {
		case OP_INCREMENT:
			memData[memPointer] += 1
		case OP_DECREMENT:
			memData[memPointer] -= 1
		case OP_FORWARD:
			memPointer += 1
		case OP_BACK:
			memPointer -= 1
		case OP_PRINT:
			fmt.Printf("%c", memData[memPointer])
		case OP_INPUT:
			memData[memPointer], _ = stdin.ReadByte()
		case OP_LOOP_BEGIN:
			if memData[memPointer] == 0 {
				var nested uint32 = 1
				for nested > 0 {
					progPointer += 1
					if progData[progPointer] == OP_LOOP_BEGIN {
						nested += 1
					} else if progData[progPointer] == OP_LOOP_END {
						nested -= 1
					}
				}
			}
		case OP_LOOP_END:
			if memData[memPointer] != 0 {
				var nested uint32 = 1
				for nested > 0 {
					progPointer -= 1
					if progData[progPointer] == OP_LOOP_BEGIN {
						nested -= 1
					} else if progData[progPointer] == OP_LOOP_END {
						nested += 1
					}
				}
			}
		}
		progPointer += 1
	}
}
