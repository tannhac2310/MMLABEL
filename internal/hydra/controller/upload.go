package controller

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/upload"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/swag/endpoint"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/swag/swagger"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type UploadController interface {
	Upload(c *gin.Context)
}

type uploadController struct {
	uploadService upload.Service
}

func RegisterUploadController(
	r *gin.RouterGroup,
	uploadService upload.Service,
) {
	g := r.Group("upload")

	var c UploadController = &uploadController{
		uploadService: uploadService,
	}

	routeutil.AddCustomEndpoint(
		g,
		"file",
		c.Upload,
		"upload file",
		endpoint.Tags(g.BasePath()),
		routeutil.Parameter(swagger.Parameter{In: "formData", Type: "file", Name: "file", Description: "file to process", Required: true}),
		routeutil.Parameter(swagger.Parameter{In: "formData", Type: "string", Name: "uploadType", Description: "upload type", Required: true}),
		endpoint.Response(http.StatusOK, &dto.UploadResponse{}, "Response"),
	)
}

func (u *uploadController) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	uploadType := c.PostForm("uploadType")
	// switch uploadType {
	//case "avatar":
	//	break
	//default:
	//	if !strings.HasPrefix(uploadType, "mobile-log-") {
	//		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage("uploadType must be avatar, mobile-log-*"))
	//		return
	//	}
	//}

	// check err
	f, err := fileHeader.Open()
	if err != nil {
		transportutil.Error(c, err)
		return
	}
	defer f.Close()

	buf := make([]byte, int(fileHeader.Size))
	_, err = f.Read(buf)
	if err != nil {
		transportutil.Error(c, err)
		return
	}
	fmt.Println(fileHeader.Filename)
	url, err := u.uploadService.Upload(c.Request.Context(), uploadType, fileHeader.Filename, bytes.NewBuffer(buf))
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.UploadResponse{
		URL: url,
	})
}
