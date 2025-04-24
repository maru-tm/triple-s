package utils

import "fmt"

func PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("  --port N       Port number (default :8080)")
	fmt.Println("  --dir S       Path to the storage directory (default ./storage)")
	fmt.Println("  --help        Show this screen.")
}
