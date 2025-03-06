package manage

import (
	commonRes "server/model/common/response"
	manageReq "server/model/manage/request"

	"github.com/gin-gonic/gin"
)

type CarLogApi struct{}

func (l *CarLogApi) GetCarLogList(c *gin.Context) {
	var carLogPageInfo manageReq.SearchCarLog
	if err := c.ShouldBindJSON(&carLogPageInfo); err != nil {
		commonRes.FailReq(err.Error(), c)
		return
	}

	if list, total, err := carLogService.GetCarLogList(c, carLogPageInfo); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
	} else {
		commonRes.OkWithDetailed(commonRes.PageResult{
			Page:     carLogPageInfo.Page,
			PageSize: carLogPageInfo.PageSize,
			List:     list,
			Total:    total,
		}, "获取成功", c)
		return
	}
}
