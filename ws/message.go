package ws

type Message struct {
	Host string `json:"host"`
	Name string `json:"name"`
	Text string `json:"text"`
}

func (self *Message) String() string {
	return self.Host + "::" + self.Name + "::" + self.Text
}
