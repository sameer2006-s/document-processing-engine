package document

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sameer2006-s/document-processing-engine/internal/auth"
)

type DocumentHandler struct {
	service *DocumentService
}

func NewDocumentHandler(service *DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

func (h *DocumentHandler) UploadFile(c *gin.Context) {
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	userId := userID.(uuid.UUID)
	file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "file is required",
        })
        return
    }

    src, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    defer src.Close()

	fileContent, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = h.service.UploadFile(&FileMetadata{
		UserID: userId,
		OriginalName: file.Filename,
		FileSize: int64(file.Size),
		ContentType: file.Header.Get("Content-Type"),
		BucketName: "documents",
		MinioKey: uuid.New().String(),
		}, fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (h *DocumentHandler) GetFile(c *gin.Context) {
	id := c.Param("id")
	fileContent, err := h.service.GetFileContent(uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	stat, err := h.service.GetFileMetadata(uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", strconv.FormatInt(stat.FileSize, 10))
	c.Header("Content-Disposition", "attachment; filename="+stat.OriginalName)
	c.Data(http.StatusOK, stat.ContentType, fileContent)
}