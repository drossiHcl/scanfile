diff --git a/fsScan/fsScan.go b/fsScan/fsScan.go
index 8c9d566..4171cbe 100644
--- a/fsScan/fsScan.go
+++ b/fsScan/fsScan.go
@@ -18,43 +18,44 @@ import (
 	"google.golang.org/grpc/credentials/insecure"
 )
 
-var FLog *os.File = nil
+// var FLog *os.File = nil
 
 var wg sync.WaitGroup
 
 func main() {
 	filesToProcess := false
 	myInit()
+
 	util.MyLog("Hello world I'm a Fs Scan ...\n")
-	mins := 1
+	mins := 0
 	for {
 
 		t1 := time.NewTimer(time.Duration(types.Env_timer_fsscan) * time.Second)
+		<-t1.C
 
-		wg.Add(1)
-
-		go func() {
+		// scan dir(s) named DATA_FILES_FOLDER
+		dirList, nbrFiles := util.ListDir(types.Env_data_files_folder)
+		util.MyLog("****** Files: %v NBR %d", dirList, nbrFiles)
 
-			defer wg.Done()
-			<-t1.C
+		// range over files
+		for _, fname := range dirList {
 
-			// scan dir(s) named DATA_FILES_FOLDER
-			dirList, nbrFiles := util.ListDir(types.Env_data_files_folder)
-			util.MyLog("****** Files: %v NBR %d", dirList, nbrFiles)
+			wg.Add(1)
+			// fname := fn
+			go func(fname string, fToProcess bool) error {
 
-			// range over files
-			for _, fn := range dirList {
+				defer wg.Done()
 
-				util.MyLog("****** File : %s %s", fn, (types.Env_data_processed_textfiles_ita_dir + fn))
+				util.MyLog("****** File : %s", fname)
 				// if destination already exixsts, skip it
-				if util.FileExists((types.Env_data_processed_textfiles_ita_dir + fn)) {
-					continue
+				if util.FileExists((types.Env_data_processed_textfiles_ita_dir + fname)) {
+					return nil
 				}
-				if util.FileExists((types.Env_data_processed_textfiles_eng_dir + fn)) {
-					continue
+				if util.FileExists((types.Env_data_processed_textfiles_eng_dir + fname)) {
+					return nil
 				}
 
-				fullFileName := types.Env_data_files_folder + fn
+				fullFileName := types.Env_data_files_folder + fname
 
 				if f1, fTyp := isFileTypeHandled(fullFileName); f1 {
 					var destination string
@@ -65,9 +66,9 @@ func main() {
 					}
 
 					if lang == "Italian" {
-						destination = (types.Env_data_textfiles_ita_dir + fn)
+						destination = (types.Env_data_textfiles_ita_dir + fname)
 					} else if lang == "English" {
-						destination = (types.Env_data_textfiles_eng_dir + fn)
+						destination = (types.Env_data_textfiles_eng_dir + fname)
 					}
 					util.MyLog("src %s dest %s lang %s", fullFileName, destination, lang)
 
@@ -81,13 +82,14 @@ func main() {
 						// log.Fatal(err)
 					}
 					if err == nil {
-						filesToProcess = true
+						fToProcess = true
 					}
 				}
 				util.MyLog("Timer 1 %d ", nbrFiles)
-			}
-		}()
+				return nil
 
+			}(fname, filesToProcess)
+		}
 		wg.Wait()
 
 		util.MyLog("****** After Wait, filesToProcess %v ", filesToProcess)
@@ -103,8 +105,20 @@ func main() {
 			util.MyLog("End %d %d ", getRes, mins)
 			filesToProcess = false
 		}
-		mins++
+
 		// TODO Every hour (or more) clean dirs
+
+		mins++
+		util.MyLog("Current Mins %d ", mins)
+		// Periodically check the log file size
+		if mins == 30 {
+			util.MyLog("Check Log Files %d ", mins)
+			rotateErr := util.CheckRotateLogFile(types.FsscanFLog, "fsScan")
+			if rotateErr != nil {
+				util.MyLog("error rotating func: %v", rotateErr)
+			}
+			mins = 0
+		}
 	}
 }
 
@@ -118,12 +132,12 @@ func myInit() {
 	}
 
 	f, err = os.Create((types.Env_log_dir + "log_fsScan.log"))
-	FLog = f
 	if err != nil {
 		log.Fatalf("error opening LOG file: %v", err)
 	}
+	types.FsscanFLog = f
 
-	log.SetOutput(f)
+	log.SetOutput(types.FsscanFLog)
 	log.SetPrefix("FsScan ")
 	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
 
