// AWS resource fetcher
package internal

import (
	"context"
	"fmt"

	//Go breaks (unused import) if this is left in during a dry run "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type FetchedResource struct {
	ID     string
	Type   string
	Region string
	Tags   map[string]string
}

func FetchAWSResources(region string) ([]FetchedResource, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	resources := []FetchedResource{}

	// Fetch EC2 instances
	ec2Client := ec2.NewFromConfig(cfg)
	ec2Resp, err := ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err == nil {
		for _, r := range ec2Resp.Reservations {
			for _, i := range r.Instances {
				tags := map[string]string{}
				for _, t := range i.Tags {
					tags[*t.Key] = *t.Value
				}
				resources = append(resources, FetchedResource{
					ID:     *i.InstanceId,
					Type:   "aws_instance",
					Region: region,
					Tags:   tags,
				})
			}
		}
	}

	// Fetch S3 buckets
	s3Client := s3.NewFromConfig(cfg)
	s3Resp, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err == nil {
		for _, b := range s3Resp.Buckets {
			resources = append(resources, FetchedResource{
				ID:     *b.Name,
				Type:   "aws_s3_bucket",
				Region: region,
				Tags:   map[string]string{}, // S3 doesn't return tags here
			})
		}
	}

	// Fetch IAM Roles
	iamClient := iam.NewFromConfig(cfg)
	iamResp, err := iamClient.ListRoles(context.TODO(), &iam.ListRolesInput{})
	if err == nil {
		for _, r := range iamResp.Roles {
			resources = append(resources, FetchedResource{
				ID:     *r.RoleName,
				Type:   "aws_iam_role",
				Region: region,
				Tags:   map[string]string{}, // IAM tags need separate call
			})
		}
	}

	return resources, nil
}
