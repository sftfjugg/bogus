package bogus

import (
	"fmt"
	"net"
	"strings"
	"log"

	"github.com/miekg/dns"
)

// ResponseReverter reverses the operations done on the question section of a packet.
// This is need because the client will otherwise disregards the response, i.e.
// dig will complain with ';; Question section mismatch: got example.org/HINFO/IN'
type ResponseReverter struct {
	dns.ResponseWriter
	originalQuestion dns.Question
	bogus            []net.IP
}

// NewResponseReverter returns a pointer to a new ResponseReverter.
func NewResponseReverter(w dns.ResponseWriter, r *dns.Msg, bogus []net.IP) *ResponseReverter {
	return &ResponseReverter{
		ResponseWriter:   w,
		originalQuestion: r.Question[0],
		bogus:            bogus,
	}
}

// WriteMsg records the status code and calls the underlying ResponseWriter's WriteMsg method.
func (r *ResponseReverter) WriteMsg(res *dns.Msg) error {
	res.Question[0] = r.originalQuestion
	for _, rr := range res.Answer {
		if rr.Header().Rrtype != dns.TypeA && rr.Header().Rrtype != dns.TypeAAAA {
			continue
		}

		ss := strings.Split(rr.String(), "\t")
		if len(ss) != 5 {
			continue
		}
		ip := net.ParseIP(ss[4])
		for _, i := range r.bogus {
			if !ip.Equal(i) {
				continue
			}
			rs := &dns.Msg{
				MsgHdr: dns.MsgHdr{
					Id:                 res.Id,
					Response:           true,
					RecursionDesired:   true,
					RecursionAvailable: true,
				},
				Question: res.Question,
				Answer:   []dns.RR{soa(r.originalQuestion.Name)},
			}
			return r.ResponseWriter.WriteMsg(rs)
		}
	}
	return r.ResponseWriter.WriteMsg(res)
}

// Write is a wrapper that records the size of the message that gets written.
func (r *ResponseReverter) Write(buf []byte) (int, error) {
	log.Println("bogus write", string(buf))
	n, err := r.ResponseWriter.Write(buf)
	return n, err
}

func soa(name string) dns.RR {
	s := fmt.Sprintf("%s 60 IN SOA ns1.%s postmaster.%s 1524370381 14400 3600 604800 60", name, name, name)
	soa, _ := dns.NewRR(s)
	return soa
}
