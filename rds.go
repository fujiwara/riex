package riex

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func (app *Riex) detectRDS(ctx context.Context) (*ReservedInstances, error) {
	pager := rds.NewDescribeReservedDBInstancesPaginator(
		app.rds, &rds.DescribeReservedDBInstancesInput{},
	)
	ris := make(ReservedInstances, 0, 100)
	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, ins := range page.ReservedDBInstances {
			ri := ReservedInstance{
				Service:      "RDS",
				InstanceType: aws.ToString(ins.DBInstanceClass),
				Name:         aws.ToString(ins.ReservedDBInstanceId),
				Description:  aws.ToString(ins.ProductDescription),
				Count:        int(ins.DBInstanceCount),
				StartTime:    aws.ToTime(ins.StartTime),
				EndTime:      ins.StartTime.Add(time.Second * time.Duration(ins.Duration)),
				State:        aws.ToString(ins.State),
			}
			if app.isPrintable(ri) {
				ris = append(ris, ri)
			}
		}
	}
	return &ris, nil
}
