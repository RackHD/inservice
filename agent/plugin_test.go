package plugins_test

import (
	"fmt"

	. "github.com/RackHD/InService/agent"
	"github.com/RackHD/InService/agent/grpc/plugin"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ServiceMock struct {
	start  error
	stop   error
	status error
}

func (s *ServiceMock) Start(_ context.Context) error {
	return s.start
}

func (s *ServiceMock) Stop(_ context.Context) error {
	return s.stop
}

func (s *ServiceMock) Status(_ context.Context) error {
	return s.status
}

var _ = Describe("Plugin", func() {
	Describe("NewPlugin", func() {
		It("should require a struct with a Service interface", func() {
			_, err := NewPlugin(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should return a Plugin struct", func() {
			plugin, err := NewPlugin(&ServiceMock{})
			Expect(err).To(Succeed())
			Expect(plugin).To(BeAssignableToTypeOf(&Plugin{}))
		})
	})

	Describe("Start", func() {
		It("should return a plugin.StartResponse on success", func() {
			p, err := NewPlugin(&ServiceMock{
				start:  nil,
				stop:   fmt.Errorf("Stop Called."),
				status: fmt.Errorf("Status Called."),
			})
			Expect(err).To(Succeed())

			response, err := p.Start(nil, &plugin.StartRequest{})
			Expect(err).To(Succeed())
			Expect(response).To(BeAssignableToTypeOf(&plugin.StartResponse{}))
		})

		It("should return an error on failure", func() {
			p, err := NewPlugin(&ServiceMock{
				start:  fmt.Errorf("Start Called."),
				stop:   nil,
				status: nil,
			})
			Expect(err).To(Succeed())

			_, err = p.Start(nil, &plugin.StartRequest{})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Stop", func() {
		It("should return a plugin.StartResponse on success", func() {
			p, err := NewPlugin(&ServiceMock{
				start:  fmt.Errorf("Start Called."),
				stop:   nil,
				status: fmt.Errorf("Status Called."),
			})
			Expect(err).To(Succeed())

			response, err := p.Stop(nil, &plugin.StopRequest{})
			Expect(err).To(Succeed())
			Expect(response).To(BeAssignableToTypeOf(&plugin.StopResponse{}))
		})

		It("should return an error on failure", func() {
			p, err := NewPlugin(&ServiceMock{
				start:  nil,
				stop:   fmt.Errorf("Stop Called."),
				status: nil,
			})
			Expect(err).To(Succeed())

			_, err = p.Stop(nil, &plugin.StopRequest{})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Status", func() {
		It("should return a plugin.StatusResponse on success", func() {
			p, err := NewPlugin(&ServiceMock{
				start:  fmt.Errorf("Start Called."),
				stop:   fmt.Errorf("Stop Called."),
				status: nil,
			})
			Expect(err).To(Succeed())

			response, err := p.Status(nil, &plugin.StatusRequest{})
			Expect(err).To(Succeed())
			Expect(response).To(BeAssignableToTypeOf(&plugin.StatusResponse{}))
		})

		It("should return an error on failure", func() {
			p, err := NewPlugin(&ServiceMock{
				start:  nil,
				stop:   nil,
				status: fmt.Errorf("Status Called."),
			})
			Expect(err).To(Succeed())

			_, err = p.Status(nil, &plugin.StatusRequest{})
			Expect(err).To(HaveOccurred())
		})
	})
})
