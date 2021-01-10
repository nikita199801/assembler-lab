package main

import (
	"fmt"
	"strconv"

	// "time"
	"bufio"
	"math"
	"os"
	"strings"
)

type ALU struct {
	OF, ZF, SF, CF int
}

// ADD instruction
func (alu *ALU) ADD(inputA, inputB int) int {
	var result int
	alu.OF, alu.SF, alu.ZF = 0, 0, 0

	result = inputA + inputB
	if result > 0xFFFFFFFF {
		alu.OF = 1
	}
	if result < 0 {
		alu.SF = 1
	}
	if result == 0 {
		alu.ZF = 1
	}
	if alu.OF == 1 {
		alu.CF = 1
	}
	fmt.Print("*********************************\n")
	fmt.Println("FLAGS")
	fmt.Println("OF:", alu.OF, "SF:", alu.SF, "ZF:", alu.ZF, "CF:", alu.CF)
	fmt.Print("*********************************\n")
	return result
}

// ADC instruction
func (alu *ALU) ADC(op1, op2 int) (result int) {
	result = op1 + op2 + alu.CF
	return
}

func (alu *ALU) SUB(op1, op2 int) (result int) {
	result = op1 - op2
	if result == 0 {
		alu.ZF = 1
	}
	if result < 0 {
		alu.SF = 1
	}
	return
}

func (alu *ALU) CMD(op1, op2 int) (result int) {
	result = op1 - op2
	if result == 0 {
		alu.ZF = 1
	}
	fmt.Print("*********************************\n")
	fmt.Println("FLAGS")
	fmt.Println("OF:", alu.OF, "SF:", alu.SF, "ZF:", alu.ZF, "CF:", alu.CF)
	fmt.Print("*********************************\n")
	return
}

// Reseting ALU flags
func (alu *ALU) ResetFlags() {
	alu.CF, alu.OF, alu.SF, alu.ZF = 0, 0, 0, 0
}

// MULtiply two values
func (alu *ALU) MUL(inputA, inputB int) (msb, lsb int) {
	var temp float64
	result := inputA * inputB
	if result > 0xFFFFFFFF {
		str := fmt.Sprintf("%b", result>>32)
		a := strings.Split(str, "")
		for i := 0; i < len(a); i++ {
			if a[i] == "1" {
				temp = temp + math.Pow(2, float64((len(a)-1-i)))
			}
		}
		msb = int(temp)
		temp = 0
		str = fmt.Sprintf("%b", result)
		b := strings.Split(str, "")
		for i := len(a); i < len(b); i++ {
			if b[i] == "1" {
				temp = temp + math.Pow(2, float64((len(b)-1-i)))
			}
		}
		lsb = int(temp)
	} else {
		msb = 0
		lsb = result
	}
	return
}

func main() {
	var EAX, EBX, ECX, EDX, ESI, EDI, idx, OPCODE int
	registers := [6]*int{&EAX, &EBX, &ECX, &EDX, &ESI, &EDI}
	var alu ALU
	i := 0
	// ********** FOR LAB_1,3
	dmem := [10]int{21, 2, 3, 4, 5}
	// cmem := [10]int{0x12060005, 0x12050000, 0x13020500, 0x20010102, 0xDF050000, 0xBF000605, 0xAF030000}
	cmem := [15]int{0xFEFEFEFE}
	// **********FOR LAB_2
	// a := []int{9456431, 9456434, 532, 44, 9456434}
	// b := []int{4531, 4531, 822, 19, 51022}
	// requiredResult := 0
	// dmem := [1024]int{}
	// cmem := [16]int{0x12050000, 0x12060005, 0x13020500, 0x20060605, 0x13030600, 0x30060605, 0x40020000, 0x20010102, 0x21040304, 0xDF050000, 0xBF000605, 0xAF030000}
	// for i := 0; i < len(a); i++ {
	// 	dmem[i] = a[i]
	// 	dmem[len(a)+i] = b[i]
	// }
	// for i := 0; i < 5; i++ {
	// 	requiredResult += a[i] * b[i]
	// }
	// fmt.Println("Required result:", fmt.Sprintf("0x%08X", requiredResult))
	// *******************
	for {
		fmt.Print("__________________\n")
		fmt.Print("| PC: ", i, " 	 |\n")
		fmt.Print("|________________|\n")
		// time.Sleep(1 * time.Second)
		OPCODE = cmem[i] >> 24
		fmt.Print("Fetching instruction: ")
		// time.Sleep(1 * time.Second)
		fmt.Printf("0x%X\n", cmem[i])
		// time.Sleep(1 * time.Second)
		fmt.Println("OPCODE:", fmt.Sprintf("%X", OPCODE))
		// time.Sleep(1 * time.Second)

		switch OPCODE {

		case 0x11:
			// Copy registers
			op1 := cmem[i] >> 16 & 0xFF
			op2 := cmem[i] >> 8 & 0xFF
			MOV(registers[op1-1], *registers[op2-1])
			RegisterState(registers)

		case 0x12:
			// LOAD a value to register
			op1 := cmem[i] >> 16 & 0xF
			MOV(registers[op1-1], cmem[i]&0xFFFF)
			RegisterState(registers)

		case 0x13:
			// LOAD from RAM by address in register
			op1 := cmem[i] >> 16 & 0xF
			op2 := registers[(cmem[i]>>8&0xF)-1]
			MOV(registers[op1-1], dmem[*op2])
			RegisterState(registers)

		case 0x20:
			// Sum two registers
			resultOp := registers[(cmem[i]>>16&0xFF)-1]
			op1 := registers[(cmem[i]>>8&0xFF)-1]
			op2 := registers[(cmem[i]&0xFF)-1]
			result := alu.ADD(*op1, *op2)
			if alu.OF == 1 {
				result = result & 0xFFFFFFFF
			}
			*resultOp = result
			RegisterState(registers)

		case 0x21:
			// Add with carry bit
			resultOp := registers[(cmem[i]>>16&0xFF)-1]
			op1 := registers[(cmem[i]>>8&0xFF)-1]
			op2 := registers[(cmem[i]&0xFF)-1]
			*resultOp = alu.ADC(*op1, *op2)
			RegisterState(registers)
			alu.ResetFlags()

		case 0x30:
			// Substract two registers
			resultOp := registers[(cmem[i]>>16&0xFF)-1]
			op1 := registers[(cmem[i]>>8&0xFF)-1]
			op2 := registers[(cmem[i]&0xFF)-1]
			*resultOp = alu.SUB(*op1, *op2)
			RegisterState(registers)

		case 0x40:
			// Multiply two values
			op1 := registers[(cmem[i]>>16&0xFF)-1]
			if (cmem[i] >> 16 & 0xFF) == 0x01 {
				msb, lsb := alu.MUL(*op1, EBX)
				EAX = lsb
				EDX = msb
				fmt.Println(msb, lsb)
			} else if (cmem[i] >> 16 & 0xFF) == 0x02 {
				msb, lsb := alu.MUL(*op1, ECX)
				EBX = lsb
				ECX = msb
			}
			RegisterState(registers)

		case 0xAF:
			//JNE instrunction
			if alu.ZF != 1 {
				op1 := cmem[i] >> 16 & 0xFF
				i = op1 - 2
			}

		case 0xBF:
			//compare instruction
			alu.ResetFlags()
			op1 := registers[(cmem[i]>>8&0xFF)-1]
			op2 := registers[(cmem[i]&0xFF)-1]
			alu.CMD(*op1, *op2)

		case 0xDF:
			//INCREMENT value
			*registers[(cmem[i]>>16&0xFF)-1]++
			RegisterState(registers)

		case 0xFE:
			// Halt executing
			idx = i
			fmt.Printf("Enter commands: ")
			for {
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString(';')
				text = strings.ToUpper(text)
				if text == "END;" {
					break
				} else {
					command := ConvertCommandIntoCode(text)
					cmem[i+1] = command
					i++
				}
			}
			i = idx
			// fmt.Println("cmem", fmt.Sprintf("%08X", cmem))
		}
		fmt.Println("===================================================")
		// time.Sleep(1 * time.Second)
		i++
	}
}

