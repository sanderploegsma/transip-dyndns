package main

import (
	"errors"
	"fmt"
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
}

func (u *Updater) UpdateDNSEntries(domainName string, dnsEntryNames []string, ttl time.Duration, entryType string, getIPAddress GetIPAddress) error {
	ipAddress, err := getIPAddress()
	if err != nil {
		return err
	}

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
				fmt.Printf("DNS entry '%s' (%s) is up-to-date\n", entryName, entryType)
				continue
			}

			fmt.Printf("Updating content of DNS entry '%s' (%s) to '%s'\n", entryName, entryType, ipAddress)
			entry.Content = ipAddress
			if err = u.DomainRepository.UpdateDNSEntry(domainName, entry); err != nil {
				fmt.Printf("Failed to update DNS entry: %v\n", err)
			}
		}

		if exists {
			continue
		}

		fmt.Printf("Creating new DNS entry '%s' (%s) with content '%s' and TTL %s\n", entryName, entryType, ipAddress, ttl)
		newEntry := domain.DNSEntry{
			Name:    entryName,
			Content: ipAddress,
			Type:    entryType,
			Expire:  int(ttl.Seconds()),
		}
		if err = u.DomainRepository.AddDNSEntry(domainName, newEntry); err != nil {
			fmt.Printf("Failed to create DNS entry: %v\n", err)
		}
	}

	if !success {
		return errors.New("failed to update one or more DNS entries")
	}

	return nil
}
