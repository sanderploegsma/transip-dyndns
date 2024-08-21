# TransIP DynDNS

A command-line utility to automatically update DNS entries in your TransIP domain with your current public IP address.
Supports both IPv4 and IPv6 addresses, and manages `A` or `AAAA` DNS records on your domain respectively.

Tip: run this utility on a schedule using `cron`, and you won't ever need to worry about keeping your DNS settings up-to-date!

## Usage

```console
$ transip-dyndns -h
NAME:
   transip-dyndns - Automatically update DNS entries in your TransIP domain with your current public IP address.

USAGE:
   transip-dyndns [global options]

VERSION:
   dev

GLOBAL OPTIONS:
   --account value, -a value                            TransIP account name [$TRANSIP_ACCOUNT_NAME]
   --private-key value, -k value                        path to TransIP API private key file [$TRANSIP_PRIVATE_KEY]
   --domain value, -d value                             domain name for which DNS entries should be synchronized
   --entry value, -e value [ --entry value, -e value ]  one or more DNS entries to synchronize
   --type value, -t value [ --type value, -t value ]    one or more DNS entry types to synchronize (options: A, AAAA)
   --ttl value                                          Time To Live (TTL) for newly created DNS entries, in seconds (default: 1h0m0s)
   --help, -h                                           show help
   --version, -v                                        print the version
```

### Usage example

Start by creating a new API key pair in your TransIP account here: https://www.transip.nl/cp/account/api/.
Make sure to store the private key somewhere safe.

Then, to automatically create or update a DNS `A` record on example.com called `foo`:

```sh
transip-dyndns --account transip-account-name --private-key /path/to/private.key --domain example.com --entry foo --type A
```

## Acknowledgements

This utility is heavily inspired by [Jerrythafast/transip-dyndns](https://github.com/Jerrythafast/transip-dyndns).
