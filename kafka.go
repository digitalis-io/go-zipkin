package zipkin

import (
	"github.com/stealthly/siesta"
	"gopkg.in/spacemonkeygo/monitor.v1/trace/gen-go/zipkin"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/elodina/pyrgus/log"
)

type KafkaCollector struct {
	producer *siesta.KafkaProducer
	topic    string
}

func (c KafkaCollector) Collect(s *zipkin.Span) {
	t := thrift.NewTMemoryBuffer()
	p := thrift.NewTBinaryProtocolTransport(t)
	err := s.Write(p)
	if err != nil {
		log.Logger.Warn("Couldn't serialize span: ", err)
		return
	}

	// TODO: latter version of Siesta provides channel to hook up on which streams sending results, so need to use that.
	// But first, will need to update libraries which will use go-zipkin to using latter Siesta as well
	c.producer.Send(&siesta.ProducerRecord{Topic: c.topic, Value: t.Buffer.Bytes()})
}
