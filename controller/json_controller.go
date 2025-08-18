package controller

import (
	"net/http"

	"sojson/dto"
	"sojson/service"

	"github.com/gin-gonic/gin"
)

// JSONController JSON处理控制器
type JSONController struct {
	jsonService *service.JSONProcessorService
}

// NewJSONController 创建JSON控制器
func NewJSONController(jsonService *service.JSONProcessorService) *JSONController {
	return &JSONController{
		jsonService: jsonService,
	}
}

// UnescapeJSON 去除转义接口
func (ctrl *JSONController) UnescapeJSON(c *gin.Context) {
	var req dto.JSONRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.JSONResponse{
			Success: false,
			Error:   "请提供要处理的文本",
		})
		return
	}

	result, err := ctrl.jsonService.UnescapeJSON(req.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.JSONResponse{
			Success: false,
			Error:   "去除转义失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Result:  result,
		Success: true,
	})
}

// FormatJSON 格式化JSON接口
func (ctrl *JSONController) FormatJSON(c *gin.Context) {
	var req dto.JSONRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.JSONResponse{
			Success: false,
			Error:   "请提供要处理的文本",
		})
		return
	}

	indent := req.Indent
	if indent == 0 {
		indent = 2
	}

	result, err := ctrl.jsonService.FormatJSON(req.Text, indent)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.JSONResponse{
			Success: false,
			Error:   "JSON格式错误: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Result:  result,
		Success: true,
	})
}

// ProcessJSON 完整处理接口
func (ctrl *JSONController) ProcessJSON(c *gin.Context) {
	var req dto.JSONRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.JSONResponse{
			Success: false,
			Error:   "请提供要处理的文本",
		})
		return
	}

	indent := req.Indent
	if indent == 0 {
		indent = 2
	}

	result, err := ctrl.jsonService.ProcessJSON(req.Text, indent)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.JSONResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.JSONResponse{
		Result:  result,
		Success: true,
	})
}

// ValidateJSON 验证JSON接口
func (ctrl *JSONController) ValidateJSON(c *gin.Context) {
	var req dto.JSONRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidateResponse{
			Valid: false,
			Error: "请提供要验证的文本",
		})
		return
	}

	if err := ctrl.jsonService.ValidateJSON(req.Text); err != nil {
		c.JSON(http.StatusOK, dto.ValidateResponse{
			Valid: false,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ValidateResponse{
		Valid:   true,
		Message: "JSON格式正确",
	})
}
