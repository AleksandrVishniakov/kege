package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Variant struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Key string `json:"key"`
}

func main() {
	ctx := context.Background()
	
	args := os.Args
	if len(args) != 2 {
		log.Fatalln("Provide task or variant id.")
	}
	
	// 25013197
	id := os.Args[1]
	
	if len(id) < 8 {
		task := findTask(ctx, id)
		fmt.Println(formatKey(task.Key))
	} else {
		variant := findVariant(ctx, id)
		for i, t := range variant.Tasks {
			fmt.Printf("%d. %s\n", i+1, formatKey(t.Key))
		}
	}
}

func findTask(ctx context.Context, id string) Task {
	cl := http.Client{}
	ctx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()
	url := taskURL(id)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	
	resp, err := cl.Do(req)
	if err != nil {
		log.Fatalf("Fail to do request to %s: %s\n", url, err.Error())
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request to %s returned %s.\n", url, resp.Status)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Read response body: %s\n", err.Error())
	}
	
	var task Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		log.Fatalf("Unmarshal response body: %s\n", err.Error())
	}
	
	return task
}

func findVariant(ctx context.Context, id string) Variant {
	cl := http.Client{}
	ctx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()
	url := variantURL(id)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	
	resp, err := cl.Do(req)
	if err != nil {
		log.Fatalf("Fail to do request to %s: %s\n", url, err.Error())
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request to %s returned %s.\n", url, resp.Status)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Read response body: %s\n", err.Error())
	}
	
	var variant Variant
	err = json.Unmarshal(body, &variant)
	if err != nil {
		log.Fatalf("Unmarshal response body: %s\n", err.Error())
	}
	
	return variant
}

func taskURL(id string) string {
	return fmt.Sprintf("https://kompege.ru/api/v1/task/%s", id)
}

func variantURL(id string) string {
	return fmt.Sprintf("https://kompege.ru/api/v1/variant/kim/%s", id)
}

func formatKey(key string) string {
	return strings.ReplaceAll(key, "\\n", "\n")
}