package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/clothing"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type StyleSearch struct{
    clothing.Style
    StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
    EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    request.PageInfo
}
