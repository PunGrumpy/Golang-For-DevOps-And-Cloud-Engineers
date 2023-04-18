package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/wardviaene/golang-for-devops-course/ssh-demo"
)

const location = "southeastasia"

func main() {
	var (
		token  azcore.TokenCredential
		pubKey string
		err    error
	)
	ctx := context.Background()
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		fmt.Println("AZURE_SUBSCRIPTION_ID is not set")
		os.Exit(1)
	}

	if pubKey, err = generateKeys(); err != nil {
		fmt.Printf("unable to generate keys: %s\n", err)
		os.Exit(1)
	}
	if token, err = getToken(); err != nil {
		fmt.Printf("unable to get token: %s\n", err)
		os.Exit(1)
	}
	if err = launchInstance(ctx, subscriptionID, token, &pubKey); err != nil {
		fmt.Printf("unable to launch instance: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully launch instance\n")
}

func generateKeys() (string, error) {
	var (
		privateKey []byte
		publicKey  []byte
		err        error
	)
	if privateKey, publicKey, err = ssh.GenerateKeys(); err != nil {
		fmt.Printf("unable to generate keys: %s\n", err)
		os.Exit(1)
	}
	if err = os.WriteFile("mykey.pem", privateKey, 0600); err != nil {
		fmt.Printf("unable to write private key: %s\n", err)
		os.Exit(1)
	}
	if err = os.WriteFile("mykey.pub", publicKey, 0644); err != nil {
		fmt.Printf("unable to write public key: %s\n", err)
		os.Exit(1)
	}

	return string(publicKey), nil
}

func getToken() (azcore.TokenCredential, error) {
	token, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return token, err
	}

	return token, nil
}

