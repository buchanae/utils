package main

import (
  "encoding/json"
  "encoding/csv"
  "fmt"
  "flag"
  "io"
  "os"
)

func main() {
  var input string

  flag.StringVar(&input, "input", input, "Input JSON file.")
  flag.Parse()

  if input == "" {
    flag.PrintDefaults()
    os.Exit(1)
  }

  var f io.Reader
  if input == "-" {
    f = os.Stdin
  } else {
    var err error
    f, err = os.Open(input)
    if err != nil {
      panic(err)
    }
  }

  dec := json.NewDecoder(f)
  uniqcolumns := map[string]struct{}{}
  var firstdata []map[string]interface{}

  i := 0
  for {
    dat := map[string]interface{}{}
    err := dec.Decode(&dat)
    if err == io.EOF {
      break
    }
    if err != nil {
      panic(err)
    }
    firstdata = append(firstdata, dat)
    for k, _ := range dat {
      uniqcolumns[k] = struct{}{}
    }
    i++

    if i > 500 {
      break
    }
  }

  var columns []string
  for k, _ := range uniqcolumns {
    columns = append(columns, k)
  }

  writer := csv.NewWriter(os.Stdout)
  defer writer.Flush()

  writer.Write(columns)

  for _, dat := range firstdata {
    var row []string

    for _, col := range columns {
      v, ok := dat[col]
      if ok {
        row = append(row, fmt.Sprintf("%s", v))
      } else {
        row = append(row, "")
      }
    }

    writer.Write(row)
  }

  for {
    dat := map[string]interface{}{}
    err := dec.Decode(&dat)
    if err == io.EOF {
      break
    }
    if err != nil {
      panic(err)
    }

    var row []string

    for _, col := range columns {
      v, ok := dat[col]
      if ok {
        row = append(row, fmt.Sprintf("%s", v))
      } else {
        row = append(row, "")
      }
    }

    writer.Write(row)
  }
}
