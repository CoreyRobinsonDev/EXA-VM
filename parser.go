package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func NewExa(filePath string) (Exa, error) {
	code := string(Unwrap(os.ReadFile(filePath)))
	name := strings.Split(strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1], ".")[0]
	exa := Exa {
		Code: code,
		Name: name,
		FileCursor: 0,
		M: make(chan string, 1),
		Marks: make(map[string]int),
	}
	vm.Exas = append(vm.Exas, exa)

	file := Unwrap(os.Open(filePath))
	defer file.Close()
	codeLines := strings.Split(code, "\n")

	for lineNum, line := range codeLines {
		if len(line) == 0 { continue }
		exa.LineNum = lineNum + 1
		words := strings.Split(line, " ")
		keyword := Unwrap(NewKeyword(words[0]))
		// Enforce HOST as the first command of a program
		if keyword != HOST && exa.Host == "" {
			err := fmt.Sprintf("%s %d: HOST must be defined at the start of program execution", name, lineNum)
			exa.Error = err
			return exa, errors.New(err)
		}
		err := keyword.Eval(&exa, words...)

		if err != nil {
			exa.Error = err.Error()
			return exa, err
		}
	}

	return exa, nil
}

func _copy(src, dst string, exa *Exa) error { 
	switch(strings.ToUpper(src)) {
	case "X":
		switch(strings.ToUpper(dst)) {
		case "X": break
		case "F": return exa.WriteFile(exa.X)
		case "T": exa.T = exa.X
		case "M": 
			exa.M <- exa.X
		default:
			return errors.New(fmt.Sprintf("Invalid COPY destination %s", dst))
		}
	case "F":
		word, err := exa.ReadFile()

		if err != nil {
			return err
		}

		switch(strings.ToUpper(dst)) {
		case "X": exa.X = word
		case "F": break
		case "T": exa.T = word
		case "M": exa.M <- word
		default:
			return errors.New(fmt.Sprintf("Invalid COPY destination %s", dst))
		}
	case "T":
		switch(strings.ToUpper(dst)) {
		case "X": exa.X = exa.T
		case "F": return exa.WriteFile(exa.T)
		case "T": break
		case "M": exa.M <- exa.T
		default:
			return errors.New(fmt.Sprintf("Invalid COPY destination %s", dst))
		}
	case "M":
		switch(strings.ToUpper(dst)) {
		case "X": exa.X = <-exa.M
		case "T": exa.T = <-exa.M
		case "F": return exa.WriteFile(<-exa.M)
		case "M": break
		default:
			return errors.New(fmt.Sprintf("Invalid COPY destination %s", dst))
		}
	default:
		srcNum, _ := strconv.Atoi(src) 

		if srcNum > 9999 {
			return errors.New(fmt.Sprintf("COPY overflow: %d exceeds 9999", srcNum))
		}

		if srcNum < -9999 {
			return errors.New(fmt.Sprintf("COPY underflwo: %d below -9999", srcNum))
		}

		switch(strings.ToUpper(dst)) {
		case "X": exa.X = src
		case "F": return exa.WriteFile(src)
		case "T": exa.T = src
		case "M": exa.M <- src
		default:
			return errors.New(fmt.Sprintf("Invalid COPY destination %s", dst))
		}
	}

	return nil 
}	
func addi() error { return nil }
func subi() error { return nil }
func muli() error { return nil }
func divi() error { return nil }
func modi() error { return nil }
func swiz() error { return nil }
func mark(name string, exa *Exa) error { 
	exa.Marks[name] = exa.LineNum
	return nil
}
func jump() error { return nil }
func tjmp() error { return nil }
func fjmp() error { return nil }
func test() error { return nil }
func repl() error { return nil }
func halt() error { return nil }
func kill() error { return nil }
func link() error { return nil }
func host(name string, exa *Exa) error {
	if name == "" {
		return errors.New("HOST: No HOST name provided")
	}
	exa.Host = name
	return nil
}
func mode() error { return nil }
func _make(name string, exa *Exa) error { 
	file := Unwrap(os.Create(name))
	defer file.Close()

	exa.F = name

	return nil
}
func grab(name string, exa *Exa) error { 
	if name == "X" {
		name = exa.X
	} else if name == "T" {
		name = exa.T
	} else if name == "M" {
		name = <- exa.M
	}

	_, err := os.ReadFile(name)

	if err != nil {
		return errors.New(fmt.Sprintf("GRAB: File %s could not be found", name))
	}

	exa.F = name

	return nil 
}
func file() error { return nil }
func seek(amount string, exa *Exa) error { 
	if exa.F == "" { return errors.New("SEEK: No file to SEEK") }
	content := string(Unwrap(os.ReadFile(exa.F)))
	contentArr := strings.Split(content, " ")

	if amount == "X" {
		if exa.X == "" {
			amount = "0"
		} else {
			amount = exa.X
		}
	} else if amount == "T" {
		if exa.T == "" {
			amount = "0"
		} else {
			amount = exa.T
		}
	} else if amount == "M" {
		m := <- exa.M
		if m == "" {
			amount = "0"
		} else {
			amount = m
		}
	}

	num, err := strconv.Atoi(amount)
	if err != nil {
		return errors.New(fmt.Sprintf("SEEK: non-numeric value %s used", amount))
	} else if num > 9999 {
		return errors.New(fmt.Sprintf("SEEK: value overflow on %s", amount))
	} else if num < -9999 {
		return errors.New(fmt.Sprintf("SEEK: value underflow on %s", amount))
	}


	exa.FileCursor += num

	if exa.FileCursor >= len(contentArr) {
		exa.FileCursor = len(contentArr) -2
	} else if exa.FileCursor < 0 { 
		exa.FileCursor = 0 
	} 

	return nil
}
func void() error { return nil }
func drop() error { return nil }
func wipe() error { return nil }
func note() error { return nil }
func noop() error { return nil }
func rand() error { return nil }

