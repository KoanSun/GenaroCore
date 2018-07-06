package main

import (
	"flag"
	"net/http"
	"io/ioutil"
	"strings"
	"github.com/tidwall/gjson"
	"github.com/GenaroNetwork/Genaro-Core/common/hexutil"
	"fmt"
	"github.com/GenaroNetwork/Genaro-Core/common"
	"log"
)

var rpcurl string

func HttpPost(url string, contentType string, body string) ([]byte, error) {
	bodyio := strings.NewReader(body)
	resp, err := http.Post(url,contentType,bodyio)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	repbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return repbody, nil
}

func GetCuBlockNum(url string) (uint64,error){
	ret,err := HttpPost(url,"application/json",`{"jsonrpc":"2.0","id":1,"method":"eth_blockNumber","params":[]}`)
	if err != nil {
		return 0,err
	}
	blockNumStr := gjson.GetBytes(ret,"result").String()
	blockNum,err := hexutil.DecodeUint64(blockNumStr)
	if err != nil {
		return 0,err
	}
	return blockNum,nil
}

func GetBlockByNumber(url string,blockNum uint64) ([]byte,error) {
	blockNumHex := hexutil.EncodeUint64(blockNum)
	ret,err := HttpPost(url,"application/json",`{"jsonrpc":"2.0","id":1,"method":"eth_getBlockByNumber","params":["`+blockNumHex+`",true]}`)
	if err != nil {
		return nil,err
	}
	return ret,err
}

func GetBlockHash(url string,blockNum uint64) (string,error){
	ret,err := GetBlockByNumber(url,blockNum)
	if err != nil {
		return "",err
	}
	blockHash := gjson.GetBytes(ret,"result.hash").String()
	return blockHash,nil
}

func SendSynState(url string,blockHash string) (string,error){
	ret,err := HttpPost(url,"application/json",`{"jsonrpc": "2.0","method": "eth_sendTransaction","params": [{"from": "0xad188b762f9e3ef76c972960b80c9dc99b9cfc73","to": "`+common.SpecialSyncAddress.String()+`","gas": "0x72bf0","gasPrice": "0x9172a","value": "0x1","extraData": "{\"msg\": \"`+blockHash+`\",\"type\": \"0xd\"}"}],"id": 1}`)
	if err != nil {
		return "",err
	}
	return gjson.ParseBytes(ret).String(),nil
}

func initarg() {
	flag.StringVar(&rpcurl, "u", "http://127.0.0.1:8545", "rpc url")
	flag.Parse()
}

func main() {
	initarg()
	cuBlockNum,err := GetCuBlockNum(rpcurl)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(cuBlockNum)
	synBlockNum := cuBlockNum/6
	if synBlockNum != 0 {
		synBlockHash,err := GetBlockHash(rpcurl,synBlockNum*6)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(synBlockHash)
		ret,err := SendSynState(rpcurl,synBlockHash)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(ret)
	}
}
