package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/JointFaaS/Storage-Center/benchmark/utils"
)

var (
	clientAddr *[]string = &[]string{}
	clients *[]*utils.Client = &[]*utils.Client{}
)

var rootCmd = &cobra.Command{
	Use:   "sc-tester",
	Short: "simple tester",
	PreRun: func(cmd *cobra.Command, args []string) {
		
	},
}

var randomCmd = &cobra.Command{
	Use: "random",
	Short: "random test",
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}

var localityCmd = &cobra.Command{
	Use: "locality",
	Short: "locality test",
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}

var worstCmd = &cobra.Command{
	Use: "worst",
	Short: "worst test",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func clientInit() {
	*clients = make([]*utils.Client, len(*clientAddr))
	for _, addr := range *clientAddr {
		*clients = append(*clients, utils.NewClient(addr))
		log.Printf("add client '%s'", addr)
	}
	log.Println("client ready")
}

func rootInit() {
	cobra.OnInitialize(clientInit)

	rootCmd.PersistentFlags().StringArrayVarP(clientAddr, "clients", "c", []string{"client1:9091", "client2:9091"}, "Tested Clients Addr")

	rootCmd.AddCommand(randomCmd, localityCmd)
}

func main() {
	rootInit()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
