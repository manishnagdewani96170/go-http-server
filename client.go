package main

import (
  "fmt"
  "log"
  "os"
  "path/filepath"
  "net/http"
  "encoding/json"
  "strconv"
  "bytes"
  // "sync"
  // "time"

)
// var wg sync.WaitGroup

func main() {
  var path string
  fmt.Println("Enter the file path")
  fmt.Scanln(&path)

  var address string
  fmt.Println("Enter the address") 
  fmt.Scanln(&address)

  current_directory := path

  if current_directory == "" {
    directory, err := os.Getwd();
    if err != nil {
      log.Fatal(err)
    }
    current_directory = directory
  }
  iterate(current_directory, address)
  
  // c := make(chan string)
    // wg.Add(1)

    // go iterate(current_directory, address, c)
  
  // wg.Wait()
  // close(c)
}

func iterate(path string, address string) {
  // defer wg.Done()
  
  // time.Sleep(1 * time.Microsecond)

  filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      log.Fatalf(err.Error())
      return nil
    }
    // fi, err := info.Stat();

    // wg.Add(1)
    // go sendRequest(path, info, address, c)
    if !info.IsDir() {
      sendRequest(path, info, address)
    }
     
    return nil
  })
}

func sendRequest(path string, info os.FileInfo, address string) {
  //defer wg.Done()
   
  fmt.Printf("File Size: %d\n", info.Size())
  fmt.Printf("File extension: %s\n", filepath.Ext(path))
  fmt.Printf("File Name: %s\n", info.Name())
  fmt.Printf("File Path: %s\n", path)

  //time.Sleep(1 * time.Microsecond)

  values := map[string]string{
    "Size": strconv.FormatInt(int64(info.Size()), 10),
    "Extension": filepath.Ext(path),
    "Name": info.Name(),
    "Path": path,
  }

  jsonValue, _ := json.Marshal(values)
  

  resp, err := http.Post(address, "application/json", bytes.NewBuffer(jsonValue))
  if err != nil {
    panic(err)
    //c <- "We could not reach:"     // pump the result into the channel
  }
  defer resp.Body.Close()

}