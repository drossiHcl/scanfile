package types

import (
	"bufio"
	"os"
	"sync"

	"github.com/pemistahl/lingua-go"
)

var (
	Env_min_word_len         int
	Env_max_words            int
	Env_max_words_percentage int

	Env_timer_fsscan    int
	Env_maxsize_logfile int64

	Env_log_dir string // const LOG_DIR = "/home/daniele/Daniele/scanfile/log/"

	Env_data_dascartare_dir string // const DATA_DASCARTARE_DIR = "/home/daniele/Daniele/myData/scanfile/daScartare/"
	Env_data_textfiles_name string // const DATA_TEXTFILES_NAME = "/home/daniele/Daniele/myData/scanfile/textFiles_"
	Env_data_files_folder   string // const DATA_FILES_FOLDER = "/home/daniele/Daniele/myData/scanfile/textFiles_input/"

	Env_data_textfiles_ita_dir           string // const DATA_TEXTFILES_ITA_DIR = "/home/daniele/Daniele/myData/scanfile/textFiles_ITA/"
	Env_data_textfiles_eng_dir           string // const DATA_TEXTFILES_ENG_DIR = "/home/daniele/Daniele/myData/scanfile/textFiles_ENG/"
	Env_name_filedascartare              string // const NAME_FILEDASCARTARE = "daScartare_"
	Env_data_processed_textfiles_ita_dir string // const DATA_PROCESSED_TEXTFILES_ITA_DIR = "/home/daniele/Daniele/myData/scanfile/textFiles_processed_ITA/"
	Env_data_processed_textfiles_eng_dir string // const DATA_PROCESSED_TEXTFILES_ENG_DIR = "/home/daniele/Daniele/myData/scanfile/textFiles_processed_ENG/"
	Env_data_output_textfiles_ita_dir    string // const DATA_OUTPUT_TEXTFILES_ITA_DIR = "/home/daniele/Daniele/myData/scanfile/textFiles_output_ITA/"
	Env_data_output_textfiles_eng_dir    string // const DATA_OUTPUT_TEXTFILES_ENG_DIR = "/home/daniele/Daniele/myData/scanfile/textFiles_output_ENG/"
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

var BackendFLog *os.File = nil
var FsscanFLog *os.File = nil
var LockLogs sync.Mutex

var Languages = []lingua.Language{
	lingua.English,
	lingua.Italian,
}
