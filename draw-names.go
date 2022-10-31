package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"gopkg.in/yaml.v3"
	"net/smtp"
	"os"
)

type Conf struct {
	Sender string `yaml:"sender"`
	Password string `yaml:"password"`
	Host string `yaml:"host"`
	Seed int64 `yaml:"seed"`
	People []struct {
		Name string `yaml:"name"`
		Email string `yaml:"email"`
	}
}

type Node struct {
	Next *Node
	Name string
	Email string
}

func readConf(filename string) (*Conf, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	conf := &Conf{}
	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return conf, err
}

func createAssignments(config *Conf) *Node {
	nodes := []*Node{}
	for _, person := range config.People {
		nodes = append(nodes, &Node{Name : person.Name, Email : person.Email})
	}
	head := nodes[0]
	unassigned := len(nodes)
	current := head
	rand.Seed(config.Seed)
	potensh := rand.Intn(len(nodes))
	for unassigned > 1 {
		for nodes[potensh].Name == current.Name || nodes[potensh].Next != nil {
			potensh = rand.Intn(len(nodes))
		}
		current.Next = nodes[potensh]
		current = current.Next
		unassigned--
	}
	current.Next = head
	return head
}

func sendEmails(configFile string) {
	config, _ := readConf(configFile)
	head := createAssignments(config)
	port := "587"
	auth := smtp.PlainAuth("", config.Sender, config.Password, config.Host)
	current := head
	headNeedsToBeSent := true
	for headNeedsToBeSent || current != head {
		msg := current.Next.Name
		body := []byte(msg)
		err := smtp.SendMail(config.Host+":"+port, auth, config.Sender, []string{current.Email}, body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		headNeedsToBeSent = false
		current = current.Next
	}
}


func main() {
	sendEmails(os.Args[1])
}
