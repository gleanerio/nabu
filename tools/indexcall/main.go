package main

import (
	"github.com/neuml/txtai.go"
)

func main() {
	embeddings := txtai.Embeddings("http://localhost:8000")

	// embeddings.Add(documents)
	embeddings.Index()

	// fmt.Println("\nBuilding an Embeddings index")
	// fmt.Printf("%-20s %s\n", "Query", "Best Match")
	// fmt.Println(strings.Repeat("-", 50))

	// for _, query := range []string{"Ridgecrest earthquake", "Lidar mapping of washington", "Island volcanos", "fault centers"} {
	// 	results := embeddings.Search(query, 1)
	// 	argmax, _ := strconv.Atoi(results[0].Id)
	// 	fmt.Printf("%-20s ::  %s\n\n", query, sections[argmax])
	// }
}
