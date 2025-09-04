package fileapi

import (
	"bytes"
	"gmc-blog-server/config"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestRouter initializes a gin router for testing and configures temporary paths.
func setupTestRouter(t *testing.T) (*gin.Engine, func()) {
	// Create temporary directories for testing
	tmpDir, err := os.MkdirTemp("", "upload-test")
	assert.NoError(t, err)

	chunkDir := filepath.Join(tmpDir, "chunks")
	fileDir := filepath.Join(tmpDir, "files")
	os.Mkdir(chunkDir, 0755)
	os.Mkdir(fileDir, 0755)

	// Override config for test
	originalChunkPath := config.AppConfig.ChunkPath
	originalFilePath := config.AppConfig.FilePath
	config.AppConfig.ChunkPath = chunkDir
	config.AppConfig.FilePath = fileDir

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Teardown function to clean up resources
	cleanup := func() {
		config.AppConfig.ChunkPath = originalChunkPath
		config.AppConfig.FilePath = originalFilePath
		os.RemoveAll(tmpDir)
	}

	return router, cleanup
}

func TestChunkUploadAndMerge(t *testing.T) {
	// 1. Setup
	router, cleanup := setupTestRouter(t)
	defer cleanup()

	// Define API endpoints
	router.POST("/upload/chunk", func(c *gin.Context) {
		err := Check(c)
		assert.NoError(t, err)
	})
	router.POST("/upload/merge", func(c *gin.Context) {
		err := MergeChunks(c)
		assert.NoError(t, err)
	})

	// Test file data
	testFileName := "test-file.txt"
	testFileContent := "this is a test file for chunked upload functionality"
	chunkSize := 5
	totalChunks := (len(testFileContent) + chunkSize - 1) / chunkSize

	// 2. Simulate Chunk Uploads
	for i := 0; i < totalChunks; i++ {
		// Calculate chunk boundaries
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(testFileContent) {
			end = len(testFileContent)
		}
		chunkContent := testFileContent[start:end]

		// Create multipart form
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", testFileName)
		assert.NoError(t, err)
		_, err = io.WriteString(part, chunkContent)
		assert.NoError(t, err)

		writer.WriteField("fileName", testFileName)
		writer.WriteField("index", strconv.Itoa(i))
		writer.WriteField("totalChunks", strconv.Itoa(totalChunks))
		writer.Close()

		// Create and serve request
		req, _ := http.NewRequest(http.MethodPost, "/upload/chunk", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 3. Simulate Merge Request
	mergeBody := new(bytes.Buffer)
	writer := multipart.NewWriter(mergeBody)
	writer.WriteField("fileName", testFileName)
	writer.WriteField("totalChunks", strconv.Itoa(totalChunks))
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/upload/merge", mergeBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions for merge
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Verify Final File
	finalFilePath := filepath.Join(config.AppConfig.FilePath, testFileName)
	_, err := os.Stat(finalFilePath)
	assert.NoError(t, err, "Final merged file should exist")

	fileContent, err := os.ReadFile(finalFilePath)
	assert.NoError(t, err)
	assert.Equal(t, testFileContent, string(fileContent), "Merged file content should match original content")

	// 5. Verify Cleanup
	tempChunkDir := filepath.Join(config.AppConfig.ChunkPath, testFileName)
	_, err = os.Stat(tempChunkDir)
	assert.True(t, os.IsNotExist(err), "Temporary chunk directory should be deleted after merge")
}
