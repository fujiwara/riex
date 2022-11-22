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
  -h, --help           Show context-sensitive help.
      --active         Show active reserved instances.
      --expired=INT    Show reserved instances expired in the last specified days.
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

## LICENSE

MIT
