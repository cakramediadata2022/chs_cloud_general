package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nfnt/resize"
)

var PublicPath string

func GetFilePath(UnitCode, FolderName string, Name string) (fullPath, fileName string, err error) {
	basePath := PublicPath
	// datePath := int(time.Now().Month())
	// yearPath := int(time.Now().Year())
	// dateYearPath := strconv.Itoa(datePath) + strconv.Itoa(yearPath)
	// Path := basePath + "/" + FolderName + "/"
	Path := fmt.Sprintf("%s/%s/%s/", basePath, FolderName, UnitCode)
	if err := CreateDirectoryIfNotExist(Path); err != nil {
		return "", "", err
	}
	now := time.Now()
	Path += fmt.Sprintf("%s_%d", Name, now.Unix())
	return Path, strings.Replace(Path, basePath, "", 1), nil
}

func CreateDirectoryIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func UploadImageUtility(file *multipart.FileHeader, UnitCode string, AWSMinioAccessKey, AWSMinioSecretKey, AWSMinioEndpoint, AWSEndpoint string) (filePath string, err error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}

	defer src.Close()

	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, src); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	var img image.Image
	if ext == ".png" {
		img, err = png.Decode(buf)
	} else if ext == ".jpg" || ext == ".jpeg" {
		img, err = jpeg.Decode(buf)
	} else {
		return "", fmt.Errorf("unsupported file format: %s", ext)
	}
	if err != nil {
		return "", err
	}

	m := resize.Resize(512, 0, img, resize.Lanczos3)
	path, imagePath, err := GetFilePath(UnitCode, "images/utility", "utility")
	if err != nil {
		return "", err
	}

	filename := path + ".png"
	imagePath += ".png"
	out, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	// Encode and save the image as PNG
	if err := png.Encode(out, m); err != nil {
		out.Close()
		return "", err
	}
	// write new image to file
	if err := png.Encode(out, m); err != nil {
		// master_data.SendResponse(global_var.ResponseCode.InvalidDataValue, fmt.Sprintf("upload file err: %s", err.Error()), nil, c)
		out.Close()
		return "", err
	}

	out.Close()
	endpoint, s3Client, err := AwsLoad(AWSMinioAccessKey, AWSMinioSecretKey, AWSMinioEndpoint, AWSEndpoint)
	if err != nil {
		fmt.Println("Error Initialize:", err.Error())
		return "", err
	}

	// Reopen the file for uploading to S3
	uploadFile, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	fileUrlPath := "pms-web/web-public" + imagePath
	object := s3.PutObjectInput{
		Bucket: aws.String("pms-storage"),
		Key:    aws.String(fileUrlPath),
		Body:   uploadFile,                // The object's contents.
		ACL:    aws.String("public-read"), // Defines Access-control List (ACL) permissions, such as private or public.
		Metadata: map[string]*string{ // Required. Defines metadata tags.
			"x-amz-meta-type":      aws.String("utility"),
			"x-amz-meta-unit-code": aws.String(UnitCode),
		},
	}

	// Step 5: Run the PutObject function with your parameters, catching for errors.
	_, err = s3Client.PutObject(&object)
	if err != nil {
		fmt.Println(err.Error())
		uploadFile.Close()
		return "", err
	}
	uploadFile.Close()
	// Delete the local file after successful upload
	err = os.Remove(filename)
	if err != nil {
		return "", fmt.Errorf("failed to delete local file: %s", err)
	}
	return endpoint + "/" + fileUrlPath, nil
}
func UploadImage(file *multipart.FileHeader, UnitCode string, FolderName string, AWSMinioAccessKey, AWSMinioSecretKey, AWSMinioEndpoint, AWSEndpoint, outputFormat string, MaxWidth uint) (filePath string, err error) {
	// Buka file yang diunggah
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Baca file ke dalam buffer
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, src); err != nil {
		return "", err
	}

	// Deteksi format file berdasarkan ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	var img image.Image

	switch ext {
	case ".png":
		img, err = png.Decode(buf)
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(buf)
	default:
		return "", fmt.Errorf("unsupported file format: %s", ext)
	}
	if err != nil {
		return "", err
	}

	// Resize gambar ke lebar 512px (tinggi menyesuaikan)
	m := resize.Resize(MaxWidth, 0, img, resize.Lanczos3)

	// Tentukan path & ekstensi file berdasarkan format output
	path, imagePath, err := GetFilePath(UnitCode, "images/"+FolderName, FolderName)
	if err != nil {
		return "", err
	}

	var filename string
	if outputFormat == "jpg" || outputFormat == "jpeg" {
		filename = path + ".jpg"
		imagePath += ".jpg"
	} else {
		filename = path + ".png"
		imagePath += ".png"
	}

	// Simpan gambar ke file sementara
	out, err := os.Create(filename)
	if err != nil {
		out.Close()
		return "", err
	}

	// Encode gambar sesuai format yang dipilih
	if outputFormat == "jpg" || outputFormat == "jpeg" {
		err = jpeg.Encode(out, m, &jpeg.Options{Quality: 80})
	} else {
		err = png.Encode(out, m)
	}
	if err != nil {
		return "", err
	}

	out.Close()
	// Inisialisasi koneksi ke AWS MinIO
	endpoint, s3Client, err := AwsLoad(AWSMinioAccessKey, AWSMinioSecretKey, AWSMinioEndpoint, AWSEndpoint)
	if err != nil {
		return "", err
	}

	// Buka kembali file untuk di-upload ke S3
	uploadFile, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	// Path penyimpanan di S3
	fileUrlPath := "pms-web/web-public" + imagePath
	object := s3.PutObjectInput{
		Bucket: aws.String("pms-storage"),
		Key:    aws.String(fileUrlPath),
		Body:   uploadFile,
		ACL:    aws.String("public-read"),
		Metadata: map[string]*string{
			"x-amz-meta-type":      aws.String(FolderName),
			"x-amz-meta-unit-code": aws.String(UnitCode),
		},
	}

	// Upload file ke S3
	_, err = s3Client.PutObject(&object)
	if err != nil {
		uploadFile.Close()
		return "", err
	}
	uploadFile.Close()
	// Hapus file lokal setelah berhasil di-upload
	if err := os.Remove(filename); err != nil {
		return "", fmt.Errorf("failed to delete local file: %s", err)
	}

	return endpoint + "/" + fileUrlPath, nil
}

