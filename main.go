package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/domain"
	"github.com/urfave/cli/v2"
)

var version = "dev"

const (
	AccountNameEnvVar       = "TRANSIP_ACCOUNT_NAME"
	AccountNameFlagFull     = "account"
	AccountNameFlagShort    = "a"
	DnsEntryFlagFull        = "entry"
	DnsEntryFlagShort       = "e"
	DnsEntryTypeFlagFull    = "type"
	DnsEntryTypeFlagShort   = "t"
	DomainNameFlagFull      = "domain"
	DomainNameFlagShort     = "d"
	DnsEntryExpiryFlagFull  = "expiry"
	PrivateKeyPathEnvVar    = "TRANSIP_PRIVATE_KEY"
	PrivateKeyPathFlagFull  = "private-key"
	PrivateKeyPathFlagShort = "k"

	DnsEntryTypeA    = "A"
	DnsEntryTypeAAAA = "AAAA"

	DefaultDnsEntryExpiry = 3600
)

var (
	supportedDnsEntryTypes      = []string{DnsEntryTypeA, DnsEntryTypeAAAA}
	getIPAddressforDnsEntryType = map[string]GetIPAddress{
		DnsEntryTypeA:    GetIPv4,
		DnsEntryTypeAAAA: GetIPv6,
	}
)

func main() {
	var accountName string
	var privateKeyPath string
	var domainName string
	var dnsEntries cli.StringSlice
	var dnsEntryTypes cli.StringSlice
	var dnsEntryExpiry int

	app := &cli.App{
		Usage:           "Automatically update DNS entries in your TransIP domain with your current public IP address.",
		Version:         version,
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        AccountNameFlagFull,
				Aliases:     []string{AccountNameFlagShort},
				Usage:       "TransIP account name",
				Required:    true,
				EnvVars:     []string{AccountNameEnvVar},
				Destination: &accountName,
			},
			&cli.StringFlag{
				Name:        PrivateKeyPathFlagFull,
				Aliases:     []string{PrivateKeyPathFlagShort},
				Usage:       "path to TransIP API private key file",
				Required:    true,
				TakesFile:   true,
				EnvVars:     []string{PrivateKeyPathEnvVar},
				Destination: &privateKeyPath,
			},
			&cli.StringFlag{
				Name:        DomainNameFlagFull,
				Aliases:     []string{DomainNameFlagShort},
				Usage:       "domain name for which DNS entries should be synchronized",
				Required:    true,
				Destination: &domainName,
			},
			&cli.StringSliceFlag{
				Name:        DnsEntryFlagFull,
				Aliases:     []string{DnsEntryFlagShort},
				Usage:       "one or more DNS entries to synchronize",
				Required:    true,
				Destination: &dnsEntries,
			},
			&cli.StringSliceFlag{
				Name:     DnsEntryTypeFlagFull,
				Aliases:  []string{DnsEntryTypeFlagShort},
				Usage:    fmt.Sprintf("one or more DNS entry types to synchronize (options: %s)", strings.Join(supportedDnsEntryTypes, ", ")),
				Required: true,
				Action: func(ctx *cli.Context, s []string) error {
					for _, v := range s {
						if !slices.Contains(supportedDnsEntryTypes, v) {
							return fmt.Errorf("unsupported DNS entry type: %s", v)
						}
					}

					return nil
				},
				Destination: &dnsEntryTypes,
			},
			&cli.IntFlag{
				Name:  DnsEntryExpiryFlagFull,
				Usage: "expiry for newly created DNS entries",
				Value: DefaultDnsEntryExpiry,
				Action: func(ctx *cli.Context, i int) error {
					if i <= 0 {
						return fmt.Errorf("invalid DNS entry expiry value: %d", i)
					}

					return nil
				},
				Destination: &dnsEntryExpiry,
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := gotransip.NewClient(gotransip.ClientConfiguration{
				AccountName:    accountName,
				PrivateKeyPath: privateKeyPath,
			})

			if err != nil {
				return err
			}

			updater := &Updater{
				DomainRepository: &domain.Repository{Client: client},
			}

			errs := make([]error, 0)

			for dnsEntryType, getIPAddress := range getIPAddressforDnsEntryType {
				if slices.Contains(dnsEntryTypes.Value(), dnsEntryType) {
					errs = append(errs, updater.UpdateDNSEntries(domainName, dnsEntries.Value(), dnsEntryExpiry, dnsEntryType, getIPAddress))
				}
			}

			return errors.Join(errs...)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
