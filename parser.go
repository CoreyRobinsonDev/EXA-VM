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
		M: make(chan string, 1),
	}
	file := Unwrap(os.Open(filePath))
	defer file.Close()
	codeLines := strings.Split(code, "\n")

	for lineNum, line := range codeLines {
		if len(line) == 0 { continue }
		words := strings.Split(line, " ")
		keyword := Unwrap(NewKeyword(words[0]))
		// Enforce HOST as the first command of a program
		if keyword != HOST && exa.Host == "" {
			return Exa{}, errors.New(fmt.Sprintf("%s %d: HOST must be defined at the start of program execution", name, lineNum))
		}
		err := keyword.Eval(&exa, words...)

		if err != nil {
			return Exa{}, err
		}
	}

	vm.Exas++
	return exa, nil
}

func _copy(src, dst string, exa *Exa) error { 
	switch(strings.ToUpper(src)) {
	case "X":
		switch(strings.ToUpper(dst)) {
		case "X":
		case "F":
			file, err := os.Open(exa.F)
			defer file.Close()

			if err != nil {
				return errors.New(fmt.Sprintf("Invalid COPY file location %s", exa.F))
			}
			Unwrap(file.WriteAt([]byte(exa.X), int64(exa.FileCursor)))
		case "T": exa.T = exa.X
		case "M": 
			exa.M <- exa.X
		default:
			return errors.New(fmt.Sprintf("Invalid COPY destination %s", dst))
		}
	case "F":
		switch(strings.ToUpper(dst)) {
		case "X": exa.X = exa.T
		case "F":
		case "T": exa.X = exa.T
		case "M": exa.M <- exa.T
		default:
			return errors.New(fmt.Sprintf("Invalid COPY destination %s", dst))
		}
	case "T":
		switch(strings.ToUpper(dst)) {
		case "X": exa.X = exa.T
		case "F":
			file, err := os.Open(exa.F)
			defer file.Close()

			if err != nil {
				return errors.New(fmt.Sprintf("Invalid COPY file location %s", exa.F))
			}
			Unwrap(file.WriteAt([]byte(exa.T), int64(exa.FileCursor)))
		case "T":
		case "M": exa.M <- exa.T
		default:
			return errors.New(fmt.Sprintf("Invalid COPY destination %s", dst))
		}
	case "M":
		switch(strings.ToUpper(dst)) {
		case "X": exa.X = <-exa.M
		case "T": exa.T = <-exa.M
		case "F":
			file, err := os.Open(exa.F)
			defer file.Close()

			if err != nil {
				return errors.New(fmt.Sprintf("Invalid COPY file location %s", exa.F))
			}
			Unwrap(file.WriteAt([]byte(<-exa.M), int64(exa.FileCursor)))
		case "M":
		default:
			return errors.New(fmt.Sprintf("Invalid COPY destination %s", dst))
		}
	default:
		srcNum, err := strconv.Atoi(src) 
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid COPY source %s", src))
		}

		if srcNum > 9999 {
			return errors.New(fmt.Sprintf("COPY overflow: %d exceeds 9999", srcNum))
		}

		if srcNum < -9999 {
			return errors.New(fmt.Sprintf("COPY underflwo: %d below -9999", srcNum))
		}

		switch(strings.ToUpper(dst)) {
		case "X": exa.X = src
		case "F":
			file, err := os.Open(exa.F)
			defer file.Close()

			if err != nil {
				return errors.New(fmt.Sprintf("Invalid COPY file location %s", exa.F))
			}
			Unwrap(file.WriteAt([]byte(src), int64(exa.FileCursor)))
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
func mark() error { return nil }
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
		return errors.New("No HOST name provided")
	}
	exa.Host = name
	return nil
}
func mode() error { return nil }
func _make(name string, exa *Exa) error { 
	file := Unwrap(os.Create("./" + name))
	defer file.Close()

	exa.F = name

	return nil
}
func grab() error { return nil }
func file() error { return nil }
func seek() error { return nil }
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
			return errors.New(fmt.Sprintf("COPY: index 2, out of range of %v array", args))
		}
		return _copy(args[1], args[2], exa)
	case ADDI: return addi()
	case SUBI: return subi()
	case MULI: return muli()
	case DIVI: return divi()
	case MODI: return modi()
	case SWIZ: return swiz()
	case MARK: return mark()
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
			return errors.New(fmt.Sprintf("HOST: index 1, out of range of %v array", args))
		}
		return host(args[1], exa)
	case MODE: return mode()
	case MAKE: 
		filename := ""
		if len(args) <= 1 {
			filename = fmt.Sprintf("%d", 400 + vm.Files)
			vm.Files++
		} else {
			filename = args[1]
		}
		return _make(filename, exa)
	case GRAB: return grab()
	case FILE: return file()
	case SEEK: return seek()
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
