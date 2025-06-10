package imagekit

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gonews/config"
	"gonews/internal/core/domain/entity"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type ImageKitAdapter interface {
	UploadImage(req *entity.FileUploadEntity) (string, error)
}

type imageKitAdapter struct {
	cfg *config.Config
}

func NewImageKitAdapter(cfg *config.Config) ImageKitAdapter {
	return &imageKitAdapter{cfg: cfg}
}

func (ik *imageKitAdapter) UploadImage(req *entity.FileUploadEntity) (string, error) {
	file, err := os.Open(req.Path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	writer.WriteField("fileName", req.Name)
	writer.WriteField("publicKey", ik.cfg.IK.PublicKey)
	writer.WriteField("useUniqueFileName", "true")
	writer.WriteField("folder", "/content")

	part, err := writer.CreateFormFile("file", req.Name)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}
	writer.Close()

	reqUrl := "https://upload.imagekit.io/api/v1/files/upload"
	reqHttp, err := http.NewRequest("POST", reqUrl, &b)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	reqHttp.Header.Set("Content-Type", writer.FormDataContentType())
	encodedKey := base64.StdEncoding.EncodeToString([]byte(ik.cfg.IK.PrivateKey + ":"))
	reqHttp.Header.Set("Authorization", "Basic "+encodedKey)

	client := &http.Client{}
	resp, err := client.Do(reqHttp)
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s", string(bodyBytes))
	}

	var result struct {
		Url string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode error: %w", err)
	}

	return result.Url, nil
}

