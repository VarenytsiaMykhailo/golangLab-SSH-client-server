package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
)
/*
127.0.0.1:2222
PASSWORD
iu9_31_05
secret
*/
/*
185.20.227.83:22
PASSWORD
iu9_31_05
123qwe
*/
func main() {
	println("ENTER THE ADDRESS TO CONNECT (ex: \"185.20.227.83:22\"):") //185.20.227.83:22
	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	addr := string(line)
	println("CHOOSE A CONNECTION METHOD \"PASSWORD\" OR \"KEY\" (ex: \"PASSWORD\"):")
	line, _, _ = bufio.NewReader(os.Stdin).ReadLine()
	connMethod := string(line)

	var config *ssh.ClientConfig
	if connMethod == "PASSWORD" {
		config = getConfigWithPass()
	} else if connMethod == "KEY" {
		config = getConfigWithKey()
	} else {
		log.Fatal("INCORRECT METHOD")
		return
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		panic(err)
	}

	println("ENTER \"EXIT\" TO QUIT")
	for {
		fmt.Println("ENTER THE COMMAND:")
		line, _, _ = bufio.NewReader(os.Stdin).ReadLine()
		command := string(line)
		fmt.Println("YOUR COMMAND:", command)
		if command == "EXIT" {
			break
		}
		sendCommandToServer(client, command)
	}
}

func getConfigWithPass() *ssh.ClientConfig {
	println("ENTER YOUR USERNAME:")
	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	username := string(line)
	println("ENTER YOUR PASSWORD:")
	line, _, _ = bufio.NewReader(os.Stdin).ReadLine()
	password := string(line)
	return &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
}

func getConfigWithKey() *ssh.ClientConfig { //пока не реализовано
	/*println("ENTER YOUR USERNAME:")
	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	username := string(line)*/

	/*
		key, err := ioutil.ReadFile("./client/key/1.ppk")
		if err != nil {
			println("TUT")
			panic(err)
		}
		println(string(key))
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			println("TUT2")
			panic(err)
		}*/
	return nil
}

func sendCommandToServer(client *ssh.Client, command string) {
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	b, err := session.CombinedOutput(command)
	fmt.Println("RESULT:")
	fmt.Println(string(b))
}
