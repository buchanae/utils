package main

import (
  "encoding/csv"
  "flag"
  "io"
  "os"
)

func main() {
  var path, keyCol, valueCol, indexCol string

  flag.StringVar(&path, "path", path, "Path to input file")
  flag.StringVar(&keyCol, "key", keyCol, "Name of key column")
  flag.StringVar(&indexCol, "index", indexCol, "Name of index column")
  flag.StringVar(&valueCol, "value", valueCol, "Name of value column")

  flag.Parse()

  if path == "" || keyCol == "" || valueCol == "" || indexCol == "" {
    flag.PrintDefaults()
    os.Exit(1)
  }

  r, err := os.Open(path)
  if err != nil {
    panic(err)
  }

  reader := csv.NewReader(r)
  reader.Comma = '\t'
  reader.ReuseRecord = true

  writer := csv.NewWriter(os.Stdout)
  //writer.Comma = '\t'
  defer writer.Flush()

  h, err := reader.Read()
  if err != nil {
    panic(err)
  }
  header := make([]string, len(h))
  copy(header, h)
  writer.Write([]string{indexCol, keyCol, valueCol})

  for {
    row, err := reader.Read()
    if err == io.EOF {
      break
    }
    if err != nil {
      panic(err)
    }

    index := row[0]
    for i := 1; i < len(header) - 1; i++ {
      key := header[i]
      value := row[i]
      writer.Write([]string{index, key, value})
    }
  }
}
