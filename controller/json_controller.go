package controller

import (
	"net/http"

	"sojson/dto"
	"sojson/service"

	"github.com/gin-gonic/gin"
)

var (
	JSONController = &jsonController{}
)

// jsonController JSON处理控制器
type jsonController struct {
}

// UnescapeJSON 去除转义接口
func (ctrl *jsonController) UnescapeJSON(c *gin.Context) {
	var req dto.JSONRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.JSONResponse{
			Success: false,
			Error:   "请提供要处理的文本",
		})
		return
	}

	result, err := service.JSONProcessorService.UnescapeJSON(req.Text)
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
func (ctrl *jsonController) FormatJSON(c *gin.Context) {
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

	result, err := service.JSONProcessorService.FormatJSON(c.Request.Context(), req.Text, indent)
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
func (ctrl *jsonController) ProcessJSON(c *gin.Context) {
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

	result, err := service.JSONProcessorService.ProcessJSON(c.Request.Context(), req.Text, indent)
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
func (ctrl *jsonController) ValidateJSON(c *gin.Context) {

	var req dto.JSONRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidateResponse{
			Valid: false,
			Error: "请提供要验证的文本",
		})
		return
	}

	if err := service.JSONProcessorService.ValidateJSON(req.Text); err != nil {
		response := dto.ValidateResponse{
			Valid: false,
			Error: err.Error(),
		}
		c.JSON(http.StatusOK, response)
		return
	}

	response := dto.ValidateResponse{
		Valid:   true,
		Message: "JSON格式正确",
	}
	c.JSON(http.StatusOK, response)
}
