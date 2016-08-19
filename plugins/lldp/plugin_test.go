package main_test

import (
	. "github.com/RackHD/InService/plugins/lldp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Describe("NewLLDPPlugin", func() {
		It("should check IP address input variables", func() {
			ifaces := []string{"eth0", "eth1"}
			_, err := NewLLDPPlugin("0", 8080, ifaces)
			Expect(err).To(HaveOccurred())
		})

		It("should check port input variable for out of range", func() {
			ifaces := []string{"eth0", "eth1"}
			_, err := NewLLDPPlugin("10.10.10.10", 0, ifaces)
			Expect(err).To(HaveOccurred())
		})

		It("should check port input variable for out of range", func() {
			ifaces := []string{"eth0", "eth1"}
			_, err := NewLLDPPlugin("10.10.10.10", 100000, ifaces)
			Expect(err).To(HaveOccurred())
		})

		It("should return an LLDP Plugin struct", func() {
			ifaces := []string{"eth0", "eth1"}
			lldp, err := NewLLDPPlugin("10.10.10.10", 8080, ifaces)
			Expect(err).To(Succeed())
			Expect(lldp).To(BeAssignableToTypeOf(&LLDPPlugin{}))
		})
	})
})
