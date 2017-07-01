package main

import (
	"flag"
	"github.com/yunify/hostnic-cni/pkg"
	"os"
	"fmt"
	"github.com/yunify/hostnic-cni/provider"
	"net"
)

var(
	cniConfig = "/etc/cni/net.d/10-hostnic.conf"
	onlyUpNic = true
)

func init() {
	flag.StringVar(&cniConfig, "cni_config", cniConfig, "the hostnic cni config file.")
	flag.BoolVar(&onlyUpNic, "only_up", onlyUpNic, "only clean up nic.")
}

func clean(n *pkg.NetConf) (error){
	nicProvider, err := provider.New(n.Provider, n.Args)
	if err != nil {
		return err
	}
	for _, vxnetID := range nicProvider.GetVxNets(){
		vxnet, err := nicProvider.GetVxNet(vxnetID)
		if err != nil {
			return err
		}
		_, network, err := net.ParseCIDR(vxnet.Network)
		if err != nil {
			return err
		}
		nics, err := pkg.ScanNicsByNetwork(network,onlyUpNic)
		if err != nil {
			return err
		}
		for _, nic := range nics {
			err := nicProvider.DeleteNic(nic)
			if err != nil {
				fmt.Printf("Delete nic %s error: %s\n", nic, err)
			}
		}
	}

}

func main() {
	flag.Parse()
	netConf, err := pkg.LoadNetConfFromFile(cniConfig)
	if err != nil {
		fmt.Println("Load net conf error:", err.Error())
		os.Exit(1)
	}
	err = clean(netConf)
	if err != nil {
		fmt.Println("Clean nic error:", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}