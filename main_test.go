package main

import (
	"fmt"
	"os"
	"testing"
)

// HOST
func TestHostEmpty(t *testing.T) {
	_, err := NewExa("./tests/HOST/HostEmpty.exa")

	if err == nil {
		t.Error("HostEmpty failed, expected error to be thrown")
	}
}

func TestHostNoName(t *testing.T) {
	_, err := NewExa("./tests/HOST/HostNoName.exa")

	if err == nil {
		t.Error("HostNoName failed, expected error to be thrown")
	}
}

func TestHostOnLineOne(t *testing.T) {
	expected := Exa { Host: "main"}
	result := Unwrap(NewExa("./tests/HOST/HostOnLineOne.exa"))

	if expected.Host != result.Host {
		t.Errorf("HostOnLineOne failed, expected %s, got %s", expected.Host, result.Host)
	}
}

func TestHostOnLineN(t *testing.T) {
	expected := Exa { Host: "main"}
	result := Unwrap(NewExa("./tests/HOST/HostOnLineN.exa"))

	if expected.Host != result.Host {
		t.Errorf("HostOnLineN failed, expected %s, got %s", expected.Host, result.Host)
	}
}

func TestHostNotFirst(t *testing.T) {
	_, err := NewExa("./tests/HOST/HostNotFirst.exa")

	if err == nil {
		t.Error("HostNotFirst failed, expected error to be thrown")
	}
}
// HOST

// COPY
func TestCopyRegToReg(t *testing.T) {
	X2M := Unwrap(NewExa("./tests/COPY/RegtoM.exa"))
	X2T := Unwrap(NewExa("./tests/COPY/RegtoT.exa"))
	_, X2FErr := NewExa("./tests/COPY/RegtoFNoFile.exa")
	_, X2XErr := NewExa("./tests/COPY/RegtoX.exa")

	X2MM := <-X2M.M
	if X2M.X != X2MM {
		t.Errorf("CopyRegToReg failed, expected %s, got %s", X2M.X, X2MM)
	}

	if X2T.X != X2T.T {
		t.Errorf("CopyRegToReg failed, expected %s, got %s", X2T.X, X2T.T)
	}

	if X2FErr == nil {
		t.Error("CopyRegToReg failed, expected error to be thrown on attempting to copy to a NULL file")
	}

	if X2XErr != nil {
		t.Errorf("CopyRegToReg failed, error (%v) shouldn't be thrown on attempting to copy to the same register", X2XErr)
	}
}

func TestCopyRegToInvalidReg(t *testing.T) {
	_, noFileErr := NewExa("./tests/COPY/RegtoFNoFile.exa")
	_, X2InvalidErr := NewExa("./tests/COPY/XtoInvalidReg.exa") 
	_, T2InvalidErr := NewExa("./tests/COPY/TtoInvalidReg.exa") 
	_, F2InvalidErr := NewExa("./tests/COPY/FNoFiletoInvalidReg.exa") 
	_, M2InvalidErr := NewExa("./tests/COPY/MtoInvalidReg.exa") 

	if noFileErr == nil ||
	X2InvalidErr == nil || 
	T2InvalidErr == nil ||
	F2InvalidErr == nil ||
	M2InvalidErr == nil {
		t.Errorf("CopyRegtoInvalidReg failed\n\tnoFileErr: %v\nX2InvalidErr: %v\nT2InvalidErr: %v\nF2InvalidErr: %v\nM2InvalidErr: %v\n",
			noFileErr,
			X2InvalidErr,
			T2InvalidErr,
			F2InvalidErr,
			M2InvalidErr,
		)
	}
}

func TestCopyNumberToInvalidReg(t *testing.T) {
	_, err := NewExa("./tests/COPY/NumtoInvalidReg.exa")

	if err == nil {
		t.Error("CopyNumberToInvalidReg failed, expected error to be thrown on attempting to copy a number to an invalid register")
	}
}

func TestCopyRegToNumber(t *testing.T) {
	_, err := NewExa("./tests/COPY/RegtoNum.exa")

	if err == nil {
		t.Error("CopyRegToNumber failed, expected error to be thrown on attempting to copy to a number")
	}
}

