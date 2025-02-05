package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	pbscan "scanfile.com/scanfile/proto"
	types "scanfile.com/scanfile/types_pkg"
	util "scanfile.com/scanfile/util_pkg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// type server struct{}
type server struct {
	pbscan.ScanStatServiceServer
}

func GrpcServer() {
	util.MyLog("Hello world I'm a Backend ...\n")

	grpcServerPort := os.Getenv("GRPC_SERVER_PORT")

	util.MyLog("... listening on grpcServerPort\n")
	lis, err := net.Listen("tcp", ":"+grpcServerPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	tls := false
	if tls {
		certfile := "/home/drossi/myTest/provaGo/ssl/server.crt"
		keyfile := "/home/drossi/myTest/provaGo/ssl/server.pem"
		creds, err := credentials.NewServerTLSFromFile(certfile, keyfile)
		if err != nil {
			log.Fatalf("Error loading certificates: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	pbscan.RegisterScanStatServiceServer(s, &server{})
	reflection.Register(s)

	util.MyLog("Call Serve ...\n")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	util.MyLog("Exited Serve ...\n")
}

func (s *server) Process(ctx context.Context, inTrig *pbscan.TriggerBackend) (*pbscan.TriggerBackendRes, error) {
	var fileName string
	var filesToRead []types.FilesToRead_T
	var fToRead types.FilesToRead_T
	filesToRead = nil

	/*
	 ** Get files to process if any
	 */
	if inTrig.Trigger == 1 {
		filesToRead_eng, nbrFiles_eng := util.ListDir(types.Env_data_textfiles_eng_dir)
		filesToRead_ita, nbrFiles_ita := util.ListDir(types.Env_data_textfiles_ita_dir)
		util.MyLog("eng: %d ita: %d\n", nbrFiles_eng, nbrFiles_ita)
		for _, fn := range filesToRead_eng {
			fileName = ""
			fileName += types.Env_data_textfiles_name + "ENG" + "/" + fn
			fToRead.FName = fn
			fToRead.FFullName = fileName
			filesToRead = append(filesToRead, fToRead)
		}
		for _, fn := range filesToRead_ita {
			fileName = ""
			fileName += types.Env_data_textfiles_name + "ITA" + "/" + fn
			fToRead.FName = fn
			fToRead.FFullName = fileName
			filesToRead = append(filesToRead, fToRead)
		}

		/*
		 ** Process files if any
		 */
		begin := time.Now()
		for _, fName := range filesToRead {
			util.MyLog("===> processing file %s\n", fName.FFullName)
			lang := util.GetLang(fName.FFullName)
			wg.Add(1)
			go processInputFiles(fName.FFullName, fName.FName, lang)
		}

		wg.Wait()

		end := time.Now()
		elapsed := end.Sub(begin)
		util.MyLog("Time elapsed %v\n", elapsed)

		filesToRead_eng = nil
		filesToRead_ita = nil
	}

	res := pbscan.TriggerBackendRes{
		TriggerRes: 1,
	}
	return &res, nil
}

func (s *server) GetFilesStats(inReq *pbscan.FsRequest, stream pbscan.ScanStatService_GetFilesStatsServer) error {
	var filesToRead []types.FilesToRead_T
	var fToRead types.FilesToRead_T
	var files_ENG, files_ITA []string
	var fileName string
	var nfe, nfi int
	filesToRead = nil

	util.MyLog("GetFilesStats server func invoked with parameters %s %s %d %s ...\n", inReq.Name, inReq.Lang, inReq.Num, inReq.List)

	if inReq.Lang == "" || inReq.Lang == "ENG" {
		files_ENG, nfe = util.ListDir(types.Env_data_output_textfiles_eng_dir)
	}
	if inReq.Lang == "" || inReq.Lang == "ITA" {
		files_ITA, nfi = util.ListDir(types.Env_data_output_textfiles_ita_dir)
	}

	if inReq.List == "listonly" {
		res := &pbscan.FilesStats{
			List: nil,
		}
		if inReq.Lang == "ENG" || inReq.Lang == "" {
			res.List = append(res.List, files_ENG...)
		}
		if inReq.Lang == "ITA" || inReq.Lang == "" {
			res.List = append(res.List, files_ITA...)
		}
		util.MyLog("*** Files list %s ***\n", res)
		stream.Send(res)
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	for _, fn := range files_ENG {
		fileName = ""
		fileName += types.Env_data_textfiles_name + "output_ENG" + "/" + fn
		fToRead.FName = fn
		fToRead.FFullName = fileName
		if (inReq.Name == "") || (strings.Contains(fToRead.FFullName, inReq.Name)) {
			filesToRead = append(filesToRead, fToRead)
		}
	}
	for _, fn := range files_ITA {
		fileName = ""
		fileName += types.Env_data_textfiles_name + "output_ITA" + "/" + fn
		fToRead.FName = fn
		fToRead.FFullName = fileName
		if (inReq.Name == "") || (strings.Contains(fToRead.FFullName, inReq.Name)) {
			filesToRead = append(filesToRead, fToRead)
		}
	}
	util.MyLog("Num files eng %d Num files ita %d ...\n", nfe, nfi)
	util.MyLog("\nFiles to READ %v ...\n", filesToRead)

	for _, fName := range filesToRead {
		util.MyLog("===> gRPC server sending file %s\n", fName.FFullName)

		res, err := prepareStreamForOneFile(fName.FFullName, inReq.Num)
		if err != nil {
			util.MyLog("gRPC server: ERROR preparing file %s\n", fName.FFullName)
		}
		stream.Send(res)
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func prepareStreamForOneFile(fileName string, maxToSend uint32) (*pbscan.FilesStats, error) {
	lang := util.GetLang(fileName)
	res := &pbscan.FilesStats{
		File: &pbscan.File{
			Language: lang,
			PathName: fileName,
		},
		List: nil,
	}
	file, err := os.Open(fileName)
	if err != nil {
		util.MyLog("gRPC Server, ERROR while opening file %v. Err=%v\n", fileName, err)
		return nil, err
	}
	defer file.Close()
	r := bufio.NewReader(file)
	numWords, err2 := r.ReadString(10) // 0x0A separator = newline
	if err2 != nil {
		util.MyLog("gRPC Server, ERROR while reading file %v . Err=%v\n", fileName, err2)
		return nil, err
	}
	nw, _ := strconv.Atoi(numWords[:(len(numWords) - 1)])
	res.File.NumWords = uint32(nw)
	nToSend := numWordsToSend(nw)
	util.MyLog("Words to Send %v ...\n", nToSend)
	for i := 0; i < nToSend && i < nw && (maxToSend == 0 || i < int(maxToSend)); i++ {
		line, err := r.ReadString(10) // 0x0A separator = newline
		if err == io.EOF {
			break
		} else if err != nil {
			util.MyLog("gRPC Server, ERROR while reading file %v . Err=%v\n", fileName, err)
			return nil, err
		}
		w, n, found := strings.Cut(line[1:(len(line)-1)], " ")
		f, _ := strconv.Atoi(n[:(len(n) - 1)])
		ww := &pbscan.Word{W: w, F: uint32(f)}
		if found {
			res.Words = append(res.Words, ww)
		}
	}
	util.MyLog("... gRPC server sending response %v ...\n", res)
	return res, nil
}

func numWordsToSend(nw int) int {
	if nw > 0 {
		if nw <= 10 {
			return nw
		}
		if nw <= 20 {
			return 15
		}
		if nw <= 50 {
			return 25
		}
		if nw <= 100 {
			return 40
		}
		if val := (nw / 100) * types.Env_max_words_percentage; val <= types.Env_max_words {
			return val
		} else {
			return types.Env_max_words
		}
	}
	return 0
}
