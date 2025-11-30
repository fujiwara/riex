# riex

AWS RI expiration detector.

## Description

riex is a AWS Reserved Instance EXpiretion detector.

riex finds reserved instances of EC2, ElastiCache, RDS, Redshift and Opensearch that will be expired in specified days.

## Installation

```console
$ brew install fujiwara/tap/riex
```

or [Binary releases](https://github.com/fujiwara/riex/releases).

## Usage

```
Usage: riex <days>

Arguments:
  <days>    Show reserved instances that will be expired within specified days.

Flags:
  -h, --help                         Show context-sensitive help.
      --active                       Show active reserved instances.
      --pending                      Show payment-pending reserved instances.
      --expired=INT                  Show reserved instances expired in the last specified days.
      --format="json"                Output format.(json, markdown, tsv)
      --dummy-output                 Dummy output for testing.
      --dummy-end-time=TIME          Endtime for testing. works only with --dummy-output.
      --ignore-tags=KEY=VALUE;...    Resource tag for ignore RI.
      --local-time                   Use local time for output.
```

`AWS_REGION` environment variable is required.


## Examples

### Find RIs that will be expired within 30 days.

```console
$ riex 30
{"service":"Redshift","name":"140aad98-3ab6-435d-bcd4-60d1e65375bc","description":"","instance_type":"ra3.xlplus","count":1,"start_time":"2021-12-21T09:17:32.937Z","end_time":"2022-12-21T09:17:32.937Z","state":"active"}
```

### Find RIs that will be expired within 30 days or whose current state is active.

```console
$ riex 30 --active
{"service":"Redshift","name":"140aad98-3ab6-435d-bcd4-60d1e65375bc","description":"","instance_type":"ra3.xlplus","count":1,"start_time":"2021-12-21T09:17:32.937Z","end_time":"2022-12-21T09:17:32.937Z","state":"active"}
{"service":"ElastiCache","name":"ri-2022-08-18-07-42-42-976","description":"redis","instance_type":"cache.r6g.large","count":2,"start_time":"2022-08-18T07:43:00.276Z","end_time":"2023-08-18T07:43:00.276Z","state":"active"}
```

### Find RIs that will be expired within 30 days or already expired last 60 days.

```
$ riex 30 --expired 60
{"service":"RDS","name":"prod-ce-8x-2","description":"aurora-mysql","instance_type":"db.r6g.8xlarge","count":1,"start_time":"2021-10-25T05:31:59.456Z","end_time":"2022-10-25T05:31:59.456Z","state":"retired"}
{"service":"Redshift","name":"140aad98-3ab6-435d-bcd4-60d1e65375bc","description":"","instance_type":"ra3.xlplus","count":1,"start_time":"2021-12-21T09:17:32.937Z","end_time":"2022-12-21T09:17:32.937Z","state":"active"}
```

### Find RIs that will be expired within 30 days, output format as markdown.

```console
$ riex 30 --format markdown
| service | name | description | instance_type | count | start_time | end_time | state |
| --- | --- | --- | --- | --- | --- | --- | --- |
| RDS | prod-ce-8x-2 | aurora-mysql | db.r6g.8xlarge | 1 | 2022-10-14T08:09:30Z | 2023-10-14T08:09:30Z | active |
| Redshift | c36868e7-5421-41d0-ab87-841a0d162d1f |  | ra3.xlplus | 1 | 2022-12-21T08:02:18Z | 2023-12-21T08:02:18Z | active |
| ElastiCache | ri-2023-08-21-01-24-15-492 | redis | cache.r6g.large | 2 | 2023-08-21T01:24:40Z | 2024-08-20T01:24:40Z | active |
```

### `--ignore-tags KEY=VALUE`

`--ignore-tags` option can be used to ignore RIs with specified tags.

```console
$ riex 30 --ignore-tags RENEWED=true
```

Reserved instances with `RENEWED=true` tag will be ignored.

Note: Reserved instances of OpenSearch and Redshift do not support tags, so `--ignore-tags` option is not work.

### GitHub Actions

`fujiwara/riex@v0` composite action can be used in GitHub Actions.
This action checks RI expiration and create an issue if RI will be expired within specified days. (default: 30 days)

```yaml
name: Check RI Expiration

on:
  schedule:
    - cron: '0 0 * * *'

jobs:
  check_ri_expiration:
    runs-on: ubuntu-latest
    steps:
      - name: Check RI expiration
        uses: fujiwara/riex@v0
        with:
          days_left: '30'
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }} # or use aws-actions/configure-aws-credentials in before step
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          AWS_REGION: 'your-aws-region'
```

## LICENSE

MIT
