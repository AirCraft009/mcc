package startup

import (
	"fmt"
	"strings"
	"sync"

	"github.com/AirCraft009/mcc/internal/fileHandling"
)

type PreProcesser struct {
	definitions []definition
	mutex       sync.Mutex
}

type definition struct {
	placeHolder string
	value       string
}

func NewPreProcesser() *PreProcesser {
	return &PreProcesser{[]definition{}, sync.Mutex{}}
}

func (pre *PreProcesser) parseDefinitions(linkable *fileHandling.Linkables) {
	var wg = &sync.WaitGroup{}

	for _, linkF := range linkable.GetFiles() {
		if linkF.FileT != fileHandling.HeaderF {

			continue
		}

		lf := linkF
		wg.Add(1)
		go func(l *fileHandling.LinkFile) {
			defer wg.Done()
			findDefinitions(l, pre)
		}(lf)

	}
	// only terminate after every goroutine is done
	wg.Wait()
}

func findDefinitions(linkF *fileHandling.LinkFile, pre *PreProcesser) {
	for _, line := range strings.Split(string(linkF.Data), "\n") {
		if !strings.HasPrefix(line, "#define ") {
			continue
		}

		line = strings.TrimPrefix(line, "#define ")
		parts := strings.Split(strings.TrimSpace(line), ":")
		if len(parts) != 2 {
			fmt.Printf("MCC-Preprocessor Warning: Error parsing definition line:%s \nfile: %s \nreason: to many sperators ':'\n", line, linkF.Path)
			continue
		}
		pre.mutex.Lock()
		pre.definitions = append(pre.definitions, definition{parts[0], parts[1]})
		pre.mutex.Unlock()
	}
}

func (pre *PreProcesser) applyDefinitions(linkable *fileHandling.Linkables) {

	var wg = &sync.WaitGroup{}

	for _, linkF := range linkable.GetFiles() {

		if linkF.FileT != fileHandling.AsmF {
			continue
		}

		lf := linkF
		wg.Add(1)
		go func(l *fileHandling.LinkFile) {
			defer wg.Done()
			replaceDefinitions(l, pre)
		}(lf)

	}
	wg.Wait()
}

func (pre *PreProcesser) Process(linkable *fileHandling.Linkables) {

	pre.parseDefinitions(linkable)
	pre.applyDefinitions(linkable)
	pre.removeComments(linkable)
}

func (pre *PreProcesser) removeComments(linkable *fileHandling.Linkables) {
	var wg = &sync.WaitGroup{}
	for _, linkF := range linkable.GetFiles() {
		// remove # lines even for header files (just makes copy size smaller)
		if linkF.FileT == fileHandling.ObjectF {
			continue
		}

		lfData := &linkF.Data
		wg.Add(1)
		go removeCommentsFromFile(lfData, wg)
	}
	wg.Wait()
}

func removeCommentsFromFile(byteData *[]byte, group *sync.WaitGroup) {
	defer group.Done()
	data := string(*byteData)
	builder := &strings.Builder{}
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		commentIndex := strings.Index(line, "#")
		if commentIndex != -1 {
			line = line[:commentIndex]
		}
		if line != "" {
			builder.WriteString(line)
			builder.WriteByte('\n')
		}
	}
	*byteData = []byte(builder.String())
}

func replaceDefinitions(linkF *fileHandling.LinkFile, pre *PreProcesser) {
	stringData := string(linkF.Data)
	for _, definition := range pre.definitions {
		stringData = strings.ReplaceAll(stringData, definition.placeHolder, definition.value)
	}
	linkF.Data = []byte(stringData)

}
