package restapi

// REST API handler for pairing
import (
	"net/http"

	"github.com/CPtung/mattercontroller/internal/matter/chiptool"
	"github.com/CPtung/mattercontroller/pkg/model"
	"github.com/gin-gonic/gin"
)

// REST API handler for pairing
func PostPairing(g *gin.Context) {
	req := model.MatterPairReqeust{}
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	controller := chiptool.New(g.Request.Context(), "")
	resp, err := controller.PairDevice(req.NodeID, req.PairCode)
	if err != nil {
		g.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	// It could be http.Created as well
	g.JSON(http.StatusOK, gin.H{"data": resp.Data})
}

// REST API handler for unpairing
func PostUnpairing(g *gin.Context) {
	deviceID := g.Param("deviceID")
	controller := chiptool.New(g.Request.Context(), "")
	resp, err := controller.UnpairDevice(deviceID)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"success": resp.Success})
}

// REST API handler for turning light on
func PutLightOnOff(g *gin.Context) {
	deviceID := g.Param("deviceID")
	body := model.MatterLightConfig{}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	controller := chiptool.New(g.Request.Context(), "")
	var err error
	if body.State == "on" {
		_, err = controller.TurnOn(deviceID)
	} else if body.State == "off" {
		_, err = controller.TurnOff(deviceID)
	} else {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := controller.GetOnOffState(deviceID)
	if err != nil || !resp.Success {
		g.JSON(http.StatusBadRequest, gin.H{"error": "failed to get state"})
		return
	}
	state, ok := resp.Data.(*model.MatterLightConfig)
	if !ok {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid state data"})
		return
	}
	g.JSON(http.StatusOK, gin.H{"isOn": state})
}

// REST API handler for getting state
func GetLightState(g *gin.Context) {
	deviceID := g.Param("deviceID")
	controller := chiptool.New(g.Request.Context(), "")
	resp, err := controller.GetOnOffState(deviceID)
	if err != nil || !resp.Success {
		g.JSON(http.StatusBadRequest, gin.H{"error": "failed to get state"})
		return
	}
	state, ok := resp.Data.(*model.MatterLightConfig)
	if !ok {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "invalid state data"})
		return
	}
	g.JSON(http.StatusOK, gin.H{"data": state})
}