func launchInstance(ctx context.Context, subscriptionID string, cred azcore.TokenCredential, pubKey *string) error {
	// Create resource group
	resourceGroupClient, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}
	resourceGroupParams := armresources.ResourceGroup{
		Location: to.Ptr(location),
	}
	resourceGroupResponse, err := resourceGroupClient.CreateOrUpdate(
		ctx,
		"azure-go-sdk",
		resourceGroupParams,
		nil,
	)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully created resource group: %s\n", *resourceGroupResponse.Name)

	// Create virtual network
	virtualNetworkClient, err := armnetwork.NewVirtualNetworksClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}

	vnetResponse, found, err := findVnet(ctx, *resourceGroupResponse.Name, "azure-go-sdk", virtualNetworkClient)
	if err != nil {
		return err
	}

	if !found {
		virtualNetworkPollerResp, err := virtualNetworkClient.BeginCreateOrUpdate(
			ctx,
			*resourceGroupResponse.Name,
			"azure-go-sdk",
			armnetwork.VirtualNetwork{
				Location: to.Ptr(location),
				Properties: &armnetwork.VirtualNetworkPropertiesFormat{
					AddressSpace: &armnetwork.AddressSpace{
						AddressPrefixes: []*string{
							to.Ptr("10.1.0.0/16"),
						},
					},
				},
			},
			nil)
		if err != nil {
			return err
		}
		virtualNetworkRepsonse, err := virtualNetworkPollerResp.PollUntilDone(ctx, nil)
		if err != nil {
			return err
		}
		vnetResponse = virtualNetworkRepsonse.VirtualNetwork
		fmt.Printf("Successfully created virtual network: %s\n", *vnetResponse.Name)
	}
	fmt.Printf("Found virtual network: %s\n", *vnetResponse.Name)

	// Create subnet
	subnetsClient, err := armnetwork.NewSubnetsClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}
	subnetsPollerResp, err := subnetsClient.BeginCreateOrUpdate(
		ctx,
		*resourceGroupResponse.Name,
		*vnetResponse.Name,
		"azure-go-sdk",
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.Ptr("10.1.0.0/24"),
			},
		},
		nil,
	)
	if err != nil {
		return err
	}
	subnetsResponse, err := subnetsPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully created subnet: %s\n", *vnetResponse.Name)

	// Create public IP address
	publicIPAddressClient, err := armnetwork.NewPublicIPAddressesClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}
	publicIPPollerResp, err := publicIPAddressClient.BeginCreateOrUpdate(
		ctx,
		*resourceGroupResponse.Name,
		"azure-go-sdk",
		armnetwork.PublicIPAddress{
			Location: to.Ptr(location),
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic),
			},
		},
		nil,
	)
	if err != nil {
		return err
	}
	publicIPAddressResponse, err := publicIPPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully created public IP address: %s\n", *vnetResponse.Name)

	// Create network security group
	networkSecurityGroupClient, err := armnetwork.NewSecurityGroupsClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}
	networkSecurityGroupPollerResp, err := networkSecurityGroupClient.BeginCreateOrUpdate(
		ctx,
		*resourceGroupResponse.Name,
		"azure-go-sdk",
		armnetwork.SecurityGroup{
			Location: to.Ptr(location),
			Properties: &armnetwork.SecurityGroupPropertiesFormat{
				SecurityRules: []*armnetwork.SecurityRule{
					{
						Name: to.Ptr("allow-ssh"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							SourceAddressPrefix:      to.Ptr("0.0.0.0/0"),
							SourcePortRange:          to.Ptr("*"),
							DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),
							DestinationPortRange:     to.Ptr("22"),
							Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
							Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
							Description:              to.Ptr("Allow All SSH on Port 22"),
							Direction:                to.Ptr(armnetwork.SecurityRuleDirectionInbound),
							Priority:                 to.Ptr(int32(1001)),
						},
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		return err
	}
	networkSecurityGroupResponse, err := networkSecurityGroupPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully created network security group: %s\n", *vnetResponse.Name)

	// Create network interface (NIC)
	interfaceClient, err := armnetwork.NewInterfacesClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}
	interfacePollerResp, err := interfaceClient.BeginCreateOrUpdate(
		ctx,
		*resourceGroupResponse.Name,
		"azure-go-sdk",
		armnetwork.Interface{
			Location: to.Ptr(location),
			Properties: &armnetwork.InterfacePropertiesFormat{
				NetworkSecurityGroup: &armnetwork.SecurityGroup{
					ID: networkSecurityGroupResponse.ID,
				},
				IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
					{
						Name: to.Ptr("azure-go-sdk"),
						Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
							PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
							Subnet: &armnetwork.Subnet{
								ID: subnetsResponse.ID,
							},
							PublicIPAddress: &armnetwork.PublicIPAddress{
								ID: publicIPAddressResponse.ID,
							},
						},
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		return err
	}

	networkInterfaceResponse, err := interfacePollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully created network interface: %s\n", *vnetResponse.Name)

	// Create VM (virtual machine)
	fmt.Println("Starting VM creation...")
	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}

	parameters := armcompute.VirtualMachine{
		Location: to.Ptr(location),
		Identity: &armcompute.VirtualMachineIdentity{
			Type: to.Ptr(armcompute.ResourceIdentityTypeNone),
		},
		Properties: &armcompute.VirtualMachineProperties{
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: &armcompute.ImageReference{
					Offer:     to.Ptr("0001-com-ubuntu-server-focal"),
					Publisher: to.Ptr("canonical"),
					SKU:       to.Ptr("20_04-lts-gen2"),
					Version:   to.Ptr("latest"),
				},
				OSDisk: &armcompute.OSDisk{
					Name:         to.Ptr("azure-go-sdk"),
					CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
					Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
					ManagedDisk: &armcompute.ManagedDiskParameters{
						StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardLRS), // OSDisk type Standard/Premium HDD/SSD
					},
					DiskSizeGB: to.Ptr[int32](50), // default 127G
				},
			},
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: to.Ptr(armcompute.VirtualMachineSizeTypes("Standard_B1s")), // VM size include vCPUs,RAM,Data Disks,Temp storage.
			},
			OSProfile: &armcompute.OSProfile{ //
				ComputerName:  to.Ptr("azure-go-sdk"),
				AdminUsername: to.Ptr("demo"),
				LinuxConfiguration: &armcompute.LinuxConfiguration{
					DisablePasswordAuthentication: to.Ptr(true),
					SSH: &armcompute.SSHConfiguration{
						PublicKeys: []*armcompute.SSHPublicKey{
							{
								Path:    to.Ptr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", "demo")),
								KeyData: pubKey,
							},
						},
					},
				},
			},
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
					{
						ID: networkInterfaceResponse.ID,
					},
				},
			},
		},
	}

	vmPollerResp, err := vmClient.BeginCreateOrUpdate(
		ctx,
		*resourceGroupResponse.Name,
		"azure-go-sdk",
		parameters,
		nil,
	)
	if err != nil {
		return err
	}
	vmResponse, err := vmPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully created VM: %s: %s (%s)\n", *vmResponse.Name, *vmResponse.ID, *vmResponse.Location)

	return nil
}

func findVnet(ctx context.Context, resourceGroupName, vnetName string, vnetClient *armnetwork.VirtualNetworksClient) (armnetwork.VirtualNetwork, bool, error) {
	vnet, err := vnetClient.Get(ctx, resourceGroupName, vnetName, nil)
	if err != nil {
		var errResponse *azcore.ResponseError
		if errors.As(err, &errResponse) && errResponse.ErrorCode == "ResourceNotFound" {
			return armnetwork.VirtualNetwork{}, false, nil
		}
		return armnetwork.VirtualNetwork{}, false, err
	}

	return vnet.VirtualNetwork, true, nil
}
