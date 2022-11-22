package riex

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
)

func (app *Riex) detectElastiCache(ctx context.Context) (*ReservedInstances, error) {
	pager := elasticache.NewDescribeReservedCacheNodesPaginator(
		app.elasticache, &elasticache.DescribeReservedCacheNodesInput{},
	)
	ris := make(ReservedInstances, 0, 100)
	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, node := range page.ReservedCacheNodes {
			ri := ReservedInstance{
				Service:      "ElastiCache",
				InstanceType: aws.ToString(node.CacheNodeType),
				Name:         aws.ToString(node.ReservedCacheNodeId),
				Description:  aws.ToString(node.ProductDescription),
				Count:        int(node.CacheNodeCount),
				StartTime:    aws.ToTime(node.StartTime),
				EndTime:      node.StartTime.Add(time.Second * time.Duration(node.Duration)),
				State:        aws.ToString(node.State),
			}
			if app.isPrintable(ri) {
				ris = append(ris, ri)
			}
		}
	}
	return &ris, nil

}
