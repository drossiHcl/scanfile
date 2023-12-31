package main

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	pbscan "scanfile.com/scanfile/proto"
	types "scanfile.com/scanfile/types_pkg"
	util "scanfile.com/scanfile/util_pkg"

	"github.com/gabriel-vasile/mimetype"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var FLog *os.File = nil

var wg sync.WaitGroup

func main() {
	filesToProcess := false
	myInit()
	util.MyLog("Hello world I'm a Fs Scan ...\n")
	mins := 1
	for {
		t1 := time.NewTimer(time.Duration(types.Env_timer_fsscan) * time.Second)
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-t1.C

			// scan dir(s) named DATA_FILES_FOLDER
			dirList, nbrFiles := util.ListDir(types.Env_data_files_folder)
			util.MyLog("****** Files: %v NBR %d", dirList, nbrFiles)

			// range over files
			for _, fn := range dirList {

				util.MyLog("****** File : %s %s", fn, (types.Env_data_processed_textfiles_ita_dir + fn))
				// if destination already exixsts, skip it
				if util.FileExists((types.Env_data_processed_textfiles_ita_dir + fn)) {
					continue
				}
				if util.FileExists((types.Env_data_processed_textfiles_eng_dir + fn)) {
					continue
				}

				fullFileName := types.Env_data_files_folder + fn

				if f1, fTyp := isFileTypeHandled(fullFileName); f1 {
					var destination string

					lang, err := util.DetectLang(fullFileName)
					if err != nil {
						log.Fatal(err)
					}

					if lang == "Italian" {
						destination = (types.Env_data_textfiles_ita_dir + fn)
					} else if lang == "English" {
						destination = (types.Env_data_textfiles_eng_dir + fn)
					}
					util.MyLog("src %s dest %s lang %s", fullFileName, destination, lang)

					if fTyp == "ASCII" {
						_, err = util.CopyFile(fullFileName, destination)
					} else {
						err = util.ConvertAndCopyFile(fullFileName, destination, fTyp)
					}
					if err != nil {
						util.MyLog("Error converting or copying file %s %v\n\n", fullFileName, err)
						// log.Fatal(err)
					}
					if err == nil {
						filesToProcess = true
					}
				}
				util.MyLog("Timer 1 %d ", nbrFiles)
			}
		}()

		wg.Wait()
		util.MyLog("****** After Wait, filesToProcess %v ", filesToProcess)

		// Trigger backEnd via gRPC
		if filesToProcess {
			resultP, cconn, err := sendRequest()
			cconn.Close()
			if err != nil {
				log.Fatalf("could not send request to backEnd: %v", err)
			}
			getRes := resultP.GetTriggerRes()
			util.MyLog("End %d %d ", getRes, mins)
			filesToProcess = false
		}
		mins++
		// TODO Every hour (or more) clean dirs
	}
}

func myInit() {
	f, err := os.Create("/home/drossi/myTest/data/log_fsScan.log")
	FLog = f
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// defer f.Close()
	// wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(f)
	log.SetPrefix("FsScan ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	myEnv := os.Environ()
	for _, e := range myEnv {
		util.MyLog("%v\n", e)
	}

	_ = util.LoadEnv("/home/drossi/myTest/data/local.env")

	log.Println("********************** FsScan Started *********************")
}

func isFileTypeHandled(fName string) (bool, string) {
	fType, err := getFileType(fName)
	if err != nil {
		util.MyLog("Error while recognizing File Type, skip this file\n")
		return false, ""
	}
	util.MyLog("File %s Type %s\n", fName, fType)
	switch fType {
	case "ASCII":
		return true, fType
	// TODO
	case "PDF":
		return true, fType
	// TODO
	case "DOC":
		return true, fType
	case "XLS":
		return true, fType
	default:
		util.MyLog("Unhandled File Type, skip this file\n")
		return false, ""
	}
}

func getFileType(fName string) (string, error) {

	// TODO
	mtype, err := mimetype.DetectFile(fName)
	if err != nil {
		util.MyLog("%s", err)
		return "", err
	}
	util.MyLog("file type %s file extension %s", mtype.String(), mtype.Extension())
	out := mtype.String()
	// ----

	if strings.Contains(string(out[:]), "ASCII") {
		return "ASCII", nil
	}
	if strings.Contains(string(out[:]), "utf-8") {
		return "ASCII", nil
	}
	if strings.Contains(string(out[:]), "text") ||
		strings.Contains(string(out[:]), "csv") ||
		strings.Contains(string(out[:]), "txt") {
		return "ASCII", nil
	}
	if strings.Contains(string(out[:]), "pdf") {
		return "PDF", nil
	}
	if strings.Contains(string(out[:]), "docx") ||
		strings.Contains(string(out[:]), "word") {
		return "DOC", nil
	}
	if strings.Contains(string(out[:]), "xlsx") ||
		strings.Contains(string(out[:]), "sheet") {
		return "XLS", nil
	}
	if strings.Contains(string(out[:]), "directory") {
		return "DIR", nil
	}
	return "", err
}

func sendRequest() (*pbscan.TriggerBackendRes, *grpc.ClientConn, error) {
	opts := []grpc.DialOption{}

	tls := false
	if tls {
		certfile := "/home/drossi/myTest/provaGo/ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certfile, "")
		if err != nil {
			log.Fatalf("Error loading CA Certificates: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.Dial("mounttest:50051", opts...)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	c := pbscan.NewScanStatServiceClient(cc)
	util.MyLog("Created client.")
	util.MyLog("fsScan does a request to the Backend to process files RPC ...")

	req := &pbscan.TriggerBackend{
		Trigger: 1,
	}

	res, err := c.Process(context.Background(), req)
	// resStream, err := c.GetFilesStats(context.Background(), req)
	if err != nil {
		util.MyLog("error while calling process RPC: %v", err)
		return nil, cc, err
	}
	return res, cc, nil
}
