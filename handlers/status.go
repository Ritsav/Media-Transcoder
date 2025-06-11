package handlers

import (
	"encoding/json"
	"fmt"
	"media_transcoder/pkg/global"
	"net/http"
)

// Returns the current queueStatus
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get Method Allowed at this endpoint", http.StatusMethodNotAllowed)
		return
	}

	// Return TaskQueue Status
	status, _ := json.Marshal(global.TaskQueue.Status())

	fmt.Fprintln(w, string(status))
}
