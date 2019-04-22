package main

import "fmt"
import "math/big"
import "time"
import "crypto/rand"
import "github.com/docopt/docopt-go"
import "gopkg.in/yaml.v2"
import "io/ioutil"
import "log"
import "strconv"
import "os"
import "encoding/csv"
import "bufio"
import "io"
func blom_key_distribution(n_string *string, t_string *string, q_string *string) bool{

    var matrix_size int
    var k_add_one int
    n := new(big.Int)
    t := new(big.Int)
    q := new(big.Int)

    n.SetString(*n_string, 10)
    t.SetString(*t_string, 10)
    q.SetString(*q_string, 16)
    k_add_one, _ = strconv.Atoi(*t_string)
    matrix_size = (k_add_one *  2) + 1
    k_add_one += 1
    // If q is not prime return false
    if !q.ProbablyPrime(20){
        return false
    }
    // If t < n - 2 exit
    if t.Cmp(n.Sub(n, big.NewInt(2)))== 1{
      return false
    }
    A := make([][]big.Int, k_add_one)
    for i := range A {
      A[i] = make([]big.Int, k_add_one)
    }
    tmp := new(big.Int)
    for i, row := range A {
      for j := range row{
        // If A[i][j] an A[j][i] are empty generate random number for them and assign it
        if A[i][j].Cmp(zero) == 0 && A[j][i].Cmp(zero) == 0{
          tmp, _ = rand.Int(rand.Reader, q)
          A[i][j] = *tmp
          A[j][i] = *tmp
        }
      }
    }

    for n := 1; n <= matrix_size; n++{
      party_polynomial := make([]big.Int, 0)
      for index := range A{
        party_polynomial = append(party_polynomial, polynomal_summation(A[index], big.NewInt(int64(n)), q))
      }
      // Out polynomial number representation as comma seperated values
      f, _ := os.Create("output/party" + strconv.Itoa(n) + ".txt" )
      for number := range party_polynomial{
        f.WriteString(fmt.Sprint(&party_polynomial[number]))
        if number+1 < len(party_polynomial){
          f.WriteString(",")
        }
      }
      f.Close()
    }
    return true

}
// Global Big ints for 0 and 1
var one = big.NewInt(1)
var zero = big.NewInt(0)

func polynomal_summation(row []big.Int, y *big.Int, q *big.Int) (big.Int){
  sum := new(big.Int)
  modresult := new(big.Int)
  temp := new(big.Int)
  for index := range row{
    // y * j mod q
    modresult.Exp(y, big.NewInt(int64(index)), q)
    // aij * (y * j mod q)
    temp = modresult.Mul(modresult, &row[index])
    temp = temp.Mod(temp, q)
    sum.Add(sum, temp)
  }
  sum.Mod(sum, q)
  // Return summation
  return *sum
}

// Defines how the yaml config is structed
type conf struct {
    N string `yaml:"n"`
    T string `yaml:"t"`
    Q string `yaml:"q"`
}

// Load the config file from a yaml file
func (c *conf) getConf(filename *string) *conf {

    yamlFile, err := ioutil.ReadFile(*filename)
    if err != nil {
        log.Printf("yamlFile.Get err   ", err)
    }
    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        log.Fatalf("Unmarshal: ", err)
    }
    return c
}

func return_big_int(number string)(big.Int){
  temp := new(big.Int)
  temp.SetString(number, 10)
  return *temp
}

func generate_session_key(filename string, other_party *big.Int, q *big.Int){
  A := make([]big.Int, 0)
  session := new(big.Int)
  csvFile, _ := os.Open(filename)
  reader := csv.NewReader(bufio.NewReader(csvFile))
  for {
      line, error := reader.Read()
      if error == io.EOF {
          break
      } else if error != nil {
            fmt.Println(error)
            break
        }
      for column := range line{
        A = append(A, return_big_int(line[column]))
      }
      *session = polynomal_summation(A, other_party, q)
      fmt.Println("This is the session key", session)
    }
}

func main(){

    usage := `bloms_scheme.

    Usage:
    bloms_scheme [--n=<n>][--t=<t>][--q=<q>][--config=<config>]
    bloms_scheme [--session_key=<session_key>][--other_party=<other_party>][--q=<q>][--config=<config>]
    bloms_scheme -h | --help
    bloms_scheme --version

    Options:
    -h --help                       Show this screen.
    --version                       Show version.
    --n=<n>                         n is the number of parties in the key exchange
    --t=<t>                         An integer such that t < new
    --q=<q>                         A prime number
    --config=<config>               A yaml file containing the config for n, t, and q
    --session_key=<session_key>     Name of the file with the existing polynomial (IE output/party4.txt)
    --other_party=<other_party>     This is the other party you want to have secure communication with  (IE 5)
    NOTE: values in the config file overwrite command line args. Args in config file can be empty`
    start := time.Now()
    arguments, _ := docopt.ParseDoc(usage)
    var c conf
    var t string
    var n string
    var q string
    // If there is a config file load it
    if arguments["--config"] != nil{
      str, _ := arguments["--config"].(string)
      c.getConf(&str)
      // If there is no value in the config file take the variable in from as an argument
      if c.N != ""{
        n = c.N
      } else {
        n = arguments["--n"].(string)
      }
      if c.T != ""{
        t = c.T
      } else {
        t = arguments["--t"].(string)
      }
      if c.Q != ""{
        q = c.Q
      } else {
        q = arguments["--q"].(string)
      }
    } else {
      n = arguments["--n"].(string)
      t = arguments["--t"].(string)
      q = arguments["--q"].(string)
    }
    if arguments["--session_key"] != nil{
      q_big_int := new(big.Int)
      q_big_int.SetString(q, 16)
      other_party := new(big.Int)
      other_party.SetString(arguments["--other_party"].(string), 10)
      generate_session_key(arguments["--session_key"].(string) ,other_party, q_big_int)
    } else{
      blom_key_distribution(&n, &t, &q)

    }
    elapsed := time.Since(start)
    fmt.Println(elapsed)


}
