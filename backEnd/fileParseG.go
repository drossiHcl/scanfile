package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"

	types "scanfile.com/scanfile/types_pkg"
	util "scanfile.com/scanfile/util_pkg"
)

var wg sync.WaitGroup

type WordsListsStats struct {
	types.WordsListsStatsT
}

func main() {

	if len(os.Args) != 1 {
		fmt.Printf("No parameters expected. Ignored \n")
	}

	myInit()

	if !(loadDaScartare("ENG")) {
		util.MyLog("Error: Unable to load file for words to skip ... \n\n")
		log.Fatalf("error Unable to load file for words to skip ENG.")
		return
	}
	//printListDs(WordsListToSkip_eng)
	types.WordsListToSkip = nil
	if !(loadDaScartare("ITA")) {
		util.MyLog("Error: Unable to load file for words to skip ... \n\n")
		log.Fatalf("error Unable to load file for words to skip ITA.")
		return
	}

	// Init and start gRPC server
	GrpcServer()
	//printListDs(WordsListToSkip_ita)
}

func myInit() {
	f, err := os.Create("/home/drossi/myTest/data/log_scanfile.log")
	types.FLog = f
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// defer f.Close()
	// wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(f)
	log.SetPrefix("BackEnd ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	_ = util.LoadEnv("/home/drossi/myTest/data/local.env")

	log.Println("********************** Backend Started *********************")
}

func processInputFiles(fileName string, fName string, lang string) {

	defer wg.Done()

	file, err := os.Open(fileName)
	if err != nil {
		util.MyLog("Error: Unable to load file to process ... \n\n")
		log.Fatal(err)
	}

	// Scan file word by word and append to a slice
	wordList := NewScanService()
	wordList.scanFile(&wordList, file, 0, lang)
	// printList(WordsList)

	// Delete duplicates from the slice
	wordList = wordList.purgeList(&wordList)

	// Sort the slice
	wordList = wordList.sortList(&wordList)

	// Move processed file to processed dir
	err = util.MoveFile(fName, lang, file)
	if err != nil {
		util.MyLog("Failed moving original file to processed dir: %s", err)
	}
	err = os.Remove(fileName)
	if err != nil {
		util.MyLog("Failed removing original file: %s", err)
	}

	file.Close()

	// Save Output to file in output dir
	util.SaveListToOutputFile(wordList.WordsListsStatsT, fName, lang)
	wordList.WordsListsStatsT = nil
}

func loadDaScartare(lang string) bool {

	fileNameDs := types.Env_data_dascartare_dir + types.Env_name_filedascartare
	if lang == "ENG" {
		fileNameDs += "ENG.txt"
	} else if lang == "ITA" {
		fileNameDs += "ITA.txt"
	} else {
		util.MyLog("\n\n Error: language not supported %s - %s \n\n\n", lang, fileNameDs)
		return false
	}

	file, err := os.Open(fileNameDs)
	if err != nil {
		util.MyLog("\n\n Error: could not open %s - %s \n\n\n", lang, fileNameDs)
		log.Fatal(err)
	}
	defer file.Close()

	// Scan file to be skipped word by word and append to a slice
	wordList := NewScanService()
	wordList.WordsListsStatsT = nil
	wordList.scanFile(&wordList, file, 1, lang)

	// Sort slice
	sort.Slice(types.WordsListToSkip, func(i, j int) bool {
		return types.WordsListToSkip[i].Word < types.WordsListToSkip[j].Word
	})
	// set proper slice
	if lang == "ENG" {
		types.WordsListToSkip_eng = types.WordsListToSkip
	} else {
		types.WordsListToSkip_ita = types.WordsListToSkip
	}
	return true
}

func (w *WordsListsStats) purgeList(wds *WordsListsStats) WordsListsStats {

	var WordsListCp WordsListsStats
	toCopy := 1
	idxCp := 0
	freq := 1
	j := 0

	for idx, w := range wds.WordsListsStatsT {

		for j = idx + 1; j < len(wds.WordsListsStatsT); j++ {

			if w.Frequency != -1 {
				if w.Word == (wds.WordsListsStatsT)[j].Word {
					// fmt.Printf("%s %s da canc \n", w.word, *wds[j].word)
					(wds.WordsListsStatsT)[j].Frequency = -1
					freq++
				}
			} else {
				toCopy = 0
			}
		}
		if toCopy == 1 {
			if (idx < len(wds.WordsListsStatsT)-1) || (wds.WordsListsStatsT)[len(wds.WordsListsStatsT)-1].Frequency != -1 {
				WordsListCp.WordsListsStatsT = append(WordsListCp.WordsListsStatsT, w)
				(WordsListCp.WordsListsStatsT)[idxCp].Frequency = freq
			}
			freq = 1
			idxCp++
		}
		toCopy = 1
	}
	wds.WordsListsStatsT = WordsListCp.WordsListsStatsT
	return WordsListCp
}

func (w *WordsListsStats) sortList(wds *WordsListsStats) WordsListsStats {
	sort.Slice(wds.WordsListsStatsT, func(i, j int) bool {
		return (wds.WordsListsStatsT)[i].Frequency > (wds.WordsListsStatsT)[j].Frequency
	})
	return *wds
}

func NewScanService() WordsListsStats {
	return WordsListsStats{}
}

func (w *WordsListsStats) scanFile(WordsList *WordsListsStats, f *os.File, Ds int, lang string) *bufio.Scanner {

	Scanner := bufio.NewScanner(f)
	Scanner.Split(bufio.ScanWords)

	for Scanner.Scan() {
		w := Scanner.Text()
		if Ds == 0 {
			w = util.ClearWord(w)
			if len(w) > types.Env_min_word_len && util.NotToBeSkipped(w, lang) {
				var Word types.WordListStatT
				Word.Frequency = 1
				Word.Word = w
				WordsList.WordsListsStatsT = append(WordsList.WordsListsStatsT, Word)
			}
		} else {
			if w[0] != '#' {
				var WordDs types.WordListToSkipT
				WordDs.Word = w
				types.WordsListToSkip = append(types.WordsListToSkip, WordDs)
			}
		}
	}

	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return Scanner
}
