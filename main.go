package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

const metaURL string = "http://169.254.169.254/latest"

type meta struct {
	token string
}

func main() {
	var (
		rr  = flag.String("r", "", "required DNS record being set")
		rt  = flag.String("t", "", "required DNS record type")
		z   = flag.String("z", "", "required Route53 hosted zone ID")
		ttl = flag.Int64("l", 15, "optional DNS record TTL in seconds")
		v2  = flag.Bool("v2", false, "optional use Instance Metadata Service version 2 (IMDSv2)")
	)
	flag.Parse()

	if *rr == "" || *rt == "" || *z == "" {
		fmt.Printf("Must provide all required options. ")
		flag.Usage()
		os.Exit(1)
	}
	*rt = strings.ToUpper(*rt)

	m := new(meta)

	if *v2 {
		if err := m.setToken(); err != nil {
			panic(err)
		}
	}

	ipv4, err := m.getIPv4()
	if err != nil {
		panic(err)
	}

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	svc := route53.New(sess)

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name:            rr,
						Type:            rt,
						TTL:             ttl,
						ResourceRecords: []*route53.ResourceRecord{{Value: &ipv4}},
					},
				},
			},
		},
		HostedZoneId: z,
	}

	result, err := svc.ChangeResourceRecordSets(params)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		switch {
		case ok:
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				panic(
					fmt.Sprintf(
						"%v: %v, status code: %v, request ID: %v",
						reqErr.Code(),
						reqErr.Message(),
						reqErr.StatusCode(),
						reqErr.RequestID(),
					),
				)
			}
			panic(
				fmt.Sprintf(
				"%v: %v, additional error: %v",
				awsErr.Code(),
				awsErr.Message(),
				awsErr.OrigErr(),
			),
		)
		default:
			panic(err)
		}
	}

	ci := result.ChangeInfo
	fmt.Printf("Route53 updated, change id: %v, status: %v\n", *ci.Id, *ci.Status)
}

func (m *meta) setToken() error {
	c := new(http.Client)
	c.Timeout = time.Second * 5
	req, err := http.NewRequest("PUT", metaURL+"/api/token", nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "21600")

	m.token, err = processRequest(c, req)
	return err
}

func (m *meta) getIPv4() (string, error) {
	c := new(http.Client)
	c.Timeout = time.Second * 5
	req, err := http.NewRequest("GET", metaURL+"/meta-data/public-ipv4", nil)
	if err != nil {
		return "", err
	}

	if m.token != "" {
		req.Header.Set("X-aws-ec2-metadata-token", m.token)
	}

	return processRequest(c, req)
}

func processRequest(c *http.Client, r *http.Request) (string, error) {
	resp, err := c.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
