package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	neutronports "github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

var openstackCloudName = "openstack"
var clientConfigFile = "clouds.yaml"

func main() {
	var err error

	// Read the clouds.yaml file.
	// That information is stored in secret cloud-credentials.
	content, err := ioutil.ReadFile(clientConfigFile)
	if err != nil {
		klog.Fatalf("Could read file %s, err: %q", clientConfigFile, err)
	}

	// Unmarshal YAML content into Clouds object.
	var clouds clientconfig.Clouds
	err = yaml.Unmarshal(content, &clouds)
	if err != nil {
		klog.Fatalf("Could not parse cloud configuration from %s, err: %q", clientConfigFile, err)
	}
	// We expect that the cloud in clouds.yaml be named "openstack".
	cloud, ok := clouds.Clouds[openstackCloudName]
	if !ok {
		klog.Fatalf("Invalid clouds.yaml file. Missing section for cloud name '%s'", openstackCloudName)
	}

	// Prepare the options.
	clientOpts := &clientconfig.ClientOpts{
		Cloud:      cloud.Cloud,
		AuthType:   cloud.AuthType,
		AuthInfo:   cloud.AuthInfo,
		RegionName: cloud.RegionName,
	}
	opts, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		klog.Fatal(err)
	}
	provider, err := openstack.NewClient(opts.IdentityEndpoint)
	if err != nil {
		klog.Fatal(err)
	}

	// Now, authenticate.
	err = openstack.Authenticate(provider, *opts)
	if err != nil {
		klog.Fatal(err)
	}

	// And another client for neutron (network).
	neutronClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		//  Region: cloud.RegionName,
	})
	if err != nil {
		klog.Fatal(err)
	}

	getPort := func() neutronports.Port {
		portListOpts := neutronports.ListOpts{
			ID: "207d5323-4858-4d8f-8122-46d32b9e2a54",
		}

		var serverPorts []neutronports.Port
		pager := neutronports.List(neutronClient, portListOpts)
		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			portList, err := neutronports.ExtractPorts(page)
			if err != nil {
				return false, err
			}
			for _, p := range portList {
				serverPorts = append(serverPorts, p)
			}
			return true, nil
		})
		if err != nil {
			klog.Fatal(err)
		}
		//	for k, p := range serverPorts[:1] {
		k := 0
		p := serverPorts[0]
		b, err := json.Marshal(p)
		if err != nil {
			klog.Fatal(err)
			//		continue
		}
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, b, "", "    "); err != nil {
			klog.Fatal(err)
			//		continue
		}
		klog.Infof("%d => %s", k, prettyJSON)
		klog.Infof("--------------")
		//}

		return p
	}

	klog.Info("Getting port")
	p := getPort()
	klog.Info("Resetting port")
	allowedPairs := []neutronports.AddressPair{}
	updateOpts := neutronports.UpdateOpts{
		AllowedAddressPairs: &allowedPairs,
		RevisionNumber:      &p.RevisionNumber,
	}
	_, err = ports.Update(neutronClient, p.ID, updateOpts).Extract()
	if err != nil {
		klog.Fatal(err)

	}

	klog.Info("Getting port")
	p = getPort()
	klog.Info("Updating port")
	allowedPairs = []neutronports.AddressPair{
		{
			IPAddress: "192.168.123.10",
		},
	}
	updateOpts = neutronports.UpdateOpts{
		AllowedAddressPairs: &allowedPairs,
		RevisionNumber:      &p.RevisionNumber,
	}
	_, err = ports.Update(neutronClient, p.ID, updateOpts).Extract()
	if err != nil {
		klog.Fatal(err)
	}

	klog.Info("Getting port")
	p = getPort()
	klog.Info("Resetting port")
	allowedPairs = []neutronports.AddressPair{}
	updateOpts = neutronports.UpdateOpts{
		AllowedAddressPairs: &allowedPairs,
		RevisionNumber:      &p.RevisionNumber,
	}
	_, err = ports.Update(neutronClient, p.ID, updateOpts).Extract()
	if err != nil {
		klog.Fatal(err)
	}

	klog.Info("Getting port")
	p = getPort()
}
