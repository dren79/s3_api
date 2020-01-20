package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func getBuckets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a session in chosen area (uses credentials from file)
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("eu-west-1")},
		)

		// Create S3 service client
		svc := s3.New(sess)

		// Use ListBuckets provided function
		result, err := svc.ListBuckets(nil)
		if err != nil {
			exitErrorf("Unable to list buckets, %v", err)
		}

		// Convert byte array to JSON using Marshal
		response, err := json.Marshal(result.Buckets)
		if err != nil {
			exitErrorf("Unable to marshal, %v", err)
		}

		writeJsonResponse(w, response, http.StatusOK)
	}
}

func getObjects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a session in chosen area (uses credentials from file)
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("eu-west-1")},
		)

		// Create S3 service client
		svc := s3.New(sess)

		// Use ListBuckets provided function
		result, err := svc.ListBuckets(nil)
		if err != nil {
			exitErrorf("Unable to list buckets, %v", err)
		}

		// Structures for the responses
		type Contents struct {
			ObjectName   string `json:"Key"`
			StorageClass string `json:"StorageClass"`
		}

		type Bucket struct {
			BucketName string     `json:"Name"`
			Contents   []Contents `json:"Contents"`
		}

		//var buffer bytes.Buffer
		var output []Bucket
		for _, b := range result.Buckets {
			var contents []Contents
			objResp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(aws.StringValue(b.Name))})
			if err != nil {
				exitErrorf("Unable to list items in bucket %q, %v", aws.StringValue(b.Name), err)
			}
			for _, obj := range objResp.Contents {
				contents = append(contents, Contents{
					ObjectName:   *obj.Key,
					StorageClass: *obj.StorageClass,
				})
			}
			if err != nil {
				exitErrorf("Unable to list items in bucket %q, %v", aws.StringValue(b.Name), err)
			}
			bucket := Bucket{
				BucketName: *b.Name,
				Contents:   contents,
			}
			output = append(output, bucket)
		}
		response, _ := json.Marshal(output)
		writeJsonResponse(w, response, http.StatusOK)

	}
}

func getAllBucketsEncryption() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.Header.Get("name")
		if name != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Only ADMIN can view encryption details"))
			return
		}
		// Create a session in chosen area (uses credentials from file)
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("eu-west-1")},
		)

		// Create S3 service client
		svc := s3.New(sess)

		result, err := svc.ListBuckets(nil)
		if err != nil {
			exitErrorf("Unable to list buckets, %v", err)
		}

		// Structures for the responses
		type Encryption struct {
			SSEAlgorithm string `json:"SSEAlgorithm"`
		}

		type Bucket struct {
			BucketName string       `json:"Name"`
			Encryption []Encryption `json:"SSEAlgorithm"`
		}

		var output []Bucket
		for _, b := range result.Buckets {
			var encryption []Encryption
			req, resp := svc.GetBucketEncryptionRequest(&s3.GetBucketEncryptionInput{
				Bucket: aws.String(aws.StringValue(b.Name)),
			})
			err := req.Send()
			if err == nil { // resp is now filled
				for _, enc := range resp.ServerSideEncryptionConfiguration.Rules {
					encryption = append(encryption, Encryption{
						SSEAlgorithm: *enc.ApplyServerSideEncryptionByDefault.SSEAlgorithm,
					})
				}
			}
			bucket := Bucket{
				BucketName: *b.Name,
				Encryption: encryption,
			}
			output = append(output, bucket)
		}
		response, _ := json.Marshal(output)
		writeJsonResponse(w, response, http.StatusOK)
	}
}

