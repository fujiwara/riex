package riex

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
)

var Version string

type Riex struct {
	config      aws.Config
	ec2         *ec2.Client
	elasticache *elasticache.Client
	rds         *rds.Client
	redshift    *redshift.Client
	opensearch  *opensearch.Client

	option    *Option
	startTime time.Time
	endTime   time.Time
}

func New(ctx context.Context, opt *Option) (*Riex, error) {
	region := os.Getenv("AWS_REGION")
	awscfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	app := &Riex{
		config:      awscfg,
		ec2:         ec2.NewFromConfig(awscfg),
		elasticache: elasticache.NewFromConfig(awscfg),
		rds:         rds.NewFromConfig(awscfg),
		redshift:    redshift.NewFromConfig(awscfg),
		opensearch:  opensearch.NewFromConfig(awscfg),
		option:      opt,
		startTime:   now.Add(time.Duration(-opt.Expired) * 24 * time.Hour),
		endTime:     now.Add(time.Duration(opt.Days) * 24 * time.Hour),
	}
	return app, nil
}

func (app *Riex) Run(ctx context.Context) error {
	funcs := []func(context.Context) (*ReservedInstances, error){
		app.detectEC2,
		app.detectRDS,
		app.detectRedshift,
		app.detectElastiCache,
		app.detectOpensearch,
	}
	var eg errgroup.Group
	var mu sync.Mutex
	ris := make(ReservedInstances, 0, 100)
	for _, fn := range funcs {
		fn := fn
		eg.Go(func() error {
			if _ris, err := fn(ctx); err != nil {
				return err
			} else {
				mu.Lock()
				ris = append(ris, *_ris...)
				mu.Unlock()
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	sort.SliceStable(ris, func(i, j int) bool {
		return ris[i].StartTime.Before(ris[j].StartTime)
	})
	return ris.Print(os.Stdout)
}

func (app *Riex) isPrintable(ri ReservedInstance) bool {
	if app.option.Active && strings.ToLower(ri.State) == "active" {
		return true
	}
	if ri.EndTime.After(app.startTime) && ri.EndTime.Before(app.endTime) {
		return true
	}
	return false
}

type ReservedInstance struct {
	Service      string    `json:"service"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	InstanceType string    `json:"instance_type"`
	Count        int       `json:"count"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	State        string    `json:"state"`
}

type ReservedInstances []ReservedInstance

func (ris ReservedInstances) Print(w io.Writer) error {
	enc := json.NewEncoder(w)
	// enc.SetIndent("", "  ")
	for _, ri := range ris {
		if err := enc.Encode(ri); err != nil {
			return err
		}
	}
	return nil
}
