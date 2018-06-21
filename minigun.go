package main

import(
	"fmt"
	"log"
	"os"
	"bufio"
	"flag"
	"sync"
	"os/exec"
	"encoding/base64"
)

var (
	command string
	hostsFile string
	args string
	output string
	json string
	stdout bool
	hostsList []string
	commandOutput []string
	wg sync.WaitGroup
)


func shot(c chan string, command string, args string, host string){
	defer wg.Done()
	defer fmt.Println("[+] Done with", host)
	fmt.Println("[+] Starting with", host)
	cmd := exec.Command(command, args, host)
	stdoutStderr, err := cmd.CombinedOutput()
	check(err)
	c <- string("\n[+] Start: " + host + "\n" + string(stdoutStderr) + "\n[+] Done: " + host + "\n") //TODO casting
}

func check(err error){
	if err != nil {
		log.Fatal(err)
        panic(err)
    }
}
func printBanner(){
	bannerBase64 := `ICAgICAgICAgICAgIF9fX19fX19fXw0KICAgICAgICAgICAgLycgICAgICAgIC98DQogICAgICAgICAgIC8gICAgICAgICAvIHxfDQogICAgICAgICAgLyAgICAgICAgIC8gIC8vfA0KICAgICAgICAgL19fX19fX19fXy8gIC8vLy98DQogICAgICAgIHwgICBfIF8gICAgfCA4by8vLy98DQogICAgICAgIHwgLycvLyApXyAgfCAgIDgvLy98DQogICAgICAgIHwvIC8vIC8vICkgfCAgIDhvLy8vfA0KICAgICAgICAvIC8vIC8vIC8vLHwgIC8gIDgvL3wNCiAgICAgICAvIC8vIC8vIC8vLyB8IC8gICA4Ly98DQogICAgICAvIC8vIC8vIC8vL19ffC8gICAgOC8vfA0KICAgICAvLihfKS8vIC8vLyB8ICAgICAgIDgvLy98DQogICAgKF8pJyBgKF8pLy98IHwgICAgICAgOC8vLy98X19fX19fX19fX18NCiAgIChfKSAvX1wgKF8pJ3wgfCAgICAgICAgOC8vLy8vLy8vLy8vLy8vLw0KICAgKF8pIFwiLyAoXyknfF98ICAgICAgICAgOC8vLy8vLy8vLy8vLy8NCiAgICAoXykuXy4oXykgZCcgSGIgICAgICAgICA4b29vb29vb29wYicNCiAgICAgIGAoXyknICBkJyAgSGBiDQogICAgICAgICAgICBkJyAgIGBiYGINCiAgICAgICAgICAgZCcgICAgIEggYGINCiAgICAgICAgICBkJyAgICAgIGBiIGBiDQogICAgICAgICBkJyAgICAgICAgICAgYGINCiAgICAgICAgZCcgICAgICAgICAgICAgYGI=`
	banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
    fmt.Println(string(banner))
    fmt.Println("Mattia Reggiani - https://github.com/mattiareggiani/minigun - info@mattiareggiani.com")
    fmt.Println()
}
func main(){
	defer fmt.Println("[+] Completed")
	printBanner()
	flag.StringVar(&command, "command", "" , "Command to run")
	flag.StringVar(&args, "args", "" , "Arguments for command")
	flag.StringVar(&hostsFile, "hosts", "" , "Hosts file list")
	flag.StringVar(&output, "output", "" , "Output file")
	flag.BoolVar(&stdout, "stdout", false , "Print to standard output")
	flag.Parse()

	// Reading hosts list
	targetList, err := os.Open(hostsFile)
	check(err)
	defer targetList.Close()
	scanner := bufio.NewScanner(targetList)
	i := 0
	for scanner.Scan() {
		hostsList = append(hostsList, scanner.Text())
		i++
	}
	check(scanner.Err())
	fmt.Println("[*] Total hosts:",i)
	
	// Starting to shot
	chanOutput := make(chan string, i)
	for host := range hostsList{
    	wg.Add(1)
    	go shot(chanOutput, command, args, hostsList[host])
	}
	wg.Wait()
	close(chanOutput)

	// Fetching output from channels
	for c := range chanOutput{
			commandOutput = append(commandOutput,c)
		}

	// Print results
	if(stdout){
		fmt.Println("[*] Output of commands:")
		for c := range commandOutput{
			fmt.Println(commandOutput[c])
		}
	}
	if(len(output) > 1) {
		fileOutput, err := os.Create(output)
		check(err)
		defer fileOutput.Close()
		for c := range commandOutput{
			_, err := fileOutput.WriteString(commandOutput[c])
    		check(err)
		}
		fileOutput.Sync()
	}
}