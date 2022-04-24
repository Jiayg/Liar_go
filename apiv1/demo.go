package apiv1

import "github.com/gogf/gf/v2/frame/g"

type DmReq struct {
	g.Meta `path:"/demo" tags:"Demo" method:"get" summary:"demo api"`
}
type DmRes struct {
	g.Meta `mime:"text/html"`
	Name   string `json:"name"`
}
