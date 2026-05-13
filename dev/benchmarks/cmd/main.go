package main

import (
	"fmt"
	"os"
	"search-benchmark/internal/handler"
	"search-benchmark/internal/service"
	"time"
)


func main() {
	evalPath := "../internal/dataset/eval_set.json"

	issueQueries, err := service.ExtractIssueQueries(evalPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error extracting issue queries from json %v", err)
		os.Exit(1)
	}

	endpoint := "http://localhost:8081/search/new"
	userId := "test-user-1"

	var positions []int
	totalQueries := 0

	start := time.Now()
	for id, queries := range issueQueries {
		for _, query := range queries {
			results, err := handler.CallSearchEndpoint(endpoint, query, userId)
			if err != nil {
				fmt.Printf("query error for %s: %v", id, err)
				continue
			}

			rank := 0
			for i, result := range results {
				if result.URL == id {
					rank = i + 1
					break
				}
			}
			positions = append(positions, rank)
			totalQueries += 1
		}
	}

	latencyMs := time.Since(start).Milliseconds()

	fmt.Printf("Benchmarking took %d ms\n", latencyMs)
	fmt.Printf("Average query took: %d ms\n", latencyMs/int64(totalQueries))
	fmt.Printf("Queries evaluated : %d\n", len(positions))
	fmt.Printf("MRR 			  : %.4f\n", service.MRR(positions))
	fmt.Printf("Recall@10 		  : %.4f\n", service.RecallAt10(positions))		
	fmt.Printf("NDCG@10 	      : %.4f\n", service.NDCGAt10(positions))
}
