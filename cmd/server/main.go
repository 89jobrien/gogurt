package server

import (
	"encoding/json"
	"fmt"
	"gogurt/internal/tools"
	"log"
	"net/http"
)

func Serve() {
    // 1. Setup tool registry
    registry := tools.NewRegistry()
    errs := registry.RegisterBatch([]*tools.Tool{
        tools.PalindromeTool,
    })
    for _, err := range errs {
        if err != nil {
            fmt.Println("Registration error:", err)
        }
    }

    registry.PrintAllDescs()

    // Define HTTP handler for tool calls
    http.HandleFunc("/tool", func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Name string          `json:"name"`
            Args json.RawMessage `json:"args"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }
        // Call the tool with provided arguments (as JSON)
        result, err := registry.Call(req.Name, string(req.Args))
        resp := map[string]interface{}{"result": result, "error": ""}
        if err != nil {
            resp["error"] = err.Error()
        }
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(resp)
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    log.Println("Server running at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
