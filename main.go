package main

import(
	"log"
	"net/http"
    "github.com/gorilla/websocket"
	"math/rand"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	// "fmt"
	"strconv"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define our message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

var clients = make(map[*websocket.Conn]bool) // connected clients
func main() {
	
	// Configure websocket route
	http.HandleFunc("/handleConnections", handleConnections)

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		msg.Message = strconv.Itoa(GetRandom(100)) 
		ws.WriteJSON(msg)
	}
}


// GetRandom is a call func
// retrun 0 ~ num-1.
func GetRandom(num int) int {
	s1 := rand.NewSource(GetSefeRandomSeed())
	r1 := rand.New(s1)
	x := r1.Intn(num)
	// fmt.Print("Rand result: = ", x, ", ")

	return x
}

// GetSefeRandomSeed is prodiver safe seed function
func GetSefeRandomSeed() int64 {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatal(err)
	}

	mnemonic, _ := bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed(mnemonic, "")

	// fmt.Println("memonic: ", mnemonic)
	// fmt.Println("seed: ", seed)

	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(account.Address.Hex())
	str := account.Address.Hex()
	subs := str[2:9]
	res, err := strconv.ParseInt(subs, 16, 64)

	if err != nil {
		log.Fatal(err)
	}

	return res
}