package gochat

import (
  	"os"
	"net"
	"fmt"
	"bufio"
	"strings"
	"container/list"
)

type Server struct {
	listener		*net.TCPListener;
	incomingMessages	chan Message;
	registerForMessages	chan *Client;
	clients			*list.List;
}

type Message struct {
	sender	string;
	message	string;
}

func (s *Server) StartServer() {
	s.incomingMessages = make(chan Message);
	s.registerForMessages = make(chan *Client);
	s.clients = list.New();

	go s.listenForConnections();
	s.startRouter();
}

func (s *Server) startRouter() {
	for {
		select {
		case msg := <-s.incomingMessages:
			s.fanOutMessage(msg)
		case newClient := <-s.registerForMessages:
			msg := Message{"**** system ****", "New user " + newClient.nickname + " has joined.\n"};
			s.fanOutMessage(msg);

			s.clients.PushFront(newClient);

			newClient.outgoingMessages <- s.clientList();
		}
	}
}

func (s *Server) clientList() (m Message) {
	clientNicknames := make([]string, s.clients.Len());
	
	i := 0;
	for e := s.clients.Front(); e != nil; e = e.Next() {
		client := e.Value.(Client);
		clientNicknames[i] = client.nickname;
		i++;
	}

	str := "Online users: " + strings.Join(clientNicknames, "   ") + "\n";

	return Message{"**** system ****", str};
}

func (s *Server) fanOutMessage(msg Message) {
	clientNicknames := make([]string, s.clients.Len());
	
	i := 0;
	for e := s.clients.Front(); e != nil; e = e.Next() {
		client := e.Value.(Client)
		ch := client.outgoingMessages

		val, ok := <- ch
		if ok {
			ch <- msg
			
		} else {
			s.clients.Remove(e);
			i--;
		}
		i++;
	}
	
}

func (s *Server) listenForConnections() {
	ip := net.ParseIP("127.0.0.1");
	addr := &net.TCPAddr{ip, 9999};
	s.listener, _ = net.ListenTCP("tcp", addr);

	for {
		s.acceptClient()
	}
}

func (s *Server) acceptClient() {
	conn, _ := s.listener.AcceptTCP();
	client := newClient(conn, s.incomingMessages);

	client.requestNick();
	s.registerForMessages <- client;
	client.sendReceiveMessages();
}

type Client struct {
	conn			*net.TCPConn;
	nickname		string;
	incomingMessages	chan Message;
	outgoingMessages	chan Message;
	reader			*bufio.Reader;
}

func newClient(conn *net.TCPConn, incoming chan Message) (c *Client) {
	c = new(Client);
	c.conn = conn;
	c.incomingMessages = incoming;
	c.outgoingMessages = make(chan Message);
	c.reader = bufio.NewReader(c.conn);
	return c;
}

func (c *Client) requestNick() {
	c.conn.Write([]byte("Please enter your nickname: "));

	nickname, _ := c.reader.ReadString('\n');

	// This is kinda stupid, but hell.
	if strings.HasSuffix(nickname, "\r\n") {
		c.nickname = nickname[0 : len(nickname)-2]
	} else {
		c.nickname = nickname[0 : len(nickname)-1]
	}
}

func (c *Client) sendReceiveMessages() {
	go c.receiveMessages();
	go c.sendMessages();
}

func (c *Client) receiveMessages() {
	for {
		bytes, err := c.reader.ReadString('\n');
		if err == nil {
			msg := Message{c.nickname, bytes};
			c.incomingMessages <- msg;
		}
		if err == os.EOF {
			close(c.outgoingMessages);
			msg := Message{"**** system ****", "User " + c.nickname + " has left.\n"};
			c.incomingMessages <- msg;
			return;
		}
	}
}

func (c *Client) sendMessages() {
	for {
		msg := <-c.outgoingMessages;
		if closed(c.outgoingMessages) {
			return
		}
		if msg.sender != c.nickname {
			c.sendMessage(msg)
		}
	}
}

func (c *Client) sendMessage(msg Message) {
	str := fmt.Sprintf("%s: %s", msg.sender, msg.message);
	c.conn.Write([]byte(str));
}
