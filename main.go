package main

import (
  "encoding/json"
  "io/ioutil"
  "strings"
  "errors"
  "github.com/kellydunn/golang-geo"
  "github.com/go-martini/martini"
)

// Zipcode struct
type Zip struct {
  City string `json:"city"`
  Lat float64 `json:"lat"`
  Lon float64 `json:"lon"`
}

// ZipMap holds all the zipcodes
type ZipMap map[string]*Zip

// Imports all zipcodes of json file into ZipMap
func (z ZipMap) importZips() error {
  raw, err := ioutil.ReadFile("./data/de.json")
  if err != nil {
    return err
  }
  err = json.Unmarshal(raw, &z)
  if err != nil {
    return err
  }
  return nil
}

// Creates new filled ZipMap
func NewFilledZipMap() (ZipMap, error) {
  m := make(ZipMap)
  err := m.importZips()
  if err != nil {
    return nil, err
  }
  return m, nil
}

// Get Zip from ZipMap by key
func (z ZipMap) SingleZip(s string) (*Zip, error) {
  if val, ok := z[s]; ok {
    return val, nil
  } else {
    return nil, errors.New("Not found")
  }
}

// Get all Zipcode that contains a string
func (z ZipMap) ZipCodesContain(s string) []string {
  found := make([]string, 0)
  for k, _ := range(z) {
    if strings.Contains(k, s) {
      found = append(found, k)
    }
  }
  return found
}

// Get Zipodes within specified area
func (z ZipMap) ZipCodesWithinRadius(s string, r float64) ([]string, error) {
  found := make([]string, 0)
  c, err := z.SingleZip(s)
  if err != nil {
    return nil, err
  }
  p1 := geo.NewPoint(c.Lat, c.Lon)
  for k, v := range z {
    p2 := geo.NewPoint(v.Lat, v.Lon)
    if p1.GreatCircleDistance(p2) < r {
      found = append(found, k)
    }
  }
  return found, nil
}

func main() {
  zips, err := NewFilledZipMap()
  if err != nil {
    errors.New("Could not create ZipMap")
  }
  m := martini.Classic()
  registerRoutes(m, &zips)
  m.RunOnAddr("0.0.0.0:8080")
}
