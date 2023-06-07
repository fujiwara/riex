package riex

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func (app *Riex) detectEC2(ctx context.Context) (*ReservedInstances, error) {
	ris := make(ReservedInstances, 0, 100)
	out, err := app.ec2.DescribeReservedInstances(ctx, &ec2.DescribeReservedInstancesInput{})
	if err != nil {
		return nil, err
	}
	for _, ins := range out.ReservedInstances {
		tags := make(map[string]string, len(ins.Tags))
		for _, tag := range ins.Tags {
			tags[*tag.Key] = *tag.Value
		}
		ri := ReservedInstance{
			Service:      "EC2",
			InstanceType: string(ins.InstanceType),
			Name:         aws.ToString(ins.ReservedInstancesId),
			Description:  string(ins.ProductDescription),
			Count:        int(*ins.InstanceCount),
			StartTime:    aws.ToTime(ins.Start),
			EndTime:      aws.ToTime(ins.End),
			State:        string(ins.State),
			Tags:         tags,
		}
		if app.isPrintable(ri) {
			ris = append(ris, ri)
		}
	}
	return &ris, nil
}
