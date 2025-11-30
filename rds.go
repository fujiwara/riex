package riex

import (
	"context"
	"log"
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
			var tags map[string]string
			listTagsOutput, err := app.rds.ListTagsForResource(ctx, &rds.ListTagsForResourceInput{
				ResourceName: ins.ReservedDBInstanceArn,
			})
			if err == nil {
				tags = make(map[string]string, len(listTagsOutput.TagList))
				for _, tag := range listTagsOutput.TagList {
					tags[*tag.Key] = *tag.Value
				}
			} else {
				log.Println("[warn] failed to get tags for", ins.ReservedDBInstanceArn, ":", err)
				tags = make(map[string]string)
			}
			ri := ReservedInstance{
				Service:      "RDS",
				InstanceType: aws.ToString(ins.DBInstanceClass),
				Name:         aws.ToString(ins.ReservedDBInstanceId),
				Description:  aws.ToString(ins.ProductDescription),
				Count:        int(aws.ToInt32(ins.DBInstanceCount)),
				StartTime:    aws.ToTime(ins.StartTime),
				EndTime:      ins.StartTime.Add(time.Second * time.Duration(int64(aws.ToInt32(ins.Duration)))),
				State:        aws.ToString(ins.State),
				Tags:         tags,
			}
			if app.isPrintable(ri) {
				ris = append(ris, ri)
			}
		}
	}
	return &ris, nil
}
