package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"cloud.google.com/go/storage"
	"github.com/gofrs/uuid"
	"github.com/joho/godotenv"
)

type Response struct {
	ImageURL string `json:"image_url"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/upload", handleUpload)
	log.Fatal(http.ListenAndServe(":9292", nil))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	f, fh, err := r.FormFile("file")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer f.Close()

	// TODO: どっかにまとめる
	bucket := os.Getenv("GCS_BUCKET")
	projectID := os.Getenv("PROJECT_ID")
	name := uuid.Must(uuid.NewV4()).String() + path.Ext(fh.Filename)
	public := true

	ctx := context.Background()
	_, objAttrs, err := upload(ctx, f, projectID, bucket, name, public)
	if err != nil {
		switch err {
		case storage.ErrBucketNotExist:
			log.Fatal("Please create the bucket first e.g. with `gsutil mb`")
		default:
			log.Fatal(err)
		}
	}

	log.Printf("URL: %s", objectURL(objAttrs))
	log.Printf("Size: %d", objAttrs.Size)
	log.Printf("MD5: %x", objAttrs.MD5)
	log.Printf("objAttrs: %+v", objAttrs)

	imageURL := Response{objectURL(objAttrs)}
	res, err := json.Marshal(imageURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func objectURL(objAttrs *storage.ObjectAttrs) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)
}

func upload(ctx context.Context, r io.Reader, projectID, bucket, name string, public bool) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	bh := client.Bucket(bucket)
	// Next check if the bucket exists
	if _, err = bh.Attrs(ctx); err != nil {
		return nil, nil, err
	}

	obj := bh.Object(name)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, r); err != nil {
		return nil, nil, err
	}
	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	if public {
		if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return nil, nil, err
		}
	}

	attrs, err := obj.Attrs(ctx)
	return obj, attrs, err
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST")
}
