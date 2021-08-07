package main

import "testing"

func TestMain(m *testing.M) {

	//This block of code allows for setup and teardowns at the package level.
	cmd := StartAndWaitForMotoServer()
	defer StopMotoServer(cmd)

	m.Run()

}
