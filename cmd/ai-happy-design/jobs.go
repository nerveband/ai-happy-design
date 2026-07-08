package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nerveband/ai-happy-design/internal/config"
	"github.com/spf13/cobra"
)

type jobRecord struct {
	ID        string                 `json:"id"`
	Kind      string                 `json:"kind"`
	Status    string                 `json:"status"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

var jobsCmd = &cobra.Command{Use: "jobs", Short: "Inspect local async job ledger"}
var jobsListCmd = &cobra.Command{Use: "list", Short: "List jobs", RunE: func(cmd *cobra.Command, args []string) error {
	jobs, err := listJobRecords()
	if err != nil {
		return err
	}
	return printJSON(map[string]interface{}{"jobs": jobs, "count": len(jobs)})
}}
var jobsGetCmd = &cobra.Command{Use: "get <id>", Short: "Get a job", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	job, err := getJobRecord(args[0])
	if err != nil {
		return err
	}
	return printJSON(job)
}}
var jobsResumeCmd = &cobra.Command{Use: "resume <id>", Short: "Mark a job resumable/pending", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	job, err := updateJobStatus(args[0], "pending")
	if err != nil {
		return err
	}
	return printJSON(job)
}}
var jobsCancelCmd = &cobra.Command{Use: "cancel <id>", Short: "Cancel a job", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	if err := cancelJobRecord(args[0]); err != nil {
		return err
	}
	return printJSON(map[string]interface{}{"id": args[0], "status": "cancelled"})
}}

func jobsPath() string { return filepath.Join(config.Dir(), "jobs.jsonl") }

func createJobRecord(kind string, payload map[string]interface{}) (jobRecord, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	job := jobRecord{ID: fmt.Sprintf("job_%d", time.Now().UnixNano()), Kind: kind, Status: "pending", CreatedAt: now, UpdatedAt: now, Payload: payload}
	if err := appendJobRecord(job); err != nil {
		return jobRecord{}, err
	}
	return job, nil
}

func appendJobRecord(job jobRecord) error {
	if err := os.MkdirAll(config.Dir(), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(jobsPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc, _ := json.Marshal(job)
	_, err = f.Write(append(enc, '\n'))
	return err
}

func listJobRecords() ([]jobRecord, error) {
	f, err := os.Open(jobsPath())
	if os.IsNotExist(err) {
		return []jobRecord{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	latest := map[string]jobRecord{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var job jobRecord
		if err := json.Unmarshal(scanner.Bytes(), &job); err == nil && job.ID != "" {
			latest[job.ID] = job
		}
	}
	out := make([]jobRecord, 0, len(latest))
	for _, job := range latest {
		out = append(out, job)
	}
	return out, scanner.Err()
}

func getJobRecord(id string) (jobRecord, error) {
	jobs, err := listJobRecords()
	if err != nil {
		return jobRecord{}, err
	}
	for _, job := range jobs {
		if job.ID == id {
			return job, nil
		}
	}
	return jobRecord{}, fmt.Errorf("job not found: %s", id)
}

func updateJobStatus(id, status string) (jobRecord, error) {
	job, err := getJobRecord(id)
	if err != nil {
		return jobRecord{}, err
	}
	job.Status = status
	job.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return job, appendJobRecord(job)
}

func cancelJobRecord(id string) error {
	_, err := updateJobStatus(id, "cancelled")
	return err
}

func init() {
	jobsCmd.AddCommand(jobsListCmd, jobsGetCmd, jobsResumeCmd, jobsCancelCmd)
	rootCmd.AddCommand(jobsCmd)
}

func parsePayload(raw string) map[string]interface{} {
	out := map[string]interface{}{}
	if strings.TrimSpace(raw) != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}
