package main

import (
"encoding/json"
"fmt"
"io/ioutil"
"log"
"net/http"
"strings"
"time"

"github.com/mantyr/pricer"
)
const (
fruit_req_api = "https://run.mocky.io/v3/c51441de-5c1a-4dc2-a44e-aab4f619926b"
veg_req_api   = "https://run.mocky.io/v3/4ec58fbc-e9e5-4ace-9ff0-4e893ef9663c"
grain_req_api = "https://run.mocky.io/v3/e6c77e5c-aec9-403f-821b-e14114220148"
)

type Fruit struct {
ID       string `json:"id"`
Name     string `json:"name"`
Quantity int    `json:"quantity"`
Price    string `json:"price"`
}

type Vegetable struct {
ID       string `json:"productId"`
Name     string `json:"productName"`
Quantity int    `json:"quantity"`
Price    string `json:"price"`
}
type Grain struct {
ID       string `json:"itemId"`
Name     string `json:"itemName"`
Quantity int    `json:"quantity"`
Price    string `json:"price"`
}

var fruit []Fruit
var vegetable []Vegetable
var grain []Grain

type JsonRequest struct {
Name     string `json:"name"`
Quantity int    `json:"quantity"`
Price    string `json:"price"`
}
type summary struct {
Fruitsum     []Fruit     `json:"Fruit"`
Vegetablesum []Vegetable `json:"Vegetables"`
Grainsum     []Grain     `json:"Grain"`
}

func showsummary(w http.ResponseWriter, r *http.Request) {
var show summary
show.Fruitsum = fruit
show.Grainsum = grain
show.Vegetablesum = vegetable
err := json.NewEncoder(w).Encode(show)
if err != nil {
fmt.Println(err)
}

}



func buyItem(w http.ResponseWriter, r *http.Request) {
FruitSupplier()
Vegetablesupplier()
Grainsupplier()
var req JsonRequest
err := json.NewDecoder(r.Body).Decode(&req)
if err != nil {
fmt.Println(err)
}
if fastbuyhelper(req) == true {
fmt.Fprintf(w, req.Name)
return
}
fmt.Fprintf(w, "NOT_FOUND")
}

func buyItemQtyPrice(w http.ResponseWriter, r *http.Request) {
var req JsonRequest
err := json.NewDecoder(r.Body).Decode(&req)
if err != nil {
fmt.Println(err)
}
//searching initially
if helper(req) == true {
fmt.Fprintf(w, req.Name)
return
}
fruit, _ = FruitSupplier()
vegetable, _ = Vegetablesupplier()
grain, _ = Grainsupplier()
if helper(req) == true {
fmt.Fprintf(w, req.Name)
return
}
fmt.Fprintf(w, "NOT_FOUND")
}

func buyItemQty(w http.ResponseWriter, r *http.Request) {
FruitSupplier()
Vegetablesupplier()
Grainsupplier()
var req JsonRequest
err := json.NewDecoder(r.Body).Decode(&req)
if err != nil {
fmt.Println(err)
}
for i, _ := range fruit {
if strings.EqualFold(fruit[i].Name, req.Name) && fruit[i].Quantity >= req.Quantity {
fmt.Fprintf(w, req.Name)
return
}
if strings.EqualFold(vegetable[i].Name, req.Name) && vegetable[i].Quantity >= req.Quantity {
fmt.Fprintf(w, req.Name)
return
}
if strings.EqualFold(grain[i].Name, req.Name) && grain[i].Quantity >= req.Quantity {
fmt.Fprintf(w, req.Name)
return
}
}
fmt.Fprintf(w, "NOT_FOUND")
}

