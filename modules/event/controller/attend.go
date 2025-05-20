package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/modules/event/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/gin-gonic/gin"
)

// Attend handles the HTTP POST request for a follower to attend an event.
// Expected route: POST /events/:event_id/attend
func (ctrl *controller) Attend(ctx *gin.Context) {
	eventID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(err)
		return
	}

	var req request.EventAttendRequest
	if err = ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}

	// Call the service to handle event attendance
	res, err := ctrl.eventService.Attend(ctx.Request.Context(), eventID, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
