package imgproxy

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/domain/image"
)

type UrlGenerator struct {
	minioBasePath string // e.g., "s3://images"
	proxyBaseURL  string // e.g., "https://imgproxy.example.com"
	keyBytes      []byte
	saltBytes     []byte
}

// NewUrlGenerator has args MinioBasePath e.g., "s3://images", ProxyBaseUrl e.g., "https://imgproxy.example.com"
func NewUrlGenerator(MinioBasePath string, ProxyBaseUrl string, KeyHex string, SaltHex string) *UrlGenerator {
	// Decode key and salt from their hex representation.
	keyBytes, err := hex.DecodeString(KeyHex)
	if err != nil {
		panic(err)
	}

	saltBytes, err := hex.DecodeString(SaltHex)
	if err != nil {
		panic(err)
	}

	return &UrlGenerator{
		minioBasePath: MinioBasePath,
		proxyBaseURL:  ProxyBaseUrl,
		keyBytes:      keyBytes,
		saltBytes:     saltBytes,
	}
}

func formatResizingType(resizingType image.ResizingType) (string, error) {
	switch resizingType {
	case image.Fit:
		return "fit", nil
	case image.Fill:
		return "fill", nil
	case image.Auto:
		return "auto", nil
	default:
		return "", errors.New("resizingType not supported")
	}
}

// GetSignedUrl generates the full signed URL for an image.
func (i UrlGenerator) GetSignedUrl(ctx context.Context, imageHash string, preset image.Preset) (string, error) {
	resizingTypeStr, err := formatResizingType(preset.ResizingType)

	// Build processing options string.
	processingOptions := fmt.Sprintf("rs:%s:%d:%d", resizingTypeStr, preset.Width, preset.Height)

	// Build the plain image URL from the Minio base path and the image hash.
	plainUrl := fmt.Sprintf("%s/%s", i.minioBasePath, imageHash)

	extension := "webp"

	// Build the path to be signed.
	// For a plain URL, the format is: "/{processing_options}/plain/{plain_url}.{extension}"
	path := fmt.Sprintf("/%s/plain/%s@%s", processingOptions, plainUrl, extension)

	// Create an HMAC hasher using the key.
	mac := hmac.New(sha256.New, i.keyBytes)
	// Write salt first, then the path.
	_, err = mac.Write(i.saltBytes)
	if err != nil {
		return "", fmt.Errorf("failed to write salt: %w", err)
	}
	_, err = mac.Write([]byte(path))
	if err != nil {
		return "", fmt.Errorf("failed to write path: %w", err)
	}

	// Compute the signature and encode it using URL-safe Base64 without padding.
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	// Build the final URL: {proxy_base_url}/{signature}{path}
	signedURL := fmt.Sprintf("%s/%s%s", i.proxyBaseURL, signature, path)
	return signedURL, nil
}

//Can be cached
//TODO Maybe use channel?

func (i UrlGenerator) GetSignedUrlBulk(ctx context.Context, imageHash string, presets []image.Preset) ([]string, error) {
	results := make([]string, len(presets))
	for j, preset := range presets {
		url, err := i.GetSignedUrl(ctx, imageHash, preset)
		if err != nil {
			return nil, fmt.Errorf("getting signed url for preset %s: %w", preset.Name, err)
		}
		results[j] = url
	}
	return results, nil
}
