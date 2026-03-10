package gormcngen_test

import (
	"testing"

	"github.com/yylego/gormcngen"
)

func TestShowDirectiveSkipGenerate(t *testing.T) {
	t.Log(gormcngen.DirectiveSkipGenerate)
}
