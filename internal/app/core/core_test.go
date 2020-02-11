package core

import (
	"context"
	"fmt"
	"git.cafebazaar.ir/bardia/lazyapi/pkg/appdetail"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCore_GetAppDetail(t *testing.T) {
	core := New()
	reply , err := core.GetAppDetail(context.Background(), &appdetail.GetAppDetailRequest{PackageName: "ir.divar"})
	if assert.NoError(t, err) {
		fmt.Printf("%v\n", reply)
	}
}