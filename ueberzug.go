package main

import (
	"encoding/json"
	"io"
	"os/exec"
)

type Ueberzug struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
}

type ueberzugCommand struct {
	Action     string `json:"action"`
	Identifier string `json:"identifier"`
	Path       string `json:"path,omitempty"`
	X          int    `json:"x,omitempty"`
	Y          int    `json:"y,omitempty"`
	MaxWidth   int    `json:"max_width,omitempty"`
	MaxHeight  int    `json:"max_height,omitempty"`
}

func NewUeberzug() (*Ueberzug, error) {
	cmd := exec.Command("ueberzug", "layer", "--silent")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Ueberzug{
		cmd:   cmd,
		stdin: stdin,
	}, nil
}

func (u *Ueberzug) Show(path string, x, y, maxWidth, maxHeight int) error {
	cmd := ueberzugCommand{
		Action:     "add",
		Identifier: "preview",
		Path:       path,
		X:          x,
		Y:          y,
		MaxWidth:   maxWidth,
		MaxHeight:  maxHeight,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = u.stdin.Write(data)
	return err
}

func (u *Ueberzug) Hide() error {
	cmd := ueberzugCommand{
		Action:     "remove",
		Identifier: "preview",
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = u.stdin.Write(data)
	return err
}

func (u *Ueberzug) Close() error {
	if u.stdin != nil {
		u.stdin.Close()
	}
	if u.cmd != nil && u.cmd.Process != nil {
		return u.cmd.Process.Kill()
	}
	return nil
}
