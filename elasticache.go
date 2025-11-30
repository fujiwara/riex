package riex

import (
	"context"
	"log"
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
			var tags map[string]string
			listTagsOutput, err := app.elasticache.ListTagsForResource(ctx, &elasticache.ListTagsForResourceInput{
				ResourceName: node.ReservationARN,
			})
			if err == nil {
				tags = make(map[string]string, len(listTagsOutput.TagList))
				for _, tag := range listTagsOutput.TagList {
					tags[*tag.Key] = *tag.Value
				}
			} else {
				log.Println("[warn] failed to get tags for", node.ReservationARN, ":", err)
				tags = make(map[string]string)
			}
			ri := ReservedInstance{
				Service:      "ElastiCache",
				InstanceType: aws.ToString(node.CacheNodeType),
				Name:         aws.ToString(node.ReservedCacheNodeId),
				Description:  aws.ToString(node.ProductDescription),
				Count:        int(aws.ToInt32(node.CacheNodeCount)),
				StartTime:    aws.ToTime(node.StartTime),
				EndTime:      node.StartTime.Add(time.Second * time.Duration(int64(aws.ToInt32(node.Duration)))),
				State:        aws.ToString(node.State),
				Tags:         tags,
			}
			if app.isPrintable(ri) {
				ris = append(ris, ri)
			}
		}
	}
	return &ris, nil

}