// UploadImageBookingEngine is a function to upload and resize images specifically for the booking engine.
func UploadImageBookingEngine(file *multipart.FileHeader, UnitCode string, FolderName string, AWSMinioAccessKey, AWSMinioSecretKey, AWSMinioEndpoint, AWSEndpoint, outputFormat string, MaxWidth uint) (filePath string, err error) {
	// Buka file yang diunggah
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Baca file ke dalam buffer
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, src); err != nil {
		return "", err
	}

	// Deteksi format file berdasarkan ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	var img image.Image

	switch ext {
	case ".png":
		img, err = png.Decode(buf)
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(buf)
	default:
		return "", fmt.Errorf("unsupported file format: %s", ext)
	}
	if err != nil {
		return "", err
	}

	// Resize gambar ke lebar 512px (tinggi menyesuaikan)
	m := resize.Resize(MaxWidth, 0, img, resize.Lanczos3)

	// Tentukan path & ekstensi file berdasarkan format output
	path, imagePath, err := GetFilePath(UnitCode, "images/"+FolderName, FolderName)
	if err != nil {
		return "", err
	}

	var filename string
	if outputFormat == "jpg" || outputFormat == "jpeg" {
		filename = path + ".jpg"
		imagePath += ".jpg"
	} else {
		filename = path + ".png"
		imagePath += ".png"
	}

	// Simpan gambar ke file sementara
	out, err := os.Create(filename)
	if err != nil {
		out.Close()
		return "", err
	}

	// Encode gambar sesuai format yang dipilih
	if outputFormat == "jpg" || outputFormat == "jpeg" {
		err = jpeg.Encode(out, m, &jpeg.Options{Quality: 80})
	} else {
		err = png.Encode(out, m)
	}
	if err != nil {
		return "", err
	}

	out.Close()
	// Inisialisasi koneksi ke AWS MinIO
	endpoint, s3Client, err := AwsLoad(AWSMinioAccessKey, AWSMinioSecretKey, AWSMinioEndpoint, AWSEndpoint)
	if err != nil {
		return "", err
	}

	// Buka kembali file untuk di-upload ke S3
	uploadFile, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	// Path penyimpanan di S3 untuk booking engine
	fileUrlPath := "cbe-web/web-public" + imagePath
	object := s3.PutObjectInput{
		Bucket: aws.String("pms-storage"),
		Key:    aws.String(fileUrlPath),
		Body:   uploadFile,
		ACL:    aws.String("public-read"),
		Metadata: map[string]*string{
			"x-amz-meta-type":      aws.String(FolderName),
			"x-amz-meta-unit-code": aws.String(UnitCode),
		},
	}

	// Upload file ke S3
	_, err = s3Client.PutObject(&object)
	if err != nil {
		uploadFile.Close()
		return "", err
	}
	uploadFile.Close()
	// Hapus file lokal setelah berhasil di-upload
	if err := os.Remove(filename); err != nil {
		return "", fmt.Errorf("failed to delete local file: %s", err)
	}

	return endpoint + "/" + fileUrlPath, nil
}
