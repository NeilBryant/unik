package aws

import (
	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/emc-advanced-dev/unik/pkg/types"
	"github.com/layer-x/layerx-commons/lxerrors"
)

const UNIK_IMAGE_ID = "UNIK_IMAGE_ID"

func (p *AwsProvider) ListImages() ([]*types.Image, error) {
	imageIds := []*string{}
	for imageId := range p.state.GetImages() {
		imageIds = append(imageIds, aws.String(imageId))
	}
	param := &ec2.DescribeImagesInput{
		ImageIds: imageIds,
	}
	output, err := p.newEC2().DescribeImages(param)
	if err != nil {
		return nil, lxerrors.New("running ec2 describe images ", err)
	}
	images := []*types.Image{}
	for _, ec2Image := range output.Images {
		imageId := *ec2Image.ImageId
		image, ok := p.state.GetImages()[imageId]
		if !ok {
			logrus.WithFields(logrus.Fields{"ec2Image": ec2Image}).Errorf("found an image that unik has no record of")
			continue
		}
		images = append(images, image)
	}
	return images, nil
}
