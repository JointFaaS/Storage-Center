package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/JointFaaS/Storage-Center/benchmark/utils"
	"github.com/spf13/cobra"
)

var (
	clientAddr *[]string        = &[]string{}
	clients    *[]*utils.Client = &[]*utils.Client{}
)

var rootCmd = &cobra.Command{
	Use:   "sc-tester",
	Short: "simple tester",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
}

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "random test",
	Run: func(cmd *cobra.Command, args []string) {
		var m map[string]string = make(map[string]string)
		var count int = 100
		var length int = len(*clients)
		var timeSlice []time.Duration = make([]time.Duration, 1000, 1000)
		var indexRecord []int = make([]int, 1000, 1000)
		start := time.Now()
		for i := 0; i < count; i++ {
			key := strconv.FormatInt(rand.Int63(), 10)
			index := rand.Int() % length
			m[key] = strconv.FormatInt(rand.Int63(), 10)
			err := (*clients)[index].Set(key, m[key])
			if err != nil {
				panic(err)
			}
			indexRecord[i*10] = index
			timeSlice[i*10] = time.Since(start)
			for j := 1; j < 10; j++ {
				index = rand.Int() % length
				indexRecord[i*10+j] = index
				s, err := (*clients)[index].Get(key)
				if err != nil {
					panic(err)
				} else if *s != m[key] {
					log.Printf("Unexpected Value %s, the record is %s\n", *s, m[key])
				}
				timeSlice[i*10+j] = time.Since(start)
			}
		}
		for i := 0; i < 1000; i++ {
			log.Println(int(timeSlice[i]))
			log.Println(indexRecord[i])
		}
	},
}

var localityCmd = &cobra.Command{
	Use:   "locality",
	Short: "locality test",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var worstCmd = &cobra.Command{
	Use:   "worst",
	Short: "worst test",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func clientInit() {
	*clients = make([]*utils.Client, 0, len(*clientAddr))
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