func FruitSupplier() ([]Fruit, error) {
client := &http.Client{Timeout: 60 * time.Second}
// Declare HTTP Method and Url
req, err := http.NewRequest("GET", fruit_req_api, nil)
if err != nil {
fmt.Println(err)
}
resp, err := client.Do(req)
if err != nil {
fmt.Println(err)
}
body, _ := ioutil.ReadAll(resp.Body)
jsonErr := json.Unmarshal([]byte(body), &fruit)
if jsonErr != nil {
log.Fatal(jsonErr)
}
//The client must close the response body when finished with it.
defer resp.Body.Close()
return fruit, nil
}

func Vegetablesupplier() ([]Vegetable, error) {
client := &http.Client{Timeout: 60 * time.Second}
// Declare HTTP Method and Url
req, err := http.NewRequest("GET", veg_req_api, nil)
if err != nil {
fmt.Println(err)
}
resp, err := client.Do(req)
if err != nil {
fmt.Println(err)
}
body, _ := ioutil.ReadAll(resp.Body)
jsonErr := json.Unmarshal([]byte(body), &vegetable)
if jsonErr != nil {
log.Fatal(jsonErr)
}
//The client must close the response body when finished with it.
defer resp.Body.Close()
return vegetable, nil
}
func Grainsupplier() ([]Grain, error) {
client := &http.Client{Timeout: 60 * time.Second}
// Declare HTTP Method and Url
req, err := http.NewRequest("GET", grain_req_api, nil)
if err != nil {
fmt.Println(err)
}
resp, err := client.Do(req)
if err != nil {
fmt.Println(err)
}
body, _ := ioutil.ReadAll(resp.Body)
jsonErr := json.Unmarshal([]byte(body), &grain)
if jsonErr != nil {
log.Fatal(jsonErr)
}
//The client must close the response body when finished with it.
defer resp.Body.Close()
return grain, nil
}
func main() {
handlers()
}
func handlers() {
http.HandleFunc("/buy-item", buyItem)
http.HandleFunc("/buy-item-qty", buyItemQty)
http.HandleFunc("/buy-item-qty-price", buyItemQtyPrice)
http.HandleFunc("/show-summary", showsummary)
http.HandleFunc("/fast-buy-item", fastBuy)
log.Fatal(http.ListenAndServe(":9090", nil))
}
func comp(item string, req string) bool {
price := pricer.NewPrice()
price.SetDefaultType("DEFAULT_TYPE")
price.Parse(item)
price.GetFloat64()
price2 := pricer.NewPrice()
price2.SetDefaultType("DEFAULT_TYPE")
price2.Parse(req)
price2.GetFloat64()
if price.GetFloat64() <= price2.GetFloat64() {
return true
}
return false
}

func fastBuy(w http.ResponseWriter, r *http.Request) {
go FruitSupplier()
go Vegetablesupplier()
go Grainsupplier()
var req JsonRequest
err := json.NewDecoder(r.Body).Decode(&req)
fmt.Println(err)
if fastbuyhelper(req) == true {
fmt.Fprintf(w, req.Name)
return
}
fmt.Fprintf(w, "NOT_FOUND")
}
func fastbuyhelper(req JsonRequest) bool {
for i, _ := range fruit {
if strings.EqualFold(fruit[i].Name, req.Name) && fruit[i].Quantity > 0 {
return true
}
if strings.EqualFold(vegetable[i].Name, req.Name) && vegetable[i].Quantity > 0 {
return true
}
if strings.EqualFold(grain[i].Name, req.Name) && grain[i].Quantity > 0 {
return true
}
}
return false
}

func helper(req JsonRequest) bool {
for i, _ := range fruit {
if strings.EqualFold(fruit[i].Name, req.Name) && fruit[i].Quantity >= req.Quantity && comp(fruit[i].Price, req.Price) {
return true
}
if strings.EqualFold(vegetable[i].Name, req.Name) && vegetable[i].Quantity >= req.Quantity && comp(vegetable[i].Price, req.Price) {
return true
}
if strings.EqualFold(grain[i].Name, req.Name) && grain[i].Quantity >= req.Quantity && comp(grain[i].Price, req.Price) {
return true
}
}
return false
}
