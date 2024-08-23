package main

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/transip/gotransip/v6/domain"
)

type DomainRepository interface {
	GetDNSEntries(domainName string) ([]domain.DNSEntry, error)
	AddDNSEntry(domainName string, dnsEntry domain.DNSEntry) error
	UpdateDNSEntry(domainName string, dnsEntry domain.DNSEntry) error
}

type GetIPAddress func() (string, error)

type Updater struct {
	DomainRepository DomainRepository
	StdOut           io.Writer
}

func (u *Updater) UpdateDNSEntries(domainName string, dnsEntryNames []string, ttl time.Duration, entryType string, getIPAddress GetIPAddress) error {
	fmt.Fprintf(u.StdOut, "Updating DNS %s records for %s\n", entryType, domainName)

	ipAddress, err := getIPAddress()
	if err != nil {
		return err
	}

	fmt.Fprintf(u.StdOut, "Current IP address: %s\n", ipAddress)

	dnsEntries, err := u.DomainRepository.GetDNSEntries(domainName)
	if err != nil {
		return err
	}

	success := true

	for _, entryName := range dnsEntryNames {
		exists := false

		for _, entry := range dnsEntries {
			if entry.Name != entryName || entry.Type != entryType {
				continue
			}

			exists = true

			if entry.Content == ipAddress {
				fmt.Fprintf(u.StdOut, "DNS %s record '%s' is up-to-date\n", entryType, entryName)
				continue
			}

			fmt.Fprintf(u.StdOut, "Updating DNS %s record '%s'\n", entryType, entryName)
			entry.Content = ipAddress
			if err = u.DomainRepository.UpdateDNSEntry(domainName, entry); err != nil {
				fmt.Fprintf(u.StdOut, "Failed to update DNS record: %v\n", err)
			}
		}

		if exists {
			continue
		}

		fmt.Fprintf(u.StdOut, "Creating DNS %s record '%s' with TTL %s\n", entryType, entryName, ttl)
		newEntry := domain.DNSEntry{
			Name:    entryName,
			Content: ipAddress,
			Type:    entryType,
			Expire:  int(ttl.Seconds()),
		}
		if err = u.DomainRepository.AddDNSEntry(domainName, newEntry); err != nil {
			fmt.Fprintf(u.StdOut, "Failed to create DNS %s record '%s': %v\n", entryType, entryName, err)
		}
	}

	if !success {
		return errors.New("failed to update one or more DNS records")
	}

	return nil
}
