package server

import (
	"fmt"
	"gogurt/internal/tools"
	"log"
	"net/http"
)



func Serve() {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	registry := tools.NewRegistry()
	errs := registry.RegisterBatch([]*tools.Tool{
		tools.PalindromeTool,
		tools.ConcatenateTool,
		tools.ReverseTool,
		tools.UppercaseTool,
		tools.AddTool,
		tools.DivideTool,
		tools.SubtractTool,
		tools.MultiplyTool,
	})
	for _, err := range errs {
		if err != nil {
			fmt.Println("Registration error:", err)
		}
	}

	stats := registry.Stats()
	fmt.Println("Tool Stats:")
	fmt.Printf("  Count: %d\n", stats.Count)
	fmt.Printf("  Available Tools: %v\n", stats.ToolNames)
	fmt.Printf("  Tool Categories: %v\n", stats.Categories)
	fmt.Printf("  Duplicates?: %v\n", stats.HasDups)
	fmt.Printf("  HasCategory: %v\n\n", stats.HasCategory)
	log.Println("Server running at :8080")
	
	handler := MiddlewareHandler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
