package riex_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/fujiwara/riex"
)

var TestRIs = riex.ReservedInstances{
	{
		Service:      "EC2",
		Name:         "m5.xlarge",
		Description:  "Linux/UNIX, Upfront",
		InstanceType: "m5.xlarge",
		Count:        2,
		StartTime:    time.Now(),
		EndTime:      time.Now().AddDate(1, 0, 0),
		State:        "active",
	},
	{
		Service:      "RDS",
		Name:         "db.m5.large",
		Description:  "MySQL, Upfront",
		InstanceType: "db.m5.large",
		Count:        1,
		StartTime:    time.Now().AddDate(0, 0, -7),
		EndTime:      time.Now().AddDate(1, 0, 0),
		State:        "retired",
	},
}

func TestPrintJSON(t *testing.T) {
	ris := TestRIs
	expectedOutput := `{"service":"EC2","name":"m5.xlarge","description":"Linux/UNIX, Upfront","instance_type":"m5.xlarge","count":2,"start_time":"` +
		ris[0].StartTime.Format(time.RFC3339) + `","end_time":"` + ris[0].EndTime.Format(time.RFC3339) +
		`","state":"active"}
{"service":"RDS","name":"db.m5.large","description":"MySQL, Upfront","instance_type":"db.m5.large","count":1,"start_time":"` +
		ris[1].StartTime.Format(time.RFC3339) + `","end_time":"` + ris[1].EndTime.Format(time.RFC3339) +
		`","state":"retired"}
`

	var buf bytes.Buffer
	app, err := riex.New(context.Background(), &riex.Option{Format: "json"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if err := app.Print(ris, &buf); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if buf.String() != expectedOutput {
		t.Errorf("Unexpected output: got '%s', expected '%s'", buf.String(), expectedOutput)
	}
}

func TestPrintTSV(t *testing.T) {
	ris := TestRIs
	expectedOutput := "service\tname\tdescription\tinstance_type\tcount\tstart_time\tend_time\tstate\n" +
		"EC2\tm5.xlarge\tLinux/UNIX, Upfront\tm5.xlarge\t2\t" + ris[0].StartTime.Format(time.RFC3339) +
		"\t" + ris[0].EndTime.Format(time.RFC3339) + "\tactive\n" +
		"RDS\tdb.m5.large\tMySQL, Upfront\tdb.m5.large\t1\t" + ris[1].StartTime.Format(time.RFC3339) +
		"\t" + ris[1].EndTime.Format(time.RFC3339) + "\tretired\n"

	var buf bytes.Buffer
	app, err := riex.New(context.Background(), &riex.Option{Format: "tsv"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if err := app.Print(ris, &buf); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if buf.String() != expectedOutput {
		t.Errorf("Unexpected output: got '%s', expected '%s'", buf.String(), expectedOutput)
	}
}

func TestPrintMarkdown(t *testing.T) {
	expectedOutput := "| service | name | description | instance_type | count | start_time | end_time | state |\n" +
		"| --- | --- | --- | --- | --- | --- | --- | --- |\n" +
		"| EC2 | m5.xlarge | Linux/UNIX, Upfront | m5.xlarge | 2 | " + TestRIs[0].StartTime.Format(time.RFC3339) +
		" | " + TestRIs[0].EndTime.Format(time.RFC3339) + " | active |\n" +
		"| RDS | db.m5.large | MySQL, Upfront | db.m5.large | 1 | " + TestRIs[1].StartTime.Format(time.RFC3339) +
		" | " + TestRIs[1].EndTime.Format(time.RFC3339) + " | retired |\n"

	var buf bytes.Buffer
	app, err := riex.New(context.Background(), &riex.Option{Format: "markdown"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if err := app.Print(TestRIs, &buf); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if buf.String() != expectedOutput {
		t.Errorf("Unexpected output: got '%s', expected '%s'", buf.String(), expectedOutput)
	}
}
