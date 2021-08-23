package main

import "./single/snet"

func main() {
	s := snet.NewConn("0.0.0.0", 8777)
	s.Start()
}