func TestCopyNumberToReg(t *testing.T) {
	N2M := Unwrap(NewExa("./tests/COPY/NumtoM.exa"))
	N2X := Unwrap(NewExa("./tests/COPY/NumtoX.exa"))
	N2T := Unwrap(NewExa("./tests/COPY/NumtoT.exa"))
	_, N2FErr := NewExa("./tests/COPY/NumtoFNoFile.exa")

	N2MM := <-N2M.M
	if "1" != N2MM {
		t.Errorf("CopyNumberToReg failed, expected 1, got %s", N2MM)
	}

	if "1" != N2X.X {
		t.Errorf("CopyNumberToReg failed, expected 1, got %s", N2X.X)
	}

	if "1" != N2T.T {
		t.Errorf("CopyNumberToReg failed, expected 1, got %s", N2T.T)
	}

	if N2FErr == nil {
		t.Error("CopyNumberToReg failed, expected error to be thrown on attempting to copy to NULL file")
	}
}

func TestCopyNumberToNumber(t *testing.T) {
	_, err := NewExa("./tests/COPY/NumtoNum.exa")

	if err == nil {
		t.Error("CopyNumberToNumber failed, expected error to be thrown on attempting to copy to a number")
	}
}

func TestCopyRegToNull(t *testing.T) {
	_, err := NewExa("./tests/COPY/XtoNull.exa")

	if err == nil {
		t.Error("CopyRegToNull failed, expected error to be thrown on attempting to copy register without a destination")
	}
}

func TestCopyNumberToNull(t *testing.T) {
	_, err := NewExa("./tests/COPY/NumtoNull.exa")

	if err == nil {
		t.Error("CopyNumberToNull failed, expected error to be thrown on attempting to copy number without a destination")
	}
}

func TestCopyNumberOverflow(t *testing.T) {
	_, err := NewExa("./tests/COPY/NumberOverflow.exa")

	if err == nil {
		t.Error("CopyNumberOverflow failed, expected error to be thrown on attempting to copy number over 9999")
	}
}

func TestCopyNumberUnderflow(t *testing.T) {
	_, err := NewExa("./tests/COPY/NumberUnderflow.exa")

	if err == nil {
		t.Error("CopyNumberUnderflow failed, expected error to be thrown on attempting to copy number under -9999")
	}
}

func TestCopyRegToFile(t *testing.T) {
	exa := Unwrap(NewExa("./tests/COPY/FtoReg.exa"))
	filename := fmt.Sprintf("./%d", 400 + vm.Files)
	Expect(os.Remove(filename))

	if exa.T != "1" {
		t.Errorf("CopyRegToFile failed, expected 1, got %s", exa.T)
	}
}
// COPY

// FILE
func TestMakeDefault(t *testing.T) {
	exa := Unwrap(NewExa("./tests/MAKE/MakeDefault.exa"))
	filename := fmt.Sprintf("./%d", 400 + vm.Files)
	_, err := os.ReadFile(filename)
	Expect(os.Remove(filename))

	if exa.F != filename || err != nil {
		t.Errorf("MakeDefault failed, expected %s, got %s", filename, exa.F)
	}
}

func TestMakeWithName(t *testing.T) {
	exa := Unwrap(NewExa("./tests/MAKE/MakeWithName.exa"))
	_, err := os.ReadFile("test.txt")
	Expect(os.Remove("test.txt"))

	if exa.F != "test.txt" || err != nil {
		t.Errorf("MakeWithName failed, expected test.txt, got %s", exa.F)
	}
}
// FILE

// SEEK
func TestSeekBackwards(t *testing.T) {
	exa := Unwrap(NewExa("./tests/SEEK/SeekBackwards.exa"))
	filename := fmt.Sprintf("./%d", 400 + vm.Files)
	Expect(os.Remove(filename))

	if exa.X != "2" {
		t.Errorf("SeekBackwards failed, expected 2, got %s", exa.X)
	}
}

func TestSeekForwards(t *testing.T) {
	exa := Unwrap(NewExa("./tests/SEEK/SeekForward.exa"))
	filename := fmt.Sprintf("./%d", 400 + vm.Files)
	Expect(os.Remove(filename))

	if exa.X != "2" {
		t.Errorf("SeekForwards failed, expected 2, got %s", exa.X)
	}
}

