package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	filePermissions = 0644
)

// create a struct that can hold all the types of Instruction pieces so that as the Instruction is read
// and  decoded I can then return the struct, read the type of Instruction and then use that to
// separate the Instructions from the data.
// Create an array of structs si I can sequentially store the structs to later print them out

type Instruction struct {
	rawInstruction string
	linevalue      uint32
	programCnt     int
	opcode         uint32
	op             string
	rd             uint8
	rn             uint8
	rm             uint8
	rt             uint8
	im             int32
	shamt          uint16
	shfcd          uint8
	field          uint32
	addr           uint16
	offset         int32
	data           int32
}

func Imm_to_32bit_converter(num uint32, bitsize uint) int32 {

	var negBitMask uint32
	var extendMask uint32

	if bitsize == 12 {
		negBitMask = 0x800 // figure out if 12 bit num is neg
		extendMask = 0xFFFFF000

	} else if bitsize == 16 {
		negBitMask = 0x8000 // figure out if 16 bit num is neg
		extendMask = 0xFFFF0000

	} else if bitsize == 19 {
		negBitMask = 0x40000 // figure out if 19 bit num is neg
		extendMask = 0xFFF80000

	} else if bitsize == 26 {
		negBitMask = 0x2000000 // figure out if 26 bit num is neg
		extendMask = 0xFC000000

	} else if bitsize == 32 {
		negBitMask = 0x10000000
		extendMask = 0x00000000

	} else {
		fmt.Println(" You ARE USING AN INVALID BIT LENGTH")
	}

	var snum int32
	snum = int32(num)
	if (negBitMask & num) > 0 { // is it?
		num = num | extendMask // if so extend with 1's
		num = num ^ 0xFFFFFFFF // 2s comp
		snum = int32(num + 1)
		snum = snum * -1 // add neg sign
	}
	return snum

}

