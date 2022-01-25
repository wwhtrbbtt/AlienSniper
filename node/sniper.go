package node

import (
	types "Alien/types"
	"fmt"
	"strings"
	"sync"
	"time"
)

const delay = 0

var (
	bearers MCbearers
)

func StartSniper(i int, payload Payload, email string) ([]types.RequestLog, float64, bool) {
	var recv []string
	var requests []types.RequestLog
	var wg sync.WaitGroup

	for g := 0; g < 2; {
		wg.Add(1)
		go func() {
			recvd := make([]byte, 4069)
			fmt.Fprintln(payload.Conns[i], payload.Payload[i])
			payload.Conns[i].Read(recvd)
			recv = append(recv, fmt.Sprintf("%v:%v", time.Now().UnixMilli(), string(recvd[9:12])))

			wg.Done()
		}()

		g++
	}

	wg.Wait()

	var sniped bool = false

	for _, status := range recv {
		if strings.Split(status, ":")[1] == "200" {
			sniped = true
		}
	}

	requests = append(requests, types.RequestLog{
		Timestamp: recv,
		Email:     email,
		Ip:        c.RemoteAddr().String(),
	})

	return requests, float64(len(recv)), sniped
}

func StartSnipe(task types.Task) {

	var l types.Log

	l.Name = task.Name
	l.Delay = pingMojang()

	accounts := task.Accounts
	droptime := task.Timestamp

	// chans := make([]chan types.Logs, len(accounts))
	var logs []types.Log
	var requests []types.RequestLog
	var wg sync.WaitGroup
	var success bool = false
	var amount float64

	bearers = bearers.AddAccounts(accounts)

	PreSleep(droptime)

	payload := bearers.CreatePayloads(task.Name)

	Sleep(droptime, delay)

	for i := range payload.AccountType {
		wg.Add(1)
		go func(i int) {
			sends, reqs, status := StartSniper(i, payload, accounts[i].Email)

			amount = amount + reqs

			if status {
				success = true
			}

			requests = append(requests, sends...)

			wg.Done()
		}(i)
	}

	wg.Wait()

	l.Requests = amount
	l.Success = success

	l.Sends = append(l.Sends, &types.Sent{
		Content: requests,
	})

	logs = append(logs, l)

	bearers = bearers.RemoveAccounts()

	var p types.Packet

	p.Type = "send_logs"
	p.Content.Logs = logs

	// TODO: this isnt a very good way of doing it, but it works for now ig.
	handleMessage(p)

}