func TestSeekJumpBeginning(t *testing.T) {
	exa := Unwrap(NewExa("./tests/SEEK/SeekJumpBeginning.exa"))
	filename := fmt.Sprintf("./%d", 400 + vm.Files)
	Expect(os.Remove(filename))

	if exa.X != "1" {
		t.Errorf("SeekJumpBeginning failed, expected 1, got %s", exa.X)
	}
}

func TestSeekJumpEnding(t *testing.T) {
	exa := Unwrap(NewExa("./tests/SEEK/SeekJumpEnding.exa"))
	filename := fmt.Sprintf("./%d", 400 + vm.Files)
	Expect(os.Remove(filename))

	if exa.X != "2" {
		t.Errorf("SeekJumpEnding failed, expected 2, got %s", exa.X)
	}
}

func TestSeekM(t *testing.T) {
	// TODO: can only be tested when REPL is implemented
}

func TestSeekT(t *testing.T) {
	exa := Unwrap(NewExa("./tests/SEEK/SeekT.exa"))
	filename := fmt.Sprintf("./%d", 400 + vm.Files)
	Expect(os.Remove(filename))

	if exa.X != "1" {
		t.Errorf("SeekT failed, expected 1, got %s", exa.X)
	}
}

func TestSeekX(t *testing.T) {
	exa := Unwrap(NewExa("./tests/SEEK/SeekX.exa"))
	filename := fmt.Sprintf("./%d", 400 + vm.Files)
	Expect(os.Remove(filename))

	if exa.X != "1" {
		t.Errorf("SeekX failed, expected 1, got %s", exa.X)
	}
}

func TestSeekOverflow(t *testing.T) {
	_, err := NewExa("./tests/SEEK/SeekOverflow.exa")
	filename := fmt.Sprintf("./%d", 400 + vm.Files)
	Expect(os.Remove(filename))

	if err == nil {
		t.Error("SeekOverflow failed, expected error to be thrown on value over 9999")
	}
}

func TestSeekUnderflows(t *testing.T) {
	_, err := NewExa("./tests/SEEK/SeekUnderflow.exa")
	filename := fmt.Sprintf("./%d", 400 + vm.Files)
	Expect(os.Remove(filename))

	if err == nil {
		t.Error("SeekUnderflow failed, expected error to be thrown on value under -9999")
	}
}
// SEEK

// GRAB
func TestGrabFakeFile(t *testing.T) {
	_, err := NewExa("./tests/GRAB/FakeFile.exa")

	if err == nil {
		t.Error("GrabFakeFile failed, expected error to be thrown on attempting to GRAB a non-existent file")
	}
}

func TestGrabM(t *testing.T) {
	// TODO: can only be tested when REPL is implemented
}

func TestGrabX(t *testing.T) {
	exa := Unwrap(NewExa("./tests/GRAB/X.exa"))

	if exa.F != "./tests/GRAB/file1" {
		t.Error(fmt.Sprintf("GrabX failed, expected ./tests/GRAB/file1, got %s", exa.F))
	}
}

func TestGrabT(t *testing.T) {
	exa := Unwrap(NewExa("./tests/GRAB/T.exa"))

	if exa.F != "./tests/GRAB/file1" {
		t.Error(fmt.Sprintf("GrabT failed, expected ./tests/GRAB/file1, got %s", exa.F))
	}
}

func TestGrabNoName(t *testing.T) {
	_, err := NewExa("./tests/GRAB/NoName.exa")

	if err == nil {
		t.Error("GrabNoName failed, expected error to be thrown on attempting to GRAB without a filename")
	}
}

func TestGrabOneFile(t *testing.T) {
	exa := Unwrap(NewExa("./tests/GRAB/OneFile.exa"))

	if exa.F != "./tests/GRAB/file1" {
		t.Error(fmt.Sprintf("GrabOneFile failed, expected ./tests/GRAB/file1, got %s", exa.F))
	}
}

func TestGrabTwoFiles(t *testing.T) {
	exa := Unwrap(NewExa("./tests/GRAB/TwoFiles.exa"))

	if exa.F != "./tests/GRAB/file2" {
		t.Error(fmt.Sprintf("GrabTwoFiles failed, expected ./tests/GRAB/file2, got %s", exa.F))
	}
}

// GRAB
