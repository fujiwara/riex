package riex

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
)

func (app *Riex) detectOpensearch(ctx context.Context) (*ReservedInstances, error) {
	pager := opensearch.NewDescribeReservedInstancesPaginator(
		app.opensearch, &opensearch.DescribeReservedInstancesInput{},
	)
	ris := make(ReservedInstances, 0, 100)
	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, node := range page.ReservedInstances {
			ri := ReservedInstance{
				Service:      "Opensearch",
				InstanceType: string(node.InstanceType),
				Name:         aws.ToString(node.ReservationName),
				Description:  aws.ToString(node.ReservedInstanceId),
				Count:        int(node.InstanceCount),
				StartTime:    aws.ToTime(node.StartTime),
				EndTime:      node.StartTime.Add(time.Second * time.Duration(node.Duration)),
				State:        aws.ToString(node.State),
				Tags:         make(map[string]string), // opeansearch reserved instance does not support tags
			}
			if app.isPrintable(ri) {
				ris = append(ris, ri)
			}
		}
	}
	return &ris, nil
}
