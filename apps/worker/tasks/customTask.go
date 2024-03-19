package tasks

import (
	"fmt"
	"time"
)

func ProcessDocument(jobId int) {
	// Fake task only
	fmt.Printf("Processing document for Job ID %v...", jobId)
	time.Sleep(5 * time.Second)
	fmt.Println("Done processing document.")
}
