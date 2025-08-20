package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)


type VM struct {
	Exas int
	Files int
}

type Exa struct {
	// The X register is a general-purpose storage register and can store a number or a keyword.
	X string 
	// The T register is a general-purpose storage register and can store a number or a keyword. It is also the destination for TEST instructions, and is the criterion for conditional jumps (TJMP and FJMP).
	T string
	// The F register allows an EXA to read and write the contents of a held file. When an EXA grabs a file, its "file cursor" will be set to the first value in the file. Reading from the F register will read this value; writing to the F register will overwrite this value. After reading or writing the F register, the file cursor will automatically advance. Writing to the end of the file will append a new value instead of overwriting.
	F string
	/* 
The M register controls an EXA's message-passing functionality. When an EXA writes to the M register the value will be stored in that EXA's outgoing message slot until another EXA reads from the M register and receives the previously written value. Both numbers and keywords can be transferred in this way. 

If an EXA writes to the M register, it will pause execution until that value is read by another EXA. If an EXA reads from the M register, it will pause execution until a value is available to be read. If two or more EXAs attempt to read from another EXA at the same time (or vice versa), one will succeed but which one succeeds will be unpredictable.

By default, an EXA can communicate with any other EXA in the same network. This can be restircted to EXAs in the same network. This can be restricted to EXAs in the same host by executing a MODE instruction. An EXA in global mode cannot communicate with an EXA in local mode, even if they are in the same host.
	*/
	M chan string
	LocalMode bool
	FileCursor int
	Host string
	Links []int
	Name string
	Code string
}

func (e *Exa) ReadFile() (string, error) {
	file, err := os.OpenFile(e.F, os.O_RDWR, 0644)
	defer file.Close()

	if err != nil {
		return "", errors.New(fmt.Sprintf("COPY: Invalid file location %s", e.F))
	}

	content := string(Unwrap(os.ReadFile(e.F)))
	contentArr := strings.Split(content, " ")

	if len(contentArr) <= e.FileCursor {
		return "", errors.New(fmt.Sprintf("COPY: file cursor position (%d) greater than %s file length", e.FileCursor, e.F))
	}

	bytes := make([]byte, len(contentArr[e.FileCursor]))

	contentOffset := ""
	for i, word := range contentArr {
		if i == e.FileCursor { break }
		contentOffset += word + " "
	}

	_, fileErr := file.ReadAt(bytes, int64(len(contentOffset)))

	if fileErr != io.EOF && fileErr != nil {
		return "", fileErr
	}
	e.FileCursor++

	return string(bytes), nil
}

func (e *Exa) WriteFile(word string) error {
	file, err := os.OpenFile(e.F, os.O_RDWR, 0644)
	defer file.Close()

	if err != nil {
		return errors.New(fmt.Sprintf("COPY: Invalid file location %s", e.F))
	}

	content := string(Unwrap(os.ReadFile(e.F)))
	contentArr := strings.Split(content, " ")

	contentOffset := ""
	for i, word := range contentArr {
		if i == e.FileCursor { break }
		contentOffset += word + " "
	}

	Unwrap(file.WriteAt([]byte(strings.Trim(word, "\"") + " "), int64(len(contentOffset))))
	e.FileCursor++

	return nil
}
