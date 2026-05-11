package minio

import (
	"fmt"
	"os"
	"context"
	"path/filepath"
	"strings"
	"github.com/minio/minio-go/v7"
)

type FaceRPCAMinIORepository struct {
	client 		*minio.Client
	bucketName 	string
}

func NewStorageRepo(client *minio.Client, bucketName string) *FaceRPCAMinIORepository {
	return &FaceRPCAMinIORepository{client: client, bucketName: bucketName}
}

func (r *FaceRPCAMinIORepository) ensureBucketExists(ctx context.Context) error {
	exists, err := r.client.BucketExists(ctx, r.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		if err := r.client.MakeBucket(ctx, r.bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", r.bucketName, err)
		}
	}
	return nil
}

func (r *FaceRPCAMinIORepository) DownloadModels(ctx context.Context, localDir string) error {
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return err
	}

	modelFiles := []string{"svm_model.pkl", "pca_transformer.pkl", "label_map.pkl"}

	for _, f := range modelFiles {
		s3Key := "models/" + f
		localPath := filepath.Join(localDir, f)

		if err := r.client.FGetObject(ctx, r.bucketName, s3Key, localPath, minio.GetObjectOptions{}); err != nil {
			return fmt.Errorf("failed to download %s: %w", f, err)
		}
	}

	return nil
}

func (r *FaceRPCAMinIORepository) UploadModel(ctx context.Context, localDir string) error {
	modelFiles := []string{"svm_model.pkl", "pca_transformer.pkl", "label_map.pkl"}

	for _, f := range modelFiles {
		localPath := filepath.Join(localDir, f)
		s3Key := "models/" + f

		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			fmt.Printf("⚠️ Local file %s not found. Did training succeed?\n", f)
			continue
		}

		if _, err := r.client.FPutObject(ctx, r.bucketName, s3Key, localPath, minio.PutObjectOptions{}); err != nil {
			return fmt.Errorf("failed to upload %s: %w", f, err)
		}
	}

	return nil
}

func (r *FaceRPCAMinIORepository) UploadDataset(ctx context.Context, datasetDir string) error {
	return filepath.Walk(datasetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if !strings.HasSuffix(info.Name(), ".jpg") {
			return nil
		}

		relPath, _ := filepath.Rel(datasetDir, path)
		s3Key := "dataset/" + filepath.ToSlash(relPath) // ensure forward slashes for S3

		_, err = r.client.FPutObject(ctx, r.bucketName, s3Key, path, minio.PutObjectOptions{})
		return err
	})
}

func (r *FaceRPCAMinIORepository) DownloadDataset(ctx context.Context, localDir string) error {
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return err
	}

	objectCh := r.client.ListObjects(ctx, r.bucketName, minio.ListObjectsOptions{
		Prefix:    "dataset/",
		Recursive: true,
	})

	count := 0
	for obj := range objectCh {
		if obj.Err != nil {
			return obj.Err
		}

		// Strip "dataset/" prefix → "1/0.jpg"
		relPath := strings.TrimPrefix(obj.Key, "dataset/")
		if relPath == "" {
			continue
		}

		localPath := filepath.Join(localDir, filepath.FromSlash(relPath))
		localSubDir := filepath.Dir(localPath)
		if err := os.MkdirAll(localSubDir, 0755); err != nil {
			return err
		}

		if err := r.client.FGetObject(ctx, r.bucketName, obj.Key, localPath, minio.GetObjectOptions{}); err != nil {
			return fmt.Errorf("failed to download %s: %w", obj.Key, err)
		}
		count++
	}

	fmt.Printf("✅ Downloaded %d files from dataset\n", count)
	return nil
}

func (r *FaceRPCAMinIORepository) DeleteUserFaces(ctx context.Context, userID int) error {
	prefix := fmt.Sprintf("dataset/%d/", userID)

	objectCh := r.client.ListObjects(ctx, r.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for obj := range objectCh {
		if obj.Err != nil {
			return obj.Err
		}
		if err := r.client.RemoveObject(ctx, r.bucketName, obj.Key, minio.RemoveObjectOptions{}); err != nil {
			return fmt.Errorf("failed to delete %s: %w", obj.Key, err)
		}
	}

	return nil
}

func (r *FaceRPCAMinIORepository) UploadUserFaces(ctx context.Context, userID int, localUserDir string) (int, error) {
	entries, err := os.ReadDir(localUserDir)
	if err != nil {
		return 0, fmt.Errorf("failed to read directory %s: %w", localUserDir, err)
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jpg") {
			continue
		}

		localPath := filepath.Join(localUserDir, entry.Name())
		s3Key := fmt.Sprintf("dataset/%d/%s", userID, entry.Name())

		if _, err := r.client.FPutObject(ctx, r.bucketName, s3Key, localPath, minio.PutObjectOptions{}); err != nil {
			return count, fmt.Errorf("failed to upload %s: %w", entry.Name(), err)
		}
		count++
	}

	return count, nil
}

func (r *FaceRPCAMinIORepository) GetModelTimestamp(ctx context.Context) (float64, error) {
	info, err := r.client.StatObject(ctx, r.bucketName, "models/svm_model.pkl", minio.StatObjectOptions{})
	if err != nil {
		return 0, nil
	}
	return float64(info.LastModified.Unix()), nil
}

func (r *FaceRPCAMinIORepository) GetDatasetTimestamp(ctx context.Context) (float64, error) {
	objectCh := r.client.ListObjects(ctx, r.bucketName, minio.ListObjectsOptions{
		Prefix:    "dataset/",
		Recursive: true,
	})

	var latest int64 = 0
	for obj := range objectCh {
		if obj.Err != nil {
			continue
		}
		if obj.LastModified.Unix() > latest {
			latest = obj.LastModified.Unix()
		}
	}
	return float64(latest), nil
}
