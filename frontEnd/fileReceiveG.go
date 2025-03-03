package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"os"
	"strconv"

	pbscan "scanfile.com/scanfile/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Word struct {
	W string // word
	F uint32 // word frequency in a File
}

type File struct {
	PathName string
	Language string
	NumWords uint32
}

type FilesStats struct {
	File  File
	Words []Word
	List  []string
}

var results_list []FilesStats
var listResult []string
var tpl *template.Template

func init() {
	fmt.Println("Exec INIT ")
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {

	// ===
	f, err := os.Create(("/tmp/log_frontEnd.log"))
	if err != nil {
		log.Fatalf("error opening LOG file: %v", err)
	}

	log.SetOutput(f)
	log.SetPrefix("FrontEndEnd ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Println("********************** Frontend Started *********************")
	// ===

	log.Println("Hello I'm the Frontend v1.0 ...")

	http.Handle("/index/", basicAuth(http.HandlerFunc(index)))
	http.Handle("/reqdata", basicAuth(http.HandlerFunc(requestData)))

	httpPort := os.Getenv("HTTP_FRONTEND_PORT")
	log.Printf("ListenAndServe on http port %s\n", httpPort)
	err = http.ListenAndServe(":"+httpPort, nil)
	// err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Printf("ListenAndServe failed %v", err)
	}
	log.Printf("Out of ListenAndServe on port 8082")
}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	// Define a map to store multiple users and their passwords
	users := map[string]string{
		os.Getenv("USERNAME"):  os.Getenv("PASSWORD"),
		os.Getenv("USERNAME2"): os.Getenv("PASSWORD2"),
		os.Getenv("USERNAME3"): os.Getenv("PASSWORD3"),
	}

	// Check if any of the environment variables are not set
	for ux, px := range users {
		if ux == "" || px == "" {
			fmt.Println("One or more environment variables for users or passwords are not set")
			return nil
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		user_ok := false
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		for u, p := range users {
			if u == user && p != pass {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			} else if u == user {
				user_ok = true
			}
		}
		if !user_ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	log.Println("Index Handler")
	err1 := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	HandleError(w, err1)
}

func requestData(w http.ResponseWriter, req *http.Request) {
	results_list = nil
	listResult = nil
	log.Println("Request data Handler")
	err := req.ParseForm()
	HandleError(w, err)
	fName := req.FormValue("NameFile")
	lang := req.FormValue("Language")
	maxn, _ := strconv.Atoi(req.FormValue("Maxnum"))
	listOnly := req.FormValue("listonly")
	log.Println(fName)
	log.Println(lang)
	log.Println(maxn)
	log.Println(listOnly)

	// Send request to backEnd
	resultsStream, cconn, err := sendRequest(fName, lang, uint32(maxn), listOnly)
	defer cconn.Close()
	if err != nil {
		log.Fatalf("could not send request to backEnd: %v", err)
	}

	// Receive response from backEnd
	err = receiveResponse(resultsStream)
	if err != nil {
		log.Fatalf("could not receive request to backEnd: %v", err)
	}

	// Send data to browser
	if listResult == nil {
		log.Println("*** *** RES *** ***", len(results_list))
		for _, f := range results_list {
			_ = tpl.ExecuteTemplate(w, "something.gohtml", f.File.PathName)
			_ = tpl.ExecuteTemplate(w, "something.gohtml", f.File.Language)
			_ = tpl.ExecuteTemplate(w, "something.gohtml", f.File.NumWords)
			for _, wrd := range f.Words {
				toHtml := fmt.Sprintf("    word %s - occourencies %d\n", wrd.W, wrd.F)
				_ = tpl.ExecuteTemplate(w, "someword.gohtml", toHtml)
			}
		}
	} else {
		for _, n := range listResult {
			_ = tpl.ExecuteTemplate(w, "something.gohtml", n)
		}
	}
	err = tpl.ExecuteTemplate(w, "buttonBackToIndex.gohtml", nil)
	HandleError(w, err)
}

func sendRequest(fileName string, language string, maxNumber uint32, list string) (pbscan.ScanStatService_GetFilesStatsClient, *grpc.ClientConn, error) {
	opts := []grpc.DialOption{}
	tls := false
	if tls {
		certfile := "../../ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certfile, "")
		if err != nil {
			log.Fatalf("Error loading CA Certificates: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	grpcServerPort := os.Getenv("GRPC_SERVER_PORT")
	log.Printf("Dial on grpc port %s\n", grpcServerPort)
	cc, err := grpc.Dial("my-backend-test:"+grpcServerPort, opts...)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	c := pbscan.NewScanStatServiceClient(cc)
	log.Println("Created client.")
	log.Println("Frontend does a request to the Backend to receive Streaming RPC ...")

	req := &pbscan.FsRequest{
		Lang: language,
		Num:  maxNumber,
		Name: fileName,
		List: list,
	}

	resStream, err := c.GetFilesStats(context.Background(), req)
	if err != nil {
		log.Printf("error while calling GetFilesStats RPC: %v", err)
		return nil, cc, err
	}
	return resStream, cc, nil
}

func receiveResponse(resultsStream pbscan.ScanStatService_GetFilesStatsClient) error {
	var results FilesStats

	var wrd Word
	for {
		msg, err := resultsStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		if msg.List == nil {
			// fmt.Printf("Response from Backend: %s\n", msg.File.GetPathName())
			results.File.PathName = msg.File.GetPathName()
			results.File.Language = msg.File.GetLanguage()
			results.File.NumWords = msg.File.GetNumWords()
			for _, wd := range msg.Words {
				// fmt.Printf("Response from Backend: %d - %s %d\n", ii, wd.GetW(), wd.GetF())
				wrd.W = wd.GetW()
				wrd.F = wd.GetF()
				results.Words = append(results.Words, wrd)
			}
			// fmt.Printf("\n")
			results_list = append(results_list, results)
			results.Words = nil
			listResult = nil
		} else {
			listResult = append(listResult, msg.List...)
		}

	}
	return nil
}

func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}
