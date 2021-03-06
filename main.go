package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	procs = flag.Int("p", runtime.NumCPU(), "how many processes to use for counting lines")
	exts  = flag.String("ext", "", "only extensions ('go,html,js')")
)

func ShouldExamine(ext string) bool {
	if *exts == "" {
		return true
	}
	ext = strings.ToLower("," + strings.TrimPrefix(ext, ".") + ",")
	return strings.Contains(*exts, ext)
}

func main() {
	flag.Parse()

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	if *exts != "" {
		*exts = "," + *exts + ","
	}

	var progress int64
	fmtprogress := func() string {
		v := atomic.LoadInt64(&progress)
		return "    files " + strconv.Itoa(int(v))
	}

	work := make(chan string, 100)
	results := make(chan CountByExt, *procs)
	go IterateDir(work, paths...)

	for i := 0; i < *procs; i++ {
		go FileWorker(work, results, &progress)
	}

	total := make(CountByExt)
	for N := *procs; N > 0; {
		select {
		case <-time.After(100 * time.Millisecond):
			line := fmtprogress()
			backspace := strings.Repeat("\r", len(line))
			fmt.Print(line, backspace)
		case result := <-results:
			for _, s := range result {
				total.Add(s)
			}
			N--
		}
	}
	fmt.Println(fmtprogress())
	fmt.Println()

	counts := make(Counts, 0, len(total))
	for _, c := range total {
		counts = append(counts, c)
	}
	sort.Sort(ByCode{counts})

	summary := Count{}

	fmt.Printf("%-12s %12s %12s %12s %12s\n", "extension", "files", "binary", "blank", "code")
	fmt.Println(strings.Repeat("-", 12*5+4))
	for _, count := range counts {
		summary.Add(count)
		fmt.Printf("%-12s %12d %12d %12d %12d\n", count.Ext, count.Files, count.Binary, count.Blank, count.Code)
	}
	fmt.Println(strings.Repeat("-", 12*5+4))
	fmt.Printf("%-12s %12d %12d %12d %12d\n", "summary", summary.Files, summary.Binary, summary.Blank, summary.Code)
}

func FileWorker(files chan string, result chan CountByExt, progress *int64) {
	total := make(CountByExt)
	defer func() { result <- total }()

	for file := range files {
		count, err := CountLines(file)
		if err != nil {
			continue
		}
		total.Add(count)
		atomic.AddInt64(progress, 1)
	}
}

func IterateDir(work chan string, roots ...string) {
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			info.Mode()
			filename := info.Name()
			if len(filename) > 1 && filename[0] == '.' {
				return filepath.SkipDir
			}
			return nil
		}

		// if a filename contains ~ we assume it's a temporary file
		if strings.Contains(filepath.Base(path), "~") {
			return nil
		}

		if !ShouldExamine(filepath.Ext(path)) {
			return nil
		}

		work <- path
		return nil
	}

	for _, root := range roots {
		err := filepath.Walk(root, walk)
		if err != nil {
			log.Println(err)
		}
	}

	close(work)
}
