package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
)

//exitErrorf() function to format error messages
func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func printErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)

}

func newSessionFromSessionOptions(options *session.Options) session.Session {
	sess, err := session.NewSessionWithOptions(*options)
	if err != nil {
		exitErrorf("Something went wrong when creating aws session object %v", err)
	}
	return *sess
}

func CheckMotoServerPort() bool {
	timeout := time.Second
	host := net.JoinHostPort("127.0.0.1", "5000")
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		fmt.Println("Connecting error:", err)
		return false
	}
	if conn != nil {
		defer conn.Close()
		fmt.Println("Opened", host)
	}
	return true
}

func WaitForMotoServer() bool {

	var index int = 0
	for !CheckMotoServerPort() && index <= 60 {
		time.Sleep(time.Second)
		index++
	}
	return index < 60
}

func StartAndWaitForMotoServer() *exec.Cmd {

	motoServerCommand := exec.Command("python", "-m", "moto.server")
	if CheckMotoServerPort() {
		fmt.Println("MotoServer was started by external means, skipping...")
		return nil
	}

	err := motoServerCommand.Start()
	if err != nil {
		printErrorf("Something went wrong when starting moto server. %v", err)
		motoServerCommand.Process.Kill()
	}

	WaitForMotoServer()

	return motoServerCommand
}

func StopMotoServer(cmd *exec.Cmd) {
	if cmd == nil {
		fmt.Println("MotoServer was started by external means, can't stop...")
		return
	}
	err := cmd.Process.Kill()
	if err != nil {
		printErrorf("Failed to stop moto server with process %v.\n %v", cmd.Process.Pid, err)
	}
}