type Keyword int

func NewKeyword(text string) (Keyword, error) {
	switch text {
	case "COPY": return COPY, nil
	case "ADDI": return ADDI, nil
	case "SUBI": return SUBI, nil
	case "MULI": return MULI, nil
	case "DIVI": return DIVI, nil
	case "MODI": return MODI, nil
	case "SWIZ": return SWIZ, nil
	case "MARK": return MARK, nil
	case "JUMP": return JUMP, nil
	case "TJMP": return TJMP, nil
	case "FJMP": return FJMP, nil
	case "TEST": return TEST, nil
	case "REPL": return REPL, nil
	case "HALT": return HALT, nil
	case "KILL": return KILL, nil
	case "LINK": return LINK, nil
	case "HOST": return HOST, nil
	case "MODE": return MODE, nil
	case "MAKE": return MAKE, nil
	case "GRAB": return GRAB, nil
	case "FILE": return FILE, nil
	case "SEEK": return SEEK, nil
	case "VOID": return VOID, nil
	case "DROP": return DROP, nil
	case "WIPE": return WIPE, nil
	case "NOTE": return NOTE, nil
	case "NOOP": return NOOP, nil
	case "RAND": return RAND, nil
	default: return -1, errors.New(fmt.Sprintf("invalid keyword: %s", text))
	}

}

const (
	COPY Keyword = iota
	ADDI
	SUBI
	MULI
	DIVI
	MODI
	SWIZ
	MARK
	JUMP
	TJMP
	FJMP
	TEST
	REPL
	HALT
	KILL
	LINK
	HOST
	MODE
	MAKE
	GRAB
	FILE
	SEEK
	VOID
	DROP
	WIPE
	NOTE
	NOOP
	RAND
)

func (k Keyword) String() string {
	switch k {
	case COPY: return "COPY"
	case ADDI: return "ADDI"
	case SUBI: return "SUBI"
	case MULI: return "MULI"
	case DIVI: return "DIVI"
	case MODI: return "MODI"
	case SWIZ: return "SWIZ"
	case MARK: return "MARK"
	case JUMP: return "JUMP"
	case TJMP: return "TJMP"
	case FJMP: return "FJMP"
	case TEST: return "TEST"
	case REPL: return "REPL"
	case HALT: return "HALT"
	case KILL: return "KILL"
	case LINK: return "LINK"
	case HOST: return "HOST"
	case MODE: return "MODE"
	case MAKE: return "MAKE"
	case GRAB: return "GRAB"
	case FILE: return "FILE"
	case SEEK: return "SEEK"
	case VOID: return "VOID"
	case DROP: return "DROP"
	case WIPE: return "WIPE"
	case NOTE: return "NOTE"
	case NOOP: return "NOOP"
	case RAND: return "RAND"
	default: return "{ERROR}"
	}
}

func (k Keyword) Eval(exa *Exa, args... string) error {
	switch k {
	case COPY: 
		if len(args) <= 2 {
			return errors.New("COPY: missing argument")
		}
		return _copy(args[1], args[2], exa)
	case ADDI: return addi()
	case SUBI: return subi()
	case MULI: return muli()
	case DIVI: return divi()
	case MODI: return modi()
	case SWIZ: return swiz()
	case MARK: 
		if len(args) <= 1 {
			return errors.New("MARK: missing label" )
		}
		return mark(args[1], exa)
	case JUMP: return jump()
	case TJMP: return tjmp()
	case FJMP: return fjmp()
	case TEST: return test()
	case REPL: return repl()
	case HALT: return halt()
	case KILL: return kill()
	case LINK: return link()
	case HOST: 
		if len(args) <= 1 {
			return errors.New("HOST: missing host name" )
		}
		return host(args[1], exa)
	case MODE: return mode()
	case MAKE: 
		filename := ""
		if len(args) <= 1 {
			vm.Files++
			filename = fmt.Sprintf("./%d", 400 + vm.Files)
		} else {
			filename = args[1]
		}
		return _make(filename, exa)
	case GRAB:
		if len(args) <= 1 {
			return errors.New("GRAB: missing file name" )
		} 
		return grab(args[1], exa)
	case FILE: return file()
	case SEEK: 
		if len(args) <= 1 {
			return errors.New("SEEK: missing cursor move amount")
		}
		return seek(args[1], exa)
	case VOID: return void()
	case DROP: return drop()
	case WIPE: return wipe()
	case NOTE: return note()
	case NOOP: return noop()
	case RAND: return rand()
	}
	return nil
}


type Register int
const (
	EOF Register = iota
	X
	T
	F
	M
	STDI
	STDO
	STDE
)

func (r Register) String() string {
	switch r {
	case EOF: return "EOF"
	case X: return "X"
	case T: return "T"
	case F: return "F"
	case M: return "M"
	case STDI: return "STDI"
	case STDO: return "STDO"
	case STDE: return "STDE"
	default: return "{ERROR}"
	}
}

