package metadata

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
)

// Metadata type is the AWS EC2 metadata API facade
type Metadata struct {
	Metadata EC2MetadataAPI
}

func New(session client.ConfigProvider) *Metadata {
	return &Metadata{
		ec2metadata.New(session),
	}
}

// EC2MetadataAPI defines a missing interface from the AWS SDK
type EC2MetadataAPI interface {
	GetInstanceIdentityDocument() (ec2metadata.EC2InstanceIdentityDocument, error)
}

// GetCurrentInstanceIDAndRegion uses EC2 Metadata to get current EC2 instance ID and AWS Region
func (m *Metadata) GetCurrentInstanceIDAndRegion() (string, string, error) {
	instanceInfo, err := m.Metadata.GetInstanceIdentityDocument()
	if err != nil {
		return "", "", err
	}

	return instanceInfo.InstanceID, instanceInfo.Region, nil
}
