package echo

import (
  "encoding/json"
  "errors"
  // "log"
  "net/http"
  "fmt"
  // "strings"

  "golang.org/x/net/context"

  "github.com/go-kit/kit/endpoint"
  httptransport "github.com/go-kit/kit/transport/http"
)

type EchoService interface {
  Echo(string) (string, error)
}

type echoService struct {}

var ErrEmpty = errors.New("empty string")

func (echoService) Echo(s string) (string, error) {
  if s == "" {
    return "", ErrEmpty
  }

  return s, nil
}

func decodeEchoRequest(r *http.Request) (interface{}, error) {
  var request echoRequest
  if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
    return nil, err
  }
  return request, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
  return json.NewEncoder(w).Encode(response)
}

func makeEchoEndpoint(svc EchoService) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(echoRequest)
    v, err := svc.Echo(req.S)
    if err != nil {
      return echoResponse{v, err.Error()}, nil
    }

    return echoResponse{v, ""}, nil

  }
}

type echoRequest struct {
  S string `json:"s"`
}

type echoResponse struct {
  V string `json:"v"`
  Err string `json:"err,omitempty"`
}

func hello(w http.ResponseWriter, r *http.Request){
    fmt.Fprint(w, "Hello 123")
}

func init() {
  ctx := context.Background()
  svc := echoService{}

  echoHandler := httptransport.NewServer(
    ctx,
    makeEchoEndpoint(svc),
    decodeEchoRequest,
    encodeResponse,
  )

  http.HandleFunc("/", hello)
  http.Handle("/echo", echoHandler)
  // log.Fatal(http.ListenAndServe(":8080", nil))
}

