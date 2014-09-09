package weatherapp9

import (
    "fmt"
  "html/template"
    "net/http"
    "encoding/json"
    "net/url"
    "log"
    "io/ioutil"
    "sync"
    "appengine"
    "appengine/urlfetch"
)
type PastWeather struct {
        Date           string
        PrecipMM       string
        TempMaxC       string
        TempMaxF       string
        TempMinC       string
        TempMinF       string
        WeatherCode    string
        WeatherDesc    []map[string]string
        WeatherIconUrl []map[string]string
  WinDir16Point  string
  WindDirDegree  string
  WindDirection  string
  WindSpeedKmph  string
  WindSpeedMiles string
}

type WeatherRequest struct {
  Query string
  Type  string
}

type Weather struct {
  Request           []WeatherRequest
  Weather           []PastWeather
}

type weatherContainer struct {
  Data Weather
}

func init() {
  http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("stylesheets"))))
  http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
    http.HandleFunc("/", handler)
  http.HandleFunc("/display", display)
}
func handler(rw http.ResponseWriter, r *http.Request) {
  fmt.Fprint(rw, rootForm)
}

const rootForm = `
<!DOCTYPE html>
  <html lang="en">
    <head>
      <meta charset="utf-8">
      <title>Go weather app</title>
      <meta name="description" content="GoWeather App">
      <meta name="author" content="Atul Bhosale">
      <!-- Mobile Specific Metas================================================== -->
      <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
      <!-- CSS ================================================== -->
      <link href="stylesheets/bootstrap.min.css" rel="stylesheet">
      <!-- Validation ================================================== -->
      <script>
        function validateForm()
        {
            var c1=document.forms["myForm"]["city1"].value;
            var c2=document.forms["myForm"]["city2"].value;
            var c3=document.forms["myForm"]["city3"].value;
            var c4=document.forms["myForm"]["city4"].value;
            var c5=document.forms["myForm"]["city5"].value;
            if ((c1==null || c1=="") || (c2==null || c2=="") || (c3==null || c3=="") || (c4==null || c4==""))
            {
                alert("Enter city name ");
                return false;
            }
        }
      </script>
    </head>
      <body>
    <!-- Primary Page Layout ================================================== -->
      <div class="container"> 
      <center>
        <h1 class="remove-bottom" style="margin-top: 40px">Go Weather App</h1>
        <h5>Version 1.0</h5>
        <hr />
        </center>
      
          <form name ="myForm" class="form-horizontal" role="form" action="/display" onsubmit="return validateForm()" method="post" accept-charset="utf-8">
          </div>
          <div class="row">
  <div class="col-md-4"><img src="images/2.jpg" align="left"></div>
  <div class="col-md-4">
              <center>
              <p>Enter the cities for which you want the maximum and minimum temperature of today</p>
                <div class="form-group">
                  <label for="regularInput1">City Name 1</label>
                  <input type="text" name="city1" id="regularInput1" class="form-control" style="text-align:center;"/>
                </div>
                <br>
                <div class="form-group">
                  <label for="regularInput2">City Name 2</label>
                  <input type="text" name="city2" id="regularInput2" class="form-control" style="text-align:center;"/>
                </div>
                <br>
                <div class="form-group">
                   <label for="regularInput3">City Name 3</label>
                   <input type="text" name="city3" id="regularInput3" class="form-control" style="text-align:center;"/>
                </div>
                <br>
                <div class="form-group">
                  <label for="regularInput4">City Name 4</label>
                  <input type="text" name="city4" id="regularInput4" class="form-control" style="text-align:center;"/>
                </div>
                <br>
                <div class="form-group">
                  <label for="regularInput5">City Name 5</label>
                  <input type="text" name="city5" id="regularInput5" class="form-control" style="text-align:center;"/>
                 </div>
                 <br>
              &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
              <button type="submit" class="btn btn-primary">Submit Form</button>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
              <button type="reset" class="btn btn=primary">Reset</button>
              </center>

      </div>
  <div class="col-md-4"><img src="images/4.jpg" align="left" alt="image"></div>
          </form>
</div>
      <!-- End Document================================================== -->
    </body>
</html>`


var dat weatherContainer
func display(w http.ResponseWriter, r *http.Request) {
 addr := []string{r.FormValue("city1"), r.FormValue("city2"), r.FormValue("city3"), r.FormValue("city4"),r.FormValue("city5")}
 var waitGroup sync.WaitGroup
 waitGroup.Add(5)
 for query:= 0; query < 5; query++ {
          go Get(w, r, query, addr,&waitGroup)
 }
    waitGroup.Wait()
    //fmt.Fprint(rw, redirect1)
}
func Get(w http.ResponseWriter, r *http.Request, query int, addr []string, waitGroup *sync.WaitGroup) {
  defer waitGroup.Done()
  s := url.QueryEscape(addr[query])
      url := fmt.Sprintf("http://api.worldweatheronline.com/free/v1/weather.ashx?q=%s&format=json&num_of_days=1&cc=no&includelocation=no&show_comments=no&key=e00de7d84c7606bca515641595683af31f8b00e6", s)
      c := appengine.NewContext(r)
  client :=urlfetch.Client(c)
      resp, err := client.Get(url)
      if err != nil {
              log.Fatal("Get: ", err)
              return
      }
      defer resp.Body.Close()
      fbody, err := ioutil.ReadAll(resp.Body)
      if err != nil {
              log.Fatal("ReadAll: ", err)
              return
      }
      json.Unmarshal(fbody, &dat)
  err1 := displayTemplate.Execute(w, dat)
    if err1 != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var displayTemplate = template.Must(template.New("display").Parse(displayTemplateHTML))

const displayTemplateHTML = `
<!DOCTYPE html>
  <html lang="en">
    <head>
      <meta charset="utf-8">
      <title>Go weather app</title>
      <link href="stylesheets/bootstrap.min.css" rel="stylesheet">  
    </head>
    <body>
    <!-- Primary Page Layout ================================================== -->
      <div class="container"> 
      <center>
        
      <table border=5 class="table">
        <thead>
          <tr>
            <th><center>City</center></th>
            <th><center>Maximum Temperature</center></th>
            <th><center>Minimum Temperature</center></th>
          </tr>
        </thead>
        <tbody>
           <tr class="success">
            <td><center>{{with .Data}}{{range .Request}}{{.Query}}{{end}}{{end}}<center></td>
            <td><center>{{with .Data}}{{range .Weather}}{{.TempMaxC}}{{end}}{{end}}<center></td>
            <td><center>{{with .Data}}{{range .Weather}}{{.TempMinC}}{{end}}{{end}}</center></td>
          </tr>
        </tbody>
      </table>
      <button type="button" class="btn btn-default"><a href="/">Start again</a></button>
      </center>
      </div>
      <br>      
      </body>
</html>`
