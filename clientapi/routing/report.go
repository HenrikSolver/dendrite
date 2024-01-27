package routing

import (
	"net/http"

	"github.com/matrix-org/dendrite/clientapi/httputil"
	"github.com/matrix-org/dendrite/roomserver/api"
	userapi "github.com/matrix-org/dendrite/userapi/api"
	"github.com/matrix-org/gomatrixserverlib/spec"
	"github.com/matrix-org/util"
	"github.com/sirupsen/logrus"
)

type ReportReq struct {
	Reason string `json:"reason"`
	Score  int    `json:"score"`
}

func Report(
	req *http.Request, rsAPI api.ClientRoomserverAPI, device *userapi.Device, roomID, eventID string,
) util.JSONResponse {

	deviceUserID, err := spec.NewUserID(device.UserID, true)
	if err != nil {
		return util.JSONResponse{
			Code: http.StatusForbidden,
			JSON: spec.Forbidden("userID doesn't have power level to change visibility"),
		}
	}
	queryReq := api.QueryMembershipForUserRequest{
		RoomID: roomID,
		UserID: *deviceUserID,
	}
	var queryRes api.QueryMembershipForUserResponse
	if err := rsAPI.QueryMembershipForUser(req.Context(), &queryReq, &queryRes); err != nil {
		util.GetLogger(req.Context()).WithError(err).Error("rsAPI.QueryMembershipsForRoom failed")
		return util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: spec.InternalServerError{},
		}
	}
	if !queryRes.IsInRoom {
		return util.JSONResponse{
			Code: http.StatusForbidden,
			JSON: spec.Forbidden("You aren't a member of this room."),
		}
	}
	var reportContent ReportReq
	resErr := httputil.UnmarshalJSONRequest(req, &reportContent)
	if resErr != nil {
		return *resErr
	}

	util.GetLogger(req.Context()).WithFields(logrus.Fields{
		"event_id": eventID,
		"room_id":  roomID,
		"reason":   reportContent.Reason,
		"score":    reportContent.Score,
	}).Info("Received report")

	return util.JSONResponse{
		Code: 200,
	}
}