func TakeApart(instructionSet []Instruction, outfile *os.File, programCnt *int) { // need tp pass in outfile!!

	numInstructions := 0
	dataElements := make([]Instruction, len(instructionSet))
	copy(dataElements, instructionSet)

	var Breaknow = false

	var programCntint int
	programCntint = *programCnt

	for idx, s := range instructionSet {

		fmt.Println(s.opcode)

		instructionSet[idx].programCnt = programCntint
		numInstructions = numInstructions + 1 // num instructions will be the count through the break instruction

		switch {

		case s.opcode == 0 && Breaknow == false:
			fmt.Println("A NOP was discovered")
			instructionSet[idx].op = "NOP"
			outputString := fmt.Sprintf("%.32s\t%d\t%s\n", instructionSet[idx].rawInstruction[0:32], instructionSet[idx].programCnt, instructionSet[idx].op)

			fmt.Println(outputString)
			fmt.Fprintf(outfile, outputString)

		case s.opcode == 2038 && Breaknow == false:
			if (s.linevalue & 0x1FFFFF) == 0x1EFFE7 {
				fmt.Println("The break just was discovered")
				instructionSet[idx].op = "BREAK"
				Breaknow = true // sets the break flag so you drop out of the instruction loop and then process the data lines
				//outputString := fmt.Sprintf("%.32s \t%d\t%s\n", instructionSet[idx].rawInstruction[0:32], instructionSet[idx].programCnt, instructionSet[idx].op)
				outputString := fmt.Sprintf("%.1s %.5s %.5s %.5s %.5s %.5s %.6s \t%d\t%s\n", instructionSet[idx].rawInstruction[0:1], instructionSet[idx].rawInstruction[1:6], instructionSet[idx].rawInstruction[6:11], instructionSet[idx].rawInstruction[11:16], instructionSet[idx].rawInstruction[16:21], instructionSet[idx].rawInstruction[21:26], instructionSet[idx].rawInstruction[26:32], instructionSet[idx].programCnt, instructionSet[idx].op)

				fmt.Println(outputString)
				fmt.Fprintf(outfile, outputString)
				fmt.Println("exiting the Instruction for loop")
				break // ensures that break will occur, not needed for other parts of case
			}

		case s.opcode >= 160 && s.opcode <= 191 && Breaknow == false:
			fmt.Println("Found an B type Instruction")
			instructionSet[idx].op = "B"
			// not sure I need to have int versus uint as input, def need it as output
			instructionSet[idx].offset = Imm_to_32bit_converter(s.linevalue&0x3FFFFFF, 26)
			fmt.Println(instructionSet[idx].offset)
			outputString := fmt.Sprintf("%.6s %.26s \t%d\t%s\t#%d\n", instructionSet[idx].rawInstruction[0:6], instructionSet[idx].rawInstruction[6:32], instructionSet[idx].programCnt, instructionSet[idx].op, instructionSet[idx].offset)

			fmt.Println(outputString)
			fmt.Fprintf(outfile, outputString)

		case s.opcode == 1112 || s.opcode == 1624 || s.opcode == 1104 || s.opcode == 1360 || s.opcode == 1872 || (s.opcode >= 1690 && s.opcode <= 1693) && Breaknow == false:
			fmt.Println("Found an R type or RShift(Logical Shift) type Instruction")
			if s.opcode == 1112 {
				instructionSet[idx].op = "ADD"
			} else if s.opcode == 1624 {
				instructionSet[idx].op = "SUB"
			} else if s.opcode == 1104 {
				instructionSet[idx].op = "AND"
			} else if s.opcode == 1360 {
				instructionSet[idx].op = "ORR"
			} else if s.opcode == 1872 {
				instructionSet[idx].op = "EOR"
			} else if s.opcode == 1690 {
				instructionSet[idx].op = "LSR"
				instructionSet[idx].shamt = uint16((s.linevalue >> 10) & 0x3F)
			} else if s.opcode == 1691 {
				instructionSet[idx].op = "LSL"
				instructionSet[idx].shamt = uint16((s.linevalue >> 10) & 0x3F)
			} else if s.opcode == 1692 {
				instructionSet[idx].op = "ASR"
				instructionSet[idx].shamt = uint16((s.linevalue >> 10) & 0x3F)
			}

			instructionSet[idx].rd = uint8(s.linevalue & 0x1F)
			instructionSet[idx].rn = uint8((s.linevalue >> 5) & 0x1F)
			instructionSet[idx].rm = uint8((s.linevalue >> 16) & 0x1F)

			if s.opcode == 1112 || s.opcode == 1624 || s.opcode == 1104 || s.opcode == 1360 || s.opcode == 1872 {
				outputString := fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s \t%d\t%s\tR%d, R%d, R%d\n", instructionSet[idx].rawInstruction[0:11], instructionSet[idx].rawInstruction[11:16], instructionSet[idx].rawInstruction[16:22], instructionSet[idx].rawInstruction[22:27], instructionSet[idx].rawInstruction[27:32], instructionSet[idx].programCnt, instructionSet[idx].op, instructionSet[idx].rd, instructionSet[idx].rn, instructionSet[idx].rm)

				fmt.Println(outputString)
				fmt.Fprintf(outfile, outputString)

			} else if s.opcode == 1690 || s.opcode == 1691 {
				outputString := fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s \t%d\t%s\tR%d, R%d, #%d\n", instructionSet[idx].rawInstruction[0:11], instructionSet[idx].rawInstruction[11:16], instructionSet[idx].rawInstruction[16:22], instructionSet[idx].rawInstruction[22:27], instructionSet[idx].rawInstruction[27:32], instructionSet[idx].programCnt, instructionSet[idx].op, instructionSet[idx].rd, instructionSet[idx].rn, instructionSet[idx].shamt)

				fmt.Println(outputString)
				fmt.Fprintf(outfile, outputString)
			}

		case s.opcode == 1160 || s.opcode == 1161 || s.opcode == 1672 || s.opcode == 1673 && Breaknow == false:
			//case 1160, 1161, 1672, 1673:
			fmt.Println("Found an I type instruction")
			instructionSet[idx].rd = uint8(s.linevalue & 0x1F)
			instructionSet[idx].rn = uint8((s.linevalue >> 5) & 0x1F)
			//instructionSet[idx].im = Imm_to_32bit_converter((s.linevalue>>10)&0xFFF, 12)
			instructionSet[idx].im = int32(s.linevalue >> 10 & 0xFFF)
			if s.opcode == 1160 || s.opcode == 1161 {
				instructionSet[idx].op = "ADDI"
			} else if s.opcode == 1672 || s.opcode == 1673 {
				instructionSet[idx].op = "SUBI"
			}

			outputString := fmt.Sprintf("%.10s %.12s %.5s %.5s \t%d\t%s\tR%d, R%d, #%d\n", instructionSet[idx].rawInstruction[0:10], instructionSet[idx].rawInstruction[10:22], instructionSet[idx].rawInstruction[22:27], instructionSet[idx].rawInstruction[27:32], instructionSet[idx].programCnt, instructionSet[idx].op, instructionSet[idx].rd, instructionSet[idx].rn, instructionSet[idx].im)
			fmt.Println(outputString)
			fmt.Fprintf(outfile, outputString)

		case s.opcode == 1984 || s.opcode == 1986 && Breaknow == false:
			//case 1984, 1986:
			fmt.Println("Found a D type instruction")
			instructionSet[idx].addr = uint16((s.linevalue >> 12) & 0x1FF)
			instructionSet[idx].rn = uint8((s.linevalue >> 5) & 0x1F)
			instructionSet[idx].rt = uint8(s.linevalue & 0x1F)
			if s.opcode == 1984 {
				instructionSet[idx].op = "STUR"
			} else if s.opcode == 1986 {
				instructionSet[idx].op = "LDUR"
			}

			//outputString := DInstr(instructionSet[idx])
			outputString := fmt.Sprintf("%.11s %.9s %.2s %.5s %.5s \t%d\t%s\tR%d, [R%d, #%d]\n", instructionSet[idx].rawInstruction[0:11], instructionSet[idx].rawInstruction[11:20], instructionSet[idx].rawInstruction[20:22], instructionSet[idx].rawInstruction[22:27], instructionSet[idx].rawInstruction[27:32], instructionSet[idx].programCnt, instructionSet[idx].op, instructionSet[idx].rt, instructionSet[idx].rn, instructionSet[idx].addr)

			fmt.Println(outputString)
			fmt.Fprintf(outfile, outputString)

		case (s.opcode >= 1684 && s.opcode <= 1687) || (s.opcode >= 1940 && s.opcode <= 1943) && Breaknow == false:
			//ase 1684, 1685, 1686, 1687, 1940, 1941, 1942, 1943:
			fmt.Println("Found an IM type Instruction")

			instructionSet[idx].field = s.linevalue >> 5 & 0xFFFF
			instructionSet[idx].rd = uint8(s.linevalue & 0x1F)

			index := uint8((s.linevalue >> 17) & 0x30)

			if s.opcode >= 1684 && s.opcode <= 1687 {
				instructionSet[idx].op = "MOVZ"
			}

			if s.opcode >= 1940 && s.opcode <= 1943 {
				instructionSet[idx].op = "MOVK"
			}

			switch index {

			case 0:
				instructionSet[idx].shfcd = 0

			case 1:
				instructionSet[idx].shfcd = 16

			case 2:
				instructionSet[idx].shfcd = 32

			case 3:
				instructionSet[idx].shfcd = 48

			}
			instructionSet[idx].shfcd = uint8((s.linevalue >> 17) & 0x30)

			outputString := fmt.Sprintf("%.9s %.2s %.16s %.5s \t%d\t%s\tR%d, %d, LSL %d\n", instructionSet[idx].rawInstruction[0:9], instructionSet[idx].rawInstruction[9:11], instructionSet[idx].rawInstruction[11:27], instructionSet[idx].rawInstruction[27:32], instructionSet[idx].programCnt, instructionSet[idx].op, instructionSet[idx].rd, instructionSet[idx].field, instructionSet[idx].shfcd)

			fmt.Println(outputString)
			fmt.Fprintf(outfile, outputString)

		case s.opcode >= 1440 && s.opcode <= 1447 && Breaknow == false:
			//case 1440,1441,1442,1443,1444,1445,1146,1447:
			fmt.Println("Found an CBZ type Instruction")
			instructionSet[idx].rd = uint8(s.linevalue & 0x1F)
			instructionSet[idx].offset = Imm_to_32bit_converter((s.linevalue>>5)&0x7FFFF, 19)

			if s.opcode >= 1440 && s.opcode <= 1447 {
				instructionSet[idx].op = "CBZ"
			}

			outputString := fmt.Sprintf("%.8s %.19s %.5s \t%d\t%s\tR%d, #%d\n", instructionSet[idx].rawInstruction[0:8], instructionSet[idx].rawInstruction[8:27], instructionSet[idx].rawInstruction[27:32], instructionSet[idx].programCnt, instructionSet[idx].op, instructionSet[idx].rd, instructionSet[idx].offset)
			fmt.Println(outputString)
			fmt.Fprintf(outfile, outputString)

		case s.opcode >= 1448 && s.opcode <= 1455 && Breaknow == false:
			//case 1448,1449,1450,1451,1452,1453,1454,1455:
			fmt.Println("Found an CBNZ type Instruction")
			instructionSet[idx].rd = uint8(s.linevalue & 0x1F)
			instructionSet[idx].offset = Imm_to_32bit_converter((s.linevalue>>5)&0x7FFFF, 19)
			if s.opcode >= 1448 && s.opcode <= 1455 {
				instructionSet[idx].op = "CBNZ"
			}

			outputString := fmt.Sprintf("%.8s %.19s %.5s \t%d\t%s\tR%d, #%d\n", instructionSet[idx].rawInstruction[0:8], instructionSet[idx].rawInstruction[8:27], instructionSet[idx].rawInstruction[27:32], instructionSet[idx].programCnt, instructionSet[idx].op, instructionSet[idx].rd, instructionSet[idx].offset)

			fmt.Println(outputString)
			fmt.Fprintf(outfile, outputString)

		}

		dataElements = dataElements[1:] // each time the cycle runs it drops the instruction that just was handled
		// when it gets to the break, the only instructions left will be data after the break
		programCntint = programCntint + 4
		if Breaknow == true {
			break
		}

	}

	for idx, d := range dataElements {

		fmt.Println("Found data")
		instructionSet[idx].op = "DATA"
		instructionSet[idx].programCnt = programCntint
		instructionSet[idx].data = Imm_to_32bit_converter(d.linevalue, 32)
		//outputString := DATA(instructionSet[idx])
		outputString := fmt.Sprintf("%.32s \t%d\t%d\n", dataElements[idx].rawInstruction[0:32], instructionSet[idx].programCnt, instructionSet[idx].data)
		fmt.Println(outputString)
		fmt.Fprintf(outfile, outputString)

		programCntint = programCntint + 4

	}

	outfile.Close()

}

