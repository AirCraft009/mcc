package pre_processor

import (
	"fmt"
	"strings"
	"sync"

	"github.com/AirCraft009/mcc/internal/linker"
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

func (pre *PreProcesser) parseDefinitions(linkable *linker.Linkables) {
	var wg = &sync.WaitGroup{}
	linkable.SetupDataGathering()
	linkF, nonNil := linkable.GetFile()
	for nonNil {
		if linkF.FileT != linker.HeaderF {

			linkF, nonNil = linkable.GetFile()
			continue
		}

		lf := linkF
		wg.Add(1)
		go func(l *linker.LinkFile) {
			defer wg.Done()
			findDefinitions(l, pre)
		}(lf)
		linkF, nonNil = linkable.GetFile()
	}
	// only terminate after every goroutine is done
	wg.Wait()
}

func findDefinitions(linkF *linker.LinkFile, pre *PreProcesser) {
	for _, line := range strings.Split(string(linkF.Data), "\n") {
		if !strings.HasPrefix(line, "#define ") {
			continue
		}

		line = strings.TrimPrefix(line, "#define ")
		parts := strings.Split(strings.TrimSpace(line), ":")
		if len(parts) != 2 {
			fmt.Printf("Error parsing definition line:%s \nfile: %s \nreason: to many sperators ':'\n", line, linkF.Path)
			continue
		}
		pre.mutex.Lock()
		pre.definitions = append(pre.definitions, definition{parts[0], parts[1]})
		pre.mutex.Unlock()
	}
}

func (pre *PreProcesser) ApplyDefinitions(linkable *linker.Linkables) {
	// parseDefinitions should already have been called
	linkable.SetupDataGathering()
	linkF, nonNil := linkable.GetFile()
	var wg = &sync.WaitGroup{}

	for nonNil {

		if linkF.FileT != linker.AsmF {
			linkF, nonNil = linkable.GetFile()
			continue
		}

		lf := linkF
		wg.Add(1)
		go func(l *linker.LinkFile) {
			defer wg.Done()
			replaceDefinitions(l, pre)
		}(lf)

		linkF, nonNil = linkable.GetFile()
	}
	wg.Wait()
}

func (pre *PreProcesser) Process(linkable *linker.Linkables) {

	pre.parseDefinitions(linkable)
	pre.ApplyDefinitions(linkable)

}

func replaceDefinitions(linkF *linker.LinkFile, pre *PreProcesser) {
	stringData := string(linkF.Data)
	for _, definition := range pre.definitions {
		stringData = strings.ReplaceAll(stringData, definition.placeHolder, definition.value)
	}
	linkF.Data = []byte(stringData)

}
