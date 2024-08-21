package main

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/transip/gotransip/v6/domain"
)

func TestDNSEntryIsUpToDate(t *testing.T) {
	// TODO
}

func TestDNSEntryUpdated(t *testing.T) {
	domainName := "example.com"
	domainRepository := &FakeDomainRepository{
		DnsEntriesByDomain: map[string][]domain.DNSEntry{
			domainName: {
				domain.DNSEntry{
					Name:    "www",
					Type:    "A",
					Expire:  3600,
					Content: "1.1.1.1",
				},
			},
		},
	}
	updater := &Updater{DomainRepository: domainRepository}

	getIPAddress := func() (string, error) {
		return "2.2.2.2", nil
	}

	err := updater.UpdateDNSEntries(domainName, []string{"www"}, 3600, "A", getIPAddress)

	if assert.Nil(t, err) {
		assert.ElementsMatch(t, domainRepository.DnsEntriesByDomain[domainName], []domain.DNSEntry{
			{
				Name:    "www",
				Type:    "A",
				Expire:  3600,
				Content: "2.2.2.2",
			},
		})
	}
}

func TestDNSEntryCreated(t *testing.T) {
	domainName := "example.com"
	domainRepository := &FakeDomainRepository{
		DnsEntriesByDomain: map[string][]domain.DNSEntry{
			domainName: {
				domain.DNSEntry{
					Name:    "www",
					Type:    "A",
					Expire:  3600,
					Content: "1.1.1.1",
				},
			},
		},
	}
	updater := &Updater{DomainRepository: domainRepository}

	getIPAddress := func() (string, error) {
		return "2.2.2.2", nil
	}

	err := updater.UpdateDNSEntries(domainName, []string{"www2"}, 3600, "A", getIPAddress)

	if assert.Nil(t, err) {
		assert.ElementsMatch(t, domainRepository.DnsEntriesByDomain[domainName], []domain.DNSEntry{
			{
				Name:    "www",
				Type:    "A",
				Expire:  3600,
				Content: "1.1.1.1",
			},
			{
				Name:    "www2",
				Type:    "A",
				Expire:  3600,
				Content: "2.2.2.2",
			},
		})
	}
}

func TestCannotDetermineIPAddress(t *testing.T) {
	domainName := "example.com"
	domainRepository := &FakeDomainRepository{
		DnsEntriesByDomain: map[string][]domain.DNSEntry{
			domainName: make([]domain.DNSEntry, 0),
		},
	}
	updater := &Updater{DomainRepository: domainRepository}

	getIPAddress := func() (string, error) {
		return "", errors.New("service unavailable")
	}

	err := updater.UpdateDNSEntries(domainName, []string{"www"}, 3600, "A", getIPAddress)
	assert.NotNil(t, err)
	assert.Empty(t, domainRepository.DnsEntriesByDomain[domainName])
}

type FakeDomainRepository struct {
	DnsEntriesByDomain map[string][]domain.DNSEntry
}

func (r *FakeDomainRepository) GetDNSEntries(domainName string) ([]domain.DNSEntry, error) {
	if entries, ok := r.DnsEntriesByDomain[domainName]; ok {
		return entries, nil
	}

	return nil, fmt.Errorf("unknown domain %s", domainName)
}

func (r *FakeDomainRepository) AddDNSEntry(domainName string, dnsEntry domain.DNSEntry) error {
	entries, ok := r.DnsEntriesByDomain[domainName]
	if !ok {
		entries = make([]domain.DNSEntry, 0)
	}

	r.DnsEntriesByDomain[domainName] = append(entries, dnsEntry)
	return nil
}

func (r *FakeDomainRepository) UpdateDNSEntry(domainName string, dnsEntry domain.DNSEntry) error {
	entries, ok := r.DnsEntriesByDomain[domainName]
	if !ok {
		return fmt.Errorf("unknown domain name %s", domainName)
	}

	for i, item := range entries {
		if item.Name == dnsEntry.Name && item.Type == dnsEntry.Type {
			entries[i] = dnsEntry
			return nil
		}
	}

	return fmt.Errorf("no existing DNS entry with name %s", dnsEntry.Name)
}
