package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Result struct {
	StatusCode int
	Err        error
}

type Report struct {
	TotalTime    time.Duration
	TotalReqs    int
	Status200    int
	StatusCodes  map[int]int
	ErrorCount   int
}

func main() {
	url := flag.String("url", "", "URL do servico a ser testado")
	requests := flag.Int("requests", 0, "Numero total de requisicoes")
	concurrency := flag.Int("concurrency", 1, "Numero de chamadas simultaneas")
	flag.Parse()

	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		fmt.Println("Uso: stress-test --url=<URL> --requests=<N> --concurrency=<C>")
		fmt.Println("  --url         URL do servico a ser testado (obrigatorio)")
		fmt.Println("  --requests    Numero total de requisicoes (obrigatorio, > 0)")
		fmt.Println("  --concurrency Numero de chamadas simultaneas (obrigatorio, > 0)")
		os.Exit(1)
	}

	report := runLoadTest(*url, *requests, *concurrency)
	printReport(report)
}

func runLoadTest(url string, totalReqs, concurrency int) Report {
	results := make(chan Result, totalReqs)
	reqsChan := make(chan struct{}, totalReqs)

	for i := 0; i < totalReqs; i++ {
		reqsChan <- struct{}{}
	}
	close(reqsChan)

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{Timeout: 10 * time.Second}
			for range reqsChan {
				resp, err := client.Get(url)
				if err != nil {
					results <- Result{Err: err}
					continue
				}
				results <- Result{StatusCode: resp.StatusCode}
				resp.Body.Close()
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)
	close(results)

	report := Report{
		TotalTime:   elapsed,
		TotalReqs:   totalReqs,
		StatusCodes: make(map[int]int),
	}

	for r := range results {
		if r.Err != nil {
			report.ErrorCount++
			continue
		}
		report.StatusCodes[r.StatusCode]++
		if r.StatusCode == http.StatusOK {
			report.Status200++
		}
	}

	return report
}

func printReport(r Report) {
	fmt.Println("========================================")
	fmt.Println("       RELATORIO DE TESTE DE CARGA      ")
	fmt.Println("========================================")
	fmt.Printf("Tempo total:          %v\n", r.TotalTime.Round(time.Millisecond))
	fmt.Printf("Total de requests:    %d\n", r.TotalReqs)
	fmt.Printf("Requests com HTTP 200: %d\n", r.Status200)

	if r.ErrorCount > 0 {
		fmt.Printf("Erros de conexao:     %d\n", r.ErrorCount)
	}

	fmt.Println("----------------------------------------")
	fmt.Println("Distribuicao de status HTTP:")
	for code, count := range r.StatusCodes {
		fmt.Printf("  HTTP %d: %d\n", code, count)
	}
	fmt.Println("========================================")
}
