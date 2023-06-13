package riex

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
)

func (app *Riex) detectRedshift(ctx context.Context) (*ReservedInstances, error) {
	pager := redshift.NewDescribeReservedNodesPaginator(
		app.redshift, &redshift.DescribeReservedNodesInput{},
	)
	ris := make(ReservedInstances, 0, 100)
	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, node := range page.ReservedNodes {
			ri := ReservedInstance{
				Service:      "Redshift",
				InstanceType: aws.ToString(node.NodeType),
				Name:         aws.ToString(node.ReservedNodeId),
				Description:  "",
				Count:        int(node.NodeCount),
				StartTime:    aws.ToTime(node.StartTime),
				EndTime:      node.StartTime.Add(time.Second * time.Duration(node.Duration)),
				State:        aws.ToString(node.State),
				Tags:         make(map[string]string), // redshift reserved instance does not support tags
			}
			if app.isPrintable(ri) {
				ris = append(ris, ri)
			}
		}
	}
	return &ris, nil

}
