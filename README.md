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

## Potential Improvements

- [ ] better error handling vs. lazy panics
- [ ] monitor change status

