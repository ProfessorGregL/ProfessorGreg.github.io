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

func main() {
	// Get our input and output filenames from the arguments
	var inputFileName string
	var outputFileName string

	var programCnt int
	programCnt = 96

	flag.StringVar(&inputFileName, "i", "", "Gets the input file name")
	flag.StringVar(&outputFileName, "o", "", "Gets the output file name")
	flag.Parse()

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
	outfile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY, filePermissions)
	if err != nil {
		if outputFileName != "" {
			fmt.Println(err)
			fmt.Println("Printing to standard output . . .")
		}
		outfile = os.Stdout
	}

	defer outfile.Close()

	//read the input file into an array of strings

	var Breaknow = false

	InputParsed := []Instruction{}
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		ns := Instruction{rawInstruction: scanner.Text()}
		InputParsed = append(InputParsed, ns)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// Now that we have all of the lines in an array, lets process the array one line at a time
	// decode the opcode and

	// Decode opcode and put in Instruction Struct
	for idx, s := range InputParsed {

		//linevalue is the base 2 version of the instruction converted from the input string

		// this is tricky since ParseUint returns uint64 value
		linevalue, err := strconv.ParseUint(s.rawInstruction, 2, 32)

		InputParsed[idx].linevalue = uint32(linevalue)

		InputParsed[idx].programCnt = programCnt

		//fmt.Printf("Binary Representation of %d is %s.\n", linevalue, IntegerToBinary(linevalue))
		if err != nil {
			fmt.Println("input = \"" + s.rawInstruction + "\"")
			fmt.Println(err)
		}
		// shift out the opcode - uint64 and convert to int and put in s=Instruction struct
		//InputParsed[idx].opcode = strconv.FormatUint(linevalue >> 21,10)

		InputParsed[idx].opcode = uint32(linevalue) >> 21

	}

	// Now that we have the opcodes parsed lets go through the InputParsed array and parse out the Instruction line.
	// Reference the Instruction cheat sheet
	// Use mask and shift to do the parseing out.

	for idx, s := range InputParsed {

		fmt.Println(s.opcode)

		InputParsed[idx].programCnt = programCnt

		switch {

		case s.opcode == 0 && Breaknow == false:
			fmt.Println("A NOP was discovered")
			InputParsed[idx].op = "NOP"
			tempString := NOP(InputParsed[idx])
			fmt.Println(tempString)
			fmt.Fprintf(outfile, tempString)

		case s.opcode == 2038 && Breaknow == false:
			if (s.linevalue & 0x1FFFFF) == 0x1EFFE7 {
				fmt.Println("The break just was discovered")
				InputParsed[idx].op = "BREAK"
				Breaknow = true
				tempString := BREAK(InputParsed[idx])
				fmt.Println(tempString)
				fmt.Fprintf(outfile, tempString)
				break
			}

		case s.opcode >= 160 && s.opcode <= 191 && Breaknow == false:
			fmt.Println("Found an B type Instruction")
			InputParsed[idx].op = "B"
			// not sure I need to have int versus uint as input, def need it as output
			InputParsed[idx].offset = Imm_to_32bit_converter(s.linevalue&0x3FFFFFF, 26)
			fmt.Println(InputParsed[idx].offset)

			tempString := BInstr(InputParsed[idx])
			fmt.Println(tempString)
			fmt.Fprintf(outfile, tempString)

		case s.opcode == 1112 || s.opcode == 1624 || s.opcode == 1104 || s.opcode == 1360 || s.opcode == 1872 || (s.opcode >= 1690 && s.opcode <= 1693) && Breaknow == false:
			//InputParsed[idx].typeofInstruction = "R"
			fmt.Println("Found an R type or RShift type Instruction")
			if s.opcode == 1112 {
				InputParsed[idx].op = "ADD"
			} else if s.opcode == 1624 {
				InputParsed[idx].op = "SUB"
			} else if s.opcode == 1104 {
				InputParsed[idx].op = "AND"
			} else if s.opcode == 1360 {
				InputParsed[idx].op = "ORR"
			} else if s.opcode == 1872 {
				InputParsed[idx].op = "EOR"
			} else if s.opcode == 1690 {
				InputParsed[idx].op = "LSR"
				InputParsed[idx].shamt = uint16((s.linevalue >> 10) & 0x3F)
			} else if s.opcode == 1691 {
				InputParsed[idx].op = "LSL"
				InputParsed[idx].shamt = uint16((s.linevalue >> 10) & 0x3F)
			} else if s.opcode == 1692 {
				InputParsed[idx].op = "ASR"
				InputParsed[idx].shamt = uint16((s.linevalue >> 10) & 0x3F)
			}

			InputParsed[idx].rd = uint8(s.linevalue & 0x1F)
			InputParsed[idx].rn = uint8((s.linevalue >> 5) & 0x1F)
			InputParsed[idx].rm = uint8((s.linevalue >> 16) & 0x1F)

			if s.opcode == 1112 || s.opcode == 1624 || s.opcode == 1104 || s.opcode == 1360 || s.opcode == 1872 {
				tempString := RInstr(InputParsed[idx])
				fmt.Println(tempString)
				fmt.Fprintf(outfile, tempString)
			} else if s.opcode == 1690 || s.opcode == 1691 {
				tempString := RShiftInstr(InputParsed[idx])
				fmt.Println(tempString)
				fmt.Fprintf(outfile, tempString)
			}

		case s.opcode == 1160 || s.opcode == 1161 || s.opcode == 1672 || s.opcode == 1673 && Breaknow == false:
			//case 1160, 1161, 1672, 1673:
			fmt.Println("Found an I type instruction")
			InputParsed[idx].rd = uint8(s.linevalue & 0x1F)
			InputParsed[idx].rn = uint8((s.linevalue >> 5) & 0x1F)
			InputParsed[idx].im = Imm_to_32bit_converter((s.linevalue>>10)&0xFFF, 12)

			if s.opcode == 1160 || s.opcode == 1161 {
				InputParsed[idx].op = "ADDI"
			} else if s.opcode == 1672 || s.opcode == 1673 {
				InputParsed[idx].op = "SUBI"
			}

			tempString := IInstr(InputParsed[idx])
			fmt.Println(tempString)
			fmt.Fprintf(outfile, tempString)

		case s.opcode == 1984 || s.opcode == 1986 && Breaknow == false:
			//case 1984, 1986:
			fmt.Println("Found a D type instruction")
			InputParsed[idx].addr = uint16((s.linevalue >> 12) & 0x1FF)
			InputParsed[idx].rn = uint8((s.linevalue >> 5) & 0x1F)
			InputParsed[idx].rt = uint8(s.linevalue & 0x1F)
			if s.opcode == 1984 {
				InputParsed[idx].op = "LDUR"
			} else if s.opcode == 1986 {
				InputParsed[idx].op = "STUR"
			}

			tempString := DInstr(InputParsed[idx])
			fmt.Println(tempString)
			fmt.Fprintf(outfile, tempString)

		case (s.opcode >= 1684 && s.opcode <= 1687) || (s.opcode >= 1940 && s.opcode <= 1943) && Breaknow == false:
			//ase 1684, 1685, 1686, 1687, 1940, 1941, 1942, 1943:
			fmt.Println("Found an IM type Instruction")

			InputParsed[idx].field = s.linevalue >> 5 & 0xFFFF
			InputParsed[idx].rd = uint8(s.linevalue & 0x1F)

			index := uint8((s.linevalue >> 17) & 0x30)

			if s.opcode >= 1684 && s.opcode <= 1687 {
				InputParsed[idx].op = "MOVZ"
			}

			if s.opcode >= 1940 && s.opcode <= 1943 {
				InputParsed[idx].op = "MOVK"
			}

			switch index {

			case 0:
				InputParsed[idx].shfcd = 0

			case 1:
				InputParsed[idx].shfcd = 16

			case 2:
				InputParsed[idx].shfcd = 32

			case 3:
				InputParsed[idx].shfcd = 48

			}
			InputParsed[idx].shfcd = uint8((s.linevalue >> 17) & 0x30)

			tempString := IMInstr(InputParsed[idx])
			fmt.Println(tempString)
			fmt.Fprintf(outfile, tempString)

		case s.opcode >= 1440 && s.opcode <= 1447 && Breaknow == false:
			//case 1440,1441,1442,1443,1444,1445,1146,1447:
			fmt.Println("Found an CBZ type Instruction")
			InputParsed[idx].rd = uint8(s.linevalue & 0x1F)
			InputParsed[idx].offset = Imm_to_32bit_converter((s.linevalue>>5)&0x7FFFF, 19)

			if s.opcode >= 1440 && s.opcode <= 1447 {
				InputParsed[idx].op = "CBZ"
			}

			tempString := CBInstr(InputParsed[idx])
			fmt.Println(tempString)
			fmt.Fprintf(outfile, tempString)

		case s.opcode >= 1448 && s.opcode <= 1455 && Breaknow == false:
			//case 1448,1449,1450,1451,1452,1453,1454,1455:
			fmt.Println("Found an CBNZ type Instruction")
			InputParsed[idx].rd = uint8(s.linevalue & 0x1F)
			InputParsed[idx].offset = Imm_to_32bit_converter((s.linevalue>>5)&0x7FFFF, 19)
			if s.opcode >= 1448 && s.opcode <= 1455 {
				InputParsed[idx].op = "CBNZ"
			}

			tempString := CBInstr(InputParsed[idx])
			fmt.Println(tempString)
			fmt.Fprintf(outfile, tempString)

		case Breaknow == true:
			fmt.Println("Found data")
			InputParsed[idx].op = "DATA"
			InputParsed[idx].data = Imm_to_32bit_converter(s.linevalue, 32)
			tempString := DATA(InputParsed[idx])
			fmt.Println(tempString)
			fmt.Fprintf(outfile, tempString)
		}

		programCnt = programCnt + 4

	}

	fmt.Println("exited the Instruction for loop")

	outfile.Close()

}
