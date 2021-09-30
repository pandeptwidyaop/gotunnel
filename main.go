package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/elliotchance/sshtunnel"
	"golang.org/x/crypto/ssh"
)

type Connection struct {
	Name        string  `json:"name"`
	Host        string  `json:"host"`
	Destination string  `json:"destination"`
	Local       string  `json:"local"`
	Username    string  `json:"username"`
	Password    *string `json:"password"`
	Key         *string `json:"key"`
	KeyPassword *string `json:"Key_password"`
}

type Config struct {
	Connection []Connection `json:"connection"`
}

var clear map[string]func()

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
		fmt.Printf("%d. %s (%s|%s -> 127.0.0.1:%s)\n", no, c2.Name, c2.Host, c2.Destination, c2.Local)
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
	CallClear()
	fmt.Println("Starting to tunnel")
	auth := LoadAuth(c)
	tunnel := sshtunnel.NewSSHTunnel(
		fmt.Sprintf("%s@%s", c.Username, c.Host),
		auth,
		c.Destination,
		c.Local,
	)
	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	err := tunnel.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func LoadAuth(c *Connection) ssh.AuthMethod {
	if c.Key != nil {
		if c.KeyPassword != nil {
			return PrivateKeyFileWithPassword(*c.Key, *c.Password)
		}
		return sshtunnel.PrivateKeyFile(*c.Key)
	}
	if c.Password == nil {
		return nil
	}
	return ssh.Password(*c.Password)
}

func PrivateKeyFileWithPassword(key string, password string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(key)
	if err != nil {
		return nil
	}

	pem, err := ssh.ParsePrivateKeyWithPassphrase(buffer, []byte(password))
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(pem)
}

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
