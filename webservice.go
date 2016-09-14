package main

import(
  "net/http"
  "encoding/json"
  "strconv"
  "github.com/go-martini/martini"
)

func registerRoutes(m *martini.ClassicMartini, z *ZipMap) {
  path := z.GetPath()
  m.Get(path, z.WebGet)
  m.Get(path+"/:id", z.WebGet)
}

// Get Path for Zip Map
func (z *ZipMap) GetPath() string {
  return "/zipcodes"
}

func (z *ZipMap) WebGet(req *http.Request, p martini.Params) (int, string) {
  query := req.URL.Query()
  if p["id"] != "" {
    zip, err := z.SingleZip(p["id"])
    if err != nil {
      return StatusCodeToJson(500)
    }
    return toJSON(zip)
  }
  if query["contain"] != nil {
    codes := z.ZipCodesContain(query["contain"][0])
    return toJSON(codes)
  }
  if query["near"] != nil {
    radius := 10.0
    if query["radius"] != nil {
      radius, _ = strconv.ParseFloat(query["radius"][0], 64)
    }
    codes, err := z.ZipCodesWithinRadius(query["near"][0], radius)
    if err != nil {
      return StatusCodeToJson(500)
    }
    return toJSON(codes)
  }
  return StatusCodeToJson(404)
}

func toJSON(v interface{}) (int, string) {
  resp, err := json.Marshal(v)
  if err != nil {
    return StatusCodeToJson(500)
  } else {
    return 200, string(resp[:])
  }
}

func StatusCodeToJson(i int) (int, string) {
  resp, _ := json.Marshal(http.StatusText(i))
  return i, string(resp[:])
}
