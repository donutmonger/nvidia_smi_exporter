package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

// name, index, temperature.gpu, utilization.gpu,
// utilization.memory, memory.total, memory.free, memory.used

func metrics(response http.ResponseWriter, request *http.Request) {
	out, err := exec.Command(
		"nvidia-smi",
		"--query-gpu=name,index,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used",
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

	metricList := []string{
		"temperature.gpu", "utilization.gpu",
		"utilization.memory", "memory.total", "memory.free", "memory.used"}

	result := ""
	for _, row := range records {
		name := fmt.Sprintf("%s[%s]", row[0], row[1])
		for idx, value := range row[2:] {
			result = fmt.Sprintf(
				"%s%s{gpu=\"%s\"} %s\n", result,
				metricList[idx], name, value)
		}
	}

	fmt.Fprintf(response, strings.Replace(result, ".", "_", -1))
}
