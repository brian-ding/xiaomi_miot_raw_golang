package main

import (
	"fmt"
	"os"

	"github.com/Brian-Ding/xiaomi_miot_raw_golang/internal/micloud"
)

func main() {
	fmt.Println("Hello, world!")
	fmt.Println(os.Args[1])
	fmt.Println(os.Args[2])
	cloud := micloud.NewMiCloud(os.Args[1], os.Args[2])
	cloud.LogIn()
	devices := cloud.GetDevices()
	fmt.Println(devices)
}
