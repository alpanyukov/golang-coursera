package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// сюда писать код

func ExecutePipeline(jobs ...job) {
	pipeline := make([]chan interface{}, len(jobs)+1)
	pipeline[0] = make(chan interface{}, 1)
	pipeline[1] = make(chan interface{}, 1)

	go func() {
		jobs[0](pipeline[0], pipeline[1])
		close(pipeline[0])
		close(pipeline[1])
	}()

	for i := 1; i < len(jobs)-1; i++ {
		pipeline[i+1] = make(chan interface{}, 1)
		worker := jobs[i]
		in := pipeline[i]
		out := pipeline[i+1]
		go func() {
			worker(in, out)
			close(out)
		}()
	}

	lastOutIndex := len(jobs) - 1
	jobs[lastOutIndex](pipeline[lastOutIndex], pipeline[lastOutIndex+1])
}

func SingleHash(in, out chan interface{}) {
	// wgMain := new(sync.WaitGroup)

	for input := range in {
		// wgMain.Add(1)
		// go func(in interface{}) {
		// 	defer wgMain.Done()
		data := fmt.Sprintf("%v", input)
		result := DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
		out <- result
		// }(input)
	}
	// wgMain.Wait()
}

func MultiHash(in, out chan interface{}) {
	for input := range in {
		data := fmt.Sprintf("%v", input)
		result := ""
		for i := 0; i < 6; i++ {
			th := DataSignerCrc32(strconv.Itoa(i) + data)
			result += th
		}
		out <- result
	}
}

func CombineResults(in, out chan interface{}) {
	result := make([]string, 0)

	for input := range in {
		data := input.(string)
		result = append(result, data)
	}

	sort.Strings(result)
	out <- strings.Join(result, "_")
}

func main() {
	ExecutePipeline([]job{
		job(func(in, out chan interface{}) {
			fmt.Println("Place in channel 1")
			out <- uint32(1)
			fmt.Println("Place in channel 3")
			out <- uint32(3)
			fmt.Println("Place in channel 4")
			out <- uint32(4)
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			fmt.Printf("Total result: %v\n", <-in)
		}),
	}...)
}