// getRegisterIndex by name
func getRegisterIndex(regiser string) (value int) {
	if regiser == "EAX" {
		value = 1
	} else if regiser == "EBX" {
		value = 2
	} else if regiser == "ECX" {
		value = 3
	} else if regiser == "EDX" {
		value = 4
	} else if regiser == "ESI" {
		value = 5
	} else if regiser == "EDI" {
		value = 6
	} else {
		value = 0
	}
	return
}

// Convert Command Into Machine Code from console input !!!!! REQUIERED FOR LAB_3
func ConvertCommandIntoCode(command string) int {
	command = strings.ReplaceAll(strings.ToUpper(command), ";", "")
	a := strings.Split(command, " ")
	var instruction int
	// fmt.Println(a)
	switch a[0] {
	case "MOV":
		if strings.Contains(command, "[") {
			index := strings.ReplaceAll(a[2], "[", "")
			index = strings.ReplaceAll(index, "]", "")
			instruction = ((0x13<<8+getRegisterIndex(a[1]))<<8 + getRegisterIndex(index)) << 8
		} else if getRegisterIndex(a[2]) == 0 {
			i, _ := strconv.Atoi(a[2])
			instruction = (0x12<<8+getRegisterIndex(a[1]))<<16 + i
		} else {
			instruction = ((0x11<<8+getRegisterIndex(a[1]))<<8 + getRegisterIndex(a[2])) << 8
		}

	case "ADD":
		instruction = 0x20
		for i := 1; i < len(a); i++ {
			instruction = instruction<<8 + getRegisterIndex(a[i])
		}

	case "JNE":
		i, _ := strconv.Atoi(a[1])
		instruction = (0xAF<<8 + i) << 16

	case "INC":
		instruction = (0xDF<<8 + getRegisterIndex(a[1])) << 16

	case "CMP":
		instruction = (0xBF<<16+getRegisterIndex(a[1]))<<8 + getRegisterIndex(a[2])
	}
	return instruction
}

// RegisterState of all register
func RegisterState(registers [6]*int) {
	fmt.Println("************Registers***************")
	h := fmt.Sprintf("0x%08X", *registers[0])
	fmt.Println("EAX:", h)
	h = fmt.Sprintf("0x%08X", *registers[1])
	fmt.Println("EBX:", h)
	h = fmt.Sprintf("0x%08X", *registers[2])
	fmt.Println("ECX:", h)
	h = fmt.Sprintf("0x%08X", *registers[3])
	fmt.Println("EDX:", h)
	h = fmt.Sprintf("0x%08X", *registers[4])
	fmt.Println("ESI:", h)
	h = fmt.Sprintf("0x%08X", *registers[5])
	fmt.Println("EDI:", h)
	fmt.Println("************************************")
}

// MOV from Register A to Register B
func MOV(register *int, data int) {
	*register = data
}
