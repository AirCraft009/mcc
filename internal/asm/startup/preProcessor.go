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
	linkable.SetupDataGathering()
	linkF, nonNil := linkable.GetFile()
	for nonNil {
		if linkF.FileT != fileHandling.HeaderF {

			linkF, nonNil = linkable.GetFile()
			continue
		}

		lf := linkF
		wg.Add(1)
		go func(l *fileHandling.LinkFile) {
			defer wg.Done()
			findDefinitions(l, pre)
		}(lf)
		linkF, nonNil = linkable.GetFile()
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
	// parseDefinitions should already have been called
	linkable.SetupDataGathering()
	linkF, nonNil := linkable.GetFile()
	var wg = &sync.WaitGroup{}

	for nonNil {

		if linkF.FileT != fileHandling.AsmF {
			linkF, nonNil = linkable.GetFile()
			continue
		}

		lf := linkF
		wg.Add(1)
		go func(l *fileHandling.LinkFile) {
			defer wg.Done()
			replaceDefinitions(l, pre)
		}(lf)

		linkF, nonNil = linkable.GetFile()
	}
	wg.Wait()
}

func (pre *PreProcesser) Process(linkable *fileHandling.Linkables) {

	pre.parseDefinitions(linkable)
	pre.applyDefinitions(linkable)

}

func replaceDefinitions(linkF *fileHandling.LinkFile, pre *PreProcesser) {
	stringData := string(linkF.Data)
	for _, definition := range pre.definitions {
		stringData = strings.ReplaceAll(stringData, definition.placeHolder, definition.value)
	}
	linkF.Data = []byte(stringData)

}
