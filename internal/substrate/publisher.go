package substrate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/parser"
	retriever "github.com/centrifuge/go-substrate-rpc-client/v4/registry/retriever"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/synternet/data-layer-sdk/pkg/options"
	"github.com/synternet/data-layer-sdk/pkg/service"
	f "github.com/vedhavyas/go-subkey"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (*BalancesTransfer) ProtoReflect() protoreflect.Message { return nil }

// TODO: check if this is needed.
// func (*BalancesTransfer) ProtoMessage() {}

type BalancesTransfer struct {
	BlockHash string `json:"block_hash"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
}

func (*ItemAdded) ProtoReflect() protoreflect.Message { return nil }

type ItemAdded struct {
	Who      string `json:"who"`
	ItemType string `json:"item_type"`
	Item     string `json:"item"`
}

func (*ItemUpdated) ProtoReflect() protoreflect.Message { return nil }

type ItemUpdated struct {
	Who      string `json:"who"`
	ItemType string `json:"item_type"`
	Item     string `json:"item"`
}

func (*ItemRead) ProtoReflect() protoreflect.Message { return nil }

type ItemRead struct {
	Item string `json:"item"`
}

func (*ExBlock) ProtoReflect() protoreflect.Message { return nil }

type ExBlock struct {
	Header     ExHeader
	Extrinsics []types.Extrinsic
}

type ExHeader struct {
	Hash types.Hash `json:"hash"`
	types.Header
}

type Publisher struct {
	*service.Service
}

func New(opts ...options.Option) *Publisher {
	ret := &Publisher{
		Service: &service.Service{},
	}

	ret.Configure(opts...)

	return ret
}

func (p *Publisher) Start() context.Context {
	ctx := context.Background()
	go func() {
		for {
			if err := p.RunSubstrate(ctx); err != nil {
				fmt.Println("Error while processing messages:", err)
				retry := 10 * time.Second
				fmt.Printf("Retrying in %f seconds...", retry.Seconds())
				time.Sleep(retry)
			} else {
				break
			}
		}
	}()
	return p.Service.Start()
}

func (p *Publisher) RunSubstrate(ctx context.Context) error {
	log.Println("Connecting to substrate API...")
	api, err := gsrpc.NewSubstrateAPI(p.RPCApi())
	if err != nil {
		return err
	}

	log.Println("Subscribing to head...")
	subNewHeads, err := api.RPC.Chain.SubscribeNewHeads()
	if err != nil {
		panic(err)
	}
	defer subNewHeads.Unsubscribe()

	retriever, err := retriever.NewDefaultEventRetriever(state.NewEventProvider(api.RPC.State), api.RPC.State)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case err := <-subNewHeads.Err():
			api.Client.Close()
			return err
		case head := <-subNewHeads.Chan():
			// Useful for testing peaq agung testnet
			// headNumber := 2045829 // ItemAdded
			// headNumber := 2043805 // ItemUpdated
			// headNumber := 2008603 // ItemRead
			// hash, err := api.RPC.Chain.GetBlockHash(uint64(headNumber))

			hash, err := api.RPC.Chain.GetBlockHash(uint64(head.Number))
			if err != nil {
				log.Fatalln(err)
			}

			events, err := retriever.GetEvents(hash)
			if err != nil {
				log.Printf("Couldn't retrieve events for block number %d: %s\n", head.Number, err)
				continue
			}

			// PeaqStorage.* events are Peaq chain specific.
			for _, e := range events {
				switch e.Name {
				case "Balances.Transfer":
					from, err := SS58AddrFromBalanceTransfer(e, 0, 0)
					if err != nil {
						log.Println(err)
						continue
					}

					to, err := SS58AddrFromBalanceTransfer(e, 1, 0)
					if err != nil {
						log.Println(err)
						continue
					}

					amount, ok := e.Fields[2].Value.(types.U128)
					if !ok {
						log.Println("Could not get amount")
						continue
					}

					divisor := big.NewInt(1e10)
					value := new(big.Float).Quo(new(big.Float).SetInt(amount.Int), new(big.Float).SetInt(divisor))
					balancesTransfer := BalancesTransfer{BlockHash: hash.Hex(), From: from, To: to, Value: value.String()}

					// Useful when debugging locally.
					// balancesTransferJson, err := json.Marshal(balancesTransfer)
					// if err != nil {
					// 	panic(err)
					// }
					// fmt.Printf("\t{prefix}.balances.transfer\n")
					// fmt.Printf("\t\t%s\n", balancesTransferJson)

					p.Publish(
						&balancesTransfer,
						"balances.transfer",
					)
				case "PeaqStorage.ItemAdded":
					who, err := SS58AddrFromBalanceTransfer(e, 0, 42)
					if err != nil {
						log.Println(err)
						continue
					}

					itemType, err := FromPeaqStorageItem(e, 1)
					if err != nil {
						log.Println(err)
						continue
					}

					item, err := FromPeaqStorageItem(e, 2)
					if err != nil {
						log.Println(err)
						continue
					}

					itemAdded := ItemAdded{Who: who, ItemType: itemType, Item: item}

					// Useful when debugging locally.
					// itemAddedJson, err := json.Marshal(itemAdded)
					// if err != nil {
					// 	panic(err)
					// }
					// fmt.Printf("\t{prefix}.peaq-storage.item-added\n")
					// fmt.Printf("\t\t%s\n", &itemAddedJson)

					p.Publish(
						&itemAdded,
						"peaq-storage.item-added",
					)
				case "PeaqStorage.ItemUpdated":
					who, err := SS58AddrFromBalanceTransfer(e, 0, 42)
					if err != nil {
						log.Println(err)
						continue
					}

					itemType, err := FromPeaqStorageItem(e, 1)
					if err != nil {
						log.Println(err)
						continue
					}
					//
					item, err := FromPeaqStorageItem(e, 2)
					if err != nil {
						log.Println(err)
						continue
					}

					itemUpdated := ItemUpdated{Who: who, ItemType: itemType, Item: item}

					// Useful when debugging locally.
					// itemUpdatedJson, err := json.Marshal(itemUpdated)
					// if err != nil {
					// 	panic(err)
					// }
					// fmt.Printf("\t{prefix}.peaq-storage.item-updated\n")
					// fmt.Printf("\t\t%s\n", &itemUpdatedJson)

					p.Publish(
						&itemUpdated,
						"peaq-storage.item-updated",
					)
				case "PeaqStorage.ItemRead":
					item, err := FromPeaqStorageItem(e, 0)
					if err != nil {
						log.Println(err)
						continue
					}

					itemRead := ItemRead{Item: item}

					// Useful when debugging locally.
					itemReadJson, err := json.Marshal(itemRead)
					if err != nil {
						panic(err)
					}
					fmt.Printf("\t{prefix}.peaq-storage.item-read\n")
					fmt.Printf("\t\t%s\n", &itemReadJson)

					p.Publish(
						&itemRead,
						"peaq-storage.item-read",
					)
				}
			}

			block, err := api.RPC.Chain.GetBlock(hash)
			if err != nil {
				log.Fatalln(err)
			}

			exBlock := ExBlock{Header: ExHeader{Hash: hash, Header: block.Block.Header}, Extrinsics: block.Block.Extrinsics}

			// Useful when debugging locally.
			// exBlock := ExBlock{Header: ExHeader{Hash: hash, Header: block.Block.Header}, Extrinsics: nil}
			// jsonBlock, err := json.Marshal(exBlock)
			// if err != nil {
			// 	log.Fatalf("Error marshalling block to JSON: %s\n", err)
			// }
			// fmt.Printf("\t{prefix}.block\n")
			// fmt.Printf("\t\t%s", jsonBlock)

			p.Publish(&exBlock, "block")
		}
	}
}

func (p *Publisher) Close() error {
	log.Println("Publisher.Close")
	p.Cancel(nil)

	var err []error

	log.Println("Waiting on publisher group")
	errGr := p.Group.Wait()
	if !errors.Is(errGr, context.Canceled) {
		err = append(err, errGr)
	}
	log.Println("Publisher.Close DONE")
	return errors.Join(err...)
}

// See https://github.com/centrifuge/go-substrate-rpc-client/issues/370
func SS58AddrFromBalanceTransfer(e *parser.Event, idx int, format uint16) (string, error) {
	if len(e.Fields) != 3 {
		return "", fmt.Errorf("Balances.Transfer should have 3 fields")
	}
	valueDecodedFields, ok := e.Fields[idx].Value.(registry.DecodedFields)
	if !ok {
		return "", fmt.Errorf("expected the value to be of type registry.DecodedFields")
	}
	v, ok := valueDecodedFields[0].Value.([]interface{})
	if !ok {
		return "", fmt.Errorf("expected to be able to cast to []interface{}")
	}
	var bytes []byte
	for _, elem := range v {
		byteVal, ok := elem.(types.U8)
		if !ok {
			return "", fmt.Errorf("element is not a byte")
		}
		bytes = append(bytes, uint8(byteVal))
	}

	// Useful for debugging.
	// for i := uint16(0); i < 65; i++ {
	// 	log.Println(i, "=>", f.SS58Encode(bytes, i))
	// }

	return f.SS58Encode(bytes, format), nil
}

func FromPeaqStorageItem(e *parser.Event, idx int) (string, error) {
	// if len(e.Fields) != 3 {
	// 	return "", fmt.Errorf("Balances.Transfer should have 3 fields")
	// }
	// valueDecodedFields, ok := e.Fields[idx].Value.(registry.DecodedFields)
	// if !ok {
	// 	return "", fmt.Errorf("expected the value to be of type registry.DecodedFields")
	// }
	v, ok := e.Fields[idx].Value.([]interface{})
	if !ok {
		return "", fmt.Errorf("expected to be able to cast to []interface{}")
	}
	var bytes []byte
	for _, elem := range v {
		byteVal, ok := elem.(types.U8)
		if !ok {
			return "", fmt.Errorf("element is not a byte")
		}
		bytes = append(bytes, uint8(byteVal))
	}

	return string(bytes), nil
}
