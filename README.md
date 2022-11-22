# riex

AWS RI expiration detector.

## Description

riex is a AWS Reserved Instance EXpiretion detector.

riex finds reserved instances of EC2, ElastiCache, RDS, Redshift and Opensearch that will be expired in specified days.

## Usage

```
Usage: riex <duration>

Arguments:
  <duration>    Show reserved instances will be expired in specified days

Flags:
  -h, --help           Show context-sensitive help.
      --active         Show active reserved instances.
      --expired=INT    Show expired reserved instances in specified days
```

`AWS_REGION` environment variable is required.

## LICENSE

MIT
