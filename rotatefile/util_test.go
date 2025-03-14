package rotatefile

import (
	"errors"
	"testing"
)

func TestPrintErrln(t *testing.T) {
	printErrln("test", nil)
	printErrln("test", errors.New("an error"))
}
