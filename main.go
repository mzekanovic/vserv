package main

import (
    "flag"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "os/exec"
    "runtime"
)

func ResponseHandler(context map[string]interface{}) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {

        playerHTML, err := Asset("tmpl/player.html")
        if err != nil {
            log.Fatal(err)
            os.Exit(1)
        }
        t, err := template.New("player").Parse(string(playerHTML))
        if err != nil {
            log.Fatal(err)
            os.Exit(1)
        }

        log.Println("[Response] ", r.Method, " ", r.URL)
        t.Execute(w, context)

    }

}

func openBrowser(url string) {
    var err error

    log.Println("openin url ", url)

    switch runtime.GOOS {
    case "linux":
        err = exec.Command("xdg-open", url).Start()
    case "windows":
        err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
    case "darwin":
        err = exec.Command("open", url).Start()
    default:
        err = fmt.Errorf("unsupported platform")
    }
    if err != nil {
        log.Fatal(err)
    }
}

func runServer(videoFileName string, mimetype string) {

    context := map[string]interface{}{"src": videoFileName, "mimetype": mimetype}

    go openBrowser("http://localhost:8080/")

    http.HandleFunc("/", ResponseHandler(context))
    http.HandleFunc("/"+videoFileName, func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, videoFileName)
    })
    log.Println("running server for ", videoFileName, " of type ", mimetype)
    http.ListenAndServe(":8080", nil)
}

func main() {

    // parse arguments
    flag.Parse()

    args := flag.Args()

    if len(args) == 0 {
        log.Fatal("file needed")
        os.Exit(1)
    }

    // get filename
    videoFileName := args[0]

    // check if file exists
    if _, err := os.Stat(videoFileName); os.IsNotExist(err) {
        log.Fatal("file ", videoFileName, " doesn't exists")
        os.Exit(1)
    }

    // open file
    fd, err := os.Open(videoFileName)
    if err != nil {
        log.Fatal("error reading file ", videoFileName)
        os.Exit(1)
    }

    defer fd.Close()

    // read atmost first 512 bytes for file format detection
    byteSlice := make([]byte, 512)
    _, err = fd.Read(byteSlice)
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

    // detect file mime type
    fileType := http.DetectContentType(byteSlice)

    // runserver if file is in right mimetype
    switch fileType {
    case "video/mp4", "video/webm", "video/ogg", "video/x-theora+ogg", "application/ogg":
        runServer(videoFileName, fileType)
    default:
        log.Fatal(videoFileName, " has wrong format ", fileType)
    }

    os.Exit(0)
}
