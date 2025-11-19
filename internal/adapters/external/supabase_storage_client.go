package external

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

const (
	supabaseStoragePath = "/storage/v1/object"
	bucketName          = "meal-photos"
	maxUploadRetries    = 3
	uploadRetryDelay    = time.Second * 2
)

// SupabaseStorageClient handles file storage operations with Supabase
type SupabaseStorageClient struct {
	projectURL string
	anonKey    string
	httpClient *http.Client
}

// NewSupabaseStorageClient creates a new Supabase storage client
func NewSupabaseStorageClient(projectURL, anonKey string) *SupabaseStorageClient {
	return &SupabaseStorageClient{
		projectURL: projectURL,
		anonKey:    anonKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// UploadImage uploads an image to Supabase storage and returns the public URL
func (c *SupabaseStorageClient) UploadImage(ctx context.Context, userID string, imageData []byte, filename string) (string, error) {
	// Generate unique path: user_id/timestamp_filename
	timestamp := time.Now().Unix()
	extension := filepath.Ext(filename)
	if extension == "" {
		extension = ".jpg"
	}
	objectPath := fmt.Sprintf("%s/%d%s", userID, timestamp, extension)

	var lastErr error
	for attempt := 0; attempt < maxUploadRetries; attempt++ {
		if attempt > 0 {
			log.Printf("[Supabase] Upload retry attempt %d/%d after error: %v", attempt+1, maxUploadRetries, lastErr)
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(uploadRetryDelay * time.Duration(attempt)):
			}
		}

		url, err := c.doUpload(ctx, objectPath, imageData)
		if err == nil {
			log.Printf("[Supabase] Successfully uploaded image: %s", objectPath)
			return url, nil
		}

		lastErr = err

		// Don't retry on context errors
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
	}

	return "", fmt.Errorf("failed to upload after %d retries: %w", maxUploadRetries, lastErr)
}

// doUpload performs the actual upload
func (c *SupabaseStorageClient) doUpload(ctx context.Context, objectPath string, imageData []byte) (string, error) {
	url := fmt.Sprintf("%s%s/%s/%s", c.projectURL, supabaseStoragePath, bucketName, objectPath)

	log.Printf("[Supabase] Uploading to: %s (size: %d bytes)", objectPath, len(imageData))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(imageData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.anonKey)
	req.Header.Set("apikey", c.anonKey)
	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("x-upsert", "true") // Overwrite if exists

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Generate public URL
	publicURL := c.GetPublicURL(objectPath)
	return publicURL, nil
}

// GetPublicURL generates a public URL for an object
func (c *SupabaseStorageClient) GetPublicURL(objectPath string) string {
	return fmt.Sprintf("%s%s/public/%s/%s", c.projectURL, supabaseStoragePath, bucketName, objectPath)
}

// DeleteImage deletes an image from Supabase storage
func (c *SupabaseStorageClient) DeleteImage(ctx context.Context, objectPath string) error {
	url := fmt.Sprintf("%s%s/%s/%s", c.projectURL, supabaseStoragePath, bucketName, objectPath)

	log.Printf("[Supabase] Deleting: %s", objectPath)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.anonKey)
	req.Header.Set("apikey", c.anonKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("[Supabase] Successfully deleted: %s", objectPath)
	return nil
}

// ListImages lists all images for a user
func (c *SupabaseStorageClient) ListImages(ctx context.Context, userID string) ([]string, error) {
	url := fmt.Sprintf("%s%s/list/%s?prefix=%s/", c.projectURL, supabaseStoragePath, bucketName, userID)

	log.Printf("[Supabase] Listing images for user: %s", userID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.anonKey)
	req.Header.Set("apikey", c.anonKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response to extract file paths
	// This is simplified - actual implementation would parse JSON response
	var files []string
	// TODO: Parse actual JSON response from Supabase
	return files, nil
}
