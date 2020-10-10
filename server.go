package main

  import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "encoding/json"
    "strconv"
    "os"
    "time"
  )

  type list struct {
    Name  []string
    Path []string
    Extension  []string
    Size [] int64
    RequestCount int64
  }

  type response struct {
    CountFileReceived   int64  `json:"no_of_file_received"`
    MaxFileSize int64 `json:"max_file_size"`
    MaxFileSizePath string `json:"max_file_size_path"`
    AvgFileSize int64 `json:"avg_file_size"`
    FileExtensions []string `json:"file_extension_list"`
    LatestFileList []string `json:"latest_file_list"`
  }

  var m list


  func httpServer(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/set_stats" && r.URL.Path != "/get_stats" {
      http.NotFound(w, r)
      return
    }
    switch r.Method {
    
    case "GET":
      for k, v := range r.URL.Query() {
        fmt.Printf("%s: %s\n", k, v)
      }

      var result int64
      
      var fpath string

      var path_length int

      var max int64

      var avg_size int64

      var file_path_index int      
      
      path_length = 10
      fpath = m.Path[0]

      
      if len(m.Size) > 0 {
        max = m.Size[0]
      }      
      
  
      for key, value := range m.Size {
        if value > max {
          max = value
          fpath = m.Path[key]
        }  
        result += value  
      }

      if m.RequestCount > 0 {   
        avg_size = result/m.RequestCount
      }  

      if len(m.Path) < 10 {
        file_path_index = (len(m.Path)-1)
        if len(m.Path) == 0 {
          file_path_index = 0
        }  
      }else {
        file_path_index = (len(m.Path)-1) - path_length
      }

      
      fmt.Printf("Number  of  files received: %d\n", m.RequestCount)
      fmt.Printf("Maximum file  size  (including  file  path): %d and %s\n", max, fpath)
      fmt.Printf("Average file  size: %d\n", avg_size)
      fmt.Printf("List  of  file  extensions: %+v\n", m.Extension)
      fmt.Printf("List  of  latest  10  file  paths received: %+v\n", m.Path[file_path_index:])

      m := response{m.RequestCount, max, fpath, avg_size, m.Extension, m.Path[file_path_index:]}

      jData, err := json.Marshal(m)
      if err != nil {
        panic(err)
      }

      w.Header().Set("Content-Type", "application/json")
      w.Write(jData)

      // w.Write([]byte("Received a GET request\n"))
 		
 		case "POST":
      reqBody, err := ioutil.ReadAll(r.Body)
      if err != nil {
        log.Fatal(err)
      }
      fmt.Printf("%s\n", reqBody)

      var data map[string]interface{}


      derr := json.Unmarshal([]byte(reqBody), &data)
      if derr != nil {
        panic(derr)
      }
         

       m.Name = append(m.Name, data["Name"].(string))
       m.Path = append(m.Path, data["Path"].(string))
       m.Extension = append(m.Extension, data["Extension"].(string))
       size, _ := strconv.ParseInt(data["Size"].(string), 10, 64)
       m.Size = append(m.Size, size)
       m.RequestCount += 1
       fmt.Println(m)

     

      w.Write([]byte("Received a POST request\n"))
  	default:
      w.WriteHeader(http.StatusNotImplemented)
      w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
    }
  }

  func getStatistics(port string) {
    address := "http://localhost:" + port + "/get_stats"
    resp, err := http.Get(address)

    if err != nil {
      panic(err)
    }
    defer resp.Body.Close()
  }

  func doEvery(d time.Duration, port string, f func(port string)) {
    for x := range time.Tick(d) {
      fmt.Println(x)
      f(port)
    }
  } 
  
  func main() {
  	var port string 
  
    // Taking input from user 
    fmt.Println("Enter the port no you want server to run on") 
    fmt.Scanln(&port) 

    _, err := strconv.ParseInt(port, 10, 0)

    if err != nil {
      fmt.Println("Port must be integer")
      os.Exit(1)
    }

    http.HandleFunc("/", httpServer)
    http.ListenAndServe(":" + port, nil)

    doEvery(20*time.Millisecond, port, getStatistics)
  }