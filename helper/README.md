This is a helper fuunction to parse an AWS billing file and insert it into dynamodb

# Configuration

## AWS credentials

You must own AWS credential to use this.

```shell
export AWS_ACCESS_KEY_ID=xxxxx
export AWS_SECRET_ACCESS_KEY=xxxxx
export AWS_REGION=us-east-1
```

## Dynamodb

you must set an environment variable that points to the table name in dynamodb. This table must exist.

```shell
export PAASTEK_TABLENAME=pricing
```


# Usage

This scripts reads a [billing report](http://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/billing-reports.html) in csv format from stdin and adds it to dynamodb.
The first line must be the header and must contains the columns names:

```
identity/LineItemId,identity/TimeInterval,bill/InvoiceId,bill/BillingEntity,bill/BillType,...
...
```

exemple:

`cat daily-report.csv | ./helper`