func main() {
	// Get our input and output filenames from the arguments
	var inputFileName string
	var outputFileName string
	var programCnt int
	programCnt = 96
	//consolePrint := true  started to add switch for console print

	// flag package - parses command line arguments based on flags
	// go run . team#_project1.go -i test1_bin.txt -o team#_out
	// so the i flag needs to return the input file name/path and the o flag returns the output file/path

	flag.StringVar(&inputFileName, "i", "", "Gets the input file name")
	flag.StringVar(&outputFileName, "o", "", "Gets the output file name")
	flag.Parse()

	fmt.Println(inputFileName)
	fmt.Println(outputFileName)

	if flag.NArg() != 0 {
		os.Exit(-1)
	}

	// Open the input file
	infile, err := os.OpenFile(inputFileName, os.O_RDONLY, filePermissions)
	if err != nil {
		fmt.Println("\n# Unable to open input file \"" + inputFileName + "\"")
	} else {
		fmt.Println("\n# input file \"" + inputFileName + "\" ")
	}
	defer infile.Close()

	//Open the output file, or set the output file to stdout
	outputFileName = outputFileName + "_dis.txt"
	outfile, err := os.OpenFile(outputFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePermissions)
	if err != nil {
		if outputFileName != "" {
			fmt.Println(err)
			fmt.Println("Printing to standard output . . .")
		}
		outfile = os.Stdout
	}

	defer outfile.Close() // keeps the outfile open

	//read the input file into an array of strings
	// This section reads in the input file.  There is no checking done on the input, just reads it

	instructionSet := []Instruction{}   // create a slice of instructions size zero
	scanner := bufio.NewScanner(infile) // create a scanner
	for scanner.Scan() {
		newInstruct := Instruction{rawInstruction: scanner.Text()} // .Text() is the "line" as a string, creates a new
		// instruction instance and puts the line int to rawInstruction variable
		instructionSet = append(instructionSet, newInstruct) // this takes the new Instruction and adds it to the
		// instructionSet slice one by one
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// Now that we have all of the instructions in a slice, lets process the slice one line at a time
	// Decode opcode out of the raw instruction and put in Instruction:opcode

	//********** Process each line in slice for opcode - store in individual instruction structs

	for idx, s := range instructionSet {

		//linevalue is the base 2 version of the instruction converted from the input string

		// this is tricky since ParseUint returns uint64 value
		// we are converting each line of text to binary number
		// had to do this because I used shift to get the opcodes
		// if I used slice instead I would not have to do this - I think

		// the difficulty here is that we need the opcode to be an int so we can use less than, but it is a string.
		// I can either convert the string to a binary and shift, or slice and the convert slice to int
		// base 2 means that you are interpreting a string that represents a binary number.  shift is much better!!!!
		linevalue, err := strconv.ParseUint(s.rawInstruction, 2, 32)

		instructionSet[idx].linevalue = uint32(linevalue) // changes to uint32

		instructionSet[idx].programCnt = programCnt // increment the program count

		if err != nil {
			fmt.Println("input = \"" + s.rawInstruction + "\"")
			fmt.Println(err)
		}

		// shifts the 32 bit integer 21 leaving the leftmost 11 bits as a number which we will only see as a base 10
		// integer not binary   - need to explain how always using 11 bits matches the different length opcodes
		instructionSet[idx].opcode = uint32(linevalue) >> 21
	}

	// Now that we have the opcodes parsed lets go through the instructionSet array and parse out the Instruction line.
	// Reference the Instruction cheat sheet
	// Use mask and shift to do the parseing out.

	TakeApart(instructionSet, outfile, &programCnt)

	outfile.Close()

}
