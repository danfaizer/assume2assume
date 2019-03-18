package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

var (
	fromAccount        = flag.String("s", "", "Soruce AWS Account ID")
	fromRole           = flag.String("r", "", "Source AWS IAM Role")
	destinationAccount = flag.String("d", "", "Destination AWS Account ID")
	destinationRole    = flag.String("t", "", "Destination AWS IAM Role")
	printSTSToken      = flag.Bool("p", false, "Print destination AWS STS Credentials")
	quietOutput        = flag.Bool("q", false, "Print only destination AWS STS Credentials if specified")
)

func mustPrintCreds(creds *sts.Credentials) {
	fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", *creds.AccessKeyId)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", *creds.SecretAccessKey)
	fmt.Printf("export AWS_SESSION_TOKEN=%s\n", *creds.SessionToken)
	fmt.Printf("export AWS_SECURITY_TOKEN=%s\n", *creds.SessionToken)
}

func assumeRole(role, sessionName string) (*sts.Credentials, error) {
	sess := session.Must(session.NewSession())

	svc := sts.New(sess)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(role),
		RoleSessionName: aws.String(sessionName),
	}

	resp, err := svc.AssumeRole(params)
	if err != nil {
		return nil, err
	}

	return resp.Credentials, nil
}

func assumeRoleWithSTSCreds(role, sessionName string, creds *sts.Credentials) (*sts.Credentials, error) {
	sess := session.Must(session.NewSession())
	conf := aws.Config{}
	conf.Credentials = credentials.NewStaticCredentials(*creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken)

	svc := sts.New(sess, &conf)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(role),
		RoleSessionName: aws.String(sessionName),
	}

	resp, err := svc.AssumeRole(params)
	if err != nil {
		return nil, err
	}

	return resp.Credentials, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: assume2assume [flags]")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *destinationAccount == "" || *destinationRole == "" || *fromAccount == "" || *fromRole == "" {
		usage()
		log.Fatal("some arguments are missing")
	}

	if !*quietOutput {
		fmt.Printf("\nassume2assume will try to assume from:\n\nAccount: %s\nRole: %s\n\nto access:\n\nAccount: %s\nRole: %s\n\n", *fromAccount, *fromRole, *destinationAccount, *destinationRole)
	}

	fromRoleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", *fromAccount, *fromRole)
	if !*quietOutput {
		fmt.Printf("Assuming fromRole: %s\n", fromRoleARN)
	}

	fromCreds, err := assumeRole(fromRoleARN, "fromSession")
	if err != nil {
		if !*quietOutput {
			fmt.Printf("error assuming fromRole: %s\n%s\n\n", fromRoleARN, err)
		}
		os.Exit(1)
	}

	destinationRoleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", *destinationAccount, *destinationRole)
	if !*quietOutput {
		fmt.Printf("Assuming destinationRole: %s\n", destinationRoleARN)
	}

	destCreds, err := assumeRoleWithSTSCreds(destinationRoleARN, "destinationSession", fromCreds)
	if err != nil {
		if !*quietOutput {
			fmt.Printf("error assuming destinationRole: %s\n%s\n\n", destinationRoleARN, err)
		}
		os.Exit(1)
	}

	if !*quietOutput {
		fmt.Printf("\nSource Role: %s in AWS Account: %s can successfully assume role %s in AWS Account %s\n\n", *fromRole, *fromAccount, *destinationRole, *destinationAccount)
	}

	if *printSTSToken == true {
		mustPrintCreds(destCreds)
	}
}
