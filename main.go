package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type ZoneRecords map[string][]string

var records = ZoneRecords{}

func addToRecords(name string, ip []string) {
	log.Printf("Register record : %s -> %s\n", name, ip)
	records[name] = ip
}

func removeFromRecords(name string) {
	log.Printf("Unregister record : %s\n", name)
	delete(records, name)
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			for _, ip := range records[q.Name] {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func StartDNSServer(host string, port int) (*dns.Server, error) {
	dns.HandleFunc(".", handleDNSRequest)
	server := &dns.Server{Addr: host + ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	if err != nil {
		return nil, errors.Wrap(err, "Error starting dns server")
	}
	return server, nil
}

func main() {

	args := os.Args[1:]
	for i := 0; i < len(args)/2; i++ {
		hostname := args[i*2]
		address := args[i*2+1]
		addToRecords(hostname, []string{address})
	}

	// start server
	port := 53

	dnsServer, err := StartDNSServer("", port)
	defer dnsServer.Shutdown()
	if err != nil {
		log.Fatalf("%e\n ", err)
		panic(err)
	}
}
