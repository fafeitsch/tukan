package tukan

import (
	"fmt"
	"sync"
)

type ScanResult struct {
	Address string
	Success bool
	Comment string
}

func (p ConnectResults) Scan() chan ScanResult {
	result := make(chan ScanResult)
	var wg sync.WaitGroup
	for connection := range p {
		wg.Add(1)
		go func(connection ConnectResult) {
			defer wg.Done()
			scanResult := ScanResult{Address: connection.Address}
			if connection.Phone != nil && connection.Error == nil {
				err := connection.Phone.Logout()
				if err != nil {
					scanResult.Success = true
					scanResult.Comment = "connection established and login successful"
				} else {
					scanResult.Comment = fmt.Sprintf("connection established and log in successful, but logout not: %v", err)
				}
			} else {
				scanResult.Comment = fmt.Sprintf("could not connect: %v", connection.Error)
			}
			result <- scanResult
		}(connection)
	}
	go func() {
		wg.Wait()
		close(result)
	}()
	return result
}
