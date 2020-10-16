// https://github.com/zhebrak/nvidia_smi_exporter
// https://play.golang.org/p/xbIw1_6rGQ
// http://golang.site/go/article/116-%EC%9C%88%EB%8F%84%EC%9A%B0%EC%A6%88-%EC%84%9C%EB%B9%84%EC%8A%A4-%ED%94%84%EB%A1%9C%EA%B7%B8%EB%9E%A8
package main

import (
    "bytes"
    "encoding/csv"
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "strings"
    "golang.org/x/sys/windows/svc"
	"time"
)


// name, index, temperature.gpu, utilization.gpu,
// utilization.memory, memory.total, memory.free, memory.used

func metrics(response http.ResponseWriter, request *http.Request) {
    out, err := exec.Command("C:/Program Files/NVIDIA Corporation/NVSMI/nvidia-smi",  
        "--query-gpu=name,index,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used",
        "--format=csv,noheader,nounits").Output()
    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

    csvReader := csv.NewReader(bytes.NewReader(out))
    csvReader.TrimLeadingSpace = true
    records, err := csvReader.ReadAll()

    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

    metricList := []string {
        "utilization.gpu",
        "utilization.memory", "memory.total", "memory.free", "memory.used", "gpu.using.pid"}

    result := ""
    result_gpu_name := ""
    result_nth := ""
    for _, row := range records {
        result_gpu_name = row[0]
        result_nth = row[1]
        name := fmt.Sprintf("%s[%s]", result_gpu_name, result_nth)
        for idx, value := range row[2:] {
            result = fmt.Sprintf(
                "%s%s{gpu=\"%s\"} %s\n", result,
                metricList[idx], name, value)
        }
    }
    // print(result)

    fmt.Fprintf(response, strings.Replace(result, ".", "_", -1))
}

func run_webserver() {
    addr := ":9101"
    if len(os.Args) > 1 {
        addr = ":" + os.Args[1]
    }

    http.HandleFunc("/metrics/", metrics)
    go http.ListenAndServe(addr, nil)
}

type WindowsNvidiaExporter struct {}

func (srv *WindowsNvidiaExporter)  Execute(args []string, req <-chan svc.ChangeRequest, stat chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
    stat <- svc.Status{State: svc.StartPending}
    
    stopChan := make(chan bool, 1)
    go runBody(stopChan)
    stat <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}
 
LOOP:
    for {
        // 서비스 변경 요청에 대해 핸들링
        switch r := <-req; r.Cmd {
        case svc.Stop, svc.Shutdown:
            stopChan <- true
            break LOOP
 
        case svc.Interrogate:
            stat <- r.CurrentStatus
            time.Sleep(100 * time.Millisecond)
            stat <- r.CurrentStatus
 
        //case svc.Pause:
        //case svc.Continue:
        }
    }
 
    stat <- svc.Status{State: svc.StopPending}
    return
}

func runBody(stopChan chan bool) {
    for {
        select {
        case <- stopChan:
            return
        default:
        }
    }
}

func main() {
    go run_webserver()

    err := svc.Run("windows_nvidia_exporter", &WindowsNvidiaExporter{})
    if err != nil {
        panic(err)
    }
}
