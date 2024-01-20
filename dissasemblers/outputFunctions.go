package main

import "fmt"

func RInstr(curInstruction Instruction) string {

	firstPart := fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s", curInstruction.rawInstruction[0:11], curInstruction.rawInstruction[11:16], curInstruction.rawInstruction[16:22], curInstruction.rawInstruction[22:27], curInstruction.rawInstruction[27:32])
	secondPart := fmt.Sprintf("\t%d\t%s\tR%d, R%d, R%d\n", curInstruction.programCnt, curInstruction.op, curInstruction.rd, curInstruction.rn, curInstruction.rm)
	return firstPart + secondPart

}

func RShiftInstr(curInstruction Instruction) string {

	firstPart := fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s", curInstruction.rawInstruction[0:11], curInstruction.rawInstruction[11:16], curInstruction.rawInstruction[16:22], curInstruction.rawInstruction[22:27], curInstruction.rawInstruction[27:32])
	secondPart := fmt.Sprintf("\t%d\t%s\tR%d, R%d, #%d\n", curInstruction.programCnt, curInstruction.op, curInstruction.rd, curInstruction.rn, curInstruction.shamt)
	return firstPart + secondPart

}

func BInstr(curInstruction Instruction) string {

	firstPart := fmt.Sprintf("%.6s %.26s   ", curInstruction.rawInstruction[0:6], curInstruction.rawInstruction[6:32])
	secondPart := fmt.Sprintf("\t%d\t%s\t#%d\n", curInstruction.programCnt, curInstruction.op, curInstruction.offset)
	return firstPart + secondPart
}

func DInstr(curInstruction Instruction) string {

	firstPart := fmt.Sprintf("%.11s %.9s %.2s %.5s %.5s", curInstruction.rawInstruction[0:11], curInstruction.rawInstruction[11:20], curInstruction.rawInstruction[20:22], curInstruction.rawInstruction[22:27], curInstruction.rawInstruction[27:32])
	secondPart := fmt.Sprintf("\t%d\t%s\tR%d, [R%d, #%d]\n", curInstruction.programCnt, curInstruction.op, curInstruction.rt, curInstruction.rn, curInstruction.addr)
	return firstPart + secondPart

}

func IInstr(curInstruction Instruction) string {

	firstPart := fmt.Sprintf("%.10s %.12s %.5s %.5s", curInstruction.rawInstruction[0:10], curInstruction.rawInstruction[10:22], curInstruction.rawInstruction[22:27], curInstruction.rawInstruction[27:32])
	secondPart := fmt.Sprintf("\t%d\t%s\tR%d, R%d, #%d\n", curInstruction.programCnt, curInstruction.op, curInstruction.rd, curInstruction.rn, curInstruction.im)
	return firstPart + secondPart

}

func CBInstr(curInstruction Instruction) string {

	firstPart := fmt.Sprintf("%.8s %.19s %.5s", curInstruction.rawInstruction[0:8], curInstruction.rawInstruction[8:27], curInstruction.rawInstruction[27:32])
	secondPart := fmt.Sprintf("\t%d\t%s\tR%d, #%d\n", curInstruction.programCnt, curInstruction.op, curInstruction.rd, curInstruction.offset)
	return firstPart + secondPart

}

func IMInstr(curInstruction Instruction) string {

	firstPart := fmt.Sprintf("%.9s %.2s %.16s %.5s", curInstruction.rawInstruction[0:9], curInstruction.rawInstruction[9:11], curInstruction.rawInstruction[11:27], curInstruction.rawInstruction[27:32])
	secondPart := fmt.Sprintf("\t%d\t%s\tR%d, %d, LSL %d\n", curInstruction.programCnt, curInstruction.op, curInstruction.rd, curInstruction.field, curInstruction.shfcd)
	return firstPart + secondPart

}

func BREAK(curInstruction Instruction) string {

	firstPart := fmt.Sprintf("%.8s %.3s %.5s %.5s %.5s %.6s", curInstruction.rawInstruction[0:8], curInstruction.rawInstruction[8:11], curInstruction.rawInstruction[11:16], curInstruction.rawInstruction[16:21], curInstruction.rawInstruction[21:27], curInstruction.rawInstruction[27:32])
	secondPart := fmt.Sprintf("\t%d\t%s\n", curInstruction.programCnt, curInstruction.op)
	return firstPart + secondPart

}

func DATA(curInstruction Instruction) string {

	firstPart := fmt.Sprintf("%.32s", curInstruction.rawInstruction[0:32])
	secondPart := fmt.Sprintf("\t%d\t%d\n", curInstruction.programCnt, curInstruction.data)
	return firstPart + secondPart

}

func NOP(curInstruction Instruction) string {

	firstPart := fmt.Sprintf("%.8s %.3s %.5s %.5s %.5s %.6s", curInstruction.rawInstruction[0:8], curInstruction.rawInstruction[8:11], curInstruction.rawInstruction[11:16], curInstruction.rawInstruction[16:21], curInstruction.rawInstruction[21:27], curInstruction.rawInstruction[27:32])
	secondPart := fmt.Sprintf("\t%d\t%s\n", curInstruction.programCnt, curInstruction.op)
	return firstPart + secondPart

}
