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
    NOTE: values in the config file overwrite command line args. Args in config file can be empty

This program allows you to generate keys for distribution or generate a session key given keys have alread been distributed.

This program assume that n and t are in decimal format and q is in hex format

All of the example tests listed in the project description are ijn the testfile directory and can be run:
`go run bloms_scheme --config testfiles/testfile1.yaml`

NOTE: You must have a output/ directory 

Example usage:
    go run bloms_scheme --n 1000 --t 50 --q 12cb639444cfb091833261b9def68e68b
    go run bloms_scheme --config testfiles/testfile1.yaml
    go run bloms_scheme --config testfiles/testfile14.yaml
    go run bloms_scheme --session_key output/party4.txt --other_party 5 --config testfiles/testfile1.yaml
    go run bloms_scheme --session_key output/party5.txt --other_party 4 --config testfiles/testfile1.txt
    
TODO:
    More Error checking. In many cases where I use the _ I am just throwing away the potential error thrown by a funciton
    Be able to specify if the values for n, t, q or the output should be in decimal or hex
    Be able to print out the matrix and polynomials to stdout
    Be able to print session keys out to file
    Fix docopt to throw error if an invalid combination of arguements is passed in
 

Design:
   When testing the speed of the program I commented out lines 77-84 which include the file I/O for writing out the files with the polynomial numbers

Dependencies:
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
You can get these package by typing "go get package" IE (go get github.com/docopt/docopt-go)

The toy example given in the project description is also it's own file. I used this example to help trouble shoot the larger issue

Time:
for testfile1.yaml 43.692506613s for 100 runs 436.925066ms per run
for testfile2.yaml 47.114040161s for 100 runs 471.140401ms per run
for testfile3.yaml 47.186525559s for 100 runs 471.865255ms per run
test 4 
