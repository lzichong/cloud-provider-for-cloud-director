package e2e

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware/cloud-provider-for-cloud-director/pkg/testingsdk"
	"github.com/vmware/cloud-provider-for-cloud-director/pkg/vcdsdk"
	"github.com/vmware/cloud-provider-for-cloud-director/tests/e2e/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"net/http"
)

const (
	testBaseName       = "e2e-ccm-automation-test"
	testServiceName    = "service-lb-test"
	testDeploymentName = "deployment-lb-test"
	httpPort           = 80

	ccmConfigMapName = "vcloud-ccm-configmap"
)

var testHttpName = "http"

var _ = Describe("Ensure Loadbalancer", func() {
	var (
		ns          *v1.Namespace
		svc         *v1.Service
		err         error
		ipamSubnet  string
		networkName string
		tc          *testingsdk.TestClient
	)

	tc, err = utils.NewTestClient(host, org, ovdcName, username, token, clusterId, false)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(tc).NotTo(BeNil())
	Expect(&tc.Cs).NotTo(BeNil())

	labels := map[string]string{
		"app": testServiceName,
	}

	httpServicePort := []v1.ServicePort{{
		Name:        testHttpName,
		Port:        httpPort,
		TargetPort:  intstr.FromInt(httpPort),
		Protocol:    v1.ProtocolTCP,
		AppProtocol: &testHttpName, // We want our virtual service to be HTTP
	}}

	ctx := context.TODO()
	ns, err = tc.CreateNameSpace(ctx, testBaseName)
	Expect(err).NotTo(HaveOccurred())

	// We will have a sample deployment so the server will return some sort of data back to us using an official e2e test image
	_, err = utils.CreateDeployment(ctx, tc, testDeploymentName, ns.Name, labels)
	Expect(err).NotTo(HaveOccurred())

	// GetConfigMap to retrieve ipamSubnet, network name for gateway manager in order to check if VCD resources are present
	ccmConfigMap, err := utils.GetConfigMap(tc.Cs, "kube-system", ccmConfigMapName)
	Expect(err).NotTo(HaveOccurred())
	Expect(ccmConfigMap).NotTo(BeNil())
	ipamSubnet, err = utils.GetIpamSubnetFromConfigMap(ccmConfigMap)
	Expect(err).NotTo(HaveOccurred())
	Expect(ipamSubnet).NotTo(BeEmpty())
	networkName, err = utils.GetNetworkNameFromConfigMap(ccmConfigMap)
	Expect(err).NotTo(HaveOccurred())
	Expect(networkName).NotTo(BeEmpty())

	// Gateway Manager is used to validate VCD networking resources
	gatewayMgr, err := vcdsdk.NewGatewayManager(context.TODO(), tc.VcdClient, networkName, ipamSubnet)
	Expect(err).NotTo(HaveOccurred())
	Expect(gatewayMgr).NotTo(BeNil())

	// Case 1. We should be able to create a LB http service on port 80 with no errors
	It("should create a load balancer service", func() {
		// Similar to Ingress setup, we will use: name=http, port=80, protocol=tcp, appProtocol=http
		By("creating a http load balancer service")
		svc, err = tc.CreateLoadBalancerService(ctx, ns.Name, testServiceName, nil, labels, httpServicePort)
		Expect(err).NotTo(HaveOccurred())
		Expect(svc).NotTo(BeNil())
	})

	// Case 2. We should have an external IP and VCD resources after creating a Loadbalancer Service
	It("should have an external IP and VCD resources", func() {
		By("getting the next available unused external IP from our edge gateway")
		availableIp, err := gatewayMgr.GetUnusedExternalIPAddress(ctx, ipamSubnet)
		Expect(err).NotTo(HaveOccurred())
		Expect(availableIp).NotTo(BeEmpty())

		By("fetching the external IP from the service")
		externalIp, err := utils.WaitForExtIP(tc.Cs, ns.Name, testServiceName)
		Expect(err).NotTo(HaveOccurred())
		Expect(externalIp).Should(Equal(availableIp))

		By("checking if resources is present in VCD Resource Set")
		virtualServiceNamePrefix := utils.GetVirtualServicePrefix(svc, tc.ClusterId)
		lbPoolNamePrefix := utils.GetLBPoolNamePrefix(svc, tc.ClusterId)
		portDetailsList := utils.GetPortDetailsList(svc)
		Expect(portDetailsList).NotTo(BeNil())

		oneArm := vcdsdk.OneArm{
			StartIP: "192.168.8.2",
			EndIP:   "192.168.8.100",
		}
		resourcesFound, err := utils.HasVCDResourcesForApplicationLB(ctx, tc, gatewayMgr, &oneArm, virtualServiceNamePrefix, lbPoolNamePrefix, portDetailsList)
		Expect(err).NotTo(HaveOccurred())
		Expect(resourcesFound).Should(BeTrue())

		By("checking virtual IP stored in in CPI vcdResourceSet matches the external IP from the load balancer service")
		externalIpExists, err := utils.IsExternalIpInVCDResourceSet(ctx, tc, externalIp, tc.ClusterId)
		Expect(err).NotTo(HaveOccurred())
		Expect(externalIpExists).To(BeTrue())
	})

	// Case 3. Check for valid external IP and connectivity for ip:port via HTTP Get Request.
	It("should have connectivity from external ip", func() {
		By("fetching the external ip of the service")
		externalIp, err := utils.WaitForExtIP(tc.Cs, ns.Name, testServiceName)
		Expect(err).NotTo(HaveOccurred())
		Expect(externalIp).NotTo(BeEmpty())

		By("checking connectivity of external ip")
		resp, err := http.Get(fmt.Sprintf("http://%s:%d", externalIp, httpPort))
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).Should(BeNumerically("==", 200))
	})

	// Case 4. We should be able to perform clean up after validation by deleting the service
	It("should ensure load balancer deleted", func() {
		By("deleting the service")
		err = tc.DeleteService(ctx, ns.Name, testServiceName)
		Expect(err).NotTo(HaveOccurred())
	})

	// Case 5. We should not have any VCD resources after the deletion of loadbalancer service
	It("should not have VCD resources after service has been deleted", func() {
		// After service is deleted, check empty resources
		By("checking VCD resources")
		virtualServiceNamePrefix := utils.GetVirtualServicePrefix(svc, clusterId)
		lbPoolNamePrefix := utils.GetLBPoolNamePrefix(svc, clusterId)
		portDetailsList := utils.GetPortDetailsList(svc)
		Expect(portDetailsList).NotTo(BeNil())

		// These will be our defaults, so it's ok to hardcode this. We can also retrieve this from ccm configmap
		oneArm := vcdsdk.OneArm{
			StartIP: "192.168.8.2",
			EndIP:   "192.168.8.100",
		}

		resourcesFound, err := utils.HasVCDResourcesForApplicationLB(context.TODO(), tc, gatewayMgr, &oneArm, virtualServiceNamePrefix, lbPoolNamePrefix, portDetailsList)
		Expect(err).NotTo(HaveOccurred())
		Expect(resourcesFound).Should(BeFalse())

		By("cleaning up the remainder of deployment and namespaces")
		err = tc.DeleteDeployment(ctx, ns.Name, testDeploymentName)
		Expect(err).NotTo(HaveOccurred())

		err = tc.DeleteNameSpace(ctx, ns.Name)
		Expect(err).NotTo(HaveOccurred())
	})
})
