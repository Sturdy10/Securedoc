package handlers

import (
	"finalCode/models"
	"finalCode/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerPort interface {
	AddActivityHandlers(c *gin.Context)
	GETuserFileHandlers(c *gin.Context)
}

type handlerAdapter struct {
	s services.ServicesPort
}

func NewHanerAdapter(s services.ServicesPort) HandlerPort {
	return &handlerAdapter{s: s}
}

func (h *handlerAdapter) AddActivityHandlers(c *gin.Context) {
	var req models.RequestActivity

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseStatusF{Status: "Error", Message: err.Error()})
		return
	}

	ScdactBinary := req.ScdactBinary
	ScdactCommands := req.ScdactCommand

	exitCode, err := h.s.AddActivityServices(req, c, ScdactBinary, ScdactCommands)
	fmt.Println("test1", exitCode)
	fmt.Println("test2", err)
	if err != nil {
		// Handling error case
		var message string
		if exitCode != 0 {
			message = ErrorCommand(exitCode)} else {message = err.Error()
			}
		c.JSON(http.StatusInternalServerError, models.ResponseStatus{Status:  "Error", Message: map[string]string{"Finalcode_result": message},})
		h.s.DeleteByScdactReqIDServices(req) 
		return
	}
	
	c.JSON(http.StatusOK, models.ResponseStatus{
		Status:  "OK",
		Message: map[string]string{"Finalcode_result": "The operation completed successfully"},
	})
}

func (h *handlerAdapter) GETuserFileHandlers(c *gin.Context) {
	scdactID := c.Param("scdact_id")

	scdactBinary, err := h.s.GETuserFileByscdactIDServices(scdactID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseStatusF{
			Status:  "Error",
			Message: err.Error(),
		})
		return
	}

	// ตรวจสอบ Content-Type ของไฟล์
	contentType := http.DetectContentType(scdactBinary[0])

	// ตั้งค่า Content-Type ใน header
	c.Header("Content-Type", contentType)

	// ส่ง binary data กลับไปที่ client
	c.Data(http.StatusOK, "application/octet-stream", scdactBinary[0])
}
