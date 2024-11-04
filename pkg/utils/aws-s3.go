package utils

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cakramediadata2022/chs_cloud_general/pkg/global_var"
)

func AwsLoad() (endpoint string, client *s3.S3, err error) {

	// Step 2: Define the parameters for the session you want to create.
	key := global_var.Config.AWS.MinioAccessKey // Access key pair.
	secret := global_var.Config.AWS.MinioSecretKey
	endpoint = global_var.Config.AWS.MinioEndpoint

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),       // Specifies your credentials.
		Endpoint:         aws.String("https://" + global_var.Config.AWS.Endpoint), // Your DigitalOcean Spaces endpoint.
		S3ForcePathStyle: aws.Bool(false),                                         // Use subdomain/virtual calling format.
		Region:           aws.String("sgp1"),                                      // Must match the region in your endpoint.
	}

	// Create a new session with the specified config
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		fmt.Println("Error creating session:", err.Error())
		return "", nil, err
	}
	s3Client := s3.New(newSession)

	return endpoint, s3Client, nil
}

func LoadReportTemplateList(Module string) ([]struct {
	Name string
	Path interface{}
}, error) {
	endpoint, s3Client, err := AwsLoad()
	if err != nil {
		fmt.Println("Error Initialize:", err.Error())
		return nil, err
	}
	// List objects in the specified bucket and prefix
	List := s3.ListObjectsV2Input{
		Bucket: aws.String("pms-storage"),
		Prefix: aws.String("pms-web/web-public/report/reports/"),
	}

	// Fetch the list of objects
	Out, err := s3Client.ListObjectsV2(&List)
	if err != nil {
		fmt.Println("Error listing objects:", err.Error())
		return nil, err
	}

	// The prefix to match files directly under 'reports/'
	reportsPrefix := "pms-web/web-public/report/reports/"
	directory := reportsPrefix + "CHS"

	// Iterate over the objects and filter out ones in subdirectories
	var templateList []struct {
		Name string
		Path interface{}
	}
	for _, item := range Out.Contents {
		// Check if the object's Key starts with the reports directory
		if strings.HasPrefix(*item.Key, directory+"/") {
			// Get the remaining part after the reports directory
			remainingPath := (*item.Key)[len(directory+"/"):]

			// Only include the file if there are no further '/' in the remaining path
			if !strings.Contains(remainingPath, "/") {
				template := strings.Replace(*item.Key, directory+"/", "", 1)
				templateList = append(templateList, struct {
					Name string
					Path interface{}
				}{
					Name: template,
					Path: fmt.Sprintf("%s%s", endpoint, *item.Key),
				})
				fmt.Println(*item.Key)
			}
		}
	}
	fmt.Println(templateList)
	return templateList, nil
}
