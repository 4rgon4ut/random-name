package ioutils

import (
	"bufio"
	"bytes"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/4rgon4ut/lightblocks-assignment/pkg/types"
)

// Reads cmds from file in goroutine and sends them to a channel
func ReadCommandsCh(scanner *bufio.Scanner) (<-chan []byte, error) {
    cmdCh := make(chan []byte)

    go func() {
        defer close(cmdCh)

        var buffer bytes.Buffer
        for scanner.Scan() {
            buffer.Write(scanner.Bytes())
        }

        if err := scanner.Err(); err != nil {
            log.Errorf("error reading file: %v", err)
            return
        }

        var commands []types.Command
        if err := json.Unmarshal(buffer.Bytes(), &commands); err != nil {
            log.Errorf("error unmarshalling JSON: %v", err)
            return
        }

        for i, cmd := range commands {
            if !cmd.IsValid() {
                log.Errorf("invalid command at index %d: %v", i, cmd)
                continue
            }

            message, err := json.Marshal(cmd)
            if err != nil {
                log.Errorf("error marshalling command at index %d: %v", i, err)
                continue
            }

            cmdCh <- message
        }
    }()

    return cmdCh, nil
}
