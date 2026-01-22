package pre_processor

import (
	"fmt"
	"mcc/Assembly-process/linker"
	"strings"
	"sync"
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
	fmt.Println("Parsing definitions")
	var wg = &sync.WaitGroup{}
	fmt.Println("setup waitgroup")
	linkable.SetupDataGathering()
	fmt.Println("setup Data")
	linkF, nonNil := linkable.GetFile()
	for nonNil {
		fmt.Println("got file linkF:", linkF)
		fmt.Println("check Header link", linkF.FileT)
		if linkF.FileT != linker.HeaderF {
			fmt.Println("is not header ")
			linkF, nonNil = linkable.GetFile()
			continue
		}

		fmt.Println("linkF-Headerfile:", linkF.FileT)
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
	fmt.Println("finding definitions-goroutine")
	for _, line := range strings.Split(string(linkF.Data), "\n") {
		fmt.Println(line)
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
		fmt.Println(parts)
		pre.definitions = append(pre.definitions, definition{parts[0], parts[1]})
		pre.mutex.Unlock()
	}
}

func (pre *PreProcesser) ApplyDefinitions(linkable *linker.Linkables) {
	// parseDefinitions should already have been called
	linkable.SetupDataGathering()
	linkF, nonNil := linkable.GetFile()
	var wg = &sync.WaitGroup{}
	fmt.Println("setup waitgroup")
	for nonNil {
		fmt.Println("got file linkF:", linkF)
		if linkF.FileT != linker.AsmF {
			linkF, nonNil = linkable.GetFile()
			continue
		}
		fmt.Println("got file linkF:", linkF.FileT)
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
	fmt.Println("started pre processor")
	pre.parseDefinitions(linkable)
	fmt.Println("finished parsing processor")
	pre.ApplyDefinitions(linkable)
	fmt.Println("finished applying definitions")
	fmt.Println(string(linkable.Files[0].Data))
}

func replaceDefinitions(linkF *linker.LinkFile, pre *PreProcesser) {
	stringData := string(linkF.Data)
	for _, definition := range pre.definitions {
		stringData = strings.ReplaceAll(stringData, definition.placeHolder, definition.value)
	}
	linkF.Data = []byte(stringData)
	fmt.Println(string(linkF.Data))
}
