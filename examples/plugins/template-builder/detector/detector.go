// Formatted with gofmt -s
package detector

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/types"
)

type fileScanner struct {
	wg       sync.WaitGroup
	mu       sync.Mutex
	stack    *types.ProjectStack
	scanners map[string]func(string) error
}

func DetectStack() *types.ProjectStack {
	stack := &types.ProjectStack{
		Dependencies: make(map[string]string),
	}

	scanner := &fileScanner{
		stack: stack,
		scanners: map[string]func(string) error{
			"package.json":     scanNodeJS,
			"requirements.txt": scanPython,
			"go.mod":           scanGo,
			".env":             scanEnv,
			"Dockerfile":       scanDocker,
		},
	}

	// Start concurrent file scanning
	scanner.scanFiles()

	// Scan k8s files separately as they use glob patterns
	scanKubernetes(stack)

	return stack
}

func (s *fileScanner) scanFiles() {
	// Create buffered channel for files
	files := make(chan string, 100)

	// Start worker pool
	for i := 0; i < 4; i++ {
		go s.worker(files)
	}

	// Walk directory and send files to workers
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			if _, ok := s.scanners[filepath.Base(path)]; ok {
				s.wg.Add(1)
				files <- path
			}
		}
		return nil
	})

	close(files)
	s.wg.Wait()
}

func (s *fileScanner) worker(files <-chan string) {
	for file := range files {
		if scanner, ok := s.scanners[filepath.Base(file)]; ok {
			if err := scanner(file); err == nil {
				s.mu.Lock()
				updateStack(s.stack, file)
				s.mu.Unlock()
			}
		}
		s.wg.Done()
	}
}

func scanNodeJS(file string) error {
	data, err := readFile(file)
	if err != nil {
		return err
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	return json.Unmarshal(data, &pkg)
}

func scanPython(file string) error {
	_, err := readFile(file)
	return err
}

func scanGo(file string) error {
	_, err := readFile(file)
	return err
}

func scanEnv(file string) error {
	_, err := readFile(file)
	return err
}

func scanDocker(file string) error {
	return nil
}

func scanKubernetes(stack *types.ProjectStack) {
	matches, _ := filepath.Glob("k8s/*.yaml")
	if len(matches) > 0 {
		stack.HasKubernetes = true
	}
}

func updateStack(stack *types.ProjectStack, file string) {
	switch filepath.Base(file) {
	case "package.json":
		stack.Language = "nodejs"
		stack.Frontend = "react"
	case "requirements.txt":
		stack.Language = "python"
		stack.Framework = "flask"
	case "go.mod":
		stack.Language = "go"
		stack.Framework = "gin"
	case "Dockerfile":
		stack.HasDocker = true
	}
}

// Optimized file reader that uses a buffer pool
var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 32*1024)
	},
}

func readFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get buffer from pool
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)

	// Read file into buffer
	var data []byte
	for {
		n, err := f.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}
