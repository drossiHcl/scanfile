package types

import (
	"bufio"
	"os"

	"github.com/pemistahl/lingua-go"
)

var (
	Env_min_word_len         int
	Env_max_words            int
	Env_max_words_percentage int

	Env_timer_fsscan int

	Env_data_dascartare_dir              string // const DATA_DASCARTARE_DIR = "/home/drossi/myTest/data/daScartare/"
	Env_data_textfiles_name              string // const DATA_TEXTFILES_NAME = "/home/drossi/myTest/data/textFiles_"
	Env_data_files_folder                string // const DATA_FILES_FOLDER = "/home/drossi/myTest/data/textFiles_input/"
	Env_data_textfiles_ita_dir           string // const DATA_TEXTFILES_ITA_DIR = "/home/drossi/myTest/data/textFiles_ITA/"
	Env_data_textfiles_eng_dir           string // const DATA_TEXTFILES_ENG_DIR = "/home/drossi/myTest/data/textFiles_ENG/"
	Env_name_filedascartare              string // const NAME_FILEDASCARTARE = "daScartare_"
	Env_data_processed_textfiles_ita_dir string // const DATA_PROCESSED_TEXTFILES_ITA_DIR = "/home/drossi/myTest/data/textFiles_processed_ITA/"
	Env_data_processed_textfiles_eng_dir string // const DATA_PROCESSED_TEXTFILES_ENG_DIR = "/home/drossi/myTest/data/textFiles_processed_ENG/"
	Env_data_output_textfiles_ita_dir    string // const DATA_OUTPUT_TEXTFILES_ITA_DIR = "/home/drossi/myTest/data/textFiles_output_ITA/"
	Env_data_output_textfiles_eng_dir    string // const DATA_OUTPUT_TEXTFILES_ENG_DIR = "/home/drossi/myTest/data/textFiles_output_ENG/"
)

type FilesToRead_T struct {
	FName     string
	FFullName string
}

type WordListStatT struct {
	Word      string
	Frequency int
}
type WordsListsStatsT []WordListStatT

type WordListToSkipT struct {
	Word string
}

type BackendScanService interface {
	scanFile(*[]WordListStatT, *os.File, int, string) *bufio.Scanner
	purgeList(*WordsListsStatsT) *WordsListsStatsT
	sortList(*WordsListsStatsT) *WordsListsStatsT
}

// var WordsList []WordListStatT
var WordsListToSkip []WordListToSkipT
var WordsListToSkip_eng []WordListToSkipT
var WordsListToSkip_ita []WordListToSkipT

var FLog *os.File = nil

var Languages = []lingua.Language{
	lingua.English,
	lingua.Italian,
}
