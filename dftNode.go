package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/zededa/zeddp/src/functionTracer/proto/dft"
)

/*type Job struct {
	ClientID	string	`json:"client_id"`
    JobID		string	`json:"job_id"`
	Tasks		[]dft.Task	`json:"tasks"`
}

type Task struct {
	TaskID		string	`json:"task_id"`
	Type		string	`json:"type"`
	Repetitions	string	`json:"repetitions"`
	Destination	string	`json:"destination"`
	Timeout		string	`json:"timeout"`
	//StartTime	string	`json:"start_time"`
}*/

/*type Result struct {
	ClientID	string	`json:"client_id"`
	JobID		string	`json:"job_id"`
	Results		string	`json:"results"`
}*/

func pingCmd(task dft.Task, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	result := "Task ID: " + task.TaskId + "\n"
	success := true
	var min, max, avg, sum, stdDev, stdDevSum float64
	min = math.MaxFloat64
	max = -1 * math.MaxFloat64
	var count int
	repetitions, _ := strconv.Atoi(task.Repetitions)
	latencies := make([]float64, repetitions)
	var IPAddr string
	for i := 0; i < repetitions; i++ {
		//Use timeout instead of gtimeout and include -4 for IPv4 or -6 for IPv6 on Linux
		cmd := "timeout " + task.Timeout + " ping -c 1 " + task.Destination
		out, cmdError := exec.Command("sh", "-c", cmd).Output()
		if cmdError != nil {
			if cmdError.Error() == "exit status 124" {
				result += "Ping Result: Timeout\n"
			} else if cmdError.Error() == "exit status 68" {
				result += "Ping Result: Unknown Host\n"
			} else if cmdError.Error() == "exit status 2" {
				result += "Ping Result: Host Unreachable\n"
			} else {
				result += "Ping Result: Other Error\n"
			}
			fmt.Printf("%s\n", cmdError.Error())
			success = false
			break
		}

		var IPAddress = regexp.MustCompile(`PING (.+) \((\d+).(\d+).(\d+).(\d+)\)`)
		var average = regexp.MustCompile(`min\/avg\/max\/(\w+) = (\d+\.\d+)\/(\d+\.\d+)\/(\d+\.\d+)\/(\d+\.\d+) ms`)
		latencyResult := average.FindAllStringSubmatch(string(out), -1)
		IPResult := IPAddress.FindAllStringSubmatch(string(out), -1)

		if len(IPResult) > 0 {
			IPAddr = IPResult[0][2] + "." + IPResult[0][3] + "." + IPResult[0][4] + "." + IPResult[0][5]
		}

		if len(latencyResult) > 0 {
			latency, _ := strconv.ParseFloat(latencyResult[0][3], 64)
			if latency > max {
				max = latency
			}
			if latency < min {
				min = latency
			}
			latencies[i] = latency
			sum += latency
			count++
		}
	}
	avg = sum / float64(count)
	for i := 0; i < count; i++ {
		stdDevSum += (latencies[i] - avg) * (latencies[i] - avg)
	}
	stdDev = math.Sqrt(stdDevSum / float64(count))
	if success {
		result += "Ping Result: Success\nDestination IP Address: " + IPAddr + "\nMin RTT: " + strconv.FormatFloat(min, 'f', -1, 64) + " ms\nAvg RTT: " + strconv.FormatFloat(avg, 'f', -1, 64) + " ms\nMax RTT: " + strconv.FormatFloat(max, 'f', -1, 64) + " ms\nStd Dev: " + strconv.FormatFloat(stdDev, 'f', -1, 64) + " ms\n\n"
	} else {
		result += "Destination IP Address: " + IPAddr + "\nMin RTT: " + strconv.FormatFloat(avg, 'f', -1, 64) + " ms\nAvg RTT: " + strconv.FormatFloat(avg, 'f', -1, 64) + " ms\nMax RTT: " + strconv.FormatFloat(avg, 'f', -1, 64) + " ms\nStd Dev: " + strconv.FormatFloat(stdDev, 'f', -1, 64) + " ms\n\n"
	}
	fmt.Println(result)
	ch <- result
}

