package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/elliotchance/sshtunnel"
	"golang.org/x/crypto/ssh"
)

type Connection struct {
	Name            string  `json:"name"`
	Host            string  `json:"host"`
	Destination     string  `json:"destination"`
	Local           string  `json:"local"`
	Username        string  `json:"username"`
	Password        *string `json:"password"`
	CertificateFile *string `json:"certificate_file"`
}

type Config struct {
	Connection []Connection `json:"connection"`
}

func (c *Config) Len() int {
	count := 0
	for i := range c.Connection {
		count = i + 1

	}
	return count
}

func (c *Config) PrintConnection() {
	if c.Len() < 1 {
		fmt.Println("Tidak ada koneksi tersedia, silakan tambahkan pada file config.json")
		os.Exit(0)
	}
	fmt.Println("Silakan pilih koneksi untuk memulai tunneling")
	for i, c2 := range c.Connection {
		no := i + 1
		fmt.Printf("%d. %s \n", no, c2.Name)
	}
}

func main() {
	var index int
	conf := LoadConfig()
	conf.PrintConnection()
	fmt.Print("Pilih : ")
	fmt.Scanf("%d", &index)
	c := conf.Connection[index-1]
	c.Start()
}

func LoadConfig() Config {
	conf := Config{}
	config, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(config, &conf)
	if err != nil {
		fmt.Println(err)
	}
	return conf
}

func (c *Connection) Start() {
	fmt.Println("Starting to tunnel")
	tunnel := sshtunnel.NewSSHTunnel(
		fmt.Sprintf("%s@%s", c.Username, c.Host),
		ssh.Password(*c.Password),
		c.Destination,
		c.Local,
	)
	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	tunnel.Start()

}
