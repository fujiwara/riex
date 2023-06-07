package riex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
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

func (app *Riex) RunForDummy(ctx context.Context, endTime time.Time) error {
	now := time.Now()
	ri := ReservedInstance{
		Service:      "DummyService",
		InstanceType: "dummy.t3.micro",
		Name:         "dummy reservation",
		Description:  "dummy description",
		Count:        1,
		StartTime:    endTime.Add(-time.Duration(1) * 24 * time.Hour * 365),
		EndTime:      endTime,
	}
	if now.Before(ri.EndTime) {
		ri.State = "active"
	} else {
		ri.State = "expired"
	}
	if app.isPrintable(ri) {
		return app.Print([]ReservedInstance{ri}, os.Stdout)
	}
	return nil
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
	return app.Print(ris, os.Stdout)
}

func (app *Riex) Print(ris ReservedInstances, w io.Writer) error {
	switch app.option.Format {
	case "json":
		return app.PrintJSON(ris, w)
	case "markdown":
		return app.PrintMarkdown(ris, w)
	case "tsv":
		return app.PrintTSV(ris, w)
	default:
		return app.PrintJSON(ris, w)
	}
}

func (app *Riex) PrintJSON(ris ReservedInstances, w io.Writer) error {
	for _, ri := range ris {
		// trucate time.Time to second
		ri.StartTime = ri.StartTime.Truncate(time.Second)
		ri.EndTime = ri.EndTime.Truncate(time.Second)
		data, err := json.Marshal(ri)
		if err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}
	return nil
}

func (app *Riex) PrintMarkdown(ris ReservedInstances, w io.Writer) error {
	if len(ris) == 0 {
		return nil
	}
	fmt.Fprintln(w, "| service | name | description | instance_type | count | start_time | end_time | state |")
	fmt.Fprintln(w, "| --- | --- | --- | --- | --- | --- | --- | --- |")
	for _, ri := range ris {
		fmt.Fprintf(w,
			"| %s | %s | %s | %s | %d | %s | %s | %s |\n",
			ri.Service, ri.Name, ri.Description, ri.InstanceType,
			ri.Count, ri.StartTime.Format(time.RFC3339), ri.EndTime.Format(time.RFC3339),
			ri.State,
		)
	}
	return nil
}

func (app *Riex) PrintTSV(ris ReservedInstances, w io.Writer) error {
	if len(ris) == 0 {
		return nil
	}
	fields := []string{"service", "name", "description", "instance_type", "count", "start_time", "end_time", "state"}
	header := strings.Join(fields, "\t")
	fmt.Fprintln(w, header)

	for _, ri := range ris {
		row := []string{
			ri.Service,
			ri.Name,
			ri.Description,
			ri.InstanceType,
			strconv.Itoa(ri.Count),
			ri.StartTime.Format(time.RFC3339),
			ri.EndTime.Format(time.RFC3339),
			ri.State,
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	return nil
}

func (app *Riex) isPrintable(ri ReservedInstance) bool {
	if !app.option.Recognized && len(app.option.RecognizedTags) > 0 {
		for key, recognizedValue := range app.option.RecognizedTags {
			if v, ok := ri.Tags[key]; ok && recognizedValue == v {
				return false
			}
		}
	}
	if app.option.Active && strings.ToLower(ri.State) == "active" {
		return true
	}
	if ri.EndTime.After(app.startTime) && ri.EndTime.Before(app.endTime) {
		return true
	}
	return false
}

type ReservedInstance struct {
	Service      string            `json:"service"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	InstanceType string            `json:"instance_type"`
	Count        int               `json:"count"`
	StartTime    time.Time         `json:"start_time"`
	EndTime      time.Time         `json:"end_time"`
	State        string            `json:"state"`
	Tags         map[string]string `json:"tags,omitempty"`
}

type ReservedInstances []ReservedInstance
