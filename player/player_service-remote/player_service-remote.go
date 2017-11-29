// Autogenerated by Thrift Compiler (0.10.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
        "flag"
        "fmt"
        "math"
        "net"
        "net/url"
        "os"
        "strconv"
        "strings"
        "git.apache.org/thrift.git/lib/go/thrift"
        "player"
)


func Usage() {
  fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
  flag.PrintDefaults()
  fmt.Fprintln(os.Stderr, "\nFunctions:")
  fmt.Fprintln(os.Stderr, "  bool ping()")
  fmt.Fprintln(os.Stderr, "  void uploadMap( gamemap)")
  fmt.Fprintln(os.Stderr, "  void uploadParamters(Args arguments)")
  fmt.Fprintln(os.Stderr, "  void assignTanks( tanks)")
  fmt.Fprintln(os.Stderr, "  void latestState(GameState state)")
  fmt.Fprintln(os.Stderr, "   getNewOrders()")
  fmt.Fprintln(os.Stderr)
  os.Exit(0)
}

func main() {
  flag.Usage = Usage
  var host string
  var port int
  var protocol string
  var urlString string
  var framed bool
  var useHttp bool
  var parsedUrl url.URL
  var trans thrift.TTransport
  _ = strconv.Atoi
  _ = math.Abs
  flag.Usage = Usage
  flag.StringVar(&host, "h", "localhost", "Specify host and port")
  flag.IntVar(&port, "p", 9090, "Specify port")
  flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
  flag.StringVar(&urlString, "u", "", "Specify the url")
  flag.BoolVar(&framed, "framed", false, "Use framed transport")
  flag.BoolVar(&useHttp, "http", false, "Use http")
  flag.Parse()
  
  if len(urlString) > 0 {
    parsedUrl, err := url.Parse(urlString)
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
    host = parsedUrl.Host
    useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http"
  } else if useHttp {
    _, err := url.Parse(fmt.Sprint("http://", host, ":", port))
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
  }
  
  cmd := flag.Arg(0)
  var err error
  if useHttp {
    trans, err = thrift.NewTHttpClient(parsedUrl.String())
  } else {
    portStr := fmt.Sprint(port)
    if strings.Contains(host, ":") {
           host, portStr, err = net.SplitHostPort(host)
           if err != nil {
                   fmt.Fprintln(os.Stderr, "error with host:", err)
                   os.Exit(1)
           }
    }
    trans, err = thrift.NewTSocket(net.JoinHostPort(host, portStr))
    if err != nil {
      fmt.Fprintln(os.Stderr, "error resolving address:", err)
      os.Exit(1)
    }
    if framed {
      trans = thrift.NewTFramedTransport(trans)
    }
  }
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating transport", err)
    os.Exit(1)
  }
  defer trans.Close()
  var protocolFactory thrift.TProtocolFactory
  switch protocol {
  case "compact":
    protocolFactory = thrift.NewTCompactProtocolFactory()
    break
  case "simplejson":
    protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
    break
  case "json":
    protocolFactory = thrift.NewTJSONProtocolFactory()
    break
  case "binary", "":
    protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
    Usage()
    os.Exit(1)
  }
  client := player.NewPlayerServiceClientFactory(trans, protocolFactory)
  if err := trans.Open(); err != nil {
    fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
    os.Exit(1)
  }
  
  switch cmd {
  case "ping":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "Ping requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.Ping())
    fmt.Print("\n")
    break
  case "uploadMap":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "UploadMap requires 1 args")
      flag.Usage()
    }
    arg20 := flag.Arg(1)
    mbTrans21 := thrift.NewTMemoryBufferLen(len(arg20))
    defer mbTrans21.Close()
    _, err22 := mbTrans21.WriteString(arg20)
    if err22 != nil { 
      Usage()
      return
    }
    factory23 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt24 := factory23.GetProtocol(mbTrans21)
    containerStruct0 := player.NewPlayerServiceUploadMapArgs()
    err25 := containerStruct0.ReadField1(jsProt24)
    if err25 != nil {
      Usage()
      return
    }
    argvalue0 := containerStruct0.Gamemap
    value0 := argvalue0
    fmt.Print(client.UploadMap(value0))
    fmt.Print("\n")
    break
  case "uploadParamters":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "UploadParamters requires 1 args")
      flag.Usage()
    }
    arg26 := flag.Arg(1)
    mbTrans27 := thrift.NewTMemoryBufferLen(len(arg26))
    defer mbTrans27.Close()
    _, err28 := mbTrans27.WriteString(arg26)
    if err28 != nil {
      Usage()
      return
    }
    factory29 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt30 := factory29.GetProtocol(mbTrans27)
    argvalue0 := player.NewArgs_()
    err31 := argvalue0.Read(jsProt30)
    if err31 != nil {
      Usage()
      return
    }
    value0 := argvalue0
    fmt.Print(client.UploadParamters(value0))
    fmt.Print("\n")
    break
  case "assignTanks":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "AssignTanks requires 1 args")
      flag.Usage()
    }
    arg32 := flag.Arg(1)
    mbTrans33 := thrift.NewTMemoryBufferLen(len(arg32))
    defer mbTrans33.Close()
    _, err34 := mbTrans33.WriteString(arg32)
    if err34 != nil { 
      Usage()
      return
    }
    factory35 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt36 := factory35.GetProtocol(mbTrans33)
    containerStruct0 := player.NewPlayerServiceAssignTanksArgs()
    err37 := containerStruct0.ReadField1(jsProt36)
    if err37 != nil {
      Usage()
      return
    }
    argvalue0 := containerStruct0.Tanks
    value0 := argvalue0
    fmt.Print(client.AssignTanks(value0))
    fmt.Print("\n")
    break
  case "latestState":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "LatestState requires 1 args")
      flag.Usage()
    }
    arg38 := flag.Arg(1)
    mbTrans39 := thrift.NewTMemoryBufferLen(len(arg38))
    defer mbTrans39.Close()
    _, err40 := mbTrans39.WriteString(arg38)
    if err40 != nil {
      Usage()
      return
    }
    factory41 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt42 := factory41.GetProtocol(mbTrans39)
    argvalue0 := player.NewGameState()
    err43 := argvalue0.Read(jsProt42)
    if err43 != nil {
      Usage()
      return
    }
    value0 := argvalue0
    fmt.Print(client.LatestState(value0))
    fmt.Print("\n")
    break
  case "getNewOrders":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetNewOrders requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetNewOrders())
    fmt.Print("\n")
    break
  case "":
    Usage()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
  }
}