func getBucketEncryption() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.Header.Get("name")
		if name != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Only ADMIN can view encryption details"))
			return
		}
		// Create a session in chosen area (uses credentials from file)
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String("eu-west-1")},
		)

		// Create S3 service client
		svc := s3.New(sess)

		keys, ok := r.URL.Query()["key"]

		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}

		// Query()["key"] will return an array of items,
		// we only want the single item.
		key := keys[0]

		//var buffer bytes.Buffer
		// Structures for the responses
		type Encryption struct {
			SSEAlgorithm string `json:"SSEAlgorithm"`
		}

		type Bucket struct {
			BucketName string       `json:"Name"`
			Encryption []Encryption `json:"SSEAlgorithm"`
		}

		var output []Bucket
		var encryption []Encryption
		req, resp := svc.GetBucketEncryptionRequest(&s3.GetBucketEncryptionInput{
			Bucket: aws.String(key),
		})
		err := req.Send()
		if err == nil { // resp is now filled
			for _, enc := range resp.ServerSideEncryptionConfiguration.Rules {
				encryption = append(encryption, Encryption{
					SSEAlgorithm: *enc.ApplyServerSideEncryptionByDefault.SSEAlgorithm,
				})
			}
			bucket := Bucket{
				BucketName: key,
				Encryption: encryption,
			}
			output = append(output, bucket)
		}
		response, _ := json.Marshal(output)
		writeJsonResponse(w, response, http.StatusOK)
	}
}

func getBucketObjectEncryption() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.Header.Get("name")
		if name != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Only ADMIN can view encryption details"))
			return
		}
		// Create a session in chosen area (uses credentials from file)
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("eu-west-1")},
		)

		// Create S3 service client
		svc := s3.New(sess)

		// Use ListBuckets provided function
		result, err := svc.ListBuckets(nil)
		if err != nil {
			exitErrorf("Unable to list buckets, %v", err)
		}
		// Structures for the responses
		type Encryption struct {
			ServerSideEncryption string `json:"ServerSideEncryption"`
		}

		type Contents struct {
			ObjectName string       `json:"Key"`
			Encryption []Encryption `json:"Encryption"`
		}

		type Bucket struct {
			BucketName string     `json:"Name"`
			Contents   []Contents `json:"Contents"`
		}

		var output []Bucket
		for _, b := range result.Buckets {
			var contents []Contents
			objResp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(aws.StringValue(b.Name))})
			if err != nil {
				exitErrorf("Unable to list items in bucket %q, %v", aws.StringValue(b.Name), err)
			}
			var encryption []Encryption
			for _, o := range objResp.Contents {
				var oName = aws.StringValue(o.Key)
				var bName = aws.StringValue(b.Name)
				res, err := svc.HeadObject(&s3.HeadObjectInput{
					Bucket: aws.String(bName),
					Key:    aws.String(oName),
				})
				if err == nil { // resp is now filled
					if len(*res.ServerSideEncryption) > 0 {
						encryption = append(encryption, Encryption{
							ServerSideEncryption: *res.ServerSideEncryption,
						})
					} else {
						encryption = append(encryption, Encryption{
							ServerSideEncryption: "Nill",
						})
					}
				}
				contents = append(contents, Contents{
					ObjectName: *o.Key,
					Encryption: encryption,
				})
			}
			output = append(output, Bucket{
				BucketName: *b.Name,
				Contents:   contents,
			})
		}

		response, _ := json.Marshal(output)
		writeJsonResponse(w, response, http.StatusOK)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func writeJsonResponse(w http.ResponseWriter, content []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(content)
}

func main() {
	router := ConfigureRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
func ConfigureRouter() *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/auth", authenticate).Methods(http.MethodPost)
	api.Handle("/list-all-buckets", authMiddleware(getBuckets())).Methods(http.MethodGet)
	api.Handle("/list-all-objects", authMiddleware(getObjects())).Methods(http.MethodGet)
	api.Handle("/get-all-buckets-encryption", authMiddleware(getAllBucketsEncryption())).Methods(http.MethodGet)
	api.Handle("/get-specific-bucket-encryption", authMiddleware(getBucketEncryption())).Methods(http.MethodGet)
	api.Handle("/get-bucket-objects-encryption", authMiddleware(getBucketObjectEncryption())).Methods(http.MethodGet)
	api.HandleFunc("", notFound)
	return api
}