package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var url = [9]string{
	"./data/buyer.0.0",
	"./data/buyer.1.1",
	"./data/good.0.0",
	"./data/good.1.1",
	"./data/good.2.2",
	"./data/order.0.0",
	"./data/order.0.3",
	"./data/order.1.1",
	"./data/order.2.2",
}

type itemList struct {
	List [][]string
}

// storeItem stroes every item of list from file into memory
func (l itemList) storeItem(i int) itemList {
	// read file
	file, err := os.Open(url[i])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// store item
	scanner := bufio.NewScanner(file)
	// In the given files of good and order, there is some bad smell.
	// transfer " " into "\t"
	if i == 0 || i == 1 {
		for scanner.Scan() {
			l.List = append(l.List, strings.Split(scanner.Text(), "\t"))
		}
	} else {
		r0 := strings.NewReplacer("    ", "\t")
		r1 := strings.NewReplacer("   ", "\t")
		r2 := strings.NewReplacer("  ", "\t")
		r3 := strings.NewReplacer(" ", "\t")
		for scanner.Scan() {
			text := r0.Replace(scanner.Text())
			text = r1.Replace(text)
			text = r2.Replace(text)
			text = r3.Replace(text)
			l.List = append(l.List, strings.Split(text, "\t"))
		}
	}
	return l
}

type order struct {
	orderID    string
	createTime int
	buyerID    string
	goodID     string
	remark     string
	amount     int
}

func normalizeOrder(storeOrder itemList) []order {
	var list []order
	for i := 0; i < len(storeOrder.List); i++ {
		var a order
		var err error
		a.orderID = storeOrder.List[i][0][8:len(storeOrder.List[i][0])]
		a.createTime, err = strconv.Atoi(storeOrder.List[i][1][11:len(storeOrder.List[i][1])])
		if err != nil {
			panic(err)
		}
		a.buyerID = storeOrder.List[i][2][8:len(storeOrder.List[i][2])]
		a.goodID = storeOrder.List[i][3][7:len(storeOrder.List[i][3])]
		// Note that "remark" is not required
		if storeOrder.List[i][4][0:6] == "remark" {
			a.remark = storeOrder.List[i][4][7:len(storeOrder.List[i][4])]
			a.amount, err = strconv.Atoi(storeOrder.List[i][5][7:len(storeOrder.List[i][5])])
			if err != nil {
				panic(err)
			}
		} else {
			a.amount, err = strconv.Atoi(storeOrder.List[i][4][7:len(storeOrder.List[i][4])])
			if err != nil {
				panic(err)
			}
		}
		list = append(list, a)
	}
	return list
}

type good struct {
	goodID      string
	salerID     string
	goodName    string
	description string
	price       float64
}

func normalizeGood(storeGood itemList) []good {
	var list []good
	for i := 0; i < len(storeGood.List); i++ {
		var a good
		var err error
		a.goodID = storeGood.List[i][0][7:len(storeGood.List[i][0])]
		a.salerID = storeGood.List[i][1][8:len(storeGood.List[i][1])]
		a.goodName = storeGood.List[i][2][10:len(storeGood.List[i][2])]
		// Note that "description" is not required.
		if storeGood.List[i][3][0:11] != "description" {
			a.price, err = strconv.ParseFloat(storeGood.List[i][3][6:len(storeGood.List[i][3])], 64)
			if err != nil {
				panic(err)
			}
		} else {
			a.description = storeGood.List[i][3][12:len(storeGood.List[i][3])]
			a.price, err = strconv.ParseFloat(storeGood.List[i][4][6:len(storeGood.List[i][4])], 64)
			if err != nil {
				panic(err)
			}
		}
		list = append(list, a)
	}
	return list
}

type buyer struct {
	buyerID      string
	contactPhone string
	address      string
	buyerName    string
}

func normalizeBuyer(storeBuyer itemList) []buyer {
	var list []buyer
	for i := 0; i < len(storeBuyer.List); i++ {
		var a buyer
		a.buyerID = storeBuyer.List[i][0][8:len(storeBuyer.List[i][0])]
		// Note that both contactPhone and address are not required.
		for j := 1; j < 4; j++ {
			switch storeBuyer.List[i][j][0:4] {
			case "cont":
				a.contactPhone = storeBuyer.List[i][j][13:len(storeBuyer.List[i][j])]
			case "addr":
				a.address = storeBuyer.List[i][j][8:len(storeBuyer.List[i][j])]
			case "buye":
				a.buyerName = storeBuyer.List[i][j][10:len(storeBuyer.List[i][j])]
				break
			default:
				continue
			}
		}
		list = append(list, a)
	}
	return list
}

func main() {

	var storeBuyer itemList
	var storeGood itemList
	var storeOrder itemList

	// store item into list
	for i := 0; i < 2; i++ {
		storeBuyer = storeBuyer.storeItem(i)
	}
	for i := 2; i < 5; i++ {
		storeGood = storeGood.storeItem(i)
	}
	for i := 5; i < 9; i++ {
		storeOrder = storeOrder.storeItem(i)
	}

	orderlist := normalizeOrder(storeOrder)
	goodlist := normalizeGood(storeGood)
	buyerlist := normalizeBuyer(storeBuyer)

	orderID := "627919339"

	for i := 0; i < len(orderlist); i++ {
		var index1, index2 int
		if orderID == orderlist[i].orderID {
			for j := 0; j < len(buyerlist); j++ {
				if buyerlist[j].buyerID == orderlist[i].buyerID {
					index1 = j
					break
				}
			}
			for j := 0; j < len(goodlist); j++ {
				if goodlist[j].goodID == orderlist[i].goodID {
					index2 = j
					break
				}
			}
			fmt.Println("OrderID = 627919339")
			fmt.Printf("Buyername: %v , Gooodname: %v \n", buyerlist[index1].buyerName, goodlist[index2].goodName)
			fmt.Printf("Amount: %v , Price: %v \n", orderlist[i].amount, goodlist[index2].price)
			fmt.Printf("Sum: %v \n", float64(orderlist[i].amount)*goodlist[index2].price)
			break
		}
	}

}
