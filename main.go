package main

import (
	"crypto/tls"
	"fmt"
	"github.com/miekg/dns"
	"io/ioutil"
	"net/http"
)

func resolveOverTCP(host string) (string, error) {
	c := new(dns.Client)
	c.Net = "tcp" // specify that we want to use TCP

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), dns.TypeA)

	// Using a public DNS server, you can change it as per your requirement
	r, _, err := c.Exchange(m, "1.1.1.1:53")
	if err != nil {
		return "", err
	}

	if len(r.Answer) == 0 {
		return "", fmt.Errorf("no answer for query")
	}

	for _, a := range r.Answer {
		if aRec, ok := a.(*dns.A); ok {
			return aRec.A.String(), nil
		}
	}

	return "", fmt.Errorf("no A record found")
}

func main() {
	// Resolve the domain name over TCP
	_, err := resolveOverTCP("catfact.ninja")
	if err != nil {
		fmt.Printf("DNS resolution error: %s\n", err)
		return
	}

	// Make the request using the hostname
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{ServerName: "catfact.ninja"},
		},
	}

	resp, err := client.Get("https://catfact.ninja/fact")
	if err != nil {
		fmt.Printf("Error making request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %s\n", err)
		return
	}

	fmt.Printf("Response: %s\n", body)
}