func curlCmd(task dft.Task, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	result := "Task ID: " + task.TaskId + "\n"
	success := true
	var min, max, avg, sum, stdDev, stdDevSum float64
	min = math.MaxFloat64
	max = -1 * math.MaxFloat64
	var count int
	repetitions, _ := strconv.Atoi(task.Repetitions)
	latencies := make([]float64, repetitions)
	var IPAddr string
	for i := 0; i < repetitions; i++ {
		//Use timeout instead of gtimeout and include -4 for IPv4 or -6 for IPv6 on Linux
		cmd := "timeout " + task.Timeout + " curl -w \"Total Time: %{time_total}\nRemote IP: %{remote_ip}\" -v " + task.Destination
		out, cmdError := exec.Command("sh", "-c", cmd).Output()
		if cmdError != nil {
			if cmdError.Error() == "exit status 124" {
				result += "Curl Result: Timeout\n"
			} else if cmdError.Error() == "exit status 6" {
				result += "Curl Result: Couldn't Resolve Host\n"
			} else if cmdError.Error() == "exit status 5" {
				result += "Curl Result: Couldn't Resolve Proxy\n"
			} else if cmdError.Error() == "exit status 7" {
				result += "Curl Result: Failed to Connect to Host\n"
			} else {
				result += "Curl Result: Other Error\n"
			}
			fmt.Printf("%s\n", cmdError.Error())
			success = false
			break
		}
		var IPAddress = regexp.MustCompile(`Remote IP: (.+)`)
		var average = regexp.MustCompile(`Total Time: (.+)`)
		latencyResult := average.FindAllStringSubmatch(string(out), -1)
		IPResult := IPAddress.FindAllStringSubmatch(string(out), -1)

		if len(IPResult) > 0 {
			IPAddr = IPResult[0][1]
		}

		if len(latencyResult) > 0 {
			latency, _ := strconv.ParseFloat(latencyResult[0][1], 64)
			latency *= 1000
			if latency > max {
				max = latency
			}
			if latency < min {
				min = latency
			}
			latencies[i] = latency
			sum += latency
			count++
		}
	}
	avg = sum / float64(count)
	for i := 0; i < count; i++ {
		stdDevSum += (latencies[i] - avg) * (latencies[i] - avg)
	}
	stdDev = math.Sqrt(stdDevSum / float64(count))
	if success {
		result += "Curl Result: Success\nDestination IP Address: " + IPAddr + "\nMin Operation Time: " + strconv.FormatFloat(min, 'f', -1, 64) + " ms\nAvg Operation Time: " + strconv.FormatFloat(avg, 'f', -1, 64) + " ms\nMax Operation Time: " + strconv.FormatFloat(max, 'f', -1, 64) + " ms\nStd Dev: " + strconv.FormatFloat(stdDev, 'f', -1, 64) + " ms\n\n"
	} else {
		result += "Destination IP Address: " + IPAddr + "\nMin Operation Time: " + strconv.FormatFloat(avg, 'f', -1, 64) + " ms\nAvg Operation Time: " + strconv.FormatFloat(avg, 'f', -1, 64) + " ms\nMax Operation Time: " + strconv.FormatFloat(avg, 'f', -1, 64) + " ms\nStd Dev: " + strconv.FormatFloat(stdDev, 'f', -1, 64) + " ms\n\n"

	}
	fmt.Println(result)
	ch <- result
}

func main() {
	ClientID := flag.String("client_id", "-1", "This Node's client id")
	ServerURL := flag.String("server_url", "-1", "The server's base URL")
	flag.Parse()
	fmt.Printf("Node Client ID: %s\n", *ClientID)
	URL := *ServerURL + "/client/get-job/" + *ClientID
	client := http.Client{Timeout: time.Second * 2}
	for {
		getResponse, err := client.Get(URL)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		body, _ := ioutil.ReadAll(getResponse.Body)
		getResponse.Body.Close()
		var job dft.Job
		_ = json.Unmarshal(body, &job)

		ch := make(chan string, len(job.Tasks)+1)
		var wg sync.WaitGroup
		results := ""
		for _, task := range job.Tasks {
			if task.Type == "ping" {
				wg.Add(1)
				go pingCmd(*task, ch, &wg)
			} else if task.Type == "curl" {
				wg.Add(1)
				go curlCmd(*task, ch, &wg)
			}

		}
		wg.Wait()
		for i := 0; i < len(job.Tasks); i++ {
			results += <-ch
		}
		resultURL := *ServerURL + "/client/post-result/" + job.JobId

		values := map[string]string{"client_id": job.ClientId, "job_id": job.JobId, "results": results}
		jsonValue, _ := json.Marshal(values)
		postResponse, _ := client.Post(resultURL, "application/json", bytes.NewBuffer(jsonValue))
		postResponse.Body.Close()
		time.Sleep(time.Second)
	}
}
