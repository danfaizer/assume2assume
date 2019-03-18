# assume2assume
Assume an AWS IAM Role in target AWS account from specified AWS IAM Role and source AWS account.

## Installation
Make sure you have a working Go environment.

To install `assume2assume` cli, simply run:
```
$ go get github.com/danfaizer/assume2assume/cmd/assume2assume
```

## Usage
```
usage: assume2assume [flags]
  -r string
    	Source AWS IAM Role
  -s string
    	Soruce AWS Account ID
  -t string
    	Destination AWS IAM Role
  -d string
    	Destination AWS Account ID
  -p	Print destination AWS STS Credentials
  -q	Print only destination AWS STS Credentials if specified
```
