package utilType

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"code.sajari.com/docconv"
	"github.com/joho/godotenv"
	"github.com/pemistahl/lingua-go"

	types "scanfile.com/scanfile/types_pkg"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func MyLog(format string, args ...interface{}) {
	l := fmt.Sprintf(format, args...)
	log.Println(l)
}

func ClearWord(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func NotToBeSkipped(str string, lang string) bool {
	if lang == "ENG" {
		for _, s := range types.WordsListToSkip_eng {
			if strings.EqualFold(s.Word, str) {
				return false
			}
		}
	} else {
		for _, s := range types.WordsListToSkip_ita {
			if strings.EqualFold(s.Word, str) {
				return false
			}
		}
	}
	return true
}

func SaveListToOutputFile(wds []types.WordListStatT, fName string, lang string) {
	var fileNameOut string

	if lang == "ITA" {
		fileNameOut = types.Env_data_output_textfiles_ita_dir
	} else {
		fileNameOut = types.Env_data_output_textfiles_eng_dir
	}
	fileNameOut += fName
	// open output file
	fout, err := os.Create(fileNameOut)
	if err != nil {
		MyLog("ERROR: cannot open output file\n")
	}
	defer fout.Close()

	w := bufio.NewWriter(fout)
	_, err = fmt.Fprintf(w, "%v\n", len(wds))
	if err != nil {
		MyLog("ERROR: cannot write output file\n")
	}
	for _, word := range wds {
		_, err = fmt.Fprintf(w, "%v\n", word)
		if err != nil {
			MyLog("ERROR: cannot write output file\n")
			break
		}
	}
	w.Flush()
}

func PrintList(wds []types.WordListStatT) {
	fmt.Println("Len of the List: ", len(wds))
	for n, word := range wds {
		fmt.Printf("%d %v ", n, word)
	}
	fmt.Printf("\n")
}
func PrintListDs(wds []types.WordListToSkipT) {
	fmt.Println("Len of the List to skip: ", len(types.WordsListToSkip))
	for n, word := range wds {
		fmt.Printf("%d %v ", n, word)
	}
	fmt.Printf("\n")
}

func GetLang(fn string) string {
	if strings.Contains(fn, "ENG") {
		return "ENG"
	} else {
		return "ITA"
	}
}

func ListDir(dir string) ([]string, int) {
	var files []string

	files = nil

	fdir, err := os.Open(dir)
	if err != nil {
		MyLog("\nError: can't open directory %s\n\n", dir)
		return files, 0
	}

	fileInfo, e := fdir.Readdir(-1)
	defer fdir.Close()
	if e != nil {
		MyLog("\nError: reading directory %s\n", dir)
		return files, 0
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, len(files)
}

func MoveFile(sourcePath, lang string, fp *os.File) error {
	var destPath string

	if lang == "ITA" {
		destPath = types.Env_data_processed_textfiles_ita_dir + sourcePath
	} else {
		destPath = types.Env_data_processed_textfiles_eng_dir + sourcePath
	}
	MyLog("Move %s to %s\n", sourcePath, destPath)

	outputFile, err := os.Create(destPath)
	if err != nil {
		MyLog("Creating failed\n")
		return fmt.Errorf("moving to processed dir failed: %s", err)
	}
	defer outputFile.Close()

	// Set the source file pointer to the beginning of the file since it is at the end
	_, err = fp.Seek(0, io.SeekStart)
	if err != nil {
		MyLog("Moving to processed dir failed\n")
		return fmt.Errorf("moving to processed dir failed: %s", err)
	}
	_, err = io.Copy(outputFile, fp)
	if err != nil {
		MyLog("Moving to processed dir failed\n")
		return fmt.Errorf("moving to processed dir failed: %s", err)
	}

	return nil
}

func FileExists(fName string) bool {
	dst_file, err := os.Open(fName)
	if err == nil {
		// destination file already exists
		dst_file.Close()
		return true
	}
	return false
}

func CopyFile(src string, dst string) (int64, error) {

	src_file, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer src_file.Close()

	src_file_stat, err := src_file.Stat()
	if err != nil {
		return 0, err
	}

	if !src_file_stat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	dst_file, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dst_file.Close()
	return io.Copy(dst_file, src_file)
}

func ConvertAndCopyFile(fullFileName string, destination string, fileType string) error {

	switch fileType {
	case "ASCII":
		return nil
	case "PDF", "DOC", "XLS":
		conv, err := ConvertPdf(fullFileName)
		if err == nil {
			fout, err := os.Create(destination)
			if err != nil {
				MyLog("ERROR: cannot open output file\n")
			}
			defer fout.Close()

			w := bufio.NewWriter(fout)
			_, err = fmt.Fprintf(w, "%s\n", conv.Body)
			if err != nil {
				MyLog("ERROR: cannot write output file\n")
			}
			w.Flush()
		}
		return nil
	// case "DOC":
	// 	return nil
	// case "XLS":
	// 	return nil
	default:
		return nil
	}
}

func GetFileContentType(fToCheck *os.File) (string, error) {

	// to sniff the content type only the first
	// 512 bytes are used.

	buf := make([]byte, 512)

	_, err := fToCheck.Read(buf)

	if err != nil {
		return "", err
	}

	// the function that actually does the trick
	contentType := http.DetectContentType(buf)

	return contentType, nil
}

func DetectLang(fileName string) (string, error) {

	detector := lingua.NewLanguageDetectorBuilder().FromLanguages(types.Languages...).Build()
	var chars string
	var numChars uint32

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanRunes)
	for Scanner.Scan() {
		chars = chars + Scanner.Text()
		if numChars++; numChars >= 100 {
			break
		}
	}
	if Scanner.Err() != nil {
		log.Fatal(Scanner.Err())
	}

	if language, exists := detector.DetectLanguageOf(chars); exists {
		/*
		 ** Returns: English, Italian
		 */
		return language.String(), nil
	}

	return "LANGNOTSUPPORTED", nil
}

func ConvertPdf(fName string) (*docconv.Response, error) {
	res, err := docconv.ConvertPath(fName)
	if err != nil {
		MyLog("\nError converting file %s %v\n", fName, err)
	}
	return res, err
}

func LoadEnv(fileEnvName string) error {
	err := godotenv.Load(fileEnvName)
	if err != nil {
		log.Fatalf("Some error occured in reading local env. Err: %s", err)
	}

	types.Env_min_word_len, err = strconv.Atoi(os.Getenv("MIN_WORD_LEN"))
	if err != nil {
		fmt.Printf("Error in Atoi 1 loading env %d \n", err)
	}
	types.Env_max_words, err = strconv.Atoi(os.Getenv("MAX_WORDS"))
	if err != nil {
		fmt.Printf("Error in Atoi 2 loading env %d \n", err)
	}
	types.Env_max_words_percentage, err = strconv.Atoi(os.Getenv("MAX_WORDS_PERCENTAGE"))
	if err != nil {
		fmt.Printf("Error in Atoi 3 loading env %d \n", err)
	}

	types.Env_timer_fsscan, err = strconv.Atoi(os.Getenv("TIMER_FSSCAN"))
	if err != nil {
		fmt.Printf("Error in Atoi 4 loading env %d \n", err)
	}
	fmt.Printf("TIMER FSSCAN loading env %d \n", types.Env_timer_fsscan)

	types.Env_log_dir = os.Getenv("LOG_DIR")
	types.Env_data_dascartare_dir = os.Getenv("DATA_DASCARTARE_DIR")
	types.Env_data_textfiles_name = os.Getenv("DATA_TEXTFILES_NAME")
	types.Env_data_files_folder = os.Getenv("DATA_FILES_FOLDER")
	types.Env_data_textfiles_ita_dir = os.Getenv("DATA_TEXTFILES_ITA_DIR")
	types.Env_data_textfiles_eng_dir = os.Getenv("DATA_TEXTFILES_ENG_DIR")
	types.Env_name_filedascartare = os.Getenv("NAME_FILEDASCARTARE")
	types.Env_data_processed_textfiles_ita_dir = os.Getenv("DATA_PROCESSED_TEXTFILES_ITA_DIR")
	types.Env_data_processed_textfiles_eng_dir = os.Getenv("DATA_PROCESSED_TEXTFILES_ENG_DIR")
	types.Env_data_output_textfiles_ita_dir = os.Getenv("DATA_OUTPUT_TEXTFILES_ITA_DIR")
	types.Env_data_output_textfiles_eng_dir = os.Getenv("DATA_OUTPUT_TEXTFILES_ENG_DIR")
	return err
}
