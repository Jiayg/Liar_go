package controller

import (
	"context"

	"github.com/jiayg/liar/apiv1"
)

var Demo = cDemo{}

type cDemo struct {
}

func (c *cDemo) Demo(ctx context.Context, req *apiv1.DmReq) (res *apiv1.DmRes, err error) {
	res = &apiv1.DmRes{Name: "赵四"}
	panic("demo wrong")
}
