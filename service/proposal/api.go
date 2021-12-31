package proposal

import (
	"app/infrastructure/log"
	"app/infrastructure/util"
	"app/service/proposal/db"
	"app/service/proposal/model"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleAPIApply(context *gin.Context) {
	var proposal model.Proposal
	err := json.NewDecoder(context.Request.Body).Decode(&proposal)
	if err != nil {
		log.Error(err.Error())
		context.String(http.StatusBadRequest, "parse request body failed")
		return
	}
	proposal.Id = uuid.New().String()
	proposal.CreateTime = util.GetCurrentShanghaiTime()
	proposal.UpdateTime = proposal.CreateTime
	proposal.Status = model.StatusChecking
	err = db.InsertProposal(&proposal)
	if err != nil {
		context.String(http.StatusInternalServerError, "insert proposal failed")
		log.Error(err.Error())
	}
	context.String(http.StatusOK, "OK")
}
