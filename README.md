# route53-updater

A simple [AWS Route53](https://aws.amazon.com/route53/) updating tool,
especially designed for EC2 instances. This tool can fetch the public IPv4
address for the EC2 Instance Metadata Service (IMDSv1 or IMDSv2), and use it
to update an existing Route53 Hosted Zone.

## Usage

<pre>
  -l int
        optional DNS record TTL (default 15)
  -r string
        required DNS record being set
  -t string
        required DNS record type
  -v2
        optional use Instance Metadata Service version 2 (IMDSv2)
  -z string
        required Route53 hosted zone ID

Example:

$ route53-updater -v2 -r node01.mydomain.com -t A -z X0X0X0X0X0X0X0X0X
Route53 updated, change id: /change/A1A1A1A1A1A1A1A1A, status: PENDING

</pre>

This expects an EC2 instance with an instance profile and IAM role attached.
See [AWS documentation](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html)
for details. Here is an example limited IAM policy to get started with
route53-updater:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "route53:ChangeResourceRecordSets"
            ],
            "Resource": [
                "arn:aws:route53:::hostedzone/X0X0X0X0X0X0X0X0X"
            ]
        }
    ]
}
```

## Potential Improvements

- [ ] better error handling vs. lazy panics
- [ ] monitor change status
- [ ] optional AWS keys as opposed to instance profile

