package main

import "reflect"
import "fmt"

type PagedLoanResponse struct {
    Items  []int
}

func printit(v interface{}) {
  fmt.Println(reflect.TypeOf(v))  
  fmt.Println(reflect.TypeOf(v).Kind())  
  fmt.Println(reflect.TypeOf(reflect.New(reflect.TypeOf(v))))

}

func main() {
  var p PagedLoanResponse
  printit(p)
  var i interface{}
  i = 0
  switch i.(type) { 
  case int:
    fmt.Println("int")
    i = "barf"
  }
 fmt.Println(i) 
}
